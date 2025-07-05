-- right now only 1 enum fo all relevant entities; probably in future I will create separate enums for customer, organization, campaign
CREATE TYPE valid_statuses AS ENUM ('active', 'cancelled', 'deleted');

CREATE TABLE IF NOT EXISTS organization (
	id VARCHAR(36) PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	status valid_statuses
);

CREATE TABLE IF NOT EXISTS customer (
	id VARCHAR(36) PRIMARY KEY,
	email VARCHAR(255) NOT NULL,
	name VARCHAR(255) NOT NULL,
    organization_id VARCHAR(36) REFERENCES organization(id),
	status valid_statuses
);

CREATE TABLE IF NOT EXISTS campaign (
	id VARCHAR(36) PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	organization_id VARCHAR(36) REFERENCES organization(id),
	customer_id VARCHAR(36) REFERENCES customer(id),
	status valid_statuses
);

CREATE TABLE IF NOT EXISTS shorturl (
	key VARCHAR(32) PRIMARY KEY,
	is_custom BOOLEAN DEFAULT FALSE,
	target VARCHAR(2048) NOT NULL,
	campaign_id VARCHAR(36) REFERENCES campaign(id),
	customer_id VARCHAR(36) REFERENCES customer(id),
	status valid_statuses
);

CREATE INDEX idx_campaign_customer ON campaign(customer_id);
CREATE INDEX idx_shorturl_customer ON shorturl(customer_id);
CREATE INDEX idx_campaign_org ON campaign(organization_id);
CREATE INDEX idx_shorturl_campaign ON shorturl(campaign_id);

