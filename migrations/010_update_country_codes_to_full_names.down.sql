-- Revert country codes back to 3-letter ISO codes

-- Update all country codes back to 3-letter ISO codes
UPDATE countries SET country_code = 'IDN' WHERE country_code = 'Indonesia';
UPDATE countries SET country_code = 'SGP' WHERE country_code = 'Singapore';
UPDATE countries SET country_code = 'MYS' WHERE country_code = 'Malaysia';
UPDATE countries SET country_code = 'THA' WHERE country_code = 'Thailand';
UPDATE countries SET country_code = 'VNM' WHERE country_code = 'Vietnam';
UPDATE countries SET country_code = 'PHL' WHERE country_code = 'Philippines';
UPDATE countries SET country_code = 'MMR' WHERE country_code = 'Myanmar';
UPDATE countries SET country_code = 'KHM' WHERE country_code = 'Cambodia';
UPDATE countries SET country_code = 'LAO' WHERE country_code = 'Laos';
UPDATE countries SET country_code = 'BRN' WHERE country_code = 'Brunei Darussalam';

-- East Asian Countries
UPDATE countries SET country_code = 'JPN' WHERE country_code = 'Japan';
UPDATE countries SET country_code = 'KOR' WHERE country_code = 'South Korea';
UPDATE countries SET country_code = 'CHN' WHERE country_code = 'China';
UPDATE countries SET country_code = 'HKG' WHERE country_code = 'Hong Kong';
UPDATE countries SET country_code = 'TWN' WHERE country_code = 'Taiwan';

-- South Asian Countries
UPDATE countries SET country_code = 'IND' WHERE country_code = 'India';
UPDATE countries SET country_code = 'PAK' WHERE country_code = 'Pakistan';
UPDATE countries SET country_code = 'BGD' WHERE country_code = 'Bangladesh';
UPDATE countries SET country_code = 'LKA' WHERE country_code = 'Sri Lanka';
UPDATE countries SET country_code = 'MDV' WHERE country_code = 'Maldives';

-- Middle East Countries
UPDATE countries SET country_code = 'SAU' WHERE country_code = 'Saudi Arabia';
UPDATE countries SET country_code = 'ARE' WHERE country_code = 'United Arab Emirates';
UPDATE countries SET country_code = 'QAT' WHERE country_code = 'Qatar';
UPDATE countries SET country_code = 'TUR' WHERE country_code = 'Turkey';
UPDATE countries SET country_code = 'ISR' WHERE country_code = 'Israel';
UPDATE countries SET country_code = 'EGY' WHERE country_code = 'Egypt';

-- European Countries
UPDATE countries SET country_code = 'GBR' WHERE country_code = 'United Kingdom';
UPDATE countries SET country_code = 'FRA' WHERE country_code = 'France';
UPDATE countries SET country_code = 'DEU' WHERE country_code = 'Germany';
UPDATE countries SET country_code = 'ITA' WHERE country_code = 'Italy';
UPDATE countries SET country_code = 'ESP' WHERE country_code = 'Spain';
UPDATE countries SET country_code = 'NLD' WHERE country_code = 'Netherlands';
UPDATE countries SET country_code = 'CHE' WHERE country_code = 'Switzerland';
UPDATE countries SET country_code = 'AUT' WHERE country_code = 'Austria';
UPDATE countries SET country_code = 'BEL' WHERE country_code = 'Belgium';
UPDATE countries SET country_code = 'SWE' WHERE country_code = 'Sweden';

-- North American Countries
UPDATE countries SET country_code = 'USA' WHERE country_code = 'United States';
UPDATE countries SET country_code = 'CAN' WHERE country_code = 'Canada';
UPDATE countries SET country_code = 'MEX' WHERE country_code = 'Mexico';

-- South American Countries
UPDATE countries SET country_code = 'BRA' WHERE country_code = 'Brazil';
UPDATE countries SET country_code = 'ARG' WHERE country_code = 'Argentina';
UPDATE countries SET country_code = 'CHL' WHERE country_code = 'Chile';
UPDATE countries SET country_code = 'PER' WHERE country_code = 'Peru';
UPDATE countries SET country_code = 'COL' WHERE country_code = 'Colombia';

-- African Countries
UPDATE countries SET country_code = 'ZAF' WHERE country_code = 'South Africa';
UPDATE countries SET country_code = 'MAR' WHERE country_code = 'Morocco';
UPDATE countries SET country_code = 'KEN' WHERE country_code = 'Kenya';
UPDATE countries SET country_code = 'TZA' WHERE country_code = 'Tanzania';

-- Oceania Countries
UPDATE countries SET country_code = 'AUS' WHERE country_code = 'Australia';
UPDATE countries SET country_code = 'NZL' WHERE country_code = 'New Zealand';
UPDATE countries SET country_code = 'PNG' WHERE country_code = 'Papua New Guinea';
UPDATE countries SET country_code = 'FJI' WHERE country_code = 'Fiji';

-- Revert the column type back to VARCHAR(3)
ALTER TABLE countries
ALTER COLUMN country_code TYPE VARCHAR(3);