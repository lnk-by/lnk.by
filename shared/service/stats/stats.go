package stats

import (
	"context"
	"fmt"
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
		return "UPDATE total_count SET total = total + 1 WHERE key = $1"
	},
	func(e Event) string {
		columnName := fmt.Sprintf("day%03d", e.Timestamp.YearDay())
		return fmt.Sprintf("UPDATE daily_count SET %[1]s = %[1]s + 1 WHERE key = $1", columnName)
	},
	func(e Event) string {
		columnName := fmt.Sprintf("hour%02d", e.Timestamp.Hour())
		return fmt.Sprintf("UPDATE hourly_count SET %[1]s = %[1]s + 1 WHERE key = $1", columnName)
	},
}

func Process(ctx context.Context, event Event) error {
	slog.Info("Processing stats", "key", event.Key, "ts", event.Timestamp)
	return db.BulkUpdateWithID(ctx, receivers, event, event.Key)
}
