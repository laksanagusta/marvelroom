-- Create master vaccine list table
CREATE TABLE master_vaccines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vaccine_code VARCHAR(20) UNIQUE NOT NULL,
    vaccine_name_id VARCHAR(255) NOT NULL,
    vaccine_name_en VARCHAR(255) NOT NULL,
    description_id TEXT,
    description_en TEXT,
    vaccine_type VARCHAR(50) NOT NULL, -- 'routine', 'travel', 'optional'
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create countries table for CDC integration
CREATE TABLE countries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    country_code VARCHAR(3) UNIQUE NOT NULL, -- ISO 3166-1 alpha-3
    country_name_id VARCHAR(255) NOT NULL,
    country_name_en VARCHAR(255) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create country vaccine requirements table (CDC data cache)
CREATE TABLE country_vaccine_requirements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    country_id UUID NOT NULL REFERENCES countries(id) ON DELETE CASCADE,
    vaccine_id UUID NOT NULL REFERENCES master_vaccines(id) ON DELETE CASCADE,
    requirement_type VARCHAR(20) NOT NULL, -- 'required', 'recommended'
    cdc_data JSONB, -- Store raw CDC response
    cached_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP NOT NULL, -- Cache expiration
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(country_id, vaccine_id)
);

-- Create indexes
CREATE INDEX idx_master_vaccines_code ON master_vaccines(vaccine_code);
CREATE INDEX idx_master_vaccines_active ON master_vaccines(is_active);
CREATE INDEX idx_master_vaccines_type ON master_vaccines(vaccine_type);

CREATE INDEX idx_countries_code ON countries(country_code);
CREATE INDEX idx_countries_active ON countries(is_active);

CREATE INDEX idx_country_vaccine_requirements_country ON country_vaccine_requirements(country_id);
CREATE INDEX idx_country_vaccine_requirements_vaccine ON country_vaccine_requirements(vaccine_id);
CREATE INDEX idx_country_vaccine_requirements_expires ON country_vaccine_requirements(expires_at);