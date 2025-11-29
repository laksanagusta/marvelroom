-- Rollback migration: Rename work_papers columns back to Indonesian names

DO $$
BEGIN
    -- Check if work_papers table exists
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'work_papers') THEN
        -- Rename columns back to Indonesian names
        IF EXISTS (
            SELECT FROM information_schema.columns
            WHERE table_schema = 'public'
            AND table_name = 'work_papers'
            AND column_name = 'organization_id'
        ) THEN
            ALTER TABLE work_papers RENAME COLUMN organization_id TO unit_id;
        END IF;

        IF EXISTS (
            SELECT FROM information_schema.columns
            WHERE table_schema = 'public'
            AND table_name = 'work_papers'
            AND column_name = 'year'
        ) THEN
            ALTER TABLE work_papers RENAME COLUMN year TO tahun;
        END IF;
    END IF;
END $$;

-- Update indexes back to Indonesian names
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'work_papers') THEN
        -- Drop new indexes
        DROP INDEX IF EXISTS idx_work_papers_organization_id;
        DROP INDEX IF EXISTS idx_work_papers_year_semester;
        DROP INDEX IF EXISTS idx_work_papers_organization_year_semester;

        -- Recreate old indexes
        CREATE INDEX IF NOT EXISTS idx_kertas_kerja_unit
        ON work_papers (unit_id)
        WHERE deleted_at IS NULL;

        CREATE INDEX IF NOT EXISTS idx_kertas_kerja_tahun_semester
        ON work_papers (tahun, semester)
        WHERE deleted_at IS NULL;
    END IF;
END $$;

-- Update foreign key constraints back to Indonesian names
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'work_papers') THEN
        -- Drop new foreign key constraint
        ALTER TABLE work_papers DROP CONSTRAINT IF EXISTS work_papers_organization_id_fkey;

        -- Add old foreign key constraint back
        ALTER TABLE work_papers
        ADD CONSTRAINT kertas_kerja_unit_id_fkey
        FOREIGN KEY (unit_id) REFERENCES units(id) ON DELETE CASCADE;
    END IF;
END $$;

-- Update unique constraints back to Indonesian names
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'work_papers') THEN
        -- Drop new unique constraint
        ALTER TABLE work_papers DROP CONSTRAINT IF EXISTS work_papers_organization_year_semester_key;

        -- Add old unique constraint back
        ALTER TABLE work_papers
        ADD CONSTRAINT kertas_kerja_unit_id_tahun_semester_key
        UNIQUE (unit_id, tahun, semester);
    END IF;
END $$;

-- Update triggers back to Indonesian names
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'work_papers') THEN
        -- Drop new trigger
        DROP TRIGGER IF EXISTS update_work_papers_updated_at ON work_papers;

        -- Add old trigger back
        CREATE TRIGGER set_timestamp_kertas_kerja
            BEFORE UPDATE ON work_papers
            FOR EACH ROW EXECUTE FUNCTION trigger_set_timestamp();
    END IF;
END $$;

-- Update primary key constraint back to Indonesian name
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'work_papers') THEN
        -- Drop new primary key constraint
        ALTER TABLE work_papers DROP CONSTRAINT IF EXISTS work_papers_pkey;

        -- Add old primary key constraint back
        ALTER TABLE work_papers ADD CONSTRAINT kertas_kerja_pkey PRIMARY KEY (id);
    END IF;
END $$;