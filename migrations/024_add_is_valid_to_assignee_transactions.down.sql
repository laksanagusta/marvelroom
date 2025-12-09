-- Rollback: Remove is_valid column from assignee_transactions table

-- Drop the index first
DROP INDEX IF EXISTS idx_assignee_transactions_is_valid;

-- Remove the is_valid column
ALTER TABLE assignee_transactions
DROP COLUMN IF EXISTS is_valid;
