package main

import (
    "bytes"
    "encoding/json"
    "log"
    "net/http"
    "sync"
    "time"

    "github.com/go-redis/redis/v8"
)

type AlertRule struct {
    ID          string
    Name        string
    Level       string
    Count       int
    Window      time.Duration
    WebhookURL  string
    Active      bool
}

type AlertState struct {
    RuleID     string
    Count      int
    LastAlert  time.Time
    Active     bool
}

type AlertManager struct {
    redisClient *redis.Client
    rules       map[string]AlertRule
    states      map[string]*AlertState
    mutex       sync.RWMutex
    running     bool
}

func NewAlertManager(client *redis.Client) *AlertManager {
    am := &AlertManager{
        redisClient: client,
        rules:       make(map[string]AlertRule),
        states:      make(map[string]*AlertState),
    }
    am.loadRules()
    return am
}

func (am *AlertManager) loadRules() {
    am.rules["error_high"] = AlertRule{
        ID:         "error_high",
        Name:       "High Error Rate",
        Level:      "ERROR",
        Count:      5,
        Window:     10 * time.Minute,
        WebhookURL: "http://localhost:3000/webhook/alert",
        Active:     true,
    }
}

func (am *AlertManager) Start() {
    am.running = true
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()

    for am.running {
        <-ticker.C
        am.cleanupOldStates()
    }
}

func (am *AlertManager) CheckAlerts(entry LogEntry) {
    am.mutex.RLock()
    defer am.mutex.RUnlock()

    for _, rule := range am.rules {
        if !rule.Active || entry.Level != rule.Level {
            continue
        }

        state, exists := am.states[rule.ID]
        if !exists {
            state = &AlertState{RuleID: rule.ID}
            am.states[rule.ID] = state
        }

        now := time.Now()
        if now.Sub(state.LastAlert) > rule.Window {
            state.Count = 0
            state.Active = false
        }

        state.Count++
        state.LastAlert = now

        if state.Count >= rule.Count && !state.Active {
            state.Active = true
            am.triggerAlert(rule, state.Count)
        }
    }
}

func (am *AlertManager) triggerAlert(rule AlertRule, count int) {
    alertData := map[string]interface{}{
        "rule_id":   rule.ID,
        "rule_name": rule.Name,
        "level":     rule.Level,
        "count":     count,
        "window":    rule.Window.String(),
        "timestamp": time.Now().ISO8601(),
    }

    jsonData, _ := json.Marshal(alertData)
    go func() {
        _, err := http.Post(rule.WebhookURL, "application/json", bytes.NewBuffer(jsonData))
        if err != nil {
            log.Printf("Failed to send alert webhook: %v", err)
        } else {
            log.Printf("Alert triggered: %s (count=%d)", rule.Name, count)
        }
    }()
}

func (am *AlertManager) cleanupOldStates() {
    am.mutex.Lock()
    defer am.mutex.Unlock()

    now := time.Now()
    for id, state := range am.states {
        if now.Sub(state.LastAlert) > 15*time.Minute {
            delete(am.states, id)
        }
    }
}

func (am *AlertManager) GetRules() []AlertRule {
    am.mutex.RLock()
    defer am.mutex.RUnlock()

    rules := make([]AlertRule, 0, len(am.rules))
    for _, rule := range am.rules {
        rules = append(rules, rule)
    }
    return rules
}