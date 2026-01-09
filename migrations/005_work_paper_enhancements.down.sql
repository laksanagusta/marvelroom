-- Rollback: Work Paper Enhancements

-- =============================================
-- 1. Revert work_paper_items changes
-- =============================================
-- Drop index on topic_id
DROP INDEX IF EXISTS idx_work_paper_items_topic_id;

-- Drop foreign key constraint
ALTER TABLE public.work_paper_items 
    DROP CONSTRAINT IF EXISTS work_paper_items_topic_id_fkey;

-- Drop topic_id column
ALTER TABLE public.work_paper_items 
    DROP COLUMN IF EXISTS topic_id;

-- Rename sequence back to sort_order
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns 
               WHERE table_name = 'work_paper_items' AND column_name = 'sequence') THEN
        ALTER TABLE public.work_paper_items RENAME COLUMN sequence TO sort_order;
    END IF;
END $$;

-- Re-add classification column
ALTER TABLE public.work_paper_items 
    ADD COLUMN IF NOT EXISTS classification text;

-- =============================================
-- 2. Revert work_papers changes
-- =============================================
-- Drop index on name
DROP INDEX IF EXISTS idx_work_papers_name;

-- Drop new unique constraint
ALTER TABLE public.work_papers 
    DROP CONSTRAINT IF EXISTS work_papers_organization_year_semester_name_key;

-- Re-add old unique constraint
ALTER TABLE public.work_papers 
    ADD CONSTRAINT work_papers_organization_year_semester_key 
    UNIQUE (organization_id, year, semester);

-- Drop name column
ALTER TABLE public.work_papers 
    DROP COLUMN IF EXISTS name;

-- =============================================
-- 3. Drop work_paper_topics table
-- =============================================
-- Drop trigger
DROP TRIGGER IF EXISTS update_work_paper_topics_updated_at ON public.work_paper_topics;

-- Drop indexes
DROP INDEX IF EXISTS idx_work_paper_topics_name;
DROP INDEX IF EXISTS idx_work_paper_topics_active;
DROP INDEX IF EXISTS idx_work_paper_topics_deleted_at;

-- Drop table
DROP TABLE IF EXISTS public.work_paper_topics;
