-- Rename tables to match new naming convention
-- This migration renames master_lakip_items to work_paper_items and paper_work_items to work_paper_notes

-- Check if old tables exist and rename them
DO $$
BEGIN
    -- Rename master_lakip_items to work_paper_items if it exists
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'master_lakip_items') THEN
        ALTER TABLE master_lakip_items RENAME TO work_paper_items;
    END IF;

    -- Rename paper_work_items to work_paper_notes if it exists
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'paper_work_items') THEN
        ALTER TABLE paper_work_items RENAME TO work_paper_notes;

        -- Rename paper_work_id to work_paper_id if column exists
        IF EXISTS (
            SELECT FROM information_schema.columns
            WHERE table_schema = 'public'
            AND table_name = 'work_paper_notes'
            AND column_name = 'paper_work_id'
        ) THEN
            ALTER TABLE work_paper_notes RENAME COLUMN paper_work_id TO work_paper_id;
        END IF;
    END IF;
END $$;

-- Update indexes to reflect the new table names
DROP INDEX IF EXISTS idx_master_lakip_items_number_active;
DROP INDEX IF EXISTS idx_master_lakip_items_parent_id;
DROP INDEX IF EXISTS idx_master_lakip_items_type;
DROP INDEX IF EXISTS idx_master_lakip_items_level;
DROP INDEX IF EXISTS idx_master_lakip_items_type_parent_number;

-- Create indexes for work_paper_items if table exists
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'work_paper_items') THEN
        -- Create regular index first (non-unique to handle potential duplicates)
        CREATE INDEX IF NOT EXISTS idx_work_paper_items_number_active
        ON work_paper_items (number)
        WHERE deleted_at IS NULL AND number IS NOT NULL;

        CREATE INDEX IF NOT EXISTS idx_work_paper_items_parent_id ON work_paper_items(parent_id) WHERE deleted_at IS NULL;
        CREATE INDEX IF NOT EXISTS idx_work_paper_items_type ON work_paper_items(type) WHERE deleted_at IS NULL;
        CREATE INDEX IF NOT EXISTS idx_work_paper_items_level ON work_paper_items(level) WHERE deleted_at IS NULL;

        -- Create regular composite index (non-unique to handle potential duplicates)
        CREATE INDEX IF NOT EXISTS idx_work_paper_items_type_parent_number
        ON work_paper_items(type, parent_id, number)
        WHERE deleted_at IS NULL;
    END IF;
END $$;

-- Update foreign key constraints in work_paper_notes
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'work_paper_notes') THEN
        -- Drop old constraint if exists
        ALTER TABLE work_paper_notes DROP CONSTRAINT IF EXISTS paper_work_items_master_item_id_fkey;
        ALTER TABLE work_paper_notes DROP CONSTRAINT IF EXISTS work_paper_notes_master_item_id_fkey;

        -- Add new foreign key constraint if work_paper_items table exists
        IF EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'work_paper_items') THEN
            ALTER TABLE work_paper_notes
            ADD CONSTRAINT work_paper_notes_master_item_id_fkey
            FOREIGN KEY (master_item_id) REFERENCES work_paper_items(id) ON DELETE RESTRICT;
        END IF;
    END IF;
END $$;

-- Update indexes for work_paper_notes
DROP INDEX IF EXISTS idx_paper_work_items_paper_work_id;
DROP INDEX IF EXISTS idx_paper_work_items_master_item_id;

DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'work_paper_notes') THEN
        -- Create index for work_paper_id if column exists
        IF EXISTS (
            SELECT FROM information_schema.columns
            WHERE table_schema = 'public'
            AND table_name = 'work_paper_notes'
            AND column_name = 'work_paper_id'
        ) THEN
            CREATE INDEX IF NOT EXISTS idx_work_paper_notes_work_paper_id
            ON work_paper_notes (work_paper_id)
            WHERE deleted_at IS NULL;
        END IF;

        CREATE INDEX IF NOT EXISTS idx_work_paper_notes_master_item_id
        ON work_paper_notes (master_item_id)
        WHERE deleted_at IS NULL;
    END IF;
END $$;

-- Update triggers
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'work_paper_items') THEN
        DROP TRIGGER IF EXISTS update_master_lakip_items_updated_at ON work_paper_items;
        DROP TRIGGER IF EXISTS update_work_paper_items_updated_at ON work_paper_items;

        CREATE TRIGGER update_work_paper_items_updated_at
            BEFORE UPDATE ON work_paper_items
            FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    END IF;

    IF EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'work_paper_notes') THEN
        DROP TRIGGER IF EXISTS update_paper_work_items_updated_at ON work_paper_notes;
        DROP TRIGGER IF EXISTS update_work_paper_notes_updated_at ON work_paper_notes;

        CREATE TRIGGER update_work_paper_notes_updated_at
            BEFORE UPDATE ON work_paper_notes
            FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
    END IF;
END $$;

-- Update function and view that reference the old table names
DROP VIEW IF EXISTS v_active_lakip_hierarchy;
DROP FUNCTION IF EXISTS get_lakip_tree_hierarchy(UUID);