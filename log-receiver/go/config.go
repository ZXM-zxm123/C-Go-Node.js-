package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	CppReceiver CppReceiverConfig `yaml:"cpp_receiver"`
	GoConsumer  GoConsumerConfig  `yaml:"go_consumer"`
	NodeAPI     NodeAPIConfig     `yaml:"node_api"`
}

type CppReceiverConfig struct {
	UDPPort     int    `yaml:"udp_port"`
	TCPPort     int    `yaml:"tcp_port"`
	QueueSize   int    `yaml:"queue_size"`
	RedisHost   string `yaml:"redis_host"`
	RedisPort   int    `yaml:"redis_port"`
	RedisStream string `yaml:"redis_stream"`
}

type GoConsumerConfig struct {
	RedisHost      string `yaml:"redis_host"`
	RedisPort      int    `yaml:"redis_port"`
	RedisStream    string `yaml:"redis_stream"`
	ConsumerGroup  string `yaml:"consumer_group"`
	ConsumerName   string `yaml:"consumer_name"`
	MetricsInterval int   `yaml:"metrics_interval"`
	GRPCPort       int    `yaml:"grpc_port"`
}

type NodeAPIConfig struct {
	HTTPPort  int    `yaml:"http_port"`
	GRPCHost  string `yaml:"grpc_host"`
	GRPCPort  int    `yaml:"grpc_port"`
	RedisHost string `yaml:"redis_host"`
	RedisPort int    `yaml:"redis_port"`
}

func LoadConfig(path string) (*Config, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for %s: %w", path, err)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("configuration file not found: %s", absPath)
		}
		return nil, fmt.Errorf("failed to read config file %s: %w", absPath, err)
	}

	var cfg Config
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)

	var lineNum int
	type nodeWithLine struct {
		*yaml.Node
	}

	for lineNum = 1; ; lineNum++ {
		var node yaml.Node
		err = decoder.Decode(&node)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, formatYAMLError(absPath, lineNum, err)
		}

		if node.Kind == yaml.DocumentNode {
			if err = node.Decode(&cfg); err != nil {
				return nil, formatYAMLError(absPath, node.Line, err)
			}
		}
	}

	if err = validateConfig(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func formatYAMLError(filePath string, lineNum int, err error) error {
	var errType string
	switch {
	case strings.Contains(err.Error(), "line"):
		errType = "YAML syntax error"
	case strings.Contains(err.Error(), "unknown field"):
		errType = "unknown configuration field"
	case strings.Contains(err.Error(), "cannot unmarshal"):
		errType = "type mismatch error"
	default:
		errType = "configuration parse error"
	}

	return fmt.Errorf(
		"%s in %s (line %d): %w",
		errType,
		filePath,
		lineNum,
		err,
	)
}

func validateConfig(cfg *Config) error {
	if cfg.GoConsumer.RedisHost == "" {
		return fmt.Errorf("validation error: go_consumer.redis_host is required")
	}
	if cfg.GoConsumer.RedisPort <= 0 || cfg.GoConsumer.RedisPort > 65535 {
		return fmt.Errorf("validation error: go_consumer.redis_port must be between 1 and 65535, got %d", cfg.GoConsumer.RedisPort)
	}
	if cfg.GoConsumer.GRPCPort <= 0 || cfg.GoConsumer.GRPCPort > 65535 {
		return fmt.Errorf("validation error: go_consumer.grpc_port must be between 1 and 65535, got %d", cfg.GoConsumer.GRPCPort)
	}
	if cfg.GoConsumer.MetricsInterval <= 0 {
		return fmt.Errorf("validation error: go_consumer.metrics_interval must be positive, got %d", cfg.GoConsumer.MetricsInterval)
	}
	return nil
}