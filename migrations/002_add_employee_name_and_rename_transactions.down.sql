-- Rename assignee_transactions table back to transactions
ALTER TABLE assignee_transactions RENAME TO transactions;

-- Remove employee_name column from assignees table
ALTER TABLE assignees DROP COLUMN employee_name;