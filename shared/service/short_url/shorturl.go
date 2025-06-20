package short_url

import (
	"errors"
	"fmt"
	"os"

	"github.com/lnk.by/shared/service"
	"github.com/lnk.by/shared/utils"
	"github.com/lnk.by/shared/utils/safemap"
)

type ShortURL struct {
	Key        string       `json:"key"`
	Custom     bool         `json:"custom"`
	Target     string       `json:"target"`
	CampaignID string       `json:"campaignId"`
	CustomerID string       `json:"customerId"`
	Status     utils.Status `json:"status"`
}

func (u *ShortURL) FieldsPtrs() []any {
	return []any{&u.Key, &u.Target, &u.CampaignID, &u.CustomerID, &u.Status}
}

func (u *ShortURL) FieldsVals() []any {
	return []any{u.Key, u.Target, u.CampaignID, u.CustomerID, u.Status}
}

var generator *service.Generator
var Custom = safemap.New[string, bool]()

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
	switch {
	case u.Target == "":
		return errors.New("target is required")
	case u.Custom && u.Key == "":
		return errors.New("custom short URL requires key")
	case !u.Custom && u.Key != "":
		return errors.New("short URL should not have a key")
	default:
		return nil
	}
}

func (u *ShortURL) Generate() {
	if u.Key == "" {
		u.Key = service.EncodeBase62(generator.NextID())
	}

	if u.Status == "" {
		u.Status = utils.StatusActive
	}
}

func (u *ShortURL) MaxAttempts() int {
	if u.Custom {
		return 1
	}

	return 10
}

const IdParam = "key"

var (
	CreateSQL   service.CreateSQL[*ShortURL]   = "INSERT INTO short_url (key, target, campaign_id, customer_id, status) VALUES ($1, $2, NULLIF($3, ''), NULLIF($4, ''), $5)"
	RetrieveSQL service.RetrieveSQL[*ShortURL] = "SELECT key, target, COALESCE(campaign_id, ''), COALESCE(customer_id, ''), status FROM short_url WHERE key = $1 AND status='active'"
	UpdateSQL   service.UpdateSQL[*ShortURL]   = "UPDATE short_url SET target = $2, campaign_id = NULLIF($3, ''), customer_id = NULLIF($4, ''), status = $5 WHERE key = $1"
	DeleteSQL   service.DeleteSQL[*ShortURL]   = "DELETE FROM short_url WHERE key = $1"
	ListSQL     service.ListSQL[*ShortURL]     = "SELECT key, target, COALESCE(campaign_id, ''), COALESCE(customer_id, ''), status FROM short_url WHERE status='active' OFFSET $1 LIMIT $2"
)
