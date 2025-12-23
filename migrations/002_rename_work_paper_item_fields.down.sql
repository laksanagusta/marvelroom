-- Rollback migration: Restore statement/explanation/filling_guide columns

ALTER TABLE public.work_paper_items 
    DROP COLUMN IF EXISTS classification,
    DROP COLUMN IF EXISTS desk_instruction;

ALTER TABLE public.work_paper_items 
    ADD COLUMN statement TEXT NOT NULL DEFAULT '',
    ADD COLUMN explanation TEXT,
    ADD COLUMN filling_guide TEXT;

-- Restore comments
COMMENT ON COLUMN public.work_paper_items.statement IS 'Pernyataan/eksistensi yang harus dicek (Kolom D)';
COMMENT ON COLUMN public.work_paper_items.explanation IS 'Penjelasan detail dari pernyataan (Kolom F)';
COMMENT ON COLUMN public.work_paper_items.filling_guide IS 'Petunjuk pengisian untuk item tersebut (Kolom G)';
