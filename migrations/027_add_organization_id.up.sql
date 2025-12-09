ALTER TABLE business_trips ADD COLUMN organization_id UUID;

CREATE INDEX idx_business_trips_organization_id ON business_trips(organization_id);
