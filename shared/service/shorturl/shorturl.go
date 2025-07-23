package shorturl

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofrs/uuid"

	"github.com/lnk.by/shared/service"
	"github.com/lnk.by/shared/service/stats"
	"github.com/lnk.by/shared/utils"
)

type ShortURL struct {
	Key        string       `json:"key"`
	Target     string       `json:"target"`
	ValidFrom  time.Time    `json:"validFrom"`
	ValidUntil time.Time    `json:"validUntil"`
	CampaignID *uuid.UUID   `json:"campaignId"`
	CustomerID *uuid.UUID   `json:"customerId"`
	Status     utils.Status `json:"status"`
	custom     bool
}

func (u *ShortURL) FieldsPtrs() []any {
	return []any{&u.Key, &u.custom, &u.Target, &u.CampaignID, &u.CustomerID, &u.Status}
}

func (u *ShortURL) FieldsVals() []any {
	return []any{u.Key, u.custom, u.Target, u.CampaignID, u.CustomerID, u.Status}
}

var generator *service.Generator

func init() {
	var err error
	generator, err = service.NewDefaultGenerator()
	if err != nil {
		fmt.Printf("Failed to initialize snowflake generator: %v\n", err)
		os.Exit(1)
	}
}

func (u *ShortURL) ParseID(idString string) (string, error) {
	return idString, nil
}

func (c *ShortURL) WithID(key string) {
	c.Key = key
}

func (u *ShortURL) Validate() error {
	// TODO: JWT: do not allow custom key, valid_from and valid_until for anonymous users
	// TODO: JWT+: in future implement limitations on custom key, valid_from and valid_until for authenticated users.
	u.custom = u.Key != ""

	if u.Target == "" {
		return errors.New("target is required")
	}

	return nil
}

func (u *ShortURL) Generate() {
	if !u.custom {
		u.Key = generator.NextBase62ID()
	}

	if u.Status == "" {
		u.Status = utils.StatusActive
	}
}

func (u *ShortURL) MaxAttempts() int {
	if u.custom {
		return 1
	}

	return 10
}

var (
	CreateSQL   service.CreateSQL[*ShortURL]   = "INSERT INTO shorturl (key, is_custom, target, campaign_id, customer_id, status) VALUES ($1, $2, $3, $4, $5, $6)"
	RetrieveSQL service.RetrieveSQL[*ShortURL] = "SELECT key, is_custom, target, campaign_id, customer_id, status FROM shorturl WHERE key = $1 AND status='active' AND now() BETWEEN valid_from AND valid_until"
	UpdateSQL   service.UpdateSQL[*ShortURL]   = "UPDATE shorturl SET target = $2, campaign_id = $3, customer_id = $4, status = $5 WHERE key = $1"
	DeleteSQL   service.DeleteSQL[*ShortURL]   = "DELETE FROM shorturl WHERE key = $1"
	ListSQL     service.ListSQL[*ShortURL]     = "SELECT key, is_custom, target, campaign_id, customer_id, status FROM shorturl WHERE status='active' AND customer_id=$1 OFFSET $2 LIMIT $3"
)

func CreateShortURL(ctx context.Context, requestBody []byte, userID *uuid.UUID) (int, string) {
	url, err := service.Parse[*ShortURL](ctx, requestBody)
	if err != nil {
		return http.StatusBadRequest, http.StatusText(http.StatusBadRequest)
	}
	if url.CustomerID == nil {
		url.CustomerID = userID
	}

	status, body := service.CreateRecord(ctx, CreateSQL, url, 0)

	e := stats.Event{Key: url.Key}
	createSQLs := []service.CreateSQL[*stats.Event]{stats.CreateTotalSQL, stats.CreateDailySQL, stats.CreateHourlySQL, stats.CreateUserAgentSQL, stats.CreateCountrySQL}
	for _, sql := range createSQLs {
		status, body := service.CreateRecord(ctx, sql, &e, 0)
		if status >= http.StatusBadRequest {
			return status, body
		}
	}

	return status, body
}
