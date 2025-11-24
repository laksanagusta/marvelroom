-- Add employee_name column to assignees table
ALTER TABLE assignees ADD COLUMN employee_name VARCHAR(255);

-- Rename transactions table to assignee_transactions
ALTER TABLE transactions RENAME TO assignee_transactions;