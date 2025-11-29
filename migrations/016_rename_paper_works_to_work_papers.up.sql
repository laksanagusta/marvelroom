-- Rename paper_works table to work_papers
-- This migration ensures the main work paper table follows the new naming convention

DO $$
BEGIN
    -- Rename paper_works to work_papers if it exists
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'paper_works') THEN
        ALTER TABLE paper_works RENAME TO work_papers;
        RAISE NOTICE 'Successfully renamed paper_works to work_papers';
    ELSIF EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'kertas_kerja') THEN
        -- If the old Indonesian name still exists, rename it to work_papers
        ALTER TABLE kertas_kerja RENAME TO work_papers;
        RAISE NOTICE 'Successfully renamed kertas_kerja to work_papers';
    END IF;
END $$;

-- Update indexes to reflect the new table name
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'work_papers') THEN
        -- Drop old indexes if they exist
        DROP INDEX IF EXISTS idx_paper_works_organization_year_semester;

        -- Create new indexes with proper names if they don't exist
        CREATE INDEX IF NOT EXISTS idx_work_papers_organization_year_semester
        ON work_papers (organization_id, year, semester)
        WHERE deleted_at IS NULL;

        CREATE INDEX IF NOT EXISTS idx_work_papers_status
        ON work_papers (status)
        WHERE deleted_at IS NULL;
    END IF;
END $$;

-- Update foreign key constraints that reference the old table name
DO $$
BEGIN
    -- Update work_paper_notes foreign key if it still references paper_works
    IF EXISTS (
        SELECT tc.constraint_name, tc.table_name
        FROM information_schema.table_constraints tc
        JOIN information_schema.key_column_usage kcu
             ON tc.constraint_name = kcu.constraint_name
        WHERE tc.constraint_schema = 'public'
        AND tc.table_name = 'work_paper_notes'
        AND tc.constraint_type = 'FOREIGN KEY'
        AND kcu.column_name = 'work_paper_id'
        AND tc.constraint_name = 'work_paper_notes_work_paper_id_fkey'
    ) THEN
        -- Check if constraint exists and needs to be updated
        ALTER TABLE work_paper_notes
        DROP CONSTRAINT IF EXISTS work_paper_notes_work_paper_id_fkey;

        ALTER TABLE work_paper_notes
        ADD CONSTRAINT work_paper_notes_work_paper_id_fkey
        FOREIGN KEY (work_paper_id) REFERENCES work_papers(id) ON DELETE CASCADE;

        RAISE NOTICE 'Updated work_paper_notes foreign key to reference work_papers table';
    END IF;
END $$;

-- Update triggers to use the new table name
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'work_papers') THEN
        -- Drop old trigger if it exists
        DROP TRIGGER IF EXISTS update_paper_works_updated_at ON work_papers;

        -- Create new trigger with proper name
        CREATE TRIGGER update_work_papers_updated_at
            BEFORE UPDATE ON work_papers
            FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

        RAISE NOTICE 'Updated triggers for work_papers table';
    END IF;
END $$;