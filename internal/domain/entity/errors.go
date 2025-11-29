package entity

import "errors"

var (
	ErrBusinessTripNotFound    = errors.New("business trip not found")
	ErrAssigneeNotFound        = errors.New("assignee not found")
	ErrTransactionNotFound     = errors.New("transaction not found")
	ErrInvalidDateRange        = errors.New("invalid date range")
	ErrDuplicateSPDNumber      = errors.New("duplicate SPD number")
	ErrUnauthorizedAccess      = errors.New("unauthorized access")

	// Desk module errors
	ErrWorkPaperItemNotFound          = errors.New("work paper item not found")
	ErrOrganizationNotFound           = errors.New("organization not found")
	ErrWorkPaperNotFound              = errors.New("work paper not found")
	ErrWorkPaperNoteNotFound          = errors.New("work paper note not found")
	ErrDuplicateWorkPaper             = errors.New("duplicate work paper for organization, year, and semester")
	ErrInvalidSemester                = errors.New("invalid semester, must be 1 or 2")
	ErrInvalidYear                    = errors.New("invalid year")
	ErrInvalidStatus                  = errors.New("invalid work paper status")
	ErrInvalidStatusTransition        = errors.New("invalid status transition")
	ErrWorkPaperItemTypeRequired      = errors.New("work paper item type is required")
	ErrWorkPaperItemNumberRequired    = errors.New("work paper item number is required")
	ErrWorkPaperItemStatementRequired = errors.New("work paper item statement is required")
	ErrInvalidWorkPaperItemType       = errors.New("invalid work paper item type, must be A, B, or C")
	ErrOrganizationIDRequired         = errors.New("organization ID is required")
	ErrWorkPaperIDRequired            = errors.New("work paper ID is required")
	ErrMasterItemIDRequired           = errors.New("master item ID is required")
	ErrWorkPaperNoteIDRequired        = errors.New("work paper note ID is required") // Legacy, keep for backward compatibility
	ErrUserIDRequired                 = errors.New("user ID is required")
	ErrUserNameRequired               = errors.New("user name is required")
	ErrInvalidSignatureType           = errors.New("invalid signature type, must be digital, manual, or approval")
	ErrSignatureNotFound              = errors.New("signature not found")
	ErrAlreadySigned                  = errors.New("signature already signed")
	ErrSignatureRejected              = errors.New("signature already rejected")
	ErrDuplicateSignature             = errors.New("signature already exists for this user and work paper")

	// Backward compatibility aliases (deprecated)
	ErrMasterLakipItemNotFound     = ErrWorkPaperItemNotFound
	ErrPaperWorkNotFound           = ErrWorkPaperNotFound
	ErrPaperWorkItemNotFound       = ErrWorkPaperNoteNotFound
	ErrDuplicatePaperWork          = ErrDuplicateWorkPaper
	ErrMasterLakipItemTypeRequired = ErrWorkPaperItemTypeRequired
	ErrMasterLakipItemNumberRequired = ErrWorkPaperItemNumberRequired
	ErrMasterLakipItemStatementRequired = ErrWorkPaperItemStatementRequired
	ErrInvalidMasterLakipItemType = ErrInvalidWorkPaperItemType
	ErrPaperWorkIDRequired         = ErrWorkPaperIDRequired
)