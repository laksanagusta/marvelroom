-- Migration: Rename work_paper_items fields from statement/explanation/filling_guide to classification/desk_instruction
-- Drop existing columns and add new ones

ALTER TABLE public.work_paper_items 
    DROP COLUMN IF EXISTS statement,
    DROP COLUMN IF EXISTS explanation,
    DROP COLUMN IF EXISTS filling_guide;

ALTER TABLE public.work_paper_items 
    ADD COLUMN classification VARCHAR(255),
    ADD COLUMN desk_instruction TEXT NOT NULL DEFAULT '';

-- Update comments
COMMENT ON COLUMN public.work_paper_items.classification IS 'Kategori/klasifikasi untuk grouping catatan';
COMMENT ON COLUMN public.work_paper_items.desk_instruction IS 'Instruksi desk yang akan dikirim sebagai prompt ke LLM';
