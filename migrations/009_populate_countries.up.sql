-- Populate countries table with popular travel destinations for Indonesian travelers

-- Southeast Asian Countries (Primary destinations)
INSERT INTO countries (country_code, country_name_id, country_name_en, is_active, created_at, updated_at) VALUES
('IDN', 'Indonesia', 'Indonesia', true, NOW(), NOW()),
('SGP', 'Singapura', 'Singapore', true, NOW(), NOW()),
('MYS', 'Malaysia', 'Malaysia', true, NOW(), NOW()),
('THA', 'Thailand', 'Thailand', true, NOW(), NOW()),
('VNM', 'Vietnam', 'Vietnam', true, NOW(), NOW()),
('PHL', 'Filipina', 'Philippines', true, NOW(), NOW()),
('MMR', 'Myanmar', 'Myanmar', true, NOW(), NOW()),
('KHM', 'Kamboja', 'Cambodia', true, NOW(), NOW()),
('LAO', 'Laos', 'Laos', true, NOW(), NOW()),
('BRN', 'Brunei Darussalam', 'Brunei Darussalam', true, NOW(), NOW()),

-- East Asian Countries
('JPN', 'Jepang', 'Japan', true, NOW(), NOW()),
('KOR', 'Korea Selatan', 'South Korea', true, NOW(), NOW()),
('CHN', 'China', 'China', true, NOW(), NOW()),
('HKG', 'Hong Kong', 'Hong Kong', true, NOW(), NOW()),
('TWN', 'Taiwan', 'Taiwan', true, NOW(), NOW()),

-- South Asian Countries
('IND', 'India', 'India', true, NOW(), NOW()),
('PAK', 'Pakistan', 'Pakistan', true, NOW(), NOW()),
('BGD', 'Bangladesh', 'Bangladesh', true, NOW(), NOW()),
('LKA', 'Sri Lanka', 'Sri Lanka', true, NOW(), NOW()),
('MDV', 'Maladewa', 'Maldives', true, NOW(), NOW()),

-- Middle East Countries
('SAU', 'Arab Saudi', 'Saudi Arabia', true, NOW(), NOW()),
('ARE', 'Uni Emirat Arab', 'United Arab Emirates', true, NOW(), NOW()),
('QAT', 'Qatar', 'Qatar', true, NOW(), NOW()),
('TUR', 'Turki', 'Turkey', true, NOW(), NOW()),
('ISR', 'Israel', 'Israel', true, NOW(), NOW()),
('EGY', 'Mesir', 'Egypt', true, NOW(), NOW()),

-- European Countries
('GBR', 'Inggris', 'United Kingdom', true, NOW(), NOW()),
('FRA', 'Prancis', 'France', true, NOW(), NOW()),
('DEU', 'Jerman', 'Germany', true, NOW(), NOW()),
('ITA', 'Italia', 'Italy', true, NOW(), NOW()),
('ESP', 'Spanyol', 'Spain', true, NOW(), NOW()),
('NLD', 'Belanda', 'Netherlands', true, NOW(), NOW()),
('CHE', 'Swiss', 'Switzerland', true, NOW(), NOW()),
('AUT', 'Austria', 'Austria', true, NOW(), NOW()),
('BEL', 'Belgia', 'Belgium', true, NOW(), NOW()),
('SWE', 'Swedia', 'Sweden', true, NOW(), NOW()),

-- North American Countries
('USA', 'Amerika Serikat', 'United States', true, NOW(), NOW()),
('CAN', 'Kanada', 'Canada', true, NOW(), NOW()),
('MEX', 'Meksiko', 'Mexico', true, NOW(), NOW()),

-- South American Countries
('BRA', 'Brasil', 'Brazil', true, NOW(), NOW()),
('ARG', 'Argentina', 'Argentina', true, NOW(), NOW()),
('CHL', 'Chile', 'Chile', true, NOW(), NOW()),
('PER', 'Peru', 'Peru', true, NOW(), NOW()),
('COL', 'Kolombia', 'Colombia', true, NOW(), NOW()),

-- African Countries
('ZAF', 'Afrika Selatan', 'South Africa', true, NOW(), NOW()),
('MAR', 'Maroko', 'Morocco', true, NOW(), NOW()),
('KEN', 'Kenya', 'Kenya', true, NOW(), NOW()),
('TZA', 'Tanzania', 'Tanzania', true, NOW(), NOW()),

-- Oceania Countries
('AUS', 'Australia', 'Australia', true, NOW(), NOW()),
('NZL', 'Selandia Baru', 'New Zealand', true, NOW(), NOW()),
('PNG', 'Papua Nugini', 'Papua New Guinea', true, NOW(), NOW()),
('FJI', 'Fiji', 'Fiji', true, NOW(), NOW());