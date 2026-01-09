-- Migration: Add topic_id to work_papers
-- This migration adds a topic_id column to the work_papers table to enforce a 1:1 relationship with topics.

-- 1. Add topic_id column
ALTER TABLE public.work_papers 
    ADD COLUMN IF NOT EXISTS topic_id uuid;

-- 2. Add foreign key constraint
ALTER TABLE public.work_papers 
    ADD CONSTRAINT work_papers_topic_id_fkey 
    FOREIGN KEY (topic_id) 
    REFERENCES public.work_paper_topics(id) 
    ON DELETE SET NULL;

-- 3. Add index on topic_id
CREATE INDEX IF NOT EXISTS idx_work_papers_topic_id ON public.work_papers USING btree (topic_id) WHERE (deleted_at IS NULL);

-- 4. Comment
COMMENT ON COLUMN public.work_papers.topic_id IS 'Reference to the specific topic this work paper covers';
