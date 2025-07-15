package stats

import (
	"context"
	"log/slog"
	"time"

	"github.com/lnk.by/shared/db"
)

type Event struct {
	Key       string    `json:"key"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	Referer   string    `json:"referer,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Language  string    `json:"language,omitempty"`
}

type StatsReceiver interface {
	Receive(Event) (string, error)
}

var receivers []StatsReceiver

func Process(ctx context.Context, event Event) error {
	slog.Info("Processing stats", "key", event.Key, "ts", event.Timestamp)

	return db.BulkUpdateWithID(ctx, receivers, func(r StatsReceiver) (string, error) {
		return r.Receive(event) // where r.Receive returns: `UPDATE stats_table SET col = col + 1 WHERE key = $1`
	}, event.Key)
}
