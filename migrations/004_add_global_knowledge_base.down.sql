-- Rollback: Remove is_global flag from knowledge_bases table

DROP INDEX IF EXISTS idx_knowledge_bases_is_global;
ALTER TABLE knowledge_bases DROP COLUMN IF EXISTS is_global;
