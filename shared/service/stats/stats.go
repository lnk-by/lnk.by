package stats

import (
	"context"
	"log/slog"
	"time"
)

type Event struct {
	Key       string    `json:"key"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	Referer   string    `json:"referer,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Language  string    `json:"language,omitempty"`
}

func Process(ctx context.Context, event Event) error {
	slog.Info("Processing stats", "key", event.Key, "ts", event.Timestamp)
	return nil
}
