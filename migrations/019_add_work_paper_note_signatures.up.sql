-- Add work paper note signatures table
CREATE TABLE work_paper_note_signatures (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    work_paper_note_id UUID NOT NULL,
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

    -- Foreign key constraint
    FOREIGN KEY (work_paper_note_id) REFERENCES work_paper_notes(id) ON DELETE CASCADE
);

-- Create indexes for work_paper_note_signatures
CREATE INDEX idx_work_paper_note_signatures_work_paper_note_id
ON work_paper_note_signatures (work_paper_note_id)
WHERE deleted_at IS NULL;

CREATE INDEX idx_work_paper_note_signatures_user_id
ON work_paper_note_signatures (user_id)
WHERE deleted_at IS NULL;

CREATE INDEX idx_work_paper_note_signatures_status
ON work_paper_note_signatures (status)
WHERE deleted_at IS NULL;

-- Create unique constraint to prevent duplicate signatures by same user on same note
CREATE UNIQUE INDEX idx_work_paper_note_signatures_unique_user_note
ON work_paper_note_signatures (work_paper_note_id, user_id)
WHERE deleted_at IS NULL AND status IN ('pending', 'signed', 'rejected');

-- Create trigger to automatically update updated_at column
CREATE TRIGGER update_work_paper_note_signatures_updated_at
    BEFORE UPDATE ON work_paper_note_signatures
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();