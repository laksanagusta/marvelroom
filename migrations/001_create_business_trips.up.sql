-- Migration: Create business trips tables
-- Description: Creates tables for business trips, assignees, and transactions

-- Create business_trips table
CREATE TABLE IF NOT EXISTS business_trips (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    activity_purpose VARCHAR(500) NOT NULL,
    destination_city VARCHAR(255) NOT NULL,
    spd_date TIMESTAMP NOT NULL,
    departure_date TIMESTAMP NOT NULL,
    return_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

-- Create assignees table
CREATE TABLE IF NOT EXISTS assignees (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    business_trip_id UUID NOT NULL REFERENCES business_trips(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    spd_number VARCHAR(100) NOT NULL,
    employee_id VARCHAR(100) NOT NULL,
    position VARCHAR(255) NOT NULL,
    rank VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL,
    UNIQUE(business_trip_id, spd_number)
);

-- Create transactions table
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assignee_id UUID NOT NULL REFERENCES assignees(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    subtype VARCHAR(50) NULL,
    amount DECIMAL(15,2) NOT NULL DEFAULT 0.00,
    total_night INTEGER NULL,
    subtotal DECIMAL(15,2) NOT NULL DEFAULT 0.00,
    description TEXT NULL,
    transport_detail TEXT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL,
    CONSTRAINT chk_transaction_type CHECK (type IN ('accommodation', 'transport', 'other', 'allowance'))
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_business_trips_deleted_at ON business_trips(deleted_at);
CREATE INDEX IF NOT EXISTS idx_business_trips_dates ON business_trips(start_date, end_date);
CREATE INDEX IF NOT EXISTS idx_business_trips_destination_city ON business_trips(destination_city);
CREATE INDEX IF NOT EXISTS idx_business_trips_created_at ON business_trips(created_at);
CREATE INDEX IF NOT EXISTS idx_business_trips_activity_purpose ON business_trips USING gin(to_tsvector('english', activity_purpose));

CREATE INDEX IF NOT EXISTS idx_assignees_deleted_at ON assignees(deleted_at);
CREATE INDEX IF NOT EXISTS idx_assignees_business_trip_id ON assignees(business_trip_id);
CREATE INDEX IF NOT EXISTS idx_assignees_employee_id ON assignees(employee_id);
CREATE INDEX IF NOT EXISTS idx_assignees_spd_number ON assignees(spd_number);

CREATE INDEX IF NOT EXISTS idx_transactions_deleted_at ON transactions(deleted_at);
CREATE INDEX IF NOT EXISTS idx_transactions_assignee_id ON transactions(assignee_id);
CREATE INDEX IF NOT EXISTS idx_transactions_type ON transactions(type);
CREATE INDEX IF NOT EXISTS idx_transactions_subtype ON transactions(subtype);

-- Create trigger to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_business_trips_updated_at BEFORE UPDATE ON business_trips
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_assignees_updated_at BEFORE UPDATE ON assignees
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_transactions_updated_at BEFORE UPDATE ON transactions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Add comments for documentation
COMMENT ON TABLE business_trips IS 'Main table for storing business trip information';
COMMENT ON TABLE assignees IS 'Table for storing employees assigned to business trips';
COMMENT ON TABLE transactions IS 'Table for storing financial transactions related to business trips';

COMMENT ON COLUMN business_trips.start_date IS 'Start date of the business trip';
COMMENT ON COLUMN business_trips.end_date IS 'End date of the business trip';
COMMENT ON COLUMN business_trips.activity_purpose IS 'Purpose of the business trip';
COMMENT ON COLUMN business_trips.destination_city IS 'Destination city for the business trip';
COMMENT ON COLUMN business_trips.spd_date IS 'SPD (Surat Perintah Dinas) date';
COMMENT ON COLUMN business_trips.departure_date IS 'Actual departure date';
COMMENT ON COLUMN business_trips.return_date IS 'Actual return date';

COMMENT ON COLUMN assignees.spd_number IS 'SPD number for the assignee';
COMMENT ON COLUMN assignees.employee_id IS 'Employee ID in the company system';
COMMENT ON COLUMN assignees.position IS 'Position of the employee';
COMMENT ON COLUMN assignees.rank IS 'Rank or grade of the employee';

COMMENT ON COLUMN transactions.type IS 'Type of transaction (accommodation, transport, other, allowance)';
COMMENT ON COLUMN transactions.subtype IS 'Subtype of transaction (hotel, flight, train, etc.)';
COMMENT ON COLUMN transactions.amount IS 'Amount per unit';
COMMENT ON COLUMN transactions.total_night IS 'Total number of nights (for accommodation)';
COMMENT ON COLUMN transactions.subtotal IS 'Total amount for this transaction';
COMMENT ON COLUMN transactions.transport_detail IS 'Details for transportation (flight number, etc.)';