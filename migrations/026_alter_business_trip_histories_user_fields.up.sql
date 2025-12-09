ALTER TABLE business_trip_histories 
DROP COLUMN changed_by_user_id,
DROP COLUMN changed_by_user_name,
ADD COLUMN changed_by VARCHAR(255);

CREATE INDEX idx_business_trip_histories_changed_by ON business_trip_histories(changed_by);
