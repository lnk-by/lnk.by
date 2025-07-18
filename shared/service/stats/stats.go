package stats

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/lnk.by/shared/db"
	"github.com/lnk.by/shared/service"
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
	updateUserAgentBasedStatistics,
}

func Process(ctx context.Context, event Event) error {
	slog.Info("Processing stats", "key", event.Key, "ts", event.Timestamp)
	return db.BulkUpdateWithID(ctx, receivers, event, event.Key)
}

func (e *Event) FieldsPtrs() []any {
	return []any{&e.Key}
}

func (e *Event) FieldsVals() []any {
	return []any{e.Key}
}

func (e *Event) Validate() error {
	if e.Key == "" {
		return errors.New("key is required")
	}

	return nil
}

func (u *Event) Generate() {
}

var (
	CreateTotalSQL  service.CreateSQL[*Event] = "INSERT INTO total_count (key) VALUES ($1)"
	CreateDailySQL  service.CreateSQL[*Event] = "INSERT INTO daily_count (key) VALUES ($1)"
	CreateHourlySQL service.CreateSQL[*Event] = "INSERT INTO hourly_count (key) VALUES ($1)"
)
