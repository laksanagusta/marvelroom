-- Drop triggers
DROP TRIGGER IF EXISTS update_master_lakip_items_updated_at ON master_lakip_items;
DROP TRIGGER IF EXISTS update_paper_works_updated_at ON paper_works;
DROP TRIGGER IF EXISTS update_paper_work_items_updated_at ON paper_work_items;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_paper_work_items_master_item_id;
DROP INDEX IF EXISTS idx_paper_work_items_paper_work_id;
DROP INDEX IF EXISTS idx_paper_works_organization_year_semester;
DROP INDEX IF EXISTS idx_master_lakip_items_number_active;

-- Drop tables in reverse order of creation (to respect foreign key constraints)
DROP TABLE IF EXISTS paper_work_items;
DROP TABLE IF EXISTS paper_works;
DROP TABLE IF EXISTS master_lakip_items;