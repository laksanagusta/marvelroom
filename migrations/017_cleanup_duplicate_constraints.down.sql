-- Rollback: Recreate the old foreign key constraint

DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'work_paper_notes') THEN
        -- Add the old constraint back
        ALTER TABLE work_paper_notes
        ADD CONSTRAINT paper_work_items_paper_work_id_fkey
        FOREIGN KEY (work_paper_id) REFERENCES work_papers(id) ON DELETE CASCADE;

        RAISE NOTICE 'Recreated foreign key constraint: paper_work_items_paper_work_id_fkey';
    END IF;
END $$;