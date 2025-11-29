-- Update master_lakip_items table to support tree structure and type hierarchy

-- First, rename Indonesian columns to English if they exist
DO $$
BEGIN
    -- Check if Indonesian columns exist and rename them to English
    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'master_lakip_items' AND column_name = 'nomor') THEN
        ALTER TABLE master_lakip_items RENAME COLUMN nomor TO number;
    END IF;

    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'master_lakip_items' AND column_name = 'pernyataan') THEN
        ALTER TABLE master_lakip_items RENAME COLUMN pernyataan TO statement;
    END IF;

    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'master_lakip_items' AND column_name = 'penjelasan') THEN
        ALTER TABLE master_lakip_items RENAME COLUMN penjelasan TO explanation;
    END IF;

    IF EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = 'master_lakip_items' AND column_name = 'petunjuk_pengisian') THEN
        ALTER TABLE master_lakip_items RENAME COLUMN petunjuk_pengisian TO filling_guide;
    END IF;
END $$;

-- Also rename kertas_kerja tables to paper_work if they exist
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'kertas_kerjas') THEN
        ALTER TABLE kertas_kerjas RENAME TO paper_works;
    END IF;

    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'kertas_kerja_items') THEN
        ALTER TABLE kertas_kerja_items RENAME TO paper_work_items;
    END IF;
END $$;

-- Add new columns for tree structure and type
ALTER TABLE master_lakip_items
ADD COLUMN type VARCHAR(10) NOT NULL DEFAULT 'A',
ADD COLUMN parent_id UUID REFERENCES master_lakip_items(id) ON DELETE SET NULL,
ADD COLUMN level INTEGER NOT NULL DEFAULT 1,
ADD COLUMN sort_order INTEGER NOT NULL DEFAULT 0;

