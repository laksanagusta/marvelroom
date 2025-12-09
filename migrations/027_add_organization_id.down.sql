DROP INDEX IF EXISTS idx_business_trips_organization_id;

ALTER TABLE business_trips DROP COLUMN organization_id;
