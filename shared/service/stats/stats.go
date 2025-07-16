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

var receivers []func(Event) string = []func(Event) string{
	func(e Event) string {
		return "UPDATE stats_table SET col = col + 1 WHERE key = $1"
	},
}

func Process(ctx context.Context, event Event) error {
	slog.Info("Processing stats", "key", event.Key, "ts", event.Timestamp)

	// where r.Receive returns: `UPDATE stats_table SET col = col + 1 WHERE key = $1`
	return db.BulkUpdateWithID(ctx, receivers, event, event.Key)
}
