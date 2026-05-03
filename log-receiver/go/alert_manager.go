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

const (
	// DefaultAlertWindow 默认告警时间窗口
	DefaultAlertWindow = 5 * time.Minute
	// DefaultAlertCooldown 默认告警冷却时间
	DefaultAlertCooldown = 10 * time.Minute
	// DefaultConsecutiveWindows 默认需要连续满足条件的窗口数
	DefaultConsecutiveWindows = 2
	// CheckInterval 状态检查间隔
	CheckInterval = 10 * time.Second
)

type AlertRule struct {
	ID                string
	Name              string
	Level             string
	Count             int
	Window            time.Duration
	WebhookURL        string
	Active            bool
	ConsecutiveCount  int           // 连续满足条件的窗口数
	AlertCooldown     time.Duration // 告警冷却时间
}

type LogEvent struct {
	Timestamp time.Time
	Source    string
	Message   string
}

type AlertState struct {
	RuleID             string
	Events             []LogEvent
	LastAlertTime      time.Time
	AlertActive        bool
	ConsecutiveSuccess int
	mutex              sync.RWMutex
}

type AlertManager struct {
	redisClient *redis.Client
	rules       map[string]AlertRule
	states      map[string]*AlertState
	mutex       sync.RWMutex
	running     bool
}

func NewAlertManager(client *redis.Client, config *Config) *AlertManager {
	am := &AlertManager{
		redisClient: client,
		rules:       make(map[string]AlertRule),
		states:      make(map[string]*AlertState),
	}
	am.loadRules(config)
	return am
}

func (am *AlertManager) loadRules(config *Config) {
	if config != nil && len(config.AlertRules) > 0 {
		log.Printf("Loading %d alert rules from configuration", len(config.AlertRules))
		for _, cfgRule := range config.AlertRules {
			rule := AlertRule{
				ID:               cfgRule.ID,
				Name:             cfgRule.Name,
				Level:            cfgRule.Level,
				Count:            cfgRule.Count,
				Window:           time.Duration(cfgRule.WindowSeconds) * time.Second,
				ConsecutiveCount: cfgRule.Consecutive,
				AlertCooldown:    time.Duration(cfgRule.CooldownSeconds) * time.Second,
				WebhookURL:       cfgRule.WebhookURL,
				Active:           cfgRule.Active,
			}
			
			// Apply defaults for missing values
			if rule.Window == 0 {
				rule.Window = DefaultAlertWindow
			}
			if rule.AlertCooldown == 0 {
				rule.AlertCooldown = DefaultAlertCooldown
			}
			if rule.ConsecutiveCount == 0 {
				rule.ConsecutiveCount = DefaultConsecutiveWindows
			}
			
			am.rules[rule.ID] = rule
			log.Printf("Loaded alert rule: %s (window=%v, count=%d, consecutive=%d, cooldown=%v)",
				rule.Name, rule.Window, rule.Count, rule.ConsecutiveCount, rule.AlertCooldown)
		}
	} else {
		log.Printf("No alert rules in configuration, using defaults")
		am.rules["error_high"] = AlertRule{
			ID:               "error_high",
			Name:             "High Error Rate",
			Level:            "ERROR",
			Count:            5,
			Window:           DefaultAlertWindow,
			ConsecutiveCount: DefaultConsecutiveWindows,
			AlertCooldown:    DefaultAlertCooldown,
			WebhookURL:       "http://localhost:3000/webhook/alert",
			Active:           true,
		}
		
		am.rules["warning_high"] = AlertRule{
			ID:               "warning_high",
			Name:             "High Warning Rate",
			Level:            "WARN",
			Count:            10,
			Window:           5 * time.Minute,
			ConsecutiveCount: 3,
			AlertCooldown:    15 * time.Minute,
			WebhookURL:       "http://localhost:3000/webhook/alert",
			Active:           true,
		}
	}
}

func (am *AlertManager) Start() {
	am.running = true
	ticker := time.NewTicker(CheckInterval)
	defer ticker.Stop()

	for am.running {
		<-ticker.C
		am.cleanupOldEvents()
	}
}

func (am *AlertManager) CheckAlerts(entry LogEntry) {
	am.mutex.RLock()
	defer am.mutex.RUnlock()

	event := LogEvent{
		Timestamp: entry.Timestamp,
		Source:    entry.Source,
		Message:   entry.Message,
	}

	for _, rule := range am.rules {
		if !rule.Active || entry.Level != rule.Level {
			continue
		}

		state, exists := am.states[rule.ID]
		if !exists {
			state = &AlertState{RuleID: rule.ID}
			am.states[rule.ID] = state
		}

		state.mutex.Lock()
		state.Events = append(state.Events, event)
		state.mutex.Unlock()

		am.checkRule(rule, state)
	}
}

