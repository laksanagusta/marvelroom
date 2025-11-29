package cdc

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type CDCClient struct {
	baseURL    string
	webBaseURL string
	httpClient *http.Client
	apiKey     string
}

type DestinationInfo struct {
	DestinationName   string            `json:"destinationName"`
	VaccineInfo       map[string]string `json:"vaccineInfo"`
	MalariaInfo       map[string]string `json:"malariaInfo"`
	HealthInfo        map[string]string `json:"healthInfo"`
	LastUpdatedDate   string            `json:"lastUpdatedDate"`
	AvoidNonessential bool              `json:"avoidNonessential"`
}

type TravelHealthNotice struct {
	NoticeTitle      string `json:"noticeTitle"`
	NoticeText       string `json:"noticeText"`
	NoticeLevel      string `json:"noticeLevel"`
	DatePosted       string `json:"datePosted"`
	DateEffective    string `json:"dateEffective"`
	DateRevised      string `json:"dateRevised"`
	DateArchived     string `json:"dateArchived"`
	Transmission     string `json:"transmission"`
	DestinationID    string `json:"destinationId"`
	DestinationName  string `json:"destinationName"`
	ExternalService  string `json:"externalService"`
	ExternalSite     string `json:"externalSite"`
	ExternalURL      string `json:"externalUrl"`
	WarningLevel     string `json:"warningLevel"`
	VaccineName      string `json:"vaccineName"`
	WhichTravelers   string `json:"whichTravelers"`
	OralAntiviralMed string `json:"oralAntiviralMed"`
	YellowFever      string `json:"yellowFever"`
	Malaria          string `json:"malaria"`
	Rabies           string `json:"rabies"`
	Measles          string `json:"measles"`
	TBE              string `json:"tbe"` // Tick-borne encephalitis
	HepatitisA       string `json:"hepatitisA"`
	HepatitisB       string `json:"hepatitisB"`
	Typhoid          string `json:"typhoid"`
	Cholera          string `json:"cholera"`
}

// VaccineExtraction represents the structured vaccine information extracted from HTML
type VaccineExtraction struct {
	CountryName         string                 `json:"countryName"`
	RequiredVaccines    []VaccineDetail        `json:"requiredVaccines"`
	RecommendedVaccines []VaccineDetail        `json:"recommendedVaccines"`
	ConsiderVaccines    []VaccineDetail        `json:"considerVaccines"`
	MalariaInfo         MalariaDetail          `json:"malariaInfo"`
	HealthNotice        string                 `json:"healthNotice"`
	LastUpdated         string                 `json:"lastUpdated"`
	RawData             map[string]interface{} `json:"rawData"`
}

type VaccineDetail struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ForWho      string `json:"forWho"`
}

type MalariaDetail struct {
	Risk        string `json:"risk"`
	Prophylaxis string `json:"prophylaxis"`
}

func NewCDCClient(baseURL, webBaseURL, apiKey string) *CDCClient {
	return &CDCClient{
		baseURL:    baseURL,
		webBaseURL: webBaseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiKey: apiKey,
	}
}

