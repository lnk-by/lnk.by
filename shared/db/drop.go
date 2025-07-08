package db

const (
	dropCampaignByCustomerIndexSQL     = `DROP INDEX IF EXISTS idx_campaign_customer;`
	dropShortURLByCustomerIndexSQL     = `DROP INDEX IF EXISTS idx_shorturl_customer;`
	dropCampaignByOrganizationIndexSQL = `DROP INDEX IF EXISTS idx_campaign_org;`
	dropShortURLByCampaignIndexSQL     = `DROP INDEX IF EXISTS idx_shorturl_campaign;`
	dropShortURLTableSQL               = `DROP TABLE IF EXISTS shorturl CASCADE;`
	dropCampaignTableSQL               = `DROP TABLE IF EXISTS campaign CASCADE;`
	dropOrganizationTableSQL           = `DROP TABLE IF EXISTS organization CASCADE;`
	dropCustomerTableSQL               = `DROP TABLE IF EXISTS customer CASCADE;`
	//dropOrganizationCustomerTableSQL = `DROP TABLE IF EXISTS organization_customer CASCADE;`
)
