package pagination

import (
	"fmt"
	"strings"
)

type QueryParams struct {
	Filters    []Filter
	Sorts      []Sort
	Pagination Pagination
}

type Filter struct {
	Field    string
	Operator string
	Value    interface{}
}

type Sort struct {
	Field string
	Order string // asc or desc
}

type Pagination struct {
	Page  int
	Limit int
}

type PagedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalItems int64       `json:"total_items"`
	TotalPages int         `json:"total_pages"`
}

type QueryBuilder struct {
	baseQuery   string
	whereClause []string
	orderClause []string
	args        []interface{}
	argCounter  int
}

func NewQueryBuilder(baseQuery string) *QueryBuilder {
	return &QueryBuilder{
		baseQuery:  baseQuery,
		argCounter: 1,
	}
}

func (qb *QueryBuilder) AddFilter(filter Filter) error {
	operator := qb.mapOperator(filter.Operator)
	if operator == "" {
		return fmt.Errorf("unsupported operator: %s", filter.Operator)
	}

	field := qb.sanitizeField(filter.Field)
	if !qb.isValidField(field) {
		return fmt.Errorf("invalid field: %s", field)
	}

	if operator == "IN" || operator == "NOT IN" {
		values, ok := filter.Value.([]interface{})
		if !ok {
			return fmt.Errorf("IN/NOT IN operator requires array value")
		}
		placeholders := make([]string, len(values))
		for i, v := range values {
			placeholders[i] = fmt.Sprintf("$%d", qb.argCounter)
			qb.args = append(qb.args, v)
			qb.argCounter++
		}
		qb.whereClause = append(qb.whereClause, fmt.Sprintf("%s %s (%s)", field, operator, strings.Join(placeholders, ",")))
	} else if operator == "LIKE" || operator == "ILIKE" {
		qb.whereClause = append(qb.whereClause, fmt.Sprintf("LOWER(%s) %s LOWER($%d)", field, operator, qb.argCounter))
		qb.args = append(qb.args, "%"+fmt.Sprintf("%v", filter.Value)+"%")
		qb.argCounter++
	} else if operator == "IS" || operator == "IS NOT" {
		// For IS NULL and IS NOT NULL, don't use parameter binding
		if filter.Value == nil {
			qb.whereClause = append(qb.whereClause, fmt.Sprintf("%s %s NULL", field, operator))
		} else {
			qb.whereClause = append(qb.whereClause, fmt.Sprintf("%s %s $%d", field, operator, qb.argCounter))
			qb.args = append(qb.args, filter.Value)
			qb.argCounter++
		}
	} else {
		qb.whereClause = append(qb.whereClause, fmt.Sprintf("%s %s $%d", field, operator, qb.argCounter))
		qb.args = append(qb.args, filter.Value)
		qb.argCounter++
	}

	return nil
}

func (qb *QueryBuilder) AddSort(sort Sort) error {
	field := qb.sanitizeField(sort.Field)
	if !qb.isValidField(field) {
		return fmt.Errorf("invalid field: %s", field)
	}

	order := strings.ToUpper(sort.Order)
	if order != "ASC" && order != "DESC" {
		order = "ASC"
	}

	qb.orderClause = append(qb.orderClause, fmt.Sprintf("%s %s", field, order))
	return nil
}

func (qb *QueryBuilder) AddPagination(pagination Pagination) {
	// Pagination will be handled separately with LIMIT and OFFSET
}

func (qb *QueryBuilder) Build() (string, []interface{}) {
	query := qb.baseQuery

	if len(qb.whereClause) > 0 {
		query += " WHERE " + strings.Join(qb.whereClause, " AND ")
	}

	if len(qb.orderClause) > 0 {
		query += " ORDER BY " + strings.Join(qb.orderClause, ", ")
	}

	return query, qb.args
}

func (qb *QueryBuilder) mapOperator(op string) string {
	operators := map[string]string{
		"eq":     "=",
		"ne":     "!=",
		"gt":     ">",
		"gte":    ">=",
		"lt":     "<",
		"lte":    "<=",
		"like":   "LIKE",
		"ilike":  "ILIKE",
		"in":     "IN",
		"nin":    "NOT IN",
		"is":     "IS",
		"is_not": "IS NOT",
	}
	return operators[op]
}

func (qb *QueryBuilder) sanitizeField(field string) string {
	// Remove any potential SQL injection attempts
	field = strings.TrimSpace(field)
	field = strings.ToLower(field)
	return field
}

func (qb *QueryBuilder) isValidField(field string) bool {
	validFields := map[string]bool{
		"id":               true,
		"start_date":       true,
		"end_date":         true,
		"activity_purpose": true,
		"destination_city": true,
		"spd_date":         true,
		"departure_date":   true,
		"return_date":      true,
		"created_at":       true,
		"updated_at":       true,
		"deleted_at":       true,
		"business_trip_id": true,
		"assignee_id":      true,
		"name":             true,
		"spd_number":       true,
		"employee_id":      true,
		"employee_name":    true,
		"employee_number":  true,
		"position":         true,
		"rank":             true,
		"type":             true,
		"subtype":          true,
		"amount":           true,
		"total_night":      true,
		"subtotal":         true,
		"description":      true,
		"transport_detail": true,

		"country_name_id": true,

		// Work paper item fields
		"number":         true,
		"statement":      true,
		"explanation":    true,
		"filling_guide":  true,
		"parent_id":      true,
		"level":          true,
		"sort_order":     true,
		"is_active":      true,
		"work_paper_id":  true,
		"master_item_id": true,
		"gdrive_link":    true,
		"is_valid":       true,
		"notes":          true,
		"last_llm_response": true,
		"organization_id": true,
		"year":           true,
		"semester":       true,
		"status":         true,
	}
	return validFields[field]
}
