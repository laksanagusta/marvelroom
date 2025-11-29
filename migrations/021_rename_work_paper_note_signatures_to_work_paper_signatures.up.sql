-- Create new work_paper_signatures table
CREATE TABLE work_paper_signatures (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    work_paper_id UUID NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    user_name VARCHAR(255) NOT NULL,
    user_email VARCHAR(255),
    user_role VARCHAR(100),
    signature_data JSONB,
    signature_type VARCHAR(50) DEFAULT 'digital' CHECK (signature_type IN ('digital', 'manual', 'approval')),
    status VARCHAR(50) DEFAULT 'pending' CHECK (status IN ('pending', 'signed', 'rejected')),
    notes TEXT,
    signed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,

    -- Foreign key constraint to work_papers
    FOREIGN KEY (work_paper_id) REFERENCES work_papers(id) ON DELETE CASCADE
);

-- Create indexes for work_paper_signatures
CREATE INDEX idx_work_paper_signatures_work_paper_id
ON work_paper_signatures (work_paper_id)
WHERE deleted_at IS NULL;

CREATE INDEX idx_work_paper_signatures_user_id
ON work_paper_signatures (user_id)
WHERE deleted_at IS NULL;

CREATE INDEX idx_work_paper_signatures_status
ON work_paper_signatures (status)
WHERE deleted_at IS NULL;

-- Create unique constraint to prevent duplicate signatures by same user on same work paper
CREATE UNIQUE INDEX idx_work_paper_signatures_unique_user_paper
ON work_paper_signatures (work_paper_id, user_id)
WHERE deleted_at IS NULL AND status IN ('pending', 'signed', 'rejected');

-- Create trigger to automatically update updated_at column
CREATE TRIGGER update_work_paper_signatures_updated_at
    BEFORE UPDATE ON work_paper_signatures
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Migrate data from work_paper_note_signatures to work_paper_signatures if needed
-- This migration moves signatures from individual notes to the work paper level
INSERT INTO work_paper_signatures (work_paper_id, user_id, user_name, user_email, user_role, signature_data, signature_type, status, notes, signed_at, created_at, updated_at)
SELECT
    wpn.work_paper_id,
    wpns.user_id,
    wpns.user_name,
    wpns.user_email,
    wpns.user_role,
    wpns.signature_data,
    wpns.signature_type,
    wpns.status,
    wpns.notes,
    wpns.signed_at,
    wpns.created_at,
    wpns.updated_at
FROM work_paper_note_signatures wpns
JOIN work_paper_notes wpn ON wpns.work_paper_note_id = wpn.id
WHERE wpns.deleted_at IS NULL
ON CONFLICT (work_paper_id, user_id) WHERE deleted_at IS NULL AND status IN ('pending', 'signed', 'rejected')
DO NOTHING;

-- Drop the old work_paper_note_signatures table and related objects
DROP TRIGGER IF EXISTS update_work_paper_note_signatures_updated_at ON work_paper_note_signatures;
DROP INDEX IF EXISTS idx_work_paper_note_signatures_work_paper_note_id;
DROP INDEX IF EXISTS idx_work_paper_note_signatures_user_id;
DROP INDEX IF EXISTS idx_work_paper_note_signatures_status;
DROP INDEX IF EXISTS idx_work_paper_note_signatures_unique_user_note;
DROP TABLE IF EXISTS work_paper_note_signatures;