-- Create indexes for tree structure
CREATE INDEX idx_master_lakip_items_parent_id ON master_lakip_items(parent_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_master_lakip_items_type ON master_lakip_items(type) WHERE deleted_at IS NULL;
CREATE INDEX idx_master_lakip_items_level ON master_lakip_items(level) WHERE deleted_at IS NULL;

-- Create unique composite index for type and numbering within same parent and level
CREATE UNIQUE INDEX idx_master_lakip_items_type_parent_number
ON master_lakip_items(type, parent_id, number)
WHERE deleted_at IS NULL;

-- Update existing data to populate type, level, and sort_order
-- Set default type to 'A' for existing records
UPDATE master_lakip_items
SET type = 'A', level = 1, sort_order = 0
WHERE type IS NULL OR type = '';

-- Insert sample LAKIP tree structure data
INSERT INTO master_lakip_items (id, type, number, statement, explanation, filling_guide, level, sort_order, created_at, updated_at) VALUES
-- Type A - Utama (Root Level)
(gen_random_uuid(), 'A', '1', 'Perjanjian Kinerja Instansi Pemerintah', 'Dokumen perjanjian kinerja yang disepakati antara pimpinan instansi pemerintah dan pimpinan lembaga penilai kinerja', 'Lampirkan perjanjian kinerja yang telah ditandatangani oleh kedua belah pihak', 1, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

(gen_random_uuid(), 'A', '2', 'Standar Pelayanan', 'Standar pelayanan yang harus dipenuhi oleh instansi pemerintah', 'Dokumentasi standar pelayanan minimal yang harus dipenuhi', 1, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

(gen_random_uuid(), 'A', '3', 'Maklumat Pelayanan', 'Informasi tentang jenis dan kualitas pelayanan yang disediakan', 'Deskripsi lengkap tentang berbagai jenis layanan yang tersedia untuk masyarakat', 1, 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

-- Type B - Kegiatan Utama (Under Type A)
(gen_random_uuid(), 'B', '1.1', 'Penyelenggaraan Pemerintahan', 'Kegiatan penyelenggaraan tata kelola pemerintahan', 'Dokumentasi kegiatan penyelenggaraan administrasi pemerintahan', 2, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

(gen_random_uuid(), 'B', '1.2', 'Pelayanan Publik', 'Kegiatan penyelenggaraan pelayanan publik kepada masyarakat', 'Bukti penyelesaian pelayanan publik yang dilaksanakan', 2, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

(gen_random_uuid(), 'B', '1.3', 'Pembinaan SDM', 'Kegiatan pembinaan sumber daya manusia aparatur sipil negara', 'Dokumentasi program pembinaan dan pengembangan kompetensi pegawai', 2, 2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

(gen_random_uuid(), 'B', '1.4', 'Pengelolaan Keuangan', 'Kegiatan pengelolaan keuangan negara yang transparan dan akuntabel', 'Laporan keuangan dan bukti penggunaan anggaran yang sah', 2, 3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

(gen_random_uuid(), 'B', '1.5', 'Pemberdayaan Masyarakat', 'Kegiatan pemberdayaan masyarakat untuk meningkatkan kemandirian', 'Program-program pemberdayaan yang telah dilaksanakan dan dampaknya', 2, 4, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

-- Type C - Kegiatan Penunjang Utama (Under Type B)
(gen_random_uuid(), 'C', '1.1.1', 'Pemantauan dan Evaluasi Internal', 'Kegiatan pemantauan dan evaluasi kinerja internal', 'Laporan hasil pemantauan dan rekomendasi perbaikan', 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

(gen_random_uuid(), 'C', '1.1.2', 'Pengelolaan Pengaduan Masyarakat', 'Kegiatan penanganan pengaduan dari masyarakat', 'Data pengaduan yang diterima dan tindak lanjut yang diambil', 3, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

(gen_random_uuid(), 'C', '1.2.1', 'Digitalisasi Layanan', 'Kegiatan transformasi digital layanan publik', 'Layanan digital yang telah dikembangkan dan diadopsi', 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

(gen_random_uuid(), 'C', '1.2.2', 'Inovasi Pelayanan', 'Kegiatan inovasi dalam penyediaan layanan publik', 'Program inovasi yang telah dilaksanakan dan hasilnya', 3, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

(gen_random_uuid(), 'C', '1.3.1', 'Pelatihan dan Pengembangan', 'Kegiatan pelatihan dan pengembangan kompetensi pegawai', 'Program pelatihan yang telah diikuti dan sertifikat yang diperoleh', 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

(gen_random_uuid(), 'C', '1.3.2', 'Kesejahteraan dan Keselamatan Kerja', 'Kegiatan penjaminan kesehatan dan keselamatan kerja pegawai', 'Program kesehatan dan keselamatan yang telah diimplementasikan', 3, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

(gen_random_uuid(), 'C', '1.4.1', 'Pengadaan Sarana dan Prasarana', 'Kegiatan pengadaan sarana pendukung operasional', 'Daftar sarana dan prasarana yang telah disediakan', 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

(gen_random_uuid(), 'C', '1.4.2', 'Pemeliharaan Aset', 'Kegiatan pemeliharaan aset milik negara', 'Jadwal dan hasil pemeliharaan aset yang dilakukan', 3, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

(gen_random_uuid(), 'C', '1.5.1', 'Kerjasama Antar Lembaga', 'Kegiatan kerjasama dengan instansi pemerintah lainnya', 'Dokumen kerjasama yang telah ditandatangani dan implementasinya', 3, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),

(gen_random_uuid(), 'C', '1.5.2', 'Partisipasi Masyarakat', 'Kegiatan partisipasi aktif masyarakat dalam pembangunan', 'Forum dan mekanisme partisipasi masyarakat yang telah dibentuk', 3, 1, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- Update parent relationships for type B items (under "Perjanjian Kinerja")
UPDATE master_lakip_items SET parent_id = (SELECT id FROM master_lakip_items WHERE number = '1' AND type = 'A' AND level = 1 AND statement = 'Perjanjian Kinerja Instansi Pemerintah' LIMIT 1)
WHERE number IN ('1.1', '1.2', '1.3', '1.4', '1.5') AND type = 'B';

-- Update parent relationships for type C items under Penyelenggaraan Pemerintahan (1.1)
UPDATE master_lakip_items SET parent_id = (SELECT id FROM master_lakip_items WHERE number = '1.1' AND type = 'B' AND level = 2 AND statement = 'Penyelenggaraan Pemerintahan' LIMIT 1)
WHERE number IN ('1.1.1', '1.1.2') AND type = 'C';

-- Update parent relationships for type C items under Pelayanan Publik (1.2)
UPDATE master_lakip_items SET parent_id = (SELECT id FROM master_lakip_items WHERE number = '1.2' AND type = 'B' AND level = 2 AND statement = 'Pelayanan Publik' LIMIT 1)
WHERE number IN ('1.2.1', '1.2.2') AND type = 'C';

-- Update parent relationships for type C items under Pembinaan SDM (1.3)
UPDATE master_lakip_items SET parent_id = (SELECT id FROM master_lakip_items WHERE number = '1.3' AND type = 'B' AND level = 2 AND statement = 'Pembinaan SDM' LIMIT 1)
WHERE number IN ('1.3.1', '1.3.2') AND type = 'C';

-- Update parent relationships for type C items under Pengelolaan Keuangan (1.4)
UPDATE master_lakip_items SET parent_id = (SELECT id FROM master_lakip_items WHERE number = '1.4' AND type = 'B' AND level = 2 AND statement = 'Pengelolaan Keuangan' LIMIT 1)
WHERE number IN ('1.4.1', '1.4.2') AND type = 'C';

-- Update parent relationships for type C items under Pemberdayaan Masyarakat (1.5)
UPDATE master_lakip_items SET parent_id = (SELECT id FROM master_lakip_items WHERE number = '1.5' AND type = 'B' AND level = 2 AND statement = 'Pemberdayaan Masyarakat' LIMIT 1)
WHERE number IN ('1.5.1', '1.5.2') AND type = 'C';

-- Create function to get hierarchical LAKIP items
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