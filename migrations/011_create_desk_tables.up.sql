-- Create master_lakip_items table
CREATE TABLE master_lakip_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    number VARCHAR(50) NOT NULL,
    statement TEXT NOT NULL,
    explanation TEXT,
    filling_guide TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Create unique index on number (for active items)
CREATE UNIQUE INDEX idx_master_lakip_items_number_active
ON master_lakip_items (number)
WHERE deleted_at IS NULL;

-- Create paper_works table
CREATE TABLE paper_works (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL,
    year INTEGER NOT NULL,
    semester INTEGER NOT NULL CHECK (semester IN (1, 2)),
    status VARCHAR(50) DEFAULT 'draft' CHECK (status IN ('draft', 'in_progress', 'completed')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

-- Create unique index on organization_id, year, semester (for active records)
CREATE UNIQUE INDEX idx_paper_works_organization_year_semester
ON paper_works (organization_id, year, semester)
WHERE deleted_at IS NULL;

-- Create paper_work_items table
CREATE TABLE paper_work_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    paper_work_id UUID NOT NULL,
    master_item_id UUID NOT NULL,
    gdrive_link TEXT,
    is_valid BOOLEAN,
    notes TEXT,
    last_llm_response JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,

    FOREIGN KEY (paper_work_id) REFERENCES paper_works(id) ON DELETE CASCADE,
    FOREIGN KEY (master_item_id) REFERENCES master_lakip_items(id) ON DELETE RESTRICT
);

-- Create indexes for paper_work_items
CREATE INDEX idx_paper_work_items_paper_work_id
ON paper_work_items (paper_work_id)
WHERE deleted_at IS NULL;

CREATE INDEX idx_paper_work_items_master_item_id
ON paper_work_items (master_item_id)
WHERE deleted_at IS NULL;

-- Create trigger to automatically update updated_at column
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for all tables
CREATE TRIGGER update_master_lakip_items_updated_at
    BEFORE UPDATE ON master_lakip_items
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_paper_works_updated_at
    BEFORE UPDATE ON paper_works
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_paper_work_items_updated_at
    BEFORE UPDATE ON paper_work_items
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();