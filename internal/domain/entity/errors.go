package entity

import "errors"

var (
	ErrBusinessTripNotFound    = errors.New("business trip not found")
	ErrAssigneeNotFound        = errors.New("assignee not found")
	ErrTransactionNotFound     = errors.New("transaction not found")
	ErrInvalidDateRange        = errors.New("invalid date range")
	ErrDuplicateSPDNumber      = errors.New("duplicate SPD number")
	ErrUnauthorizedAccess      = errors.New("unauthorized access")
)