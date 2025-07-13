package shorturl

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/lnk.by/shared/service"
	"github.com/lnk.by/shared/utils"
)

type ShortURL struct {
	Key        string       `json:"key"`
	Target     string       `json:"target"`
	ValidFrom  time.Time    `json:"validFrom"`
	ValidUntil time.Time    `json:"validUntil"`
	CampaignID string       `json:"campaignId"`
	CustomerID string       `json:"customerId"`
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

func (u *ShortURL) WithId(key string) {
	u.Key = key
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
	CreateSQL   service.CreateSQL[*ShortURL]   = "INSERT INTO shorturl (key, is_custom, target, campaign_id, customer_id, status) VALUES ($1, $2, $3, NULLIF($4, ''), NULLIF($5, ''), $6)"
	RetrieveSQL service.RetrieveSQL[*ShortURL] = "SELECT key, is_custom, target, COALESCE(campaign_id, ''), COALESCE(customer_id, ''), status FROM shorturl WHERE key = $1 AND status='active' AND now() >= valid_from AND now() < valid_until"
	UpdateSQL   service.UpdateSQL[*ShortURL]   = "UPDATE shorturl SET target = $2, campaign_id = NULLIF($3, ''), customer_id = NULLIF($4, ''), status = $5 WHERE key = $1"
	DeleteSQL   service.DeleteSQL[*ShortURL]   = "DELETE FROM shorturl WHERE key = $1"
	ListSQL     service.ListSQL[*ShortURL]     = "SELECT key, is_custom, target, COALESCE(campaign_id, ''), COALESCE(customer_id, ''), status FROM shorturl WHERE status='active' OFFSET $1 LIMIT $2"
)
