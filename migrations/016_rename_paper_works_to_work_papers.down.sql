-- Rollback migration: Rename work_papers back to paper_works

DO $$
BEGIN
    -- Rename work_papers back to paper_works if it exists
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'work_papers') THEN
        ALTER TABLE work_papers RENAME TO paper_works;
        RAISE NOTICE 'Successfully rolled back: work_papers renamed to paper_works';
    END IF;
END $$;

-- Update indexes to reflect the old table name
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'paper_works') THEN
        -- Drop new indexes
        DROP INDEX IF EXISTS idx_work_papers_organization_year_semester;
        DROP INDEX IF EXISTS idx_work_papers_status;

        -- Create old indexes
        CREATE INDEX IF NOT EXISTS idx_paper_works_organization_year_semester
        ON paper_works (organization_id, year, semester)
        WHERE deleted_at IS NULL;
    END IF;
END $$;

-- Update foreign key constraints back to reference paper_works
DO $$
BEGIN
    -- Update work_paper_notes foreign key back to reference paper_works
    IF EXISTS (
        SELECT tc.constraint_name
        FROM information_schema.table_constraints tc
        JOIN information_schema.key_column_usage kcu
             ON tc.constraint_name = kcu.constraint_name
        WHERE tc.constraint_schema = 'public'
        AND tc.table_name = 'work_paper_notes'
        AND tc.constraint_type = 'FOREIGN KEY'
        AND kcu.column_name = 'work_paper_id'
    ) THEN
        ALTER TABLE work_paper_notes
        DROP CONSTRAINT IF EXISTS work_paper_notes_work_paper_id_fkey;

        ALTER TABLE work_paper_notes
        ADD CONSTRAINT work_paper_notes_work_paper_id_fkey
        FOREIGN KEY (work_paper_id) REFERENCES paper_works(id) ON DELETE CASCADE;
    END IF;
END $$;

-- Update triggers to use the old table name
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'paper_works') THEN
        -- Drop new trigger
        DROP TRIGGER IF EXISTS update_work_papers_updated_at ON paper_works;

        -- Create old trigger
        CREATE TRIGGER update_paper_works_updated_at
            BEFORE UPDATE ON paper_works
            FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    END IF;
END $$;