package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/repository"
	"sandbox/internal/infrastructure/cdc"
)

type CDCService struct {
	vaccinesRepo     repository.VaccinesRepository
	cdcClient        *cdc.CDCClient
	vaccineExtractor VaccineExtractor
}

func NewCDCService(vaccinesRepo repository.VaccinesRepository, cdcClient *cdc.CDCClient, vaccineExtractor VaccineExtractor) *CDCService {
	return &CDCService{
		vaccinesRepo:     vaccinesRepo,
		cdcClient:        cdcClient,
		vaccineExtractor: vaccineExtractor,
	}
}

type CountryVaccineRecommendation struct {
	CountryCode         string                  `json:"country_code"`
	CountryName         string                  `json:"country_name"`
	RequiredVaccines    []*entity.MasterVaccine `json:"required_vaccines"`
	RecommendedVaccines []*entity.MasterVaccine `json:"recommended_vaccines"`
	ConsiderVaccines    []*entity.MasterVaccine `json:"consider_vaccines"`
	MalariaRisk         string                  `json:"malaria_risk"`
	MalariaProphylaxis  string                  `json:"malaria_prophylaxis"`
	HealthNotice        string                  `json:"health_notice"`
	LastUpdated         time.Time               `json:"last_updated"`
}

func (s *CDCService) GetVaccineRecommendationsByCountry(ctx context.Context, countryCode string, language string) (*CountryVaccineRecommendation, error) {
	// Convert country code/name to full country name for database lookup
	fullCountryName := s.convertToFullCountryName(countryCode)

	// Get country information from our database
	country, err := s.vaccinesRepo.GetCountryByCode(ctx, fullCountryName)
	if err != nil {
		return nil, fmt.Errorf("country not found: %w", err)
	}

	// CDC client must be available
	if s.cdcClient == nil {
		return nil, fmt.Errorf("CDC client not initialized for %s. Unable to retrieve vaccine recommendations", countryCode)
	}

	// Get HTML from CDC website
	htmlData, err := s.cdcClient.GetDestinationHTML(ctx, countryCode)
	if err != nil {
		return nil, fmt.Errorf("failed to get CDC HTML for %s: %w", countryCode, err)
	}

	// Process HTML with Gemini to extract vaccine information
	vaccineData, err := s.vaccineExtractor.ExtractVaccineRecommendations(ctx, htmlData)
	if err != nil {
		return nil, fmt.Errorf("failed to process CDC HTML with Gemini for %s: %w", countryCode, err)
	}

	// Get all active vaccines from our database for mapping
	allVaccines, err := s.vaccinesRepo.ListActiveMasterVaccines(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get vaccines: %w", err)
	}

	// Convert Gemini response to our format
	recommendation := s.convertGeminiResponseToRecommendation(vaccineData, country, language, allVaccines)

	return recommendation, nil
}

