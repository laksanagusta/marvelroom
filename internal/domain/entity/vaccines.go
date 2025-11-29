package entity

import (
	"errors"
	"strings"
	"time"
)

type VaccineType string

const (
	VaccineTypeRoutine  VaccineType = "routine"
	VaccineTypeTravel   VaccineType = "travel"
	VaccineTypeOptional VaccineType = "optional"
)

type RequirementType string

const (
	RequirementTypeRequired      RequirementType = "required"
	RequirementTypeRecommended RequirementType = "recommended"
)

// MasterVaccine represents a master vaccine entity
type MasterVaccine struct {
	ID              string      `db:"id"`
	VaccineCode     string      `db:"vaccine_code"`
	VaccineNameID   string      `db:"vaccine_name_id"`
	VaccineNameEN   string      `db:"vaccine_name_en"`
	DescriptionID   *string     `db:"description_id"`
	DescriptionEN   *string     `db:"description_en"`
	VaccineType     VaccineType `db:"vaccine_type"`
	IsActive        bool        `db:"is_active"`
	CreatedAt       time.Time   `db:"created_at"`
	UpdatedAt       time.Time   `db:"updated_at"`
}

// Country represents a country entity
type Country struct {
	ID           string    `db:"id"`
	CountryCode  string    `db:"country_code"`
	CountryNameID string   `db:"country_name_id"`
	CountryNameEN string   `db:"country_name_en"`
	IsActive     bool      `db:"is_active"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// CountryVaccineRequirement represents vaccine requirements for a country
type CountryVaccineRequirement struct {
	ID              string           `db:"id"`
	CountryID       string           `db:"country_id"`
	VaccineID       string           `db:"vaccine_id"`
	RequirementType RequirementType  `db:"requirement_type"`
	CDCData         *string          `db:"cdc_data"`
	CachedAt        time.Time        `db:"cached_at"`
	ExpiresAt       time.Time        `db:"expires_at"`
	CreatedAt       time.Time        `db:"created_at"`
}

// NewMasterVaccine creates a new master vaccine with validation
func NewMasterVaccine(vaccineCode, vaccineNameID, vaccineNameEN string, vaccineType VaccineType, descriptionID, descriptionEN *string) (*MasterVaccine, error) {
	// Validation
	if strings.TrimSpace(vaccineCode) == "" {
		return nil, errors.New("vaccine code is required")
	}

	if strings.TrimSpace(vaccineNameID) == "" {
		return nil, errors.New("vaccine name (Indonesian) is required")
	}

	if strings.TrimSpace(vaccineNameEN) == "" {
		return nil, errors.New("vaccine name (English) is required")
	}

	if !isValidVaccineType(vaccineType) {
		return nil, errors.New("invalid vaccine type")
	}

	return &MasterVaccine{
		VaccineCode:   strings.ToUpper(strings.TrimSpace(vaccineCode)),
		VaccineNameID: strings.TrimSpace(vaccineNameID),
		VaccineNameEN: strings.TrimSpace(vaccineNameEN),
		DescriptionID: trimStringPtr(descriptionID),
		DescriptionEN: trimStringPtr(descriptionEN),
		VaccineType:   vaccineType,
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}, nil
}

// NewCountry creates a new country with validation
func NewCountry(countryCode, countryNameID, countryNameEN string) (*Country, error) {
	// Validation
	if strings.TrimSpace(countryCode) == "" {
		return nil, errors.New("country code is required")
	}

	if len(strings.TrimSpace(countryCode)) < 3 || len(strings.TrimSpace(countryCode)) > 100 {
		return nil, errors.New("country code must be between 3 and 100 characters")
	}

	if strings.TrimSpace(countryNameID) == "" {
		return nil, errors.New("country name (Indonesian) is required")
	}

	if strings.TrimSpace(countryNameEN) == "" {
		return nil, errors.New("country name (English) is required")
	}

	return &Country{
		CountryCode:   strings.TrimSpace(countryCode), // Don't force uppercase for full names
		CountryNameID: strings.TrimSpace(countryNameID),
		CountryNameEN: strings.TrimSpace(countryNameEN),
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}, nil
}

// NewCountryVaccineRequirement creates a new country vaccine requirement with validation
func NewCountryVaccineRequirement(countryID, vaccineID string, requirementType RequirementType, cdcData *string, cacheDuration time.Duration) (*CountryVaccineRequirement, error) {
	// Validation
	if strings.TrimSpace(countryID) == "" {
		return nil, errors.New("country ID is required")
	}

	if strings.TrimSpace(vaccineID) == "" {
		return nil, errors.New("vaccine ID is required")
	}

	if !isValidRequirementType(requirementType) {
		return nil, errors.New("invalid requirement type")
	}

	now := time.Now()
	return &CountryVaccineRequirement{
		CountryID:       strings.TrimSpace(countryID),
		VaccineID:       strings.TrimSpace(vaccineID),
		RequirementType: requirementType,
		CDCData:         cdcData,
		CachedAt:        now,
		ExpiresAt:       now.Add(cacheDuration),
		CreatedAt:       now,
	}, nil
}

// Update updates master vaccine
func (mv *MasterVaccine) Update(vaccineNameID, vaccineNameEN string, vaccineType VaccineType, descriptionID, descriptionEN *string) error {
	if strings.TrimSpace(vaccineNameID) == "" {
		return errors.New("vaccine name (Indonesian) is required")
	}

	if strings.TrimSpace(vaccineNameEN) == "" {
		return errors.New("vaccine name (English) is required")
	}

	if !isValidVaccineType(vaccineType) {
		return errors.New("invalid vaccine type")
	}

	mv.VaccineNameID = strings.TrimSpace(vaccineNameID)
	mv.VaccineNameEN = strings.TrimSpace(vaccineNameEN)
	mv.DescriptionID = trimStringPtr(descriptionID)
	mv.DescriptionEN = trimStringPtr(descriptionEN)
	mv.VaccineType = vaccineType
	mv.UpdatedAt = time.Now()

	return nil
}

// Update updates country
func (c *Country) Update(countryNameID, countryNameEN string) error {
	if strings.TrimSpace(countryNameID) == "" {
		return errors.New("country name (Indonesian) is required")
	}

	if strings.TrimSpace(countryNameEN) == "" {
		return errors.New("country name (English) is required")
	}

	c.CountryNameID = strings.TrimSpace(countryNameID)
	c.CountryNameEN = strings.TrimSpace(countryNameEN)
	c.UpdatedAt = time.Now()

	return nil
}

// RefreshCache refreshes the cache for country vaccine requirement
func (cvr *CountryVaccineRequirement) RefreshCache(cdcData *string, cacheDuration time.Duration) {
	now := time.Now()
	cvr.CDCData = cdcData
	cvr.CachedAt = now
	cvr.ExpiresAt = now.Add(cacheDuration)
}

// IsExpired checks if the cached data is expired
func (cvr *CountryVaccineRequirement) IsExpired() bool {
	return time.Now().After(cvr.ExpiresAt)
}

// Activate sets the entity to active status
func (mv *MasterVaccine) Activate() {
	mv.IsActive = true
	mv.UpdatedAt = time.Now()
}

// Deactivate sets the entity to inactive status
func (mv *MasterVaccine) Deactivate() {
	mv.IsActive = false
	mv.UpdatedAt = time.Now()
}

// Activate sets the country to active status
func (c *Country) Activate() {
	c.IsActive = true
	c.UpdatedAt = time.Now()
}

// Deactivate sets the country to inactive status
func (c *Country) Deactivate() {
	c.IsActive = false
	c.UpdatedAt = time.Now()
}

// GetDisplayName returns the display name based on language preference
func (mv *MasterVaccine) GetDisplayName(language string) string {
	switch strings.ToLower(language) {
	case "id", "indonesia":
		return mv.VaccineNameID
	case "en", "english":
		fallthrough
	default:
		return mv.VaccineNameEN
	}
}

// GetDescription returns the description based on language preference
func (mv *MasterVaccine) GetDescription(language string) *string {
	switch strings.ToLower(language) {
	case "id", "indonesia":
		return mv.DescriptionID
	case "en", "english":
		fallthrough
	default:
		return mv.DescriptionEN
	}
}

// GetDisplayName returns the display name based on language preference
func (c *Country) GetDisplayName(language string) string {
	switch strings.ToLower(language) {
	case "id", "indonesia":
		return c.CountryNameID
	case "en", "english":
		fallthrough
	default:
		return c.CountryNameEN
	}
}

// IsTravelVaccine checks if vaccine is for travel purposes
func (mv *MasterVaccine) IsTravelVaccine() bool {
	return mv.VaccineType == VaccineTypeTravel
}

// IsRoutineVaccine checks if vaccine is a routine vaccine
func (mv *MasterVaccine) IsRoutineVaccine() bool {
	return mv.VaccineType == VaccineTypeRoutine
}

// IsOptionalVaccine checks if vaccine is optional
func (mv *MasterVaccine) IsOptionalVaccine() bool {
	return mv.VaccineType == VaccineTypeOptional
}

// IsRequired checks if the requirement is mandatory
func (cvr *CountryVaccineRequirement) IsRequired() bool {
	return cvr.RequirementType == RequirementTypeRequired
}

// IsRecommended checks if the requirement is recommended
func (cvr *CountryVaccineRequirement) IsRecommended() bool {
	return cvr.RequirementType == RequirementTypeRecommended
}

// Helper functions

// trimStringPtr trims whitespace from string pointer
func trimStringPtr(s *string) *string {
	if s == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*s)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

// isValidVaccineType checks if the vaccine type is valid
func isValidVaccineType(vaccineType VaccineType) bool {
	switch vaccineType {
	case VaccineTypeRoutine, VaccineTypeTravel, VaccineTypeOptional:
		return true
	default:
		return false
	}
}

// isValidRequirementType checks if the requirement type is valid
func isValidRequirementType(requirementType RequirementType) bool {
	switch requirementType {
	case RequirementTypeRequired, RequirementTypeRecommended:
		return true
	default:
		return false
	}
}


// Getters

func (mv *MasterVaccine) GetID() string            { return mv.ID }
func (mv *MasterVaccine) GetVaccineCode() string   { return mv.VaccineCode }
func (mv *MasterVaccine) GetVaccineNameID() string { return mv.VaccineNameID }
func (mv *MasterVaccine) GetVaccineNameEN() string { return mv.VaccineNameEN }
func (mv *MasterVaccine) GetVaccineType() VaccineType { return mv.VaccineType }
func (mv *MasterVaccine) GetIsActive() bool        { return mv.IsActive }
func (mv *MasterVaccine) GetCreatedAt() time.Time  { return mv.CreatedAt }
func (mv *MasterVaccine) GetUpdatedAt() time.Time  { return mv.UpdatedAt }

func (c *Country) GetID() string           { return c.ID }
func (c *Country) GetCountryCode() string  { return c.CountryCode }
func (c *Country) GetCountryNameID() string { return c.CountryNameID }
func (c *Country) GetCountryNameEN() string { return c.CountryNameEN }
func (c *Country) GetIsActive() bool       { return c.IsActive }
func (c *Country) GetCreatedAt() time.Time { return c.CreatedAt }
func (c *Country) GetUpdatedAt() time.Time { return c.UpdatedAt }

func (cvr *CountryVaccineRequirement) GetID() string              { return cvr.ID }
func (cvr *CountryVaccineRequirement) GetCountryID() string       { return cvr.CountryID }
func (cvr *CountryVaccineRequirement) GetVaccineID() string       { return cvr.VaccineID }
func (cvr *CountryVaccineRequirement) GetRequirementType() RequirementType { return cvr.RequirementType }
func (cvr *CountryVaccineRequirement) GetCachedAt() time.Time     { return cvr.CachedAt }
func (cvr *CountryVaccineRequirement) GetExpiresAt() time.Time    { return cvr.ExpiresAt }
func (cvr *CountryVaccineRequirement) GetCreatedAt() time.Time    { return cvr.CreatedAt }