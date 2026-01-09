-- Migration: Smart Document Linking
-- This migration adds fields for automatic Google Drive folder linking:
-- 1. expected_folder_name on work_paper_items (master data)
-- 2. source_folder_link and last_folder_sync_at on work_papers
-- 3. file_status on work_paper_notes for tracking folder status

-- =============================================
-- 1. Add expected_folder_name to work_paper_items
-- =============================================
ALTER TABLE public.work_paper_items 
ADD COLUMN IF NOT EXISTS expected_folder_name VARCHAR(255);

COMMENT ON COLUMN public.work_paper_items.expected_folder_name IS 
'Expected subfolder name in Google Drive to check for files';

-- Add index for faster lookups
CREATE INDEX IF NOT EXISTS idx_work_paper_items_expected_folder_name 
ON public.work_paper_items USING btree (expected_folder_name) 
WHERE (deleted_at IS NULL);

-- =============================================
-- 2. Add source_folder_link and last_folder_sync_at to work_papers
-- =============================================
ALTER TABLE public.work_papers 
ADD COLUMN IF NOT EXISTS source_folder_link TEXT,
ADD COLUMN IF NOT EXISTS last_folder_sync_at TIMESTAMP WITH TIME ZONE;

COMMENT ON COLUMN public.work_papers.source_folder_link IS 
'Root Google Drive folder link for automatic document linking';

COMMENT ON COLUMN public.work_papers.last_folder_sync_at IS 
'Timestamp of last folder synchronization';

-- =============================================
-- 3. Add file_status to work_paper_notes
-- =============================================
ALTER TABLE public.work_paper_notes 
ADD COLUMN IF NOT EXISTS file_status VARCHAR(20) DEFAULT 'pending';

COMMENT ON COLUMN public.work_paper_notes.file_status IS 
'Status of file in folder: pending, found, missing, linked';

-- Add check constraint for valid status values
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'work_paper_notes_file_status_check') THEN
        ALTER TABLE public.work_paper_notes 
        ADD CONSTRAINT work_paper_notes_file_status_check 
        CHECK (file_status IN ('pending', 'found', 'missing', 'linked'));
    END IF;
END $$;

-- Add index for filtering by file_status
CREATE INDEX IF NOT EXISTS idx_work_paper_notes_file_status 
ON public.work_paper_notes USING btree (file_status) 
WHERE (deleted_at IS NULL);
