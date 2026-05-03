package main

import (
    "context"
    "log"
    "net"
    "time"

    "github.com/go-redis/redis/v8"
    "google.golang.org/grpc"

    pb "log_consumer/proto"
)

const (
    redisAddr     = "localhost:6379"
    streamName    = "log_stream"
    consumerGroup = "log_consumers"
    consumerName  = "consumer_1"
    grpcPort      = ":50051"
    metricsInterval = 60 * time.Second
)

var (
    redisClient *redis.Client
    logParser   *LogParser
    metrics     *MetricsManager
    alertManager *AlertManager
)

func main() {
    redisClient = redis.NewClient(&redis.Options{
        Addr: redisAddr,
    })

    ctx := context.Background()
    _, err := redisClient.Ping(ctx).Result()
    if err != nil {
        log.Fatalf("Failed to connect to Redis: %v", err)
    }

    createConsumerGroup(ctx)

    logParser = NewLogParser()
    metrics = NewMetricsManager(metricsInterval)
    alertManager = NewAlertManager(redisClient)

    go consumeStream(ctx)
    go metrics.Start()
    go alertManager.Start()

    lis, err := net.Listen("tcp", grpcPort)
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    s := grpc.NewServer()
    pb.RegisterLogServiceServer(s, &server{})

    log.Printf("gRPC server listening on %s", grpcPort)
    if err := s.Serve(lis); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}

func createConsumerGroup(ctx context.Context) {
    err := redisClient.XGroupCreateMkStream(ctx, streamName, consumerGroup, "0").Err()
    if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
        log.Printf("Failed to create consumer group: %v", err)
    }
}

func consumeStream(ctx context.Context) {
    for {
        entries, err := redisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
            Group:    consumerGroup,
            Consumer: consumerName,
            Streams:  []string{streamName, ">"},
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
                redisClient.XAck(ctx, streamName, consumerGroup, msg.ID)
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