// convertGeminiResponseToRecommendation converts Gemini response to our recommendation format
func (s *CDCService) convertGeminiResponseToRecommendation(vaccineData map[string]interface{}, country *entity.Country, language string, allVaccines []*entity.MasterVaccine) *CountryVaccineRecommendation {
	recommendation := &CountryVaccineRecommendation{
		CountryCode: country.CountryCode,
		CountryName: country.GetDisplayName(language),
	}

	// Create vaccine lookup map
	vaccineMap := make(map[string]*entity.MasterVaccine)
	for _, vaccine := range allVaccines {
		vaccineMap[strings.ToLower(vaccine.VaccineNameEN)] = vaccine
		vaccineMap[strings.ToLower(vaccine.VaccineNameID)] = vaccine
	}

	// Extract country name if available
	if countryName, ok := vaccineData["countryName"].(string); ok && countryName != "" {
		recommendation.CountryName = countryName
	}

	// Extract required vaccines
	if requiredVaccines, ok := vaccineData["requiredVaccines"].([]interface{}); ok {
		for _, req := range requiredVaccines {
			if reqMap, ok := req.(map[string]interface{}); ok {
				if name, ok := reqMap["name"].(string); ok {
					if vaccine := s.findVaccineByName(name, vaccineMap); vaccine != nil {
						recommendation.RequiredVaccines = append(recommendation.RequiredVaccines, vaccine)
					}
				}
			}
		}
	}

	// Extract recommended vaccines
	if recommendedVaccines, ok := vaccineData["recommendedVaccines"].([]interface{}); ok {
		for _, rec := range recommendedVaccines {
			if recMap, ok := rec.(map[string]interface{}); ok {
				if name, ok := recMap["name"].(string); ok {
					if vaccine := s.findVaccineByName(name, vaccineMap); vaccine != nil {
						recommendation.RecommendedVaccines = append(recommendation.RecommendedVaccines, vaccine)
					}
				}
			}
		}
	}

	// Extract consider vaccines
	if considerVaccines, ok := vaccineData["considerVaccines"].([]interface{}); ok {
		for _, con := range considerVaccines {
			if conMap, ok := con.(map[string]interface{}); ok {
				if name, ok := conMap["name"].(string); ok {
					if vaccine := s.findVaccineByName(name, vaccineMap); vaccine != nil {
						recommendation.ConsiderVaccines = append(recommendation.ConsiderVaccines, vaccine)
					}
				}
			}
		}
	}

	// Extract malaria information
	if malariaInfo, ok := vaccineData["malariaInfo"].(map[string]interface{}); ok {
		if risk, ok := malariaInfo["risk"].(string); ok {
			recommendation.MalariaRisk = risk
		}
		if prophylaxis, ok := malariaInfo["prophylaxis"].(string); ok {
			recommendation.MalariaProphylaxis = prophylaxis
		}
	}

	// Extract health notice
	if healthNotice, ok := vaccineData["healthNotice"].(string); ok {
		recommendation.HealthNotice = healthNotice
	}

	// Extract last updated date
	if lastUpdated, ok := vaccineData["lastUpdated"].(string); ok {
		// Try to parse various date formats
		if parsed, err := time.Parse("2006-01-02", lastUpdated); err == nil {
			recommendation.LastUpdated = parsed
		} else if parsed, err := time.Parse("01/02/2006", lastUpdated); err == nil {
			recommendation.LastUpdated = parsed
		}
	}

	return recommendation
}

// findVaccineByName finds vaccine by name using fuzzy matching
func (s *CDCService) findVaccineByName(name string, vaccineMap map[string]*entity.MasterVaccine) *entity.MasterVaccine {
	lowerName := strings.ToLower(name)

	// Direct match
	if vaccine, exists := vaccineMap[lowerName]; exists {
		return vaccine
	}

	// Common vaccine name mappings
	vaccineMappings := map[string]string{
		"hepatitis a":           "Hepatitis A",
		"hepatitis b":           "Hepatitis B",
		"typhoid":               "Typhoid",
		"yellow fever":          "Yellow Fever",
		"rabies":                "Rabies",
		"meningitis":            "Meningococcal",
		"meningococcal":         "Meningococcal",
		"cholera":               "Cholera",
		"japanese encephalitis": "Japanese Encephalitis",
		"influenza":             "Influenza",
		"flu":                   "Influenza",
		"mmr":                   "MMR",
		"measles":               "MMR",
		"mumps":                 "MMR",
		"rubella":               "MMR",
		"tetanus":               "Tetanus",
		"diphtheria":            "Diphtheria",
		"pertussis":             "Pertussis",
		"polio":                 "Polio",
		"covid":                 "COVID-19",
		"covid-19":              "COVID-19",
	}

	// Try mapping
	if mappedName, ok := vaccineMappings[lowerName]; ok {
		if vaccine, exists := vaccineMap[strings.ToLower(mappedName)]; exists {
			return vaccine
		}
	}

	// Try partial matching
	for mappingKey, mappedName := range vaccineMappings {
		if strings.Contains(lowerName, mappingKey) || strings.Contains(mappingKey, lowerName) {
			if vaccine, exists := vaccineMap[strings.ToLower(mappedName)]; exists {
				return vaccine
			}
		}
	}

	// Try fuzzy matching with existing vaccine names
	for existingName, vaccine := range vaccineMap {
		if strings.Contains(existingName, lowerName) || strings.Contains(lowerName, existingName) {
			return vaccine
		}
	}

	return nil
}

