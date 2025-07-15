CREATE TABLE IF NOT EXISTS organization (
	id VARCHAR(36) PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	status VARCHAR(16) CHECK (status IN ('active', 'cancelled', 'deleted')),
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS customer (
	id VARCHAR(36) PRIMARY KEY,
	email VARCHAR(255) NOT NULL,
	name VARCHAR(255) NOT NULL,
    organization_id VARCHAR(36) REFERENCES organization(id),
	status VARCHAR(16) CHECK (status IN ('active', 'cancelled', 'deleted')),
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS campaign (
	id VARCHAR(36) PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
    valid_from TIMESTAMPTZ NOT NULL DEFAULT now(),
    valid_until TIMESTAMPTZ NOT NULL DEFAULT '2050-01-01 00:00:00+00',
	organization_id VARCHAR(36) REFERENCES organization(id),
	customer_id VARCHAR(36) REFERENCES customer(id),
	status VARCHAR(16) CHECK (status IN ('active', 'cancelled', 'deleted')),
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS shorturl (
	key VARCHAR(32) PRIMARY KEY,
	is_custom BOOLEAN DEFAULT FALSE,
	target VARCHAR(2048) NOT NULL,
    valid_from TIMESTAMPTZ NOT NULL DEFAULT now(),
    valid_until TIMESTAMPTZ NOT NULL DEFAULT '2050-01-01 00:00:00+00',
	campaign_id VARCHAR(36) REFERENCES campaign(id),
	customer_id VARCHAR(36) REFERENCES customer(id),
	status VARCHAR(16) CHECK (status IN ('active', 'cancelled', 'deleted')),
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_customer_by_organization ON customer(organization_id);

CREATE INDEX IF NOT EXISTS idx_campaign_customer ON campaign(customer_id);

CREATE INDEX IF NOT EXISTS idx_campaign_org ON campaign(organization_id);

CREATE INDEX IF NOT EXISTS idx_shorturl_customer ON shorturl(customer_id);

CREATE INDEX IF NOT EXISTS idx_shorturl_campaign ON shorturl(campaign_id);

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS set_updated_at_organization_trigger ON organization;
CREATE TRIGGER set_updated_at_organization_trigger BEFORE UPDATE ON organization FOR EACH ROW EXECUTE FUNCTION set_updated_at();


DROP TRIGGER IF EXISTS set_updated_at_customer_trigger ON customer;
CREATE TRIGGER set_updated_at_customer_trigger BEFORE UPDATE ON customer FOR EACH ROW EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS set_updated_at_campaign_trigger ON campaign;
CREATE TRIGGER set_updated_at_campaign_trigger BEFORE UPDATE ON campaign FOR EACH ROW EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS set_updated_at_shorturl_trigger ON shorturl;
CREATE TRIGGER set_updated_at_shorturl_trigger BEFORE UPDATE ON shorturl FOR EACH ROW EXECUTE FUNCTION set_updated_at();
