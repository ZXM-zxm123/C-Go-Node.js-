package main

import (
    "regexp"
    "strings"
    "time"
)

type LogEntry struct {
    Timestamp time.Time
    Level     string
    Source    string
    Message   string
    Fields    map[string]string
}

type LogParser struct {
    patterns []*regexp.Regexp
    filters  []FilterRule
}

type FilterRule struct {
    Field    string
    Pattern  *regexp.Regexp
    Include  bool
}

func NewLogParser() *LogParser {
    patterns := []*regexp.Regexp{
        regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d+)?(?:Z|[+-]\d{2}:\d{2})?)\s+(\w+)\s+(.*)$`),
        regexp.MustCompile(`^(\d{4}/\d{2}/\d{2}\s+\d{2}:\d{2}:\d{2})\s+(\w+)\s+(.*)$`),
        regexp.MustCompile(`^\[?(\d{2}/\w{3}/\d{4}:\d{2}:\d{2}:\d{2})\]?\s*(\w+)?\s*(.*)$`),
    }

    filters := []FilterRule{
        {Field: "Level", Pattern: regexp.MustCompile(`^(DEBUG|INFO|WARN|ERROR|FATAL)$`), Include: true},
    }

    return &LogParser{
        patterns: patterns,
        filters:  filters,
    }
}

func (p *LogParser) Parse(line string) LogEntry {
    entry := LogEntry{
        Fields: make(map[string]string),
        Level:  "INFO",
        Message: line,
    }

    for _, pattern := range p.patterns {
        matches := pattern.FindStringSubmatch(line)
        if len(matches) >= 3 {
            entry.Timestamp = parseTimestamp(matches[1])
            if matches[2] != "" {
                entry.Level = strings.ToUpper(matches[2])
            }
            entry.Message = matches[3]
            p.extractFields(entry.Message, entry.Fields)
            break
        }
    }

    if entry.Timestamp.IsZero() {
        entry.Timestamp = time.Now()
    }

    return entry
}

func parseTimestamp(ts string) time.Time {
    formats := []string{
        time.RFC3339,
        "2006/01/02 15:04:05",
        "02/Jan/2006:15:04:05",
        time.RFC3339Nano,
    }

    for _, format := range formats {
        if t, err := time.Parse(format, ts); err == nil {
            return t
        }
    }

    return time.Time{}
}

func (p *LogParser) extractFields(message string, fields map[string]string) {
    keyValuePattern := regexp.MustCompile(`(\w+)=("[^"]*"|\S+)`)
    matches := keyValuePattern.FindAllStringSubmatch(message, -1)

    for _, match := range matches {
        key := match[1]
        value := strings.Trim(match[2], "\"")
        fields[key] = value
    }
}

func (p *LogParser) Filter(entry LogEntry) bool {
    for _, rule := range p.filters {
        var value string
        switch rule.Field {
        case "Level":
            value = entry.Level
        case "Source":
            value = entry.Source
        default:
            value = entry.Fields[rule.Field]
        }

        matched := rule.Pattern.MatchString(value)
        if rule.Include && !matched {
            return false
        }
        if !rule.Include && matched {
            return false
        }
    }
    return true
}