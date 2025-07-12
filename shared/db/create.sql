CREATE TABLE IF NOT EXISTS organization (
	id UUID PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	status VARCHAR(16) CHECK (status IN ('active', 'cancelled', 'deleted'))
);

CREATE TABLE IF NOT EXISTS customer (
	id UUID PRIMARY KEY,
	email VARCHAR(255) NOT NULL,
	name VARCHAR(255) NOT NULL,
    organization_id UUID REFERENCES organization(id),
	status VARCHAR(16) CHECK (status IN ('active', 'cancelled', 'deleted'))
);

CREATE TABLE IF NOT EXISTS campaign (
	id UUID PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	organization_id UUID REFERENCES organization(id),
	customer_id UUID REFERENCES customer(id),
	status VARCHAR(16) CHECK (status IN ('active', 'cancelled', 'deleted'))
);

CREATE TABLE IF NOT EXISTS shorturl (
	key VARCHAR(32) PRIMARY KEY,
	is_custom BOOLEAN DEFAULT FALSE,
	target VARCHAR(2048) NOT NULL,
	campaign_id UUID REFERENCES campaign(id),
	customer_id UUID REFERENCES customer(id),
	status VARCHAR(16) CHECK (status IN ('active', 'cancelled', 'deleted'))
);

CREATE INDEX IF NOT EXISTS idx_customer_by_organization ON customer(organization_id);

CREATE INDEX IF NOT EXISTS idx_campaign_customer ON campaign(customer_id);

CREATE INDEX IF NOT EXISTS idx_campaign_org ON campaign(organization_id);

CREATE INDEX IF NOT EXISTS idx_shorturl_customer ON shorturl(customer_id);

CREATE INDEX IF NOT EXISTS idx_shorturl_campaign ON shorturl(campaign_id);