func (am *AlertManager) checkRule(rule AlertRule, state *AlertState) {
	state.mutex.Lock()
	defer state.mutex.Unlock()

	now := time.Now()

	// 1. 检查是否在告警冷却期
	if !state.LastAlertTime.IsZero() && now.Sub(state.LastAlertTime) < rule.AlertCooldown {
		return
	}

	// 2. 清理窗口外的事件
	cutoff := now.Add(-rule.Window)
	var recentEvents []LogEvent
	for _, event := range state.Events {
		if event.Timestamp.After(cutoff) {
			recentEvents = append(recentEvents, event)
		}
	}
	state.Events = recentEvents

	// 3. 统计当前窗口的事件数
	currentCount := len(state.Events)
	conditionMet := currentCount >= rule.Count

	// 4. 更新连续满足条件计数
	if conditionMet {
		state.ConsecutiveSuccess++
		log.Printf("Rule %s: current=%d, threshold=%d, consecutive=%d/%d",
			rule.ID, currentCount, rule.Count, state.ConsecutiveSuccess, rule.ConsecutiveCount)
	} else {
		if state.ConsecutiveSuccess > 0 {
			log.Printf("Rule %s: reset consecutive counter (current=%d < threshold=%d)",
				rule.ID, currentCount, rule.Count)
		}
		state.ConsecutiveSuccess = 0
	}

	// 5. 检查是否满足连续窗口条件并触发告警
	if state.ConsecutiveSuccess >= rule.ConsecutiveCount && !state.AlertActive {
		state.AlertActive = true
		state.LastAlertTime = now
		am.triggerAlert(rule, currentCount)
	}
}

func (am *AlertManager) triggerAlert(rule AlertRule, count int) {
	alertData := map[string]interface{}{
		"rule_id":           rule.ID,
		"rule_name":         rule.Name,
		"level":             rule.Level,
		"count":             count,
		"threshold":         rule.Count,
		"window_seconds":    rule.Window.Seconds(),
		"window":            rule.Window.String(),
		"consecutive":       rule.ConsecutiveCount,
		"cooldown_seconds":  rule.AlertCooldown.Seconds(),
		"cooldown":          rule.AlertCooldown.String(),
		"timestamp":         time.Now().Format(time.RFC3339),
	}

	jsonData, _ := json.Marshal(alertData)
	
	log.Printf("=== ALERT TRIGGERED ===")
	log.Printf("Rule: %s (%s)", rule.Name, rule.ID)
	log.Printf("Count: %d (Threshold: %d)", count, rule.Count)
	log.Printf("Window: %s, Consecutive Checks: %d", rule.Window, rule.ConsecutiveCount)
	log.Printf("Cooldown: %s", rule.AlertCooldown)
	log.Printf("=== END ALERT ===")
	
	go func() {
		_, err := http.Post(rule.WebhookURL, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Failed to send alert webhook: %v", err)
		} else {
			log.Printf("Alert webhook sent successfully")
		}
	}()
}

func (am *AlertManager) cleanupOldEvents() {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	now := time.Now()
	for ruleID, state := range am.states {
		rule, exists := am.rules[ruleID]
		if !exists {
			delete(am.states, ruleID)
			continue
		}

		state.mutex.Lock()
		defer state.mutex.Unlock()

		// 清理超过窗口大小的旧事件
		cutoff := now.Add(-rule.Window)
		var filtered []LogEvent
		for _, event := range state.Events {
			if event.Timestamp.After(cutoff) {
				filtered = append(filtered, event)
			}
		}
		state.Events = filtered

		// 如果告警活跃且超过冷却期，重置状态
		if state.AlertActive && now.Sub(state.LastAlertTime) > rule.AlertCooldown {
			log.Printf("Rule %s: alert cooldown expired, resetting state", ruleID)
			state.AlertActive = false
			state.ConsecutiveSuccess = 0
		}
	}
}

func (am *AlertManager) AddRule(rule AlertRule) error {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	if rule.Window == 0 {
		rule.Window = DefaultAlertWindow
	}
	if rule.AlertCooldown == 0 {
		rule.AlertCooldown = DefaultAlertCooldown
	}
	if rule.ConsecutiveCount == 0 {
		rule.ConsecutiveCount = DefaultConsecutiveWindows
	}

	am.rules[rule.ID] = rule
	log.Printf("Added/updated alert rule: %s", rule.ID)
	return nil
}

func (am *AlertManager) RemoveRule(ruleID string) {
	am.mutex.Lock()
	defer am.mutex.Unlock()

	delete(am.rules, ruleID)
	delete(am.states, ruleID)
	log.Printf("Removed alert rule: %s", ruleID)
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

func (am *AlertManager) GetState(ruleID string) (*AlertState, bool) {
	am.mutex.RLock()
	defer am.mutex.RUnlock()

	state, exists := am.states[ruleID]
	return state, exists
}

func (s *AlertState) GetCurrentCount() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.Events)
}

func (s *AlertState) IsAlertActive() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.AlertActive
}
