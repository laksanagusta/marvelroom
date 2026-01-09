-- Rollback: Smart Document Linking

-- Remove file_status from work_paper_notes
ALTER TABLE public.work_paper_notes 
DROP CONSTRAINT IF EXISTS work_paper_notes_file_status_check;

DROP INDEX IF EXISTS idx_work_paper_notes_file_status;

ALTER TABLE public.work_paper_notes 
DROP COLUMN IF EXISTS file_status;

-- Remove source_folder_link and last_folder_sync_at from work_papers
ALTER TABLE public.work_papers 
DROP COLUMN IF EXISTS source_folder_link,
DROP COLUMN IF EXISTS last_folder_sync_at;

-- Remove expected_folder_name from work_paper_items
DROP INDEX IF EXISTS idx_work_paper_items_expected_folder_name;

ALTER TABLE public.work_paper_items 
DROP COLUMN IF EXISTS expected_folder_name;
