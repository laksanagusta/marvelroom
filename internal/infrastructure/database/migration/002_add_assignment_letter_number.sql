-- Add assignment_letter_number column to business_trips table
ALTER TABLE business_trips ADD COLUMN IF NOT EXISTS assignment_letter_number VARCHAR(100);
