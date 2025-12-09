DROP INDEX IF EXISTS idx_business_trip_histories_changed_by;

ALTER TABLE business_trip_histories 
DROP COLUMN changed_by,
ADD COLUMN changed_by_user_id VARCHAR(100),
ADD COLUMN changed_by_user_name VARCHAR(255);

CREATE INDEX idx_business_trip_histories_changed_by_user_id ON business_trip_histories(changed_by_user_id);
