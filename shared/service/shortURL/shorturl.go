package shortURL

type ShortURL struct {
	Key        string `db:"key" json:"key"`
	Target     string `db:"target" json:"target"`
	CampaignID string `db:"campaign_id" json:"campaignId"`
}
