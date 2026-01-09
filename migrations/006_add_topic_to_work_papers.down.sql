-- Migration Down: Add topic_id to work_papers

ALTER TABLE public.work_papers 
    DROP COLUMN IF EXISTS topic_id;