func (c *CDCClient) GetDestinationInfo(ctx context.Context, countryCode string) (*DestinationInfo, error) {
	// CDC API endpoint for destination information
	endpoint := fmt.Sprintf("/destinations/regulation/json/%s", url.QueryEscape(countryCode))

	log.Println("CDC API endpoint:", c.baseURL+endpoint)

	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.apiKey != "" {
		req.Header.Set("X-Api-Key", c.apiKey)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("CDC API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var destination DestinationInfo
	if err := json.Unmarshal(body, &destination); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &destination, nil
}

func (c *CDCClient) GetTravelHealthNotices(ctx context.Context, countryCode string) ([]*TravelHealthNotice, error) {
	// CDC API endpoint for travel health notices
	endpoint := fmt.Sprintf("/travelHealthNotices/json/%s", url.QueryEscape(countryCode))

	log.Println("CDC API endpoint:", c.baseURL+endpoint)
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.apiKey != "" {
		req.Header.Set("X-Api-Key", c.apiKey)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("CDC API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var notices []*TravelHealthNotice
	if err := json.Unmarshal(body, &notices); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return notices, nil
}

// Alternative method to get data from CDC's structured data format
// This uses the CDC's public JSON endpoints which may have different structure
func (c *CDCClient) GetRawDestinationData(ctx context.Context, countryCode string) (map[string]interface{}, error) {
	endpoint := fmt.Sprintf("/destinations/json/%s", url.QueryEscape(countryCode))

	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("CDC API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return data, nil
}

// GetDestinationHTML fetches the HTML content from CDC travel destination page
func (c *CDCClient) GetDestinationHTML(ctx context.Context, countryCode string) (string, error) {
	// Construct CDC website URL for the country
	url := fmt.Sprintf("%s/travel/destinations/traveler/none/%s", c.webBaseURL, strings.ToLower(countryCode))

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers to mimic a browser request
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Connection", "keep-alive")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("CDC website returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Extract only the relevant vaccine table content to optimize processing
	return c.extractVaccineTableContent(string(body)), nil
}

// extractVaccineTableContent extracts only the relevant vaccine-related tables from HTML
func (c *CDCClient) extractVaccineTableContent(html string) string {
	// Look for the specific vaccine table mentioned by the user
	destVMARegex := regexp.MustCompile(`(?i)<table[^>]*class="[^"]*disease[^"]*hidden-one[^"]*"[^>]*id="dest-vm-a"[^>]*>.*?</table>`)

	// Also look for other vaccine-related tables and sections
	patterns := []string{
		`(?i)<table[^>]*class="[^"]*disease[^"]*"[^>]*>.*?</table>`, // Any disease table
		`(?i)<div[^>]*class="[^"]*vaccine[^"]*"[^>]*>.*?</div>`,     // Vaccine-related divs
		`(?i)<h2[^>]*class="[^"]*vaccine[^"]*"[^>]*>.*?</h2>`,       // Vaccine sections with class
		`(?i)<h3[^>]*class="[^"]*vaccine[^"]*"[^>]*>.*?</h3>`,       // Vaccine subsections with class
		`(?i)<div[^>]*id="[^"]*vaccine[^"]*"[^>]*>.*?</div>`,        // Elements with vaccine IDs
	}

	var extractedContent []string

	// First try to find the specific table mentioned
	if match := destVMARegex.FindString(html); match != "" {
		extractedContent = append(extractedContent, match)
		log.Printf("Found dest-vm-a vaccine table with %d characters", len(match))
	}

	// Then try other patterns
	for _, pattern := range patterns {
		regex := regexp.MustCompile(pattern + `(?s)`)
		matches := regex.FindAllString(html, -1)
		for _, match := range matches {
			// Avoid duplicates and very long matches
			if len(match) < 50000 && !c.containsMatch(extractedContent, match) {
				extractedContent = append(extractedContent, match)
			}
		}
	}

	// If no vaccine content found, try to find any health-related content
	if len(extractedContent) == 0 {
		healthPatterns := []string{
			`(?i)<div[^>]*class="[^"]*health[^"]*"[^>]*>.*?</div>`,
			`(?i)<table[^>]*class="[^"]*medicine[^"]*"[^>]*>.*?</table>`,
			`(?i)<h2[^>]*class="[^"]*health[^"]*"[^>]*>.*?</h2>`,
		}

		for _, pattern := range healthPatterns {
			regex := regexp.MustCompile(pattern + `(?s)`)
			if matches := regex.FindAllString(html, 3); len(matches) > 0 {
				extractedContent = append(extractedContent, matches...)
				break
			}
		}
	}

	if len(extractedContent) == 0 {
		log.Printf("Warning: No vaccine-specific content found, falling back to truncated HTML")
		// Fallback: return just the beginning of the HTML with a reasonable limit
		if len(html) > 100000 {
			return html[:100000]
		}
		return html
	}

	result := strings.Join(extractedContent, "\n\n")
	log.Printf("Extracted %d vaccine-related sections, total characters: %d (original: %d)",
		len(extractedContent), len(result), len(html))

	return result
}

// containsMatch checks if a match already exists in the extracted content
func (c *CDCClient) containsMatch(content []string, match string) bool {
	for _, existing := range content {
		if existing == match {
			return true
		}
	}
	return false
}
