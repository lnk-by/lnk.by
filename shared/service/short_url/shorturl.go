package short_url

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/lnk.by/shared/service"
	"github.com/lnk.by/shared/utils"
)

type ShortURL struct {
	Key        string       `json:"key"`
	Target     string       `json:"target"`
	CampaignID string       `json:"campaignId"`
	CustomerID string       `json:"customerId"`
	Custom     bool         `json:"custom"`
	Status     utils.Status `json:"status"`
}

func (u *ShortURL) FieldsPtrs() []any {
	return []any{&u.Key, &u.Target, &u.CampaignID, &u.CustomerID, &u.Custom, &u.Status}
}

func (u *ShortURL) FieldsVals() []any {
	return []any{u.Key, u.Target, u.CampaignID, u.CustomerID, u.Custom, u.Status}
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

func (u *ShortURL) WithId(key string) {
	u.Key = key
}

func (u *ShortURL) Validate() error {
	if u.Target == "" {
		return errors.New("Target is required")
	}

	if !u.Custom {
		u.Key = service.EncodeBase62(generator.NextID())
	}
	if u.Status == "" {
		u.Status = utils.StatusActive
	}

	return nil
}

func Create(ctx context.Context, requestedShortURL ShortURL) (int, string) {
	attempts := 10
	requestedShortURL.Custom = false
	if requestedShortURL.Key != "" {
		attempts = 1
		requestedShortURL.Custom = true
	}
	return service.CreateWithRetries(ctx, CreateSQL, &requestedShortURL, attempts)
}

const IdParam = "key"

var (
	CreateSQL   service.CreateSQL[*ShortURL]   = "INSERT INTO short_url (key, target, campaign_id, customer_id, is_custom, status) VALUES ($1, $2, NULLIF($3, ''), NULLIF($4, ''), $5, $6)"
	RetrieveSQL service.RetrieveSQL[*ShortURL] = "SELECT key, target, COALESCE(campaign_id, ''), COALESCE(customer_id, ''), is_custom, status FROM short_url WHERE key = $1 AND status='active'"
	UpdateSQL   service.UpdateSQL[*ShortURL]   = "UPDATE short_url SET target = $2, campaign_id = NULLIF($3, ''), customer_id = NULLIF($4, ''), is_custom = $5, status = $6 WHERE key = $1"
	DeleteSQL   service.DeleteSQL[*ShortURL]   = "DELETE FROM short_url WHERE key = $1"
	ListSQL     service.ListSQL[*ShortURL]     = "SELECT key, target, COALESCE(campaign_id, ''), COALESCE(customer_id, ''), is_custom, status FROM short_url WHERE status='active' OFFSET $1 LIMIT $2"
)