func (s *CDCService) parseVaccineInfo(vaccineInfo map[string]string, allVaccines []*entity.MasterVaccine) (required, recommended, consider []*entity.MasterVaccine) {
	vaccineMap := make(map[string]*entity.MasterVaccine)
	for _, vaccine := range allVaccines {
		vaccineMap[strings.ToLower(vaccine.VaccineNameEN)] = vaccine
		vaccineMap[strings.ToLower(vaccine.VaccineNameID)] = vaccine
	}

	// Common vaccine mappings
	vaccineMappings := map[string]string{
		"hepatitis a":           "Hepatitis A",
		"hepatitis b":           "Hepatitis B",
		"typhoid":               "Typhoid",
		"yellow fever":          "Yellow Fever",
		"rabies":                "Rabies",
		"meningitis":            "Meningococcal",
		"cholera":               "Cholera",
		"japanese encephalitis": "Japanese Encephalitis",
		"influenza":             "Influenza",
		"mmr":                   "MMR",
		"tetanus":               "Tetanus",
		"diphtheria":            "Diphtheria",
		"pertussis":             "Pertussis",
		"polio":                 "Polio",
	}

	for key, value := range vaccineInfo {
		lowerKey := strings.ToLower(key)
		lowerValue := strings.ToLower(value)

		// Check if this is a required vaccine
		if strings.Contains(lowerValue, "required") || strings.Contains(lowerValue, "mandatory") {
			if vaccineName := s.findVaccineName(lowerKey, vaccineMappings); vaccineName != nil {
				if vaccine, exists := vaccineMap[*vaccineName]; exists {
					required = append(required, vaccine)
				}
			}
		}

		// Check if this is a recommended vaccine
		if strings.Contains(lowerValue, "recommended") || strings.Contains(lowerValue, "suggested") {
			if vaccineName := s.findVaccineName(lowerKey, vaccineMappings); vaccineName != nil {
				if vaccine, exists := vaccineMap[*vaccineName]; exists {
					recommended = append(recommended, vaccine)
				}
			}
		}

		// Check if this is to be considered
		if strings.Contains(lowerValue, "consider") || strings.Contains(lowerValue, "if staying") {
			if vaccineName := s.findVaccineName(lowerKey, vaccineMappings); vaccineName != nil {
				if vaccine, exists := vaccineMap[*vaccineName]; exists {
					consider = append(consider, vaccine)
				}
			}
		}
	}

	return required, recommended, consider
}

func (s *CDCService) findVaccineName(key string, mappings map[string]string) *string {
	if name, ok := mappings[key]; ok {
		return &name
	}

	// Try partial matching
	for mappingKey, name := range mappings {
		if strings.Contains(key, mappingKey) || strings.Contains(mappingKey, key) {
			return &name
		}
	}

	return nil
}

