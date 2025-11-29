-- Clean up duplicate foreign key constraints in work_paper_notes
-- This migration removes the old constraint name and keeps the correct one

DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'work_paper_notes') THEN
        -- Drop the old constraint if it exists
        ALTER TABLE work_paper_notes DROP CONSTRAINT IF EXISTS paper_work_items_paper_work_id_fkey;

        RAISE NOTICE 'Dropped old foreign key constraint: paper_work_items_paper_work_id_fkey';
    END IF;
END $$;