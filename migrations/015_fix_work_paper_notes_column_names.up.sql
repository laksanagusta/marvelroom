-- Fix work_paper_notes column naming to match entity expectations
-- This migration renames Indonesian column names to English names expected by the Go code

DO $$
BEGIN
    -- Check if work_paper_notes table exists
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'work_paper_notes') THEN
        -- Rename columns from Indonesian to English names if they exist with Indonesian names

        -- Rename link_gdrive to gdrive_link
        IF EXISTS (
            SELECT FROM information_schema.columns
            WHERE table_schema = 'public'
            AND table_name = 'work_paper_notes'
            AND column_name = 'link_gdrive'
        ) THEN
            ALTER TABLE work_paper_notes RENAME COLUMN link_gdrive TO gdrive_link;
        END IF;

        -- Rename hasil_valid to is_valid
        IF EXISTS (
            SELECT FROM information_schema.columns
            WHERE table_schema = 'public'
            AND table_name = 'work_paper_notes'
            AND column_name = 'hasil_valid'
        ) THEN
            ALTER TABLE work_paper_notes RENAME COLUMN hasil_valid TO is_valid;
        END IF;

        -- Rename catatan to notes
        IF EXISTS (
            SELECT FROM information_schema.columns
            WHERE table_schema = 'public'
            AND table_name = 'work_paper_notes'
            AND column_name = 'catatan'
        ) THEN
            ALTER TABLE work_paper_notes RENAME COLUMN catatan TO notes;
        END IF;

        -- Drop checked_at column if it exists (not needed in entity)
        IF EXISTS (
            SELECT FROM information_schema.columns
            WHERE table_schema = 'public'
            AND table_name = 'work_paper_notes'
            AND column_name = 'checked_at'
        ) THEN
            ALTER TABLE work_paper_notes DROP COLUMN checked_at;
        END IF;

        -- Drop kertas_kerja_id column if it exists (duplicate of work_paper_id)
        IF EXISTS (
            SELECT FROM information_schema.columns
            WHERE table_schema = 'public'
            AND table_name = 'work_paper_notes'
            AND column_name = 'kertas_kerja_id'
        ) THEN
            ALTER TABLE work_paper_notes DROP COLUMN kertas_kerja_id;
        END IF;

        RAISE NOTICE 'Successfully renamed work_paper_notes columns to match entity expectations';
    END IF;
END $$;

-- Update indexes to reflect the new column names
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'work_paper_notes') THEN
        -- Drop old indexes if they exist
        DROP INDEX IF EXISTS idx_kertas_kerja_items_checked_at;
        DROP INDEX IF EXISTS idx_kertas_kerja_items_kertas_kerja;
        DROP INDEX IF EXISTS idx_kertas_kerja_items_master;
        DROP INDEX IF EXISTS idx_kertas_kerja_items_valid;

        -- Create new indexes with proper names if they don't exist
        CREATE INDEX IF NOT EXISTS idx_work_paper_notes_is_valid
        ON work_paper_notes (is_valid)
        WHERE deleted_at IS NULL;

        CREATE INDEX IF NOT EXISTS idx_work_paper_notes_gdrive_link
        ON work_paper_notes (gdrive_link)
        WHERE deleted_at IS NULL AND gdrive_link IS NOT NULL;
    END IF;
END $$;

-- Update foreign key constraints to use proper names
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'work_paper_notes') THEN
        -- Drop old constraint if it exists
        ALTER TABLE work_paper_notes DROP CONSTRAINT IF EXISTS kertas_kerja_items_kertas_kerja_id_fkey;

        -- Ensure the work_paper_id foreign key constraint exists
        IF NOT EXISTS (
            SELECT tc.constraint_name
            FROM information_schema.table_constraints tc
            WHERE tc.constraint_schema = 'public'
            AND tc.table_name = 'work_paper_notes'
            AND tc.constraint_name = 'work_paper_notes_work_paper_id_fkey'
        ) THEN
            -- Check if paper_works table exists
            IF EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'paper_works') THEN
                ALTER TABLE work_paper_notes
                ADD CONSTRAINT work_paper_notes_work_paper_id_fkey
                FOREIGN KEY (work_paper_id) REFERENCES paper_works(id) ON DELETE CASCADE;
            ELSIF EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'work_papers') THEN
                ALTER TABLE work_paper_notes
                ADD CONSTRAINT work_paper_notes_work_paper_id_fkey
                FOREIGN KEY (work_paper_id) REFERENCES work_papers(id) ON DELETE CASCADE;
            END IF;
        END IF;
    END IF;
END $$;