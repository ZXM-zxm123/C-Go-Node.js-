package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/go-redis/redis/v8"
	"google.golang.org/grpc"

	pb "log_consumer/proto"
)

var (
	cfg          *Config
	redisClient  *redis.Client
	logParser    *LogParser
	metrics      *MetricsManager
	alertManager *AlertManager
)

func main() {
	configPath := getConfigPath()
	var err error
	cfg, err := LoadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}

	redisAddr := fmt.Sprintf("%s:%d", cfg.GoConsumer.RedisHost, cfg.GoConsumer.RedisPort)
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	ctx := context.Background()
	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis at %s: %v", redisAddr, err)
	}

	createConsumerGroup(ctx)

	logParser := NewLogParser()
	metrics := NewMetricsManager(time.Duration(cfg.GoConsumer.MetricsInterval) * time.Second)
	alertManager := NewAlertManager(redisClient, cfg)

	go consumeStream(ctx)
	go metrics.Start()
	go alertManager.Start()

	grpcPortStr := fmt.Sprintf(":%d", cfg.GoConsumer.GRPCPort)
	lis, err := net.Listen("tcp", grpcPortStr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", grpcPortStr, err)
	}

	s := grpc.NewServer()
	pb.RegisterLogServiceServer(s, &server{})

	log.Printf("gRPC server listening on %s", grpcPortStr)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func getConfigPath() string {
	if path := os.Getenv("CONFIG_PATH"); path != "" {
		return path
	}
	return filepath.Join("..", "config", "config.yaml")
}

func createConsumerGroup(ctx context.Context) {
	err := redisClient.XGroupCreateMkStream(ctx, cfg.GoConsumer.RedisStream, cfg.GoConsumer.ConsumerGroup, "0").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		log.Printf("Failed to create consumer group: %v", err)
	}
}

func consumeStream(ctx context.Context) {
	for {
		entries, err := redisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    cfg.GoConsumer.ConsumerGroup,
			Consumer: cfg.GoConsumer.ConsumerName,
			Streams:  []string{cfg.GoConsumer.RedisStream, ">"},
			Count:    100,
			Block:    0,
		}).Result()

		if err != nil {
			log.Printf("Failed to read stream: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		for _, stream := range entries {
			for _, msg := range stream.Messages {
				processMessage(msg)
				redisClient.XAck(ctx, cfg.GoConsumer.RedisStream, cfg.GoConsumer.ConsumerGroup, msg.ID)
			}
		}
	}
}

func processMessage(msg redis.XMessage) {
    data, ok := msg.Values["data"].(string)
    if !ok {
        return
    }

    source, _ := msg.Values["source"].(string)

    parsed := logParser.Parse(data)
    parsed.Source = source

    if logParser.Filter(parsed) {
        metrics.Record(parsed)
        alertManager.CheckAlerts(parsed)
    }
}