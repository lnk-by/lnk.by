package stats

import (
	"context"
	"log/slog"
	"time"
)

type StatsEvent struct {
	Key       string    `json:"key"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	Referer   string    `json:"referer,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Language  string    `json:"language,omitempty"`
}

func ProcessStatistics(ctx context.Context, event StatsEvent) error {
	slog.Info("Processing stats", "key", event.Key, "ts", event.Timestamp)
	return nil
}
