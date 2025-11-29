-- Rollback migration - rename tables back to original names

-- Rename work_paper_items back to master_lakip_items
ALTER TABLE work_paper_items RENAME TO master_lakip_items;

-- Rename work_paper_notes back to paper_work_items
ALTER TABLE work_paper_notes RENAME TO paper_work_items;

-- Update indexes to reflect the old table names
DROP INDEX IF EXISTS idx_work_paper_items_number_active;
DROP INDEX IF EXISTS idx_work_paper_items_parent_id;
DROP INDEX IF EXISTS idx_work_paper_items_type;
DROP INDEX IF EXISTS idx_work_paper_items_level;
DROP INDEX IF EXISTS idx_work_paper_items_type_parent_number;

CREATE UNIQUE INDEX idx_master_lakip_items_number_active
ON master_lakip_items (number)
WHERE deleted_at IS NULL;

CREATE INDEX idx_master_lakip_items_parent_id ON master_lakip_items(parent_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_master_lakip_items_type ON master_lakip_items(type) WHERE deleted_at IS NULL;
CREATE INDEX idx_master_lakip_items_level ON master_lakip_items(level) WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX idx_master_lakip_items_type_parent_number
ON master_lakip_items(type, parent_id, number)
WHERE deleted_at IS NULL;

-- Update foreign key constraints in paper_work_items
ALTER TABLE paper_work_items DROP CONSTRAINT IF EXISTS work_paper_notes_master_item_id_fkey;
ALTER TABLE paper_work_items ADD CONSTRAINT paper_work_items_master_item_id_fkey
FOREIGN KEY (master_item_id) REFERENCES master_lakip_items(id) ON DELETE RESTRICT;

-- Update indexes for paper_work_items
DROP INDEX IF EXISTS idx_work_paper_notes_paper_work_id;
DROP INDEX IF EXISTS idx_work_paper_notes_master_item_id;

CREATE INDEX idx_paper_work_items_paper_work_id
ON paper_work_items (paper_work_id)
WHERE deleted_at IS NULL;

CREATE INDEX idx_paper_work_items_master_item_id
ON paper_work_items (master_item_id)
WHERE deleted_at IS NULL;

-- Update triggers
DROP TRIGGER IF EXISTS update_work_paper_items_updated_at ON master_lakip_items;
CREATE TRIGGER update_master_lakip_items_updated_at
    BEFORE UPDATE ON master_lakip_items
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_work_paper_notes_updated_at ON paper_work_items;
CREATE TRIGGER update_paper_work_items_updated_at
    BEFORE UPDATE ON paper_work_items
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Restore function and view with original names
DROP FUNCTION IF EXISTS get_work_paper_tree_hierarchy(UUID);
DROP VIEW IF EXISTS v_active_work_paper_hierarchy;

CREATE OR REPLACE FUNCTION get_lakip_tree_hierarchy(parent_id UUID)
RETURNS TABLE (
    id UUID,
    type VARCHAR(10),
    number TEXT,
    statement TEXT,
    explanation TEXT,
    filling_guide TEXT,
    parent_id UUID,
    level INTEGER,
    sort_order INTEGER,
    is_active BOOLEAN,
    path TEXT
) AS $$
WITH RECURSIVE lakip_items AS (
    -- Base case: get root items
    SELECT
        id, type, number, statement, explanation, filling_guide,
        parent_id, level, sort_order, is_active,
        number::text as path
    FROM master_lakip_items
    WHERE (parent_id = $1 OR ($1 IS NULL AND parent_id IS NULL))
    AND deleted_at IS NULL
    AND is_active = true

    UNION ALL

    -- Recursive case: get children
    SELECT
        child.id, child.type, child.number, child.statement, child.explanation, child.filling_guide,
        child.parent_id, child.level, child.sort_order, child.is_active,
        parent.path || '.' || child.number::text as path
    FROM master_lakip_items child
    INNER JOIN lakip_items parent ON child.parent_id = parent.id
    WHERE child.deleted_at IS NULL
    AND child.is_active = true
)
SELECT * FROM lakip_items
ORDER BY level, sort_order, number;
$$ LANGUAGE sql SECURITY DEFINER;

-- Create view for active LAKIP hierarchy
CREATE OR REPLACE VIEW v_active_lakip_hierarchy AS
SELECT * FROM get_lakip_tree_hierarchy(NULL);