// convertToFullCountryName converts country codes or names to full country names
func (s *CDCService) convertToFullCountryName(countryCode string) string {
	// Normalize the input
	lowerCode := strings.ToLower(strings.TrimSpace(countryCode))

	// Common country code mappings to full names
	countryMappings := map[string]string{
		// ISO-3 to full names
		"jpn": "Japan",
		"idn": "Indonesia",
		"mys": "Malaysia",
		"tha": "Thailand",
		"vnm": "Vietnam",
		"phl": "Philippines",
		"sgp": "Singapore",
		"chn": "China",
		"kor": "South Korea",
		"ind": "India",
		"usa": "United States",
		"gbr": "United Kingdom",
		"deu": "Germany",
		"fra": "France",
		"ita": "Italy",
		"esp": "Spain",
		"nld": "Netherlands",
		"aus": "Australia",
		"can": "Canada",
		"bra": "Brazil",
		"arg": "Argentina",
		"mex": "Mexico",
		"sau": "Saudi Arabia",
		"are": "United Arab Emirates",
		"egy": "Egypt",
		"zaf": "South Africa",
		"tur": "Turkey",
		"isr": "Israel",
		"mmr": "Myanmar",
		"khm": "Cambodia",
		"lao": "Laos",
		"brn": "Brunei Darussalam",
		"hkg": "Hong Kong",
		"twn": "Taiwan",
		"pak": "Pakistan",
		"bgd": "Bangladesh",
		"lka": "Sri Lanka",
		"mdv": "Maldives",
		"qat": "Qatar",
		"che": "Switzerland",
		"aut": "Austria",
		"bel": "Belgium",
		"swe": "Sweden",
		"chl": "Chile",
		"per": "Peru",
		"col": "Colombia",
		"mar": "Morocco",
		"ken": "Kenya",
		"tza": "Tanzania",
		"nzl": "New Zealand",
		"png": "Papua New Guinea",
		"fji": "Fiji",

		// ISO-2 to full names
		"jp": "Japan",
		"id": "Indonesia",
		"my": "Malaysia",
		"th": "Thailand",
		"vn": "Vietnam",
		"ph": "Philippines",
		"sg": "Singapore",
		"cn": "China",
		"kr": "South Korea",
		"in": "India",
		"us": "United States",
		"gb": "United Kingdom",
		"de": "Germany",
		"fr": "France",
		"it": "Italy",
		"es": "Spain",
		"nl": "Netherlands",
		"au": "Australia",
		"ca": "Canada",
		"br": "Brazil",
		"ar": "Argentina",
		"mx": "Mexico",
		"sa": "Saudi Arabia",
		"ae": "United Arab Emirates",
		"eg": "Egypt",
		"za": "South Africa",
		"tr": "Turkey",
		"il": "Israel",
		"mm": "Myanmar",
		"kh": "Cambodia",
		"la": "Laos",
		"bn": "Brunei Darussalam",
		"hk": "Hong Kong",
		"tw": "Taiwan",
		"pk": "Pakistan",
		"bd": "Bangladesh",
		"lk": "Sri Lanka",
		"mv": "Maldives",
		"qa": "Qatar",
		"ch": "Switzerland",
		"at": "Austria",
		"be": "Belgium",
		"se": "Sweden",
		"cl": "Chile",
		"pe": "Peru",
		"co": "Colombia",
		"ma": "Morocco",
		"ke": "Kenya",
		"tz": "Tanzania",
		"nz": "New Zealand",
		"pg": "Papua New Guinea",
		"fj": "Fiji",

		// Common name mappings
		"japan":                "Japan",
		"indonesia":            "Indonesia",
		"malaysia":             "Malaysia",
		"thailand":             "Thailand",
		"vietnam":              "Vietnam",
		"philippines":          "Philippines",
		"singapore":            "Singapore",
		"china":                "China",
		"south korea":          "South Korea",
		"korea":                "South Korea",
		"india":                "India",
		"united states":        "United States",
		"america":              "United States",
		"uk":                   "United Kingdom",
		"united kingdom":       "United Kingdom",
		"england":              "United Kingdom",
		"germany":              "Germany",
		"france":               "France",
		"italy":                "Italy",
		"spain":                "Spain",
		"netherlands":          "Netherlands",
		"holland":              "Netherlands",
		"australia":            "Australia",
		"canada":               "Canada",
		"brazil":               "Brazil",
		"argentina":            "Argentina",
		"mexico":               "Mexico",
		"saudi arabia":         "Saudi Arabia",
		"united arab emirates": "United Arab Emirates",
		"uae":                  "United Arab Emirates",
		"egypt":                "Egypt",
		"south africa":         "South Africa",
		"turkey":               "Turkey",
		"israel":               "Israel",
		"myanmar":              "Myanmar",
		"cambodia":             "Cambodia",
		"laos":                 "Laos",
		"brunei":               "Brunei Darussalam",
		"hong kong":            "Hong Kong",
		"taiwan":               "Taiwan",
		"pakistan":             "Pakistan",
		"bangladesh":           "Bangladesh",
		"sri lanka":            "Sri Lanka",
		"maldives":             "Maldives",
		"qatar":                "Qatar",
		"switzerland":          "Switzerland",
		"austria":              "Austria",
		"belgium":              "Belgium",
		"sweden":               "Sweden",
		"chile":                "Chile",
		"peru":                 "Peru",
		"colombia":             "Colombia",
		"morocco":              "Morocco",
		"kenya":                "Kenya",
		"tanzania":             "Tanzania",
		"new zealand":          "New Zealand",
		"papua new guinea":     "Papua New Guinea",
		"fiji":                 "Fiji",
	}

	// Look for exact match in mappings
	if fullName, exists := countryMappings[lowerCode]; exists {
		return fullName
	}

	// Try to find partial matches for country names
	for countryName, fullName := range countryMappings {
		if strings.Contains(lowerCode, countryName) || strings.Contains(countryName, lowerCode) {
			return fullName
		}
	}

	// If input looks like it might already be a full country name, capitalize it properly
	if len(strings.Fields(lowerCode)) > 1 || len(lowerCode) > 5 {
		return s.titleCaseCountryName(lowerCode)
	}

	// Last resort: return the input as-is (might already be a full country name)
	return strings.TrimSpace(countryCode)
}

// titleCaseCountryName converts a country name to title case
func (s *CDCService) titleCaseCountryName(name string) string {
	words := strings.Fields(strings.ToLower(name))
	for i, word := range words {
		if word == "and" || word == "of" || word == "the" {
			continue // keep these lowercase unless they're the first word
		}
		if i == 0 || word != "and" && word != "of" && word != "the" {
			words[i] = strings.Title(word)
		}
	}
	return strings.Join(words, " ")
}
