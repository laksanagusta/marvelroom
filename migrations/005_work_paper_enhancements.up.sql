-- Migration: Work Paper Enhancements
-- This migration adds:
-- 1. work_paper_topics table for classifying work paper items with AI context
-- 2. name field to work_papers table with updated unique constraint
-- 3. Updates work_paper_items to use topic_id and sequence

-- =============================================
-- 1. Create work_paper_topics table
-- =============================================
CREATE TABLE IF NOT EXISTS public.work_paper_topics (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name character varying(255) NOT NULL,
    description text,
    template_path text,
    template_version integer DEFAULT 1 NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at timestamp with time zone,
    CONSTRAINT work_paper_topics_pkey PRIMARY KEY (id),
    CONSTRAINT work_paper_topics_name_unique UNIQUE (name)
);

-- Add indexes for work_paper_topics
CREATE INDEX idx_work_paper_topics_name ON public.work_paper_topics USING btree (name) WHERE (deleted_at IS NULL);
CREATE INDEX idx_work_paper_topics_active ON public.work_paper_topics USING btree (is_active) WHERE (deleted_at IS NULL);
CREATE INDEX idx_work_paper_topics_deleted_at ON public.work_paper_topics USING btree (deleted_at);

-- Add comments for work_paper_topics
COMMENT ON TABLE public.work_paper_topics IS 'Topics/classifications for work paper items with AI prompt context';
COMMENT ON COLUMN public.work_paper_topics.name IS 'Unique name of the topic';
COMMENT ON COLUMN public.work_paper_topics.description IS 'Description used as context in AI checks';
COMMENT ON COLUMN public.work_paper_topics.template_path IS 'Path to Excel template file for embedding notes';
COMMENT ON COLUMN public.work_paper_topics.template_version IS 'Version number of the template for tracking changes';

-- Add trigger for updated_at on work_paper_topics
CREATE TRIGGER update_work_paper_topics_updated_at 
    BEFORE UPDATE ON public.work_paper_topics 
    FOR EACH ROW 
    EXECUTE FUNCTION public.update_updated_at_column();

-- =============================================
-- 2. Add name field to work_papers
-- =============================================
ALTER TABLE public.work_papers 
    ADD COLUMN IF NOT EXISTS name character varying(255);

-- Set default values for existing records (if any)
UPDATE public.work_papers 
SET name = CONCAT('Work Paper ', year, ' S', semester) 
WHERE name IS NULL;

-- Make name NOT NULL after setting defaults
ALTER TABLE public.work_papers 
    ALTER COLUMN name SET NOT NULL;

-- Drop old unique constraint
ALTER TABLE public.work_papers 
    DROP CONSTRAINT IF EXISTS work_papers_organization_year_semester_key;

-- Add new unique constraint including name
ALTER TABLE public.work_papers 
    ADD CONSTRAINT work_papers_organization_year_semester_name_key 
    UNIQUE (organization_id, year, semester, name);

-- Add index on name for better query performance
CREATE INDEX IF NOT EXISTS idx_work_papers_name ON public.work_papers USING btree (name) WHERE (deleted_at IS NULL);

-- Add comment
COMMENT ON COLUMN public.work_papers.name IS 'Name of the work paper for identification';

-- =============================================
-- 3. Update work_paper_items to use topic_id and sequence
-- =============================================
-- Note: This assumes existing data has been cleared (TRUNCATE) as per user approval

-- Drop old classification column if exists
ALTER TABLE public.work_paper_items 
    DROP COLUMN IF EXISTS classification;

-- Rename sort_order to sequence if not already done
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.columns 
               WHERE table_name = 'work_paper_items' AND column_name = 'sort_order') THEN
        ALTER TABLE public.work_paper_items RENAME COLUMN sort_order TO sequence;
    END IF;
END $$;

-- Add topic_id column if not exists
ALTER TABLE public.work_paper_items 
    ADD COLUMN IF NOT EXISTS topic_id uuid;

-- Add foreign key constraint to work_paper_topics if not exists
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints 
                   WHERE constraint_name = 'work_paper_items_topic_id_fkey') THEN
        ALTER TABLE public.work_paper_items 
            ADD CONSTRAINT work_paper_items_topic_id_fkey 
            FOREIGN KEY (topic_id) 
            REFERENCES public.work_paper_topics(id) 
            ON DELETE SET NULL;
    END IF;
END $$;

-- Add index on topic_id if not exists
CREATE INDEX IF NOT EXISTS idx_work_paper_items_topic_id ON public.work_paper_items USING btree (topic_id) WHERE (deleted_at IS NULL);

-- Update comments
COMMENT ON COLUMN public.work_paper_items.topic_id IS 'Reference to work paper topic for classification';
COMMENT ON COLUMN public.work_paper_items.sequence IS 'Sequence number for ordering items';
