package main

import (
    "context"
    "time"

    pb "log_consumer/proto"
)

type server struct {
    pb.UnimplementedLogServiceServer
}

func (s *server) GetMetrics(ctx context.Context, req *pb.MetricsRequest) (*pb.MetricsResponse, error) {
    m := metrics.GetMetrics()

    levelCounts := make(map[string]int64)
    for k, v := range m.LevelCounts {
        levelCounts[k] = v
    }

    sourceCounts := make(map[string]int64)
    for k, v := range m.SourceCounts {
        sourceCounts[k] = v
    }

    return &pb.MetricsResponse{
        TotalCount:   m.TotalCount,
        LevelCounts:  levelCounts,
        SourceCounts: sourceCounts,
        ErrorRate:    m.ErrorRate,
        AvgLatency:   m.AvgLatency,
        Timestamp:    time.Now().Format(time.RFC3339),
    }, nil
}

func (s *server) GetAlerts(ctx context.Context, req *pb.AlertsRequest) (*pb.AlertsResponse, error) {
    rules := alertManager.GetRules()
    pbRules := make([]*pb.AlertRule, 0, len(rules))

    for _, rule := range rules {
        pbRules = append(pbRules, &pb.AlertRule{
            Id:         rule.ID,
            Name:       rule.Name,
            Level:      rule.Level,
            Count:      int32(rule.Count),
            Window:     int32(rule.Window.Minutes()),
            WebhookUrl: rule.WebhookURL,
            Active:     rule.Active,
        })
    }

    return &pb.AlertsResponse{
        Rules: pbRules,
    }, nil
}