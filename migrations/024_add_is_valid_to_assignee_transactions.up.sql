-- Migration: Add is_valid column to assignee_transactions table
-- Description: Adds a boolean column to track document authenticity validation

ALTER TABLE assignee_transactions
ADD COLUMN is_valid BOOLEAN DEFAULT true;

-- Create index for better query performance on is_valid
CREATE INDEX IF NOT EXISTS idx_assignee_transactions_is_valid ON assignee_transactions(is_valid);

-- Add comment for documentation
COMMENT ON COLUMN assignee_transactions.is_valid IS 'Indicates whether the transaction document has been validated as authentic (true) or potentially fraudulent (false)';
