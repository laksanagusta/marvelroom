-- Drop work paper note signatures table
DROP TABLE IF EXISTS work_paper_note_signatures;

-- Drop indexes for work_paper_note_signatures
DROP INDEX IF EXISTS idx_work_paper_note_signatures_work_paper_note_id;
DROP INDEX IF EXISTS idx_work_paper_note_signatures_user_id;
DROP INDEX IF EXISTS idx_work_paper_note_signatures_status;
DROP INDEX IF EXISTS idx_work_paper_note_signatures_unique_user_note;

-- Drop trigger for work_paper_note_signatures
DROP TRIGGER IF EXISTS update_work_paper_note_signatures_updated_at ON work_paper_note_signatures;