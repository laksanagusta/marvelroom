-- Update country codes from 3-letter ISO codes to full country names

-- First, alter the country_code column to accommodate longer names
ALTER TABLE countries
ALTER COLUMN country_code TYPE VARCHAR(100);

-- Update all country codes to full names (English)
UPDATE countries SET country_code = 'Indonesia' WHERE country_code = 'IDN';
UPDATE countries SET country_code = 'Singapore' WHERE country_code = 'SGP';
UPDATE countries SET country_code = 'Malaysia' WHERE country_code = 'MYS';
UPDATE countries SET country_code = 'Thailand' WHERE country_code = 'THA';
UPDATE countries SET country_code = 'Vietnam' WHERE country_code = 'VNM';
UPDATE countries SET country_code = 'Philippines' WHERE country_code = 'PHL';
UPDATE countries SET country_code = 'Myanmar' WHERE country_code = 'MMR';
UPDATE countries SET country_code = 'Cambodia' WHERE country_code = 'KHM';
UPDATE countries SET country_code = 'Laos' WHERE country_code = 'LAO';
UPDATE countries SET country_code = 'Brunei Darussalam' WHERE country_code = 'BRN';

-- East Asian Countries
UPDATE countries SET country_code = 'Japan' WHERE country_code = 'JPN';
UPDATE countries SET country_code = 'South Korea' WHERE country_code = 'KOR';
UPDATE countries SET country_code = 'China' WHERE country_code = 'CHN';
UPDATE countries SET country_code = 'Hong Kong' WHERE country_code = 'HKG';
UPDATE countries SET country_code = 'Taiwan' WHERE country_code = 'TWN';

-- South Asian Countries
UPDATE countries SET country_code = 'India' WHERE country_code = 'IND';
UPDATE countries SET country_code = 'Pakistan' WHERE country_code = 'PAK';
UPDATE countries SET country_code = 'Bangladesh' WHERE country_code = 'BGD';
UPDATE countries SET country_code = 'Sri Lanka' WHERE country_code = 'LKA';
UPDATE countries SET country_code = 'Maldives' WHERE country_code = 'MDV';

-- Middle East Countries
UPDATE countries SET country_code = 'Saudi Arabia' WHERE country_code = 'SAU';
UPDATE countries SET country_code = 'United Arab Emirates' WHERE country_code = 'ARE';
UPDATE countries SET country_code = 'Qatar' WHERE country_code = 'QAT';
UPDATE countries SET country_code = 'Turkey' WHERE country_code = 'TUR';
UPDATE countries SET country_code = 'Israel' WHERE country_code = 'ISR';
UPDATE countries SET country_code = 'Egypt' WHERE country_code = 'EGY';

-- European Countries
UPDATE countries SET country_code = 'United Kingdom' WHERE country_code = 'GBR';
UPDATE countries SET country_code = 'France' WHERE country_code = 'FRA';
UPDATE countries SET country_code = 'Germany' WHERE country_code = 'DEU';
UPDATE countries SET country_code = 'Italy' WHERE country_code = 'ITA';
UPDATE countries SET country_code = 'Spain' WHERE country_code = 'ESP';
UPDATE countries SET country_code = 'Netherlands' WHERE country_code = 'NLD';
UPDATE countries SET country_code = 'Switzerland' WHERE country_code = 'CHE';
UPDATE countries SET country_code = 'Austria' WHERE country_code = 'AUT';
UPDATE countries SET country_code = 'Belgium' WHERE country_code = 'BEL';
UPDATE countries SET country_code = 'Sweden' WHERE country_code = 'SWE';

-- North American Countries
UPDATE countries SET country_code = 'United States' WHERE country_code = 'USA';
UPDATE countries SET country_code = 'Canada' WHERE country_code = 'CAN';
UPDATE countries SET country_code = 'Mexico' WHERE country_code = 'MEX';

-- South American Countries
UPDATE countries SET country_code = 'Brazil' WHERE country_code = 'BRA';
UPDATE countries SET country_code = 'Argentina' WHERE country_code = 'ARG';
UPDATE countries SET country_code = 'Chile' WHERE country_code = 'CHL';
UPDATE countries SET country_code = 'Peru' WHERE country_code = 'PER';
UPDATE countries SET country_code = 'Colombia' WHERE country_code = 'COL';

-- African Countries
UPDATE countries SET country_code = 'South Africa' WHERE country_code = 'ZAF';
UPDATE countries SET country_code = 'Morocco' WHERE country_code = 'MAR';
UPDATE countries SET country_code = 'Kenya' WHERE country_code = 'KEN';
UPDATE countries SET country_code = 'Tanzania' WHERE country_code = 'TZA';

-- Oceania Countries
UPDATE countries SET country_code = 'Australia' WHERE country_code = 'AUS';
UPDATE countries SET country_code = 'New Zealand' WHERE country_code = 'NZL';
UPDATE countries SET country_code = 'Papua New Guinea' WHERE country_code = 'PNG';
UPDATE countries SET country_code = 'Fiji' WHERE country_code = 'FJI';