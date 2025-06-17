package shortURL

import "github.com/lnk.by/shared/service"

type ShortURL struct {
	Key        string `json:"key"`
	Target     string `json:"target"`
	CampaignID string `json:"campaignId"`
}

func (u *ShortURL) FieldsPtrs() []any {
	return []any{&u.Key, &u.Target, &u.CampaignID}
}

func (u *ShortURL) FieldsVals() []any {
	return []any{u.Key, u.Target, u.CampaignID}
}

var (
	CreateSQL   service.CreateSQL[*ShortURL]   = "INSERT INTO short_url (key, target, campaign_id) VALUES ($1, $2, $3)"
	RetrieveSQL service.RetrieveSQL[*ShortURL] = "SELECT key, target, campaign_id FROM short_url WHERE key = $1"
	UpdateSQL   service.UpdateSQL[*ShortURL]   = "UPDATE short_url SET target = $2, campaign_id = $3 WHERE key = $1"
	DeleteSQL   service.DeleteSQL[*ShortURL]   = "DELETE FROM short_url WHERE key = $1"
	ListSQL     service.ListSQL[*ShortURL]     = "SELECT key, target, campaign_id FROM short_url OFFSET $1 LIMIT $2"
)
