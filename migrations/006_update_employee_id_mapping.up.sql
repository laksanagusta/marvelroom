-- Migration: Update employee_id to store external API user ID
-- Description:
-- - employee_id will now store the user ID from external API (response.id)
-- - employee_number will store the NIP/employee number from external API (response.employee_id)
-- - For existing data, we'll keep current employee_id as employee_number (NIP)
-- - employee_id will be NULL until data is fetched from external API

-- Add a comment to clarify the new meaning of columns
COMMENT ON COLUMN assignees.employee_id IS 'User ID from external API system (response.id)';
COMMENT ON COLUMN assignees.employee_number IS 'Employee NIP/number from external API system (response.employee_id)';

-- For existing data, move current employee_id to employee_number if employee_number is NULL
UPDATE assignees
SET employee_number = employee_id
WHERE employee_number IS NULL OR employee_number = '';

-- Note: employee_id will be populated later when data is fetched from external API
-- For now, existing employee_id values will remain and can be used as employee_number (NIP) for API calls