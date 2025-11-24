-- Down migration: Revert employee_id mapping changes
-- This migration restores the original meaning where employee_id stores NIP

-- Move employee_number back to employee_id for records where employee_id is NULL
UPDATE assignees
SET employee_id = employee_number
WHERE employee_id IS NULL OR employee_id = '';

-- Update comments to reflect original meaning
COMMENT ON COLUMN assignees.employee_id IS 'Employee ID in the company system (NIP)';
COMMENT ON COLUMN assignees.employee_number IS 'Employee number (NIP) - deprecated in this migration';