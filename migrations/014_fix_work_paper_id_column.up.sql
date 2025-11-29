-- Fix work_paper_id column in work_paper_notes table
-- This migration ensures the work_paper_id column exists and is properly configured

DO $$
BEGIN
    -- Check if work_paper_notes table exists
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'work_paper_notes') THEN
        -- Check if work_paper_id column exists, if not add it
        IF NOT EXISTS (
            SELECT FROM information_schema.columns
            WHERE table_schema = 'public'
            AND table_name = 'work_paper_notes'
            AND column_name = 'work_paper_id'
        ) THEN
            -- Check if paper_work_id exists (old name) and rename it
            IF EXISTS (
                SELECT FROM information_schema.columns
                WHERE table_schema = 'public'
                AND table_name = 'work_paper_notes'
                AND column_name = 'paper_work_id'
            ) THEN
                ALTER TABLE work_paper_notes RENAME COLUMN paper_work_id TO work_paper_id;
            ELSE
                -- If neither column exists, add the work_paper_id column
                ALTER TABLE work_paper_notes ADD COLUMN work_paper_id UUID NOT NULL DEFAULT gen_random_uuid();

                -- Note: This is a fallback scenario. The default value should be replaced with actual foreign key references.
                -- You may need to run an UPDATE script to populate this column with proper work_paper_id values.
                RAISE WARNING 'Added work_paper_id column with default UUID values. Please update these with proper foreign key references.';
            END IF;
        END IF;

        -- Add foreign key constraint if it doesn't exist
        IF NOT EXISTS (
            SELECT tc.constraint_name
            FROM information_schema.table_constraints tc
            WHERE tc.constraint_schema = 'public'
            AND tc.table_name = 'work_paper_notes'
            AND tc.constraint_name = 'work_paper_notes_work_paper_id_fkey'
        ) THEN
            -- Only add constraint if paper_works table exists
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

        -- Create index for work_paper_id if it doesn't exist
        CREATE INDEX IF NOT EXISTS idx_work_paper_notes_work_paper_id
        ON work_paper_notes (work_paper_id)
        WHERE deleted_at IS NULL;

    ELSE
        -- If work_paper_notes doesn't exist, check if paper_work_items exists and rename it
        IF EXISTS (SELECT FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'paper_work_items') THEN
            ALTER TABLE paper_work_items RENAME TO work_paper_notes;

            -- Rename paper_work_id to work_paper_id if it exists
            IF EXISTS (
                SELECT FROM information_schema.columns
                WHERE table_schema = 'public'
                AND table_name = 'work_paper_notes'
                AND column_name = 'paper_work_id'
            ) THEN
                ALTER TABLE work_paper_notes RENAME COLUMN paper_work_id TO work_paper_id;
            END IF;
        END IF;
    END IF;
END $$;