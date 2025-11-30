-- Migration: Add verificators to business trips
-- Description: Adds support for assigning multiple verificators to business trips

-- Create business_trip_verificators table
CREATE TABLE IF NOT EXISTS business_trip_verificators (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    business_trip_id UUID NOT NULL REFERENCES business_trips(id) ON DELETE CASCADE,
    user_id VARCHAR(100) NOT NULL,
    user_name VARCHAR(255) NOT NULL,
    employee_number VARCHAR(50) NOT NULL,
    position VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    verified_at TIMESTAMP NULL,
    verification_notes TEXT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL,
    UNIQUE(business_trip_id, user_id),
    CONSTRAINT chk_verificator_status CHECK (status IN ('pending', 'approved', 'rejected'))
);

-- Add indexes for better performance
CREATE INDEX IF NOT EXISTS idx_business_trip_verificators_deleted_at ON business_trip_verificators(deleted_at);
CREATE INDEX IF NOT EXISTS idx_business_trip_verificators_business_trip_id ON business_trip_verificators(business_trip_id);
CREATE INDEX IF NOT EXISTS idx_business_trip_verificators_user_id ON business_trip_verificators(user_id);
CREATE INDEX IF NOT EXISTS idx_business_trip_verificators_employee_number ON business_trip_verificators(employee_number);
CREATE INDEX IF NOT EXISTS idx_business_trip_verificators_status ON business_trip_verificators(status);

-- Create trigger to automatically update updated_at timestamp
CREATE TRIGGER update_business_trip_verificators_updated_at BEFORE UPDATE ON business_trip_verificators
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Add comments for documentation
COMMENT ON TABLE business_trip_verificators IS 'Table for storing users assigned to verify business trips';
COMMENT ON COLUMN business_trip_verificators.business_trip_id IS 'Reference to the business trip being verified';
COMMENT ON COLUMN business_trip_verificators.user_id IS 'User ID in the system';
COMMENT ON COLUMN business_trip_verificators.user_name IS 'Full name of the verificator';
COMMENT ON COLUMN business_trip_verificators.employee_number IS 'Employee number/NIP of the verificator';
COMMENT ON COLUMN business_trip_verificators.position IS 'Position of the verificator';
COMMENT ON COLUMN business_trip_verificators.status IS 'Verification status (pending, approved, rejected)';
COMMENT ON COLUMN business_trip_verificators.verified_at IS 'Timestamp when verification was completed';
COMMENT ON COLUMN business_trip_verificators.verification_notes IS 'Notes or comments from verificator';