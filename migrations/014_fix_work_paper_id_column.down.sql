-- Rollback the work_paper_id column fix
-- This migration reverses the changes made in the up migration

DO $$
BEGIN
    -- Remove foreign key constraint if it exists
    ALTER TABLE work_paper_notes DROP CONSTRAINT IF EXISTS work_paper_notes_work_paper_id_fkey;

    -- Remove index if it exists
    DROP INDEX IF EXISTS idx_work_paper_notes_work_paper_id;

    -- Note: We won't rename the column back to paper_work_id or drop it in rollback
    -- as this could cause data loss. The column should be kept for compatibility.

END $$;