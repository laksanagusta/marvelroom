package pagination

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type QueryParser struct{}

func NewQueryParser() *QueryParser {
	return &QueryParser{}
}

func (qp *QueryParser) Parse(params map[string]string) (*QueryParams, error) {
	result := &QueryParams{
		Filters:    []Filter{},
		Sorts:      []Sort{},
		Pagination: Pagination{Page: 1, Limit: 20},
	}

	for key, value := range params {
		if key == "page" {
			if page, err := strconv.Atoi(value); err == nil && page > 0 {
				result.Pagination.Page = page
			}
			continue
		}

		if key == "limit" {
			if limit, err := strconv.Atoi(value); err == nil && limit > 0 && limit <= 100 {
				result.Pagination.Limit = limit
			}
			continue
		}

		if key == "sort" {
			sorts := qp.parseSorts(value)
			result.Sorts = append(result.Sorts, sorts...)
			continue
		}

		log.Println(key, value)

		filter, err := qp.parseFilter(key, value)
		if err == nil {
			result.Filters = append(result.Filters, filter)
		}
	}

	return result, nil
}

func (qp *QueryParser) parseFilter(key, value string) (Filter, error) {
	parts := strings.SplitN(value, " ", 2)
	if len(parts) != 2 {
		return Filter{}, fmt.Errorf("invalid filter format for %s: %s", key, value)
	}

	operator := parts[0]
	filterValue := parts[1]

	typedValue := qp.convertValue(key, filterValue)

	return Filter{
		Field:    key,
		Operator: operator,
		Value:    typedValue,
	}, nil
}

func (qp *QueryParser) parseSorts(value string) []Sort {
	var sorts []Sort
	sortPairs := strings.Split(value, ",")

	for _, pair := range sortPairs {
		parts := strings.SplitN(strings.TrimSpace(pair), " ", 2)
		field := parts[0]
		order := "asc"
		if len(parts) == 2 {
			order = parts[1]
		}
		sorts = append(sorts, Sort{Field: field, Order: order})
	}

	return sorts
}

func (qp *QueryParser) convertValue(field, value string) interface{} {
	// Handle array values for IN/NOT IN operators
	if strings.Contains(value, ",") {
		values := strings.Split(value, ",")
		result := make([]interface{}, len(values))
		for i, v := range values {
			result[i] = qp.convertSingleValue(field, strings.TrimSpace(v))
		}
		return result
	}

	return qp.convertSingleValue(field, value)
}

func (qp *QueryParser) convertSingleValue(field, value string) interface{} {
	switch field {
	case "page", "limit", "total_night":
		if v, err := strconv.Atoi(value); err == nil {
			return v
		}
	case "amount", "subtotal":
		if v, err := strconv.ParseFloat(value, 64); err == nil {
			return v
		}
	case "start_date", "end_date", "spd_date", "departure_date", "return_date", "created_at", "updated_at":
		if v, err := time.Parse("2006-01-02", value); err == nil {
			return v
		}
		if v, err := time.Parse(time.RFC3339, value); err == nil {
			return v
		}
	}
	return value
}
