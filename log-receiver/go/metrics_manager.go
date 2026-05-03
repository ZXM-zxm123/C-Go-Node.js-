package main

import (
    "log"
    "sync"
    "time"
)

type Metrics struct {
    TotalCount     int64
    LevelCounts    map[string]int64
    SourceCounts   map[string]int64
    ErrorRate      float64
    AvgLatency     float64
}

type MetricsManager struct {
    interval        time.Duration
    currentMetrics  Metrics
    lastMetrics     Metrics
    mutex           sync.RWMutex
    running         bool
}

func NewMetricsManager(interval time.Duration) *MetricsManager {
    return &MetricsManager{
        interval: interval,
        currentMetrics: Metrics{
            LevelCounts:  make(map[string]int64),
            SourceCounts: make(map[string]int64),
        },
        lastMetrics: Metrics{
            LevelCounts:  make(map[string]int64),
            SourceCounts: make(map[string]int64),
        },
    }
}

func (m *MetricsManager) Record(entry LogEntry) {
    m.mutex.Lock()
    defer m.mutex.Unlock()

    m.currentMetrics.TotalCount++
    m.currentMetrics.LevelCounts[entry.Level]++
    m.currentMetrics.SourceCounts[entry.Source]++
}

func (m *MetricsManager) Start() {
    m.running = true
    ticker := time.NewTicker(m.interval)
    defer ticker.Stop()

    for m.running {
        <-ticker.C
        m.calculateMetrics()
        log.Printf("Metrics updated: Total=%d, Errors=%d", 
            m.lastMetrics.TotalCount, 
            m.lastMetrics.LevelCounts["ERROR"])
    }
}

func (m *MetricsManager) calculateMetrics() {
    m.mutex.Lock()
    defer m.mutex.Unlock()

    m.lastMetrics = m.currentMetrics

    totalErrors := m.lastMetrics.LevelCounts["ERROR"]
    if m.lastMetrics.TotalCount > 0 {
        m.lastMetrics.ErrorRate = float64(totalErrors) / float64(m.lastMetrics.TotalCount) * 100
    }

    m.currentMetrics = Metrics{
        LevelCounts:  make(map[string]int64),
        SourceCounts: make(map[string]int64),
    }
}

func (m *MetricsManager) GetMetrics() Metrics {
    m.mutex.RLock()
    defer m.mutex.RUnlock()
    return m.lastMetrics
}