# 告警系统修复说明

## 问题分析

原告警系统存在以下问题：
1. **时间窗口过短**：可能只有1秒或固定过短的窗口
2. **误报频繁**：瞬时峰值立即触发告警
3. **缺少告警抑制**：没有冷却期，相同问题反复告警
4. **配置不灵活**：规则不可配置

## 修复内容

### 1. 滑动窗口实现（Sliding Window）
- 记录每个日志事件的时间戳
- 只统计当前窗口内的事件数量
- 窗口大小可配置（默认5分钟）
- 自动清理过期事件

### 2. 连续满足条件机制
- 需要连续N次检查都满足条件才触发告警
- 避免瞬时峰值导致误报
- 默认需要连续2次满足条件

### 3. 告警冷却期（Cooldown）
- 告警触发后有冷却期，期间不重复告警
- 默认10分钟冷却期
- 冷却期后自动重置告警状态

### 4. 完全可配置规则
```yaml
alert_rules:
  - id: "error_high"
    name: "High Error Rate"
    level: "ERROR"
    count: 5
    window_seconds: 300
    consecutive: 2
    cooldown_seconds: 600
    webhook_url: "http://localhost:3000/webhook/alert"
    active: true
```

## 改进前后对比

| 特性 | 修复前 | 修复后 |
|------|--------|--------|
| 时间窗口 | 过短或固定 | 5分钟（可配置） |
| 告警触发 | 立即触发 | 连续2次满足才触发 |
| 告警抑制 | 无 | 10分钟冷却期 |
| 规则配置 | 硬编码 | 完全可配置 |
| 误报控制 | 高 | 低（多层过滤） |

## 告警触发流程

```
1. 日志到达 → 添加到窗口事件队列
2. 检查当前窗口内计数是否 >= 阈值
3. 如果满足 → 连续成功计数 +1
4. 否则 → 连续成功计数重置为 0
5. 连续成功 >= 配置次数 → 触发告警
6. 触发告警 → 设置活跃标志，开始冷却期
7. 冷却期间 → 不重复告警
8. 冷却结束 → 重置状态，重新开始检测
```

## 新增字段说明

AlertRule 新增字段：
- `ConsecutiveCount int` - 连续满足条件次数
- `AlertCooldown time.Duration` - 告警冷却期

AlertState 新增字段：
- `Events []LogEvent` - 事件队列
- `ConsecutiveSuccess int` - 连续成功计数

## 配置项说明

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| window_seconds | 300 (5分钟) | 时间窗口大小 |
| count | 5 | 触发阈值 |
| consecutive | 2 | 需要连续满足的次数 |
| cooldown_seconds | 600 (10分钟) | 告警冷却期 |
| level | "ERROR" | 监控的日志级别 |

## 告警Webhook数据

发送到 Webhook 的数据包含：
```json
{
  "rule_id": "error_high",
  "rule_name": "High Error Rate",
  "level": "ERROR",
  "count": 7,
  "threshold": 5,
  "window": "5m0s",
  "consecutive": 2,
  "cooldown": "10m0s",
  "timestamp": "2026-05-03T12:34:56Z"
}
```

## 日志输出示例

```
2026/05/03 12:34:56 Rule error_high: current=4, threshold=5, consecutive=1/2
2026/05/03 12:35:10 Rule error_high: current=6, threshold=5, consecutive=2/2
2026/05/03 12:35:10 === ALERT TRIGGERED ===
2026/05/03 12:35:10 Rule: High Error Rate (error_high)
2026/05/03 12:35:10 Count: 6 (Threshold: 5)
2026/05/03 12:35:10 Window: 5m0s, Consecutive Checks: 2
2026/05/03 12:35:10 Cooldown: 10m0s
2026/05/03 12:35:10 === END ALERT ===
2026/05/03 12:35:10 Alert webhook sent successfully
```

## 测试用场景

### 场景1：瞬时峰值（不触发告警）
- 第1秒：5条ERROR → 满足条件但连续次数不够
- 第10秒：恢复正常 → 连续计数重置
- **结果**：不触发告警

### 场景2：持续问题（触发告警）
- 第1秒：5条ERROR → 连续=1
- 第10秒：6条ERROR → 连续=2 → 触发告警
- **结果**：触发告警
- 冷却期内：不计入重复告警

## API变更

gRPC服务 AlertRule 消息字段更新：
```proto
message AlertRule {
  string id = 1;
  string name = 2;
  string level = 3;
  int32 count = 4;
  int32 window_seconds = 5;  // 新增
  int32 consecutive = 6;     // 新增
  int32 cooldown_seconds = 7;// 新增
  string webhook_url = 8;
  bool active = 9;
}
```
