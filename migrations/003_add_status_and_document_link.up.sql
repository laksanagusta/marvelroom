-- Migration: Add status and document_link columns to business_trips table
-- Description: Adds status tracking and document link functionality

-- Add status column with default value 'draft'
ALTER TABLE business_trips ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'draft';

-- Add document_link column (optional)
ALTER TABLE business_trips ADD COLUMN document_link TEXT NULL;

-- Add constraint to ensure only valid status values
ALTER TABLE business_trips ADD CONSTRAINT chk_business_trip_status
    CHECK (status IN ('draft', 'ongoing', 'completed', 'canceled'));

-- Create index for status column for better query performance
CREATE INDEX IF NOT EXISTS idx_business_trips_status ON business_trips(status);

-- Add comments for documentation
COMMENT ON COLUMN business_trips.status IS 'Status of the business trip: draft, ongoing, completed, or canceled';
COMMENT ON COLUMN business_trips.document_link IS 'Link to Google Drive or other document storage for this business trip';