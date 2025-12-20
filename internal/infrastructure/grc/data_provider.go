package grc

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"sandbox/internal/domain/entity"
)

const (
	// Google Spreadsheet CSV export URL
	SpreadsheetURL = "https://docs.google.com/spreadsheets/d/1-m9ZMTkyG1UubyLSDZS_17yxuayA7E6G/export?format=csv"
	// Cache duration - refresh data every 5 minutes
	CacheDuration = 5 * time.Minute
)

// DataProvider fetches and caches GRC data from Google Spreadsheet
type DataProvider struct {
	mu          sync.RWMutex
	units       []entity.GRCUnit
	average     entity.GRCUnit
	lastFetched time.Time
	httpClient  *http.Client
}

// NewDataProvider creates a new data provider
func NewDataProvider() *DataProvider {
	return &DataProvider{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// FetchAll fetches units and average from spreadsheet in a single call (no cache)
func (dp *DataProvider) FetchAll() ([]entity.GRCUnit, entity.GRCUnit, error) {
	units, err := dp.fetchAndCache()
	if err != nil {
		return nil, entity.GRCUnit{}, err
	}

	dp.mu.RLock()
	defer dp.mu.RUnlock()
	return units, dp.average, nil
}

// GetAllUnits returns all GRC units, fetching from spreadsheet directly (no cache)
func (dp *DataProvider) GetAllUnits() ([]entity.GRCUnit, error) {
	// Always fetch fresh data (cache disabled)
	return dp.fetchAndCache()
}

// GetAverageScores returns the average scores
func (dp *DataProvider) GetAverageScores() (entity.GRCUnit, error) {
	// Always fetch fresh data (cache disabled)
	_, err := dp.fetchAndCache()
	if err != nil {
		return entity.GRCUnit{}, err
	}

	dp.mu.RLock()
	defer dp.mu.RUnlock()
	return dp.average, nil
}

// RefreshData forces a refresh of the data from the spreadsheet
func (dp *DataProvider) RefreshData() error {
	_, err := dp.fetchAndCache()
	return err
}

// fetchAndCache fetches data from Google Spreadsheet (no cache)
func (dp *DataProvider) fetchAndCache() ([]entity.GRCUnit, error) {
	dp.mu.Lock()
	defer dp.mu.Unlock()

	log.Println("Fetching fresh data from spreadsheet...")

	resp, err := dp.httpClient.Get(SpreadsheetURL)
	if err != nil {
		// Return cached data if available
		if len(dp.units) > 0 {
			return dp.units, nil
		}
		return nil, fmt.Errorf("failed to fetch spreadsheet: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if len(dp.units) > 0 {
			return dp.units, nil
		}
		return nil, fmt.Errorf("spreadsheet returned status %d", resp.StatusCode)
	}

	units, avg, err := dp.parseCSV(resp.Body)
	if err != nil {
		if len(dp.units) > 0 {
			return dp.units, nil
		}
		return nil, fmt.Errorf("failed to parse CSV: %w", err)
	}

	dp.units = units
	dp.average = avg
	dp.lastFetched = time.Now()

	return units, nil
}

// parseCSV parses the CSV data from the spreadsheet
func (dp *DataProvider) parseCSV(r io.Reader) ([]entity.GRCUnit, entity.GRCUnit, error) {
	reader := csv.NewReader(r)

	// Skip header row
	_, err := reader.Read()
	if err != nil {
		return nil, entity.GRCUnit{}, err
	}

	var units []entity.GRCUnit
	var average entity.GRCUnit

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue // Skip malformed rows
		}

		// Skip empty rows or rows without enough columns
		if len(record) < 16 {
			continue
		}

		// Check if this is the average row (No is empty or non-numeric)
		id, err := strconv.Atoi(strings.TrimSpace(record[0]))
		if err != nil {
			// This might be the "Rata-Rata Satker" row
			if strings.Contains(record[1], "Rata") {
				average = dp.parseAverageRow(record)
			}
			continue
		}

		unit := dp.parseUnitRow(id, record)
		units = append(units, unit)
	}

	return units, average, nil
}

// parseUnitRow parses a single unit row from CSV
func (dp *DataProvider) parseUnitRow(id int, record []string) entity.GRCUnit {
	name := strings.TrimSpace(record[1])

	// Parse percentage scores
	pski := parseFloat(record[9])
	pspip := parseFloat(record[10])
	pmri := parseFloat(record[11])
	piepk := parseFloat(record[12])
	pwbkrb := parseFloat(record[13])
	psakip := parseFloat(record[14])

	// Calculate average from 6 percentage components
	average := (pski + pspip + pmri + piepk + pwbkrb + psakip) / 6.0

	return entity.GRCUnit{
		ID:       id,
		Name:     name,
		Category: dp.inferCategory(name),
		SKI:      parseFloat(record[2]),
		SPIP:     parseFloat(record[3]),
		MRI:      parseFloat(record[4]),
		IEPK:     parseFloat(record[5]),
		SPIPT:    parseFloat(record[6]),
		WBKRB:    parseFloat(record[7]),
		SAKIP:    parseFloat(record[8]),
		PSKI:     pski,
		PSPIP:    pspip,
		PMRI:     pmri,
		PIEPK:    piepk,
		PWBKRB:   pwbkrb,
		PSAKIP:   psakip,
		Average:  average,
	}
}

// parseAverageRow parses the average row from CSV
func (dp *DataProvider) parseAverageRow(record []string) entity.GRCUnit {
	// Parse percentage scores
	pski := parseFloat(record[9])
	pspip := parseFloat(record[10])
	pmri := parseFloat(record[11])
	piepk := parseFloat(record[12])
	pwbkrb := parseFloat(record[13])
	psakip := parseFloat(record[14])

	// Calculate average from 6 percentage components
	average := (pski + pspip + pmri + piepk + pwbkrb + psakip) / 6.0

	return entity.GRCUnit{
		Name:    "Rata-Rata Satker",
		SKI:     parseFloat(record[2]),
		SPIP:    parseFloat(record[3]),
		MRI:     parseFloat(record[4]),
		IEPK:    parseFloat(record[5]),
		SPIPT:   parseFloat(record[6]),
		WBKRB:   parseFloat(record[7]),
		SAKIP:   parseFloat(record[8]),
		PSKI:    pski,
		PSPIP:   pspip,
		PMRI:    pmri,
		PIEPK:   piepk,
		PWBKRB:  pwbkrb,
		PSAKIP:  psakip,
		Average: average,
	}
}

// inferCategory determines the category based on unit name
func (dp *DataProvider) inferCategory(name string) string {
	nameLower := strings.ToLower(name)

	if strings.Contains(nameLower, "balai besar") {
		return "Balai Besar"
	}
	if strings.Contains(nameLower, "kelas i ") || strings.HasSuffix(nameLower, "kelas i") {
		return "Kelas I"
	}
	if strings.Contains(nameLower, "kelas ii") {
		return "Kelas II"
	}
	if strings.Contains(nameLower, "loka") {
		return "Loka"
	}
	if strings.Contains(nameLower, "direktorat") || strings.Contains(nameLower, "sekretariat") {
		return "Direktorat"
	}

	return "Lainnya"
}

// parseFloat safely parses a float from string
func parseFloat(s string) float64 {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, ",", ".")
	val, _ := strconv.ParseFloat(s, 64)
	return val
}
