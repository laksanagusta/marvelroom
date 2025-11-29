-- Drop views and functions
DROP VIEW IF EXISTS v_active_lakip_hierarchy;
DROP FUNCTION IF EXISTS get_lakip_tree_hierarchy(UUID);

-- Drop indexes
DROP INDEX IF EXISTS idx_master_lakip_items_type_parent_number;
DROP INDEX IF EXISTS idx_master_lakip_items_level;
DROP INDEX IF EXISTS idx_master_lakip_items_type;
DROP INDEX IF EXISTS idx_master_lakip_items_parent_id;

-- Remove foreign key constraint
ALTER TABLE master_lakip_items DROP CONSTRAINT IF EXISTS master_lakip_items_parent_id_fkey;

-- Remove columns for tree structure
ALTER TABLE master_lakip_items
DROP COLUMN IF EXISTS type,
DROP COLUMN IF EXISTS parent_id,
DROP COLUMN IF EXISTS level,
DROP COLUMN IF EXISTS sort_order;

-- Restore original structure (if needed)
ALTER TABLE master_lakip_items
ADD COLUMN IF NOT EXISTS number VARCHAR(50) NOT NULL,
ADD COLUMN IF NOT EXISTS statement TEXT NOT NULL,
ADD COLUMN IF NOT EXISTS explanation TEXT,
ADD COLUMN IF NOT EXISTS filling_guide TEXT,
ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT true,
ADD COLUMN IF NOT EXISTS created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP NULL;

-- Restore unique index
CREATE UNIQUE INDEX IF NOT EXISTS idx_master_lakip_items_number_active
ON master_lakip_items (number)
WHERE deleted_at IS NULL;