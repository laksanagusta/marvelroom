-- Rollback column name changes in work_paper_notes table
-- This migration reverts English column names back to Indonesian names

DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'work_paper_notes') THEN
        -- Rename columns back to Indonesian names if they exist with English names

        -- Rename gdrive_link back to link_gdrive
        IF EXISTS (
            SELECT FROM information_schema.columns
            WHERE table_schema = 'public'
            AND table_name = 'work_paper_notes'
            AND column_name = 'gdrive_link'
        ) THEN
            ALTER TABLE work_paper_notes RENAME COLUMN gdrive_link TO link_gdrive;
        END IF;

        -- Rename is_valid back to hasil_valid
        IF EXISTS (
            SELECT FROM information_schema.columns
            WHERE table_schema = 'public'
            AND table_name = 'work_paper_notes'
            AND column_name = 'is_valid'
        ) THEN
            ALTER TABLE work_paper_notes RENAME COLUMN is_valid TO hasil_valid;
        END IF;

        -- Rename notes back to catatan
        IF EXISTS (
            SELECT FROM information_schema.columns
            WHERE table_schema = 'public'
            AND table_name = 'work_paper_notes'
            AND column_name = 'notes'
        ) THEN
            ALTER TABLE work_paper_notes RENAME COLUMN notes TO catatan;
        END IF;

        -- Add back checked_at column if it doesn't exist
        IF NOT EXISTS (
            SELECT FROM information_schema.columns
            WHERE table_schema = 'public'
            AND table_name = 'work_paper_notes'
            AND column_name = 'checked_at'
        ) THEN
            ALTER TABLE work_paper_notes ADD COLUMN checked_at timestamp with time zone;
        END IF;

        -- Add back kertas_kerja_id column if it doesn't exist
        IF NOT EXISTS (
            SELECT FROM information_schema.columns
            WHERE table_schema = 'public'
            AND table_name = 'work_paper_notes'
            AND column_name = 'kertas_kerja_id'
        ) THEN
            ALTER TABLE work_paper_notes ADD COLUMN kertas_kerja_id UUID NOT NULL DEFAULT gen_random_uuid();
        END IF;

        RAISE NOTICE 'Successfully reverted work_paper_notes columns to Indonesian names';
    END IF;
END $$;

-- Update indexes back to Indonesian names
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'work_paper_notes') THEN
        -- Drop English indexes
        DROP INDEX IF EXISTS idx_work_paper_notes_is_valid;
        DROP INDEX IF EXISTS idx_work_paper_notes_gdrive_link;

        -- Recreate Indonesian indexes
        CREATE INDEX IF NOT EXISTS idx_kertas_kerja_items_checked_at
        ON work_paper_notes (checked_at);

        CREATE INDEX IF NOT EXISTS idx_kertas_kerja_items_kertas_kerja
        ON work_paper_notes (kertas_kerja_id)
        WHERE deleted_at IS NULL;

        CREATE INDEX IF NOT EXISTS idx_kertas_kerja_items_master
        ON work_paper_notes (master_item_id)
        WHERE deleted_at IS NULL;

        CREATE INDEX IF NOT EXISTS idx_kertas_kerja_items_valid
        ON work_paper_notes (hasil_valid)
        WHERE deleted_at IS NULL;
    END IF;
END $$;