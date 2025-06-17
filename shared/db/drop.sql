DROP INDEX IF EXISTS idx_campaign_customer;
DROP INDEX IF EXISTS idx_shorturl_customer;
DROP INDEX IF EXISTS idx_campaign_org;
DROP INDEX IF EXISTS idx_shorturl_campaign;

DROP TABLE IF EXISTS short_url CASCADE;
DROP TABLE IF EXISTS campaign CASCADE;
DROP TABLE IF EXISTS organization_customer CASCADE;
DROP TABLE IF EXISTS organization CASCADE;
DROP TABLE IF EXISTS customer CASCADE;

DROP TYPE valid_statuses CASCADE;
