-- Migration: Create business_trip_tables
-- Description: Create tables for business trip management system

-- Create business_trips table
CREATE TABLE IF NOT EXISTS business_trips (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    activity_purpose VARCHAR(255) NOT NULL,
    destination_city VARCHAR(255) NOT NULL,
    spd_date DATE NOT NULL,
    departure_date DATE NOT NULL,
    return_date DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- Validations
    CONSTRAINT check_business_trip_dates CHECK (
        end_date >= start_date AND
        return_date >= departure_date AND
        spd_date <= departure_date
    )
);

-- Create assignees table
CREATE TABLE IF NOT EXISTS assignees (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    business_trip_id UUID NOT NULL REFERENCES business_trips(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    spd_number VARCHAR(100) NOT NULL,
    employee_id VARCHAR(50) NOT NULL,
    position VARCHAR(255) NOT NULL,
    rank VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- Unique constraint for spd_number within a business trip
    CONSTRAINT unique_spd_number_per_trip UNIQUE(business_trip_id, spd_number)
);

-- Create transactions table
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    assignee_id UUID NOT NULL REFERENCES assignees(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('accommodation', 'transport', 'other', 'allowance')),
    subtype VARCHAR(100),
    amount DECIMAL(12,2) NOT NULL CHECK (amount >= 0),
    total_night INTEGER DEFAULT 1 CHECK (total_night >= 0),
    subtotal DECIMAL(12,2) NOT NULL CHECK (subtotal >= 0),
    description TEXT,
    transport_detail TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_business_trips_dates ON business_trips(start_date, end_date);
CREATE INDEX IF NOT EXISTS idx_business_trips_destination ON business_trips(destination_city);
CREATE INDEX IF NOT EXISTS idx_assignees_business_trip_id ON assignees(business_trip_id);
CREATE INDEX IF NOT EXISTS idx_assignees_employee_id ON assignees(employee_id);
CREATE INDEX IF NOT EXISTS idx_transactions_assignee_id ON transactions(assignee_id);
CREATE INDEX IF NOT EXISTS idx_transactions_type ON transactions(type);

-- Create trigger to update updated_at columns
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
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
COMMENT ON TABLE transactions IS 'Table for storing transaction details for each assignee';

COMMENT ON COLUMN business_trips.activity_purpose IS 'Purpose of the business trip activity';
COMMENT ON COLUMN business_trips.destination_city IS 'Destination city for the business trip';
COMMENT ON COLUMN business_trips.spd_date IS 'Date when SPD (Surat Perjalanan Dinas) was issued';
COMMENT ON COLUMN assignees.spd_number IS 'SPD number for the assignee';
COMMENT ON COLUMN assignees.employee_id IS 'Employee ID (NIP) for the assignee';
COMMENT ON COLUMN assignees.rank IS 'Employee rank/golongan';
COMMENT ON COLUMN transactions.type IS 'Transaction type: accommodation, transport, other, or allowance';
COMMENT ON COLUMN transactions.subtype IS 'Transaction subtype: hotel, flight, train, taxi, daily_allowance, etc.';
COMMENT ON COLUMN transactions.amount IS 'Base amount for the transaction';
COMMENT ON COLUMN transactions.total_night IS 'Number of nights for accommodation type transactions';
COMMENT ON COLUMN transactions.subtotal IS 'Calculated subtotal (amount * total_night for accommodation, or amount for other types)';
COMMENT ON COLUMN transactions.description IS 'Description of what the transaction is for';
COMMENT ON COLUMN transactions.transport_detail IS 'Detailed transport information (routes, addresses, etc.)';