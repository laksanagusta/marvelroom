-- Fix work_papers column names to match application expectations
-- This migration renames Indonesian column names to English names expected by the Go code

DO $$
BEGIN
    -- Check if work_papers table exists
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'work_papers') THEN
        -- Rename columns from Indonesian to English names if they exist

        -- Rename unit_id to organization_id
        IF EXISTS (
            SELECT FROM information_schema.columns
            WHERE table_schema = 'public'
            AND table_name = 'work_papers'
            AND column_name = 'unit_id'
        ) THEN
            ALTER TABLE work_papers RENAME COLUMN unit_id TO organization_id;
            RAISE NOTICE 'Renamed unit_id to organization_id';
        END IF;

        -- Rename tahun to year
        IF EXISTS (
            SELECT FROM information_schema.columns
            WHERE table_schema = 'public'
            AND table_name = 'work_papers'
            AND column_name = 'tahun'
        ) THEN
            ALTER TABLE work_papers RENAME COLUMN tahun TO year;
            RAISE NOTICE 'Renamed tahun to year';
        END IF;
    END IF;
END $$;

-- Update indexes to reflect the new column names
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'work_papers') THEN
        -- Drop old indexes with Indonesian names
        DROP INDEX IF EXISTS idx_kertas_kerja_unit;
        DROP INDEX IF EXISTS idx_kertas_kerja_tahun_semester;

        -- Create new indexes with English names
        CREATE INDEX IF NOT EXISTS idx_work_papers_organization_id
        ON work_papers (organization_id)
        WHERE deleted_at IS NULL;

        CREATE INDEX IF NOT EXISTS idx_work_papers_year_semester
        ON work_papers (year, semester)
        WHERE deleted_at IS NULL;

        CREATE INDEX IF NOT EXISTS idx_work_papers_organization_year_semester
        ON work_papers (organization_id, year, semester)
        WHERE deleted_at IS NULL;

        RAISE NOTICE 'Updated indexes for work_papers table';
    END IF;
END $$;

-- Update foreign key constraints to use new column names
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'work_papers') THEN
        -- Drop old foreign key constraint
        ALTER TABLE work_papers DROP CONSTRAINT IF EXISTS kertas_kerja_unit_id_fkey;

        -- Add new foreign key constraint with new column name
        ALTER TABLE work_papers
        ADD CONSTRAINT work_papers_organization_id_fkey
        FOREIGN KEY (organization_id) REFERENCES units(id) ON DELETE CASCADE;

        RAISE NOTICE 'Updated foreign key constraints for work_papers table';
    END IF;
END $$;

-- Update unique constraints to use new column names
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'work_papers') THEN
        -- Drop old unique constraint
        ALTER TABLE work_papers DROP CONSTRAINT IF EXISTS kertas_kerja_unit_id_tahun_semester_key;

        -- Add new unique constraint with English names
        ALTER TABLE work_papers
        ADD CONSTRAINT work_papers_organization_year_semester_key
        UNIQUE (organization_id, year, semester);

        RAISE NOTICE 'Updated unique constraints for work_papers table';
    END IF;
END $$;

-- Update triggers to use new table name
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'work_papers') THEN
        -- Drop old trigger
        DROP TRIGGER IF EXISTS set_timestamp_kertas_kerja ON work_papers;

        -- Ensure the correct trigger exists
        DROP TRIGGER IF EXISTS update_work_papers_updated_at ON work_papers;
        CREATE TRIGGER update_work_papers_updated_at
            BEFORE UPDATE ON work_papers
            FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

        RAISE NOTICE 'Updated triggers for work_papers table';
    END IF;
END $$;

-- Update primary key constraint name
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'work_papers') THEN
        -- Drop old primary key constraint
        ALTER TABLE work_papers DROP CONSTRAINT IF EXISTS kertas_kerja_pkey;

        -- Add new primary key constraint with English name
        ALTER TABLE work_papers ADD CONSTRAINT work_papers_pkey PRIMARY KEY (id);

        RAISE NOTICE 'Updated primary key constraint for work_papers table';
    END IF;
END $$;