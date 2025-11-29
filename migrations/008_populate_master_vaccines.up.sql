-- Populate master vaccines with common vaccines available in Indonesia

-- Routine Vaccines (Imunisasi Anak)
INSERT INTO master_vaccines (vaccine_code, vaccine_name_id, vaccine_name_en, description_id, description_en, vaccine_type, is_active, created_at, updated_at) VALUES
('BCG', 'BCG (Tuberkulosis)', 'BCG (Tuberculosis)', 'Vaksin untuk mencegah penyakit tuberkulosis parah pada anak', 'Vaccine to prevent severe tuberculosis in children', 'routine', true, NOW(), NOW()),
('POLIO', 'Polio (Poliomielitis)', 'Polio (Poliomyelitis)', 'Vaksin untuk mencegah penyakit polio yang menyebabkan lumpuh', 'Vaccine to prevent polio which causes paralysis', 'routine', true, NOW(), NOW()),
('DPT', 'DPT (Difteri, Pertusis, Tetanus)', 'DPT (Diphtheria, Pertussis, Tetanus)', 'Vaksin kombinasi untuk mencegah difteri, pertusis (batuk rejan), dan tetanus', 'Combined vaccine to prevent diphtheria, pertussis (whooping cough), and tetanus', 'routine', true, NOW(), NOW()),
('MMR', 'MMR (Mumps, Measles, Rubella)', 'MMR (Mumps, Measles, Rubella)', 'Vaksin kombinasi untuk mencegah gondongan, campak, dan rubella', 'Combined vaccine to prevent mumps, measles, and rubella', 'routine', true, NOW(), NOW()),
('HEP_B', 'Hepatitis B', 'Hepatitis B', 'Vaksin untuk mencegah penyakit hepatitis B', 'Vaccine to prevent hepatitis B', 'routine', true, NOW(), NOW()),
('HIB', 'Haemophilus influenzae tipe B', 'Haemophilus influenzae type B', 'Vaksin untuk mencegah pneumonia dan meningitis akibat HIB', 'Vaccine to prevent pneumonia and meningitis caused by HIB', 'routine', true, NOW(), NOW()),
('PCV', 'Pneumokokus (PCV13)', 'Pneumococcal (PCV13)', 'Vaksin untuk mencegah pneumonia, meningitis, dan infeksi pneumokokus lainnya', 'Vaccine to prevent pneumonia, meningitis, and other pneumococcal infections', 'routine', true, NOW(), NOW()),
('RV', 'Rotavirus', 'Rotavirus', 'Vaksin untuk mencegah diare berat akibat rotavirus pada anak', 'Vaccine to prevent severe diarrhea caused by rotavirus in children', 'routine', true, NOW(), NOW()),
('VAR', 'Varisela (Cacar Air)', 'Varicella (Chickenpox)', 'Vaksin untuk mencegah penyakit cacar air', 'Vaccine to prevent chickenpox', 'routine', true, NOW(), NOW()),

-- Travel Vaccines
('HEP_A', 'Hepatitis A', 'Hepatitis A', 'Vaksin untuk mencegah penyakit hepatitis A yang sering terjadi di negara berkembang', 'Vaccine to prevent hepatitis A common in developing countries', 'travel', true, NOW(), NOW()),
('TYPHOID', 'Tifoid', 'Typhoid', 'Vaksin untuk mencegah penyakit tifoid yang disebabkan oleh bakteri Salmonella typhi', 'Vaccine to prevent typhoid fever caused by Salmonella typhi bacteria', 'travel', true, NOW(), NOW()),
('JE', 'Japanese Encephalitis', 'Japanese Encephalitis', 'Vaksin untuk mencegah penyakit radang otak yang disebabkan oleh virus Japanese Encephalitis', 'Vaccine to prevent encephalitis caused by Japanese Encephalitis virus', 'travel', true, NOW(), NOW()),
('YELLOW_FEVER', 'Yellow Fever', 'Yellow Fever', 'Vaksin untuk mencegah demam kuning yang wajib untuk beberapa negara di Afrika dan Amerika Selatan', 'Vaccine to prevent yellow fever required for some countries in Africa and South America', 'travel', true, NOW(), NOW()),
('RABIES', 'Rabies', 'Rabies', 'Vaksin untuk mencegah penyakit rabies (gila anjing) setelah gigitan hewan', 'Vaccine to prevent rabies after animal bite', 'travel', true, NOW(), NOW()),
('MENING', 'Meningokokus (ACWY)', 'Meningococcal (ACWY)', 'Vaksin untuk mencegah meningitis dan sepsis yang disebabkan oleh bakteri meningokokus', 'Vaccine to prevent meningitis and sepsis caused by meningococcal bacteria', 'travel', true, NOW(), NOW()),
('CHOLERA', 'Kolera', 'Cholera', 'Vaksin untuk mencegah penyakit kolera yang menyebabkan diare akut', 'Vaccine to prevent cholera which causes acute diarrhea', 'travel', true, NOW(), NOW()),

-- Optional/Additional Vaccines
('HPV', 'HPV (Human Papillomavirus)', 'HPV (Human Papillomavirus)', 'Vaksin untuk mencegah kanker serviks dan kanker lain yang disebabkan oleh HPV', 'Vaccine to prevent cervical cancer and other cancers caused by HPV', 'optional', true, NOW(), NOW()),
('INFLUENZA', 'Influenza (Flu)', 'Influenza (Flu)', 'Vaksin untuk mencegah penyakit flu musiman', 'Vaccine to prevent seasonal flu', 'optional', true, NOW(), NOW()),
('COVID19', 'COVID-19', 'COVID-19', 'Vaksin untuk mencegah penyakit COVID-19', 'Vaccine to prevent COVID-19', 'optional', true, NOW(), NOW()),
('DTAP', 'DTAP (Difteri, Tetanus, Pertusis)', 'DTAP (Diphtheria, Tetanus, Pertussis)', 'Vaksin booster untuk difteri, tetanus, dan pertusis', 'Booster vaccine for diphtheria, tetanus, and pertussis', 'optional', true, NOW(), NOW()),
('TD', 'Td (Tetanus, Difteri)', 'Td (Tetanus, Diphtheria)', 'Vaksin booster untuk tetanus dan difteri', 'Booster vaccine for tetanus and diphtheria', 'optional', true, NOW(), NOW());