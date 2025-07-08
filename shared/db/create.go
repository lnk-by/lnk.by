package db

const (
	createOrganizationTableSQL = `CREATE TABLE IF NOT EXISTS organization (
	id VARCHAR(36) PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	status VARCHAR(16) CHECK (status IN ('active', 'cancelled', 'deleted'))
);`

	createCustomerTableSQL = `CREATE TABLE IF NOT EXISTS customer (
	id VARCHAR(36) PRIMARY KEY,
	email VARCHAR(255) NOT NULL,
	name VARCHAR(255) NOT NULL,
    organization_id VARCHAR(36) REFERENCES organization(id),
	status VARCHAR(16) CHECK (status IN ('active', 'cancelled', 'deleted'))
);`

	createCampaignTableSQL = `CREATE TABLE IF NOT EXISTS campaign (
	id VARCHAR(36) PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	organization_id VARCHAR(36) REFERENCES organization(id),
	customer_id VARCHAR(36) REFERENCES customer(id),
	status VARCHAR(16) CHECK (status IN ('active', 'cancelled', 'deleted'))
);`

	createShortURLTableSQL = `CREATE TABLE IF NOT EXISTS shorturl (
	key VARCHAR(32) PRIMARY KEY,
	is_custom BOOLEAN DEFAULT FALSE,
	target VARCHAR(2048) NOT NULL,
	campaign_id VARCHAR(36) REFERENCES campaign(id),
	customer_id VARCHAR(36) REFERENCES customer(id),
	status VARCHAR(16) CHECK (status IN ('active', 'cancelled', 'deleted'))
);`

	createCustomerByOrganizationIndexSQL = `CREATE INDEX IF NOT EXISTS idx_customer_by_organization ON customer(organization_id);`
	createCampaignByCustomerIndexSQL     = `CREATE INDEX IF NOT EXISTS idx_campaign_customer ON campaign(customer_id);`
	createCampaignByOrganizationIndexSQL = `CREATE INDEX IF NOT EXISTS idx_campaign_org ON campaign(organization_id);`
	createShortURLByCustomerIndexSQL     = `CREATE INDEX IF NOT EXISTS idx_shorturl_customer ON shorturl(customer_id);`
	createShortURLByCampaignIndexSQL     = `CREATE INDEX IF NOT EXISTS idx_shorturl_campaign ON shorturl(campaign_id);`
)
