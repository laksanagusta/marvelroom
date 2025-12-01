package business_trip

import (
	"time"

	"sandbox/internal/domain/entity"
	"sandbox/pkg/nullable"

	"github.com/invopop/validation"
)

// Global date validation function for easier use
func validateDateFormat(dateStr string) bool {
	if dateStr == "" {
		return true // Let Required handle empty validation
	}
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}

// BusinessTripRequest represents the request body for creating/updating a business trip
type BusinessTripRequest struct {
	BusinessTripNumber string               `json:"business_trip_number,omitempty"`
	StartDate          string               `json:"start_date"`
	EndDate            string               `json:"end_date"`
	ActivityPurpose    string               `json:"activity_purpose"`
	DestinationCity    string               `json:"destination_city"`
	SPDDate            string               `json:"spd_date"`
	DepartureDate      string               `json:"departure_date"`
	ReturnDate         string               `json:"return_date"`
	Status             string               `json:"status"`
	DocumentLink       string               `json:"document_link"`
	Verificators       []VerificatorRequest `json:"verificators"`
	Assignees          []AssigneeRequest    `json:"assignees"`
}

func (r BusinessTripRequest) Validate() error {
	// Basic validation with invopop
	err := validation.ValidateStruct(&r,
		validation.Field(&r.StartDate, validation.Required, validation.Length(1, 50)),
		validation.Field(&r.EndDate, validation.Required, validation.Length(1, 50)),
		validation.Field(&r.ActivityPurpose, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.DestinationCity, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.SPDDate, validation.Required, validation.Length(1, 50)),
		validation.Field(&r.DepartureDate, validation.Required, validation.Length(1, 50)),
		validation.Field(&r.ReturnDate, validation.Required, validation.Length(1, 50)),
		validation.Field(&r.Status, validation.Length(0, 20)),
		validation.Field(&r.DocumentLink, validation.Length(0, 500)),
		validation.Field(&r.Assignees, validation.Required, validation.Length(1, 50), validation.Each()),
		validation.Field(&r.Verificators, validation.Each()),
	)
	if err != nil {
		return err
	}

	// Manual validation for date format (custom validation)
	if !validateDateFormat(r.StartDate) {
		return validation.NewError("startDate", "must be a valid date in YYYY-MM-DD format")
	}
	if !validateDateFormat(r.EndDate) {
		return validation.NewError("endDate", "must be a valid date in YYYY-MM-DD format")
	}
	if !validateDateFormat(r.SPDDate) {
		return validation.NewError("spdDate", "must be a valid date in YYYY-MM-DD format")
	}
	if !validateDateFormat(r.DepartureDate) {
		return validation.NewError("departureDate", "must be a valid date in YYYY-MM-DD format")
	}
	if !validateDateFormat(r.ReturnDate) {
		return validation.NewError("returnDate", "must be a valid date in YYYY-MM-DD format")
	}

	// Validate status if provided
	if r.Status != "" {
		validStatuses := map[string]bool{
			"draft":           true,
			"ready_to_verify": true,
			"ongoing":         true,
			"completed":       true,
			"canceled":        true,
		}
		if !validStatuses[r.Status] {
			return validation.NewError("status", "must be one of: draft, ready_to_verify, ongoing, completed, canceled")
		}
	}

	// Validate nested verificators manually as they have their own validation
	for _, vr := range r.Verificators {
		if err := vr.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (r BusinessTripRequest) ToEntity() (*entity.BusinessTrip, error) {
	startDate, err := time.Parse("2006-01-02", r.StartDate)
	if err != nil {
		return nil, err
	}

	endDate, err := time.Parse("2006-01-02", r.EndDate)
	if err != nil {
		return nil, err
	}

	spdDate, err := time.Parse("2006-01-02", r.SPDDate)
	if err != nil {
		return nil, err
	}

	departureDate, err := time.Parse("2006-01-02", r.DepartureDate)
	if err != nil {
		return nil, err
	}

	returnDate, err := time.Parse("2006-01-02", r.ReturnDate)
	if err != nil {
		return nil, err
	}

	bt, err := entity.NewBusinessTrip(startDate, endDate, spdDate, departureDate, returnDate, r.ActivityPurpose, r.DestinationCity)
	if err != nil {
		return nil, err
	}

	// Set status if provided (must be valid status)
	if r.Status != "" {
		status := entity.BusinessTripStatus(r.Status)
		if status != entity.BusinessTripStatusDraft {
			// If status is not draft, we need to validate the transition
			if err := bt.UpdateStatus(status); err != nil {
				return nil, err
			}
		}
	}

	// Set document link
	if r.DocumentLink != "" {
		bt.UpdateDocumentLink(r.DocumentLink)
	}

	// Add verificators
	for _, vr := range r.Verificators {
		_, err := bt.AddVerificator(vr.UserID, vr.UserName, vr.EmployeeNumber, vr.Position)
		if err != nil {
			return nil, err
		}
	}

	// Add assignees
	for _, assigneeReq := range r.Assignees {
		assignee, err := bt.AddAssignee(assigneeReq.Name, assigneeReq.SPDNumber, assigneeReq.EmployeeID, assigneeReq.EmployeeName, assigneeReq.EmployeeNumber, assigneeReq.Position, assigneeReq.Rank)
		if err != nil {
			return nil, err
		}

		// Add transactions for each assignee
		for _, transactionReq := range assigneeReq.Transactions {
			txType := entity.TransactionType(transactionReq.Type)
			subtype := entity.TransactionSubtype(transactionReq.Subtype)

			transaction, err := entity.NewTransaction(
				transactionReq.Name,
				txType,
				subtype,
				transactionReq.Amount,
				transactionReq.Amount, // Will be calculated in NewTransaction
				transactionReq.TotalNight,
				transactionReq.Description,
				transactionReq.TransportDetail,
			)
			if err != nil {
				return nil, err
			}

			err = assignee.AddTransaction(transaction)
			if err != nil {
				return nil, err
			}
		}
	}

	return bt, nil
}

// AssigneeRequest represents the request body for an assignee
type AssigneeRequest struct {
	Name           string               `json:"name"`
	SPDNumber      string               `json:"spd_number"`
	EmployeeID     string               `json:"employee_id"`
	EmployeeName   string               `json:"employee_name"`
	EmployeeNumber string               `json:"employee_number"`
	Position       string               `json:"position"`
	Rank           string               `json:"rank"`
	Transactions   []TransactionRequest `json:"transactions"`
}

func (r AssigneeRequest) Validate() error {
	err := validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.SPDNumber, validation.Required, validation.Length(1, 100)),
		validation.Field(&r.EmployeeNumber, validation.Required, validation.Length(1, 50)),
		validation.Field(&r.Position, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.Rank, validation.Required, validation.Length(1, 100)),
		validation.Field(&r.Transactions, validation.Each()),
	)
	if err != nil {
		return err
	}

	// Validate nested transactions manually as they have their own validation
	for _, tx := range r.Transactions {
		if err := tx.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// TransactionRequest represents the request body for a transaction
type TransactionRequest struct {
	Name            string  `json:"name"`
	Type            string  `json:"type"`
	Subtype         string  `json:"subtype"`
	Amount          float64 `json:"amount"`
	TotalNight      *int    `json:"total_night"`
	Description     string  `json:"description"`
	TransportDetail string  `json:"transport_detail"`
}

// VerificatorRequest represents the request body for a verificator
type VerificatorRequest struct {
	UserID         string `json:"user_id"`
	UserName       string `json:"user_name"`
	EmployeeNumber string `json:"employee_number"`
	Position       string `json:"position"`
}

func (r VerificatorRequest) Validate() error {
	// Basic validation with invopop
	err := validation.ValidateStruct(&r,
		validation.Field(&r.UserID, validation.Required, validation.Length(1, 100)),
		validation.Field(&r.UserName, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.EmployeeNumber, validation.Required, validation.Length(1, 50)),
		validation.Field(&r.Position, validation.Required, validation.Length(1, 255)),
	)
	if err != nil {
		return err
	}

	return nil
}

func (r TransactionRequest) Validate() error {
	// Basic validation with invopop
	err := validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.Type, validation.Required, validation.Length(1, 50)),
		validation.Field(&r.Subtype, validation.Length(0, 50)),
		validation.Field(&r.Amount, validation.Required, validation.Min(0.0)),
		validation.Field(&r.TotalNight, validation.Min(0)),
		validation.Field(&r.Description, validation.Length(0, 1000)),
		validation.Field(&r.TransportDetail, validation.Length(0, 1000)),
	)
	if err != nil {
		return err
	}

	// Manual validation for enum values (custom validation)
	validTypes := map[string]bool{
		"accommodation": true,
		"transport":     true,
		"other":         true,
		"allowance":     true,
	}
	if !validTypes[r.Type] {
		return validation.NewError("type", "must be one of: accommodation, transport, other, allowance")
	}

	if r.Subtype != "" {
		validSubtypes := map[string]bool{
			"hotel":           true,
			"flight":          true,
			"train":           true,
			"taxi":            true,
			"daily_allowance": true,
			"rental_car":      true,
			"meal":            true,
			"other":           true,
		}
		if !validSubtypes[r.Subtype] {
			return validation.NewError("subtype", "must be one of: hotel, flight, train, taxi, daily_allowance, rental_car, meal, other")
		}
	}

	return nil
}

// UpdateBusinessTripRequest represents the request body for updating a business trip
type UpdateBusinessTripRequest struct {
	BusinessTripID     string              `params:"tripId" json:"tripId"`
	BusinessTripNumber nullable.NullString `json:"business_trip_number"`
	StartDate          nullable.NullString `json:"start_date"`
	EndDate            nullable.NullString `json:"end_date"`
	ActivityPurpose    nullable.NullString `json:"activity_purpose"`
	DestinationCity    nullable.NullString `json:"destination_city"`
	SPDDate            nullable.NullString `json:"spd_date"`
	DepartureDate      nullable.NullString `json:"departure_date"`
	ReturnDate         nullable.NullString `json:"return_date"`
	Status             nullable.NullString `json:"status"`
	DocumentLink       nullable.NullString `json:"document_link"`
}

// UpdateBusinessTripWithAssigneesRequest represents the request body for updating a business trip with full replace of assignees and transactions
type UpdateBusinessTripWithAssigneesRequest struct {
	BusinessTripID     string               `params:"tripId" json:"tripId"`
	BusinessTripNumber string               `json:"business_trip_number"`
	StartDate          string               `json:"start_date"`
	EndDate            string               `json:"end_date"`
	ActivityPurpose    string               `json:"activity_purpose"`
	DestinationCity    string               `json:"destination_city"`
	SPDDate            string               `json:"spd_date"`
	DepartureDate      string               `json:"departure_date"`
	ReturnDate         string               `json:"return_date"`
	Status             string               `json:"status"`
	DocumentLink       string               `json:"document_link"`
	Verificators       []VerificatorRequest `json:"verificators"`
	Assignees          []AssigneeRequest    `json:"assignees"`
}

func (r UpdateBusinessTripRequest) Validate() error {
	// Validate BusinessTripID
	if r.BusinessTripID == "" {
		return validation.NewError("tripId", "is required")
	}

	// For nullable fields, we need manual validation as invopop doesn't directly support nullable types
	if r.StartDate.IsSet() && !validateDateFormat(r.StartDate.String) {
		return validation.NewError("startDate", "must be a valid date in YYYY-MM-DD format")
	}
	if r.EndDate.IsSet() && !validateDateFormat(r.EndDate.String) {
		return validation.NewError("endDate", "must be a valid date in YYYY-MM-DD format")
	}
	if r.SPDDate.IsSet() && !validateDateFormat(r.SPDDate.String) {
		return validation.NewError("spdDate", "must be a valid date in YYYY-MM-DD format")
	}
	if r.DepartureDate.IsSet() && !validateDateFormat(r.DepartureDate.String) {
		return validation.NewError("departureDate", "must be a valid date in YYYY-MM-DD format")
	}
	if r.ReturnDate.IsSet() && !validateDateFormat(r.ReturnDate.String) {
		return validation.NewError("returnDate", "must be a valid date in YYYY-MM-DD format")
	}

	if r.ActivityPurpose.IsSet() && len(r.ActivityPurpose.String) > 255 {
		return validation.NewError("activityPurpose", "must be less than 255 characters")
	}
	if r.DestinationCity.IsSet() && len(r.DestinationCity.String) > 255 {
		return validation.NewError("destinationCity", "must be less than 255 characters")
	}

	// Validate status if provided
	if r.Status.IsSet() && r.Status.String != "" {
		validStatuses := map[string]bool{
			"draft":     true,
			"ongoing":   true,
			"completed": true,
			"canceled":  true,
		}
		if !validStatuses[r.Status.String] {
			return validation.NewError("status", "must be one of: draft, ongoing, completed, canceled")
		}
	}

	// Validate document link length if provided
	if r.DocumentLink.IsSet() && len(r.DocumentLink.String) > 500 {
		return validation.NewError("documentLink", "must be less than 500 characters")
	}

	return nil
}

func (r UpdateBusinessTripWithAssigneesRequest) Validate() error {
	// Validate BusinessTripID
	if r.BusinessTripID == "" {
		return validation.NewError("tripId", "is required")
	}

	// Basic validation with invopop
	err := validation.ValidateStruct(&r,
		validation.Field(&r.StartDate, validation.Required, validation.Length(1, 50)),
		validation.Field(&r.EndDate, validation.Required, validation.Length(1, 50)),
		validation.Field(&r.ActivityPurpose, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.DestinationCity, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.SPDDate, validation.Required, validation.Length(1, 50)),
		validation.Field(&r.DepartureDate, validation.Required, validation.Length(1, 50)),
		validation.Field(&r.ReturnDate, validation.Required, validation.Length(1, 50)),
		validation.Field(&r.Status, validation.Length(0, 20)),
		validation.Field(&r.DocumentLink, validation.Length(0, 500)),
		validation.Field(&r.Assignees, validation.Required, validation.Length(1, 50), validation.Each()),
		validation.Field(&r.Verificators, validation.Each()),
	)
	if err != nil {
		return err
	}

	// Manual validation for date format (custom validation)
	if !validateDateFormat(r.StartDate) {
		return validation.NewError("startDate", "must be a valid date in YYYY-MM-DD format")
	}
	if !validateDateFormat(r.EndDate) {
		return validation.NewError("endDate", "must be a valid date in YYYY-MM-DD format")
	}
	if !validateDateFormat(r.SPDDate) {
		return validation.NewError("spdDate", "must be a valid date in YYYY-MM-DD format")
	}
	if !validateDateFormat(r.DepartureDate) {
		return validation.NewError("departureDate", "must be a valid date in YYYY-MM-DD format")
	}
	if !validateDateFormat(r.ReturnDate) {
		return validation.NewError("returnDate", "must be a valid date in YYYY-MM-DD format")
	}

	// Validate status if provided
	if r.Status != "" {
		validStatuses := map[string]bool{
			"draft":           true,
			"ready_to_verify": true,
			"ongoing":         true,
			"completed":       true,
			"canceled":        true,
		}
		if !validStatuses[r.Status] {
			return validation.NewError("status", "must be one of: draft, ready_to_verify, ongoing, completed, canceled")
		}
	}

	return nil
}

func (r UpdateBusinessTripWithAssigneesRequest) ToEntity(businessTripID string) (*entity.BusinessTrip, error) {
	startDate, err := time.Parse("2006-01-02", r.StartDate)
	if err != nil {
		return nil, err
	}

	endDate, err := time.Parse("2006-01-02", r.EndDate)
	if err != nil {
		return nil, err
	}

	spdDate, err := time.Parse("2006-01-02", r.SPDDate)
	if err != nil {
		return nil, err
	}

	departureDate, err := time.Parse("2006-01-02", r.DepartureDate)
	if err != nil {
		return nil, err
	}

	returnDate, err := time.Parse("2006-01-02", r.ReturnDate)
	if err != nil {
		return nil, err
	}

	bt, err := entity.NewBusinessTrip(startDate, endDate, spdDate, departureDate, returnDate, r.ActivityPurpose, r.DestinationCity)
	if err != nil {
		return nil, err
	}

	// Set the ID for update
	bt.ID = businessTripID

	// Set document link
	if r.DocumentLink != "" {
		bt.UpdateDocumentLink(r.DocumentLink)
	}

	// Set status if provided (must be valid status)
	if r.Status != "" {
		status := entity.BusinessTripStatus(r.Status)
		if status != entity.BusinessTripStatusDraft {
			// If status is not draft, we need to validate the transition
			if err := bt.UpdateStatus(status); err != nil {
				return nil, err
			}
		}
	}

	// Add verificators
	for _, vr := range r.Verificators {
		_, err := bt.AddVerificator(vr.UserID, vr.UserName, vr.EmployeeNumber, vr.Position)
		if err != nil {
			return nil, err
		}
	}

	// Add assignees
	for _, assigneeReq := range r.Assignees {
		assignee, err := bt.AddAssignee(assigneeReq.Name, assigneeReq.SPDNumber, assigneeReq.EmployeeID, assigneeReq.EmployeeName, assigneeReq.EmployeeNumber, assigneeReq.Position, assigneeReq.Rank)
		if err != nil {
			return nil, err
		}

		// Add transactions for each assignee
		for _, transactionReq := range assigneeReq.Transactions {
			txType := entity.TransactionType(transactionReq.Type)
			subtype := entity.TransactionSubtype(transactionReq.Subtype)

			transaction, err := entity.NewTransaction(
				transactionReq.Name,
				txType,
				subtype,
				transactionReq.Amount,
				transactionReq.Amount, // Will be calculated in NewTransaction
				transactionReq.TotalNight,
				transactionReq.Description,
				transactionReq.TransportDetail,
			)
			if err != nil {
				return nil, err
			}

			err = assignee.AddTransaction(transaction)
			if err != nil {
				return nil, err
			}
		}
	}

	return bt, nil
}

// BusinessTripResponse represents the response body for a business trip
type BusinessTripResponse struct {
	ID                 string                `json:"id"`
	BusinessTripNumber string                `json:"business_trip_number"`
	StartDate          string                `json:"start_date"`
	EndDate            string                `json:"end_date"`
	ActivityPurpose    string                `json:"activity_purpose"`
	DestinationCity    string                `json:"destination_city"`
	SPDDate            string                `json:"spd_date"`
	DepartureDate      string                `json:"departure_date"`
	ReturnDate         string                `json:"return_date"`
	Status             string                `json:"status"`
	DocumentLink       string                `json:"document_link"`
	TotalCost          float64               `json:"total_cost"`
	Verificators       []VerificatorResponse `json:"verificators"`
	Assignees          []AssigneeResponse    `json:"assignees"`
	CreatedAt          string                `json:"created_at"`
	UpdatedAt          string                `json:"updated_at"`
}

// VerificatorResponse represents the response body for a verificator
type VerificatorResponse struct {
	ID                string  `json:"id"`
	BusinessTripID    string  `json:"business_trip_id"`
	UserID            string  `json:"user_id"`
	UserName          string  `json:"user_name"`
	EmployeeNumber    string  `json:"employee_number"`
	Position          string  `json:"position"`
	Status            string  `json:"status"`
	VerifiedAt        *string `json:"verified_at"`
	VerificationNotes string  `json:"verification_notes"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
}

// AssigneeResponse represents the response body for an assignee
type AssigneeResponse struct {
	ID             string                `json:"id"`
	Name           string                `json:"name"`
	SPDNumber      string                `json:"spd_number"`
	EmployeeID     string                `json:"employee_id"`
	EmployeeName   string                `json:"employee_name"`
	EmployeeNumber string                `json:"employee_number"`
	Position       string                `json:"position"`
	Rank           string                `json:"rank"`
	TotalCost      float64               `json:"total_cost"`
	Transactions   []TransactionResponse `json:"transactions"`
	CreatedAt      string                `json:"created_at"`
	UpdatedAt      string                `json:"updated_at"`
}

// TransactionResponse represents the response body for a transaction
type TransactionResponse struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Type            string  `json:"type"`
	Subtype         string  `json:"subtype"`
	Amount          float64 `json:"amount"`
	TotalNight      *int    `json:"total_night,omitempty"`
	Subtotal        float64 `json:"subtotal"`
	Description     string  `json:"description,omitempty"`
	TransportDetail string  `json:"transport_detail,omitempty"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

// BusinessTripListResponse represents the response for business trip list
type BusinessTripListResponse struct {
	BusinessTrips []BusinessTripResponse `json:"business_trips"`
	Total         int                    `json:"total"`
	Page          int                    `json:"page"`
	Limit         int                    `json:"limit"`
	TotalPages    int                    `json:"total_pages"`
}

// QueryParams represents query parameters for filtering and pagination
type QueryParams struct {
	Page            int    `query:"page"`
	Limit           int    `query:"limit"`
	Search          string `query:"search"`
	DestinationCity string `query:"destination_city"`
	StartDate       string `query:"start_date"`
	EndDate         string `query:"end_date"`
	Status          string `query:"status"`
	SortBy          string `query:"sort_by"`
	SortDirection   string `query:"sort_direction"`
}

func (p QueryParams) Validate() error {
	// Basic validation with invopop
	err := validation.ValidateStruct(&p,
		validation.Field(&p.Page, validation.Min(1)),
		validation.Field(&p.Limit, validation.Min(1), validation.Max(100)),
		validation.Field(&p.StartDate, validation.Length(0, 50)),
		validation.Field(&p.EndDate, validation.Length(0, 50)),
		validation.Field(&p.Status, validation.Length(0, 20)),
		validation.Field(&p.SortDirection, validation.Length(0, 10)),
	)
	if err != nil {
		return err
	}

	// Manual validation for date format and sort direction (custom validation)
	if p.StartDate != "" && !validateDateFormat(p.StartDate) {
		return validation.NewError("start_date", "must be a valid date in YYYY-MM-DD format")
	}
	if p.EndDate != "" && !validateDateFormat(p.EndDate) {
		return validation.NewError("end_date", "must be a valid date in YYYY-MM-DD format")
	}
	if p.SortDirection != "" && p.SortDirection != "asc" && p.SortDirection != "desc" {
		return validation.NewError("sort_direction", "must be 'asc' or 'desc'")
	}

	// Validate status if provided
	if p.Status != "" {
		validStatuses := map[string]bool{
			"draft":           true,
			"ready_to_verify": true,
			"ongoing":         true,
			"completed":       true,
			"canceled":        true,
		}
		if !validStatuses[p.Status] {
			return validation.NewError("status", "must be one of: draft, ready_to_verify, ongoing, completed, canceled")
		}
	}

	return nil
}

func (p QueryParams) SetDefaults() QueryParams {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 {
		p.Limit = 10
	}
	if p.SortBy == "" {
		p.SortBy = "created_at"
	}
	if p.SortDirection == "" {
		p.SortDirection = "desc"
	}
	return p
}

func FromEntity(bt *entity.BusinessTrip) *BusinessTripResponse {
	assignees := make([]AssigneeResponse, len(bt.GetAssignees()))

	for i, assignee := range bt.GetAssignees() {
		transactions := make([]TransactionResponse, len(assignee.GetTransactions()))

		for j, tx := range assignee.GetTransactions() {
			transactions[j] = TransactionResponse{
				ID:              tx.GetID(),
				Name:            tx.GetName(),
				Type:            string(tx.GetType()),
				Subtype:         string(tx.GetSubtype()),
				Amount:          tx.GetAmount(),
				TotalNight:      tx.GetTotalNight(),
				Subtotal:        tx.GetSubtotal(),
				Description:     tx.GetDescription(),
				TransportDetail: tx.GetTransportDetail(),
				CreatedAt:       tx.CreatedAt.Format(time.RFC3339),
				UpdatedAt:       tx.UpdatedAt.Format(time.RFC3339),
			}
		}

		assignees[i] = AssigneeResponse{
			ID:             assignee.GetID(),
			Name:           assignee.GetName(),
			SPDNumber:      assignee.GetSPDNumber(),
			EmployeeID:     assignee.GetEmployeeID(),
			EmployeeName:   assignee.GetEmployeeName(),
			EmployeeNumber: assignee.GetEmployeeNumber(),
			Position:       assignee.GetPosition(),
			Rank:           assignee.GetRank(),
			TotalCost:      assignee.GetTotalCost(),
			Transactions:   transactions,
			CreatedAt:      assignee.CreatedAt.Format(time.RFC3339),
			UpdatedAt:      assignee.UpdatedAt.Format(time.RFC3339),
		}
	}

	// Create verificators response
	verificators := make([]VerificatorResponse, len(bt.GetVerificators()))
	for i, verificator := range bt.GetVerificators() {
		var verifiedAt *string
		if verificator.GetVerifiedAt() != nil {
			verified := verificator.GetVerifiedAt().Format(time.RFC3339)
			verifiedAt = &verified
		}
		verificators[i] = VerificatorResponse{
			ID:                verificator.GetID(),
			BusinessTripID:    verificator.GetBusinessTripID(),
			UserID:            verificator.GetUserID(),
			UserName:          verificator.GetUserName(),
			EmployeeNumber:    verificator.GetEmployeeNumber(),
			Position:          verificator.GetPosition(),
			Status:            string(verificator.GetStatus()),
			VerifiedAt:        verifiedAt,
			VerificationNotes: verificator.GetVerificationNotes(),
			CreatedAt:         verificator.CreatedAt.Format(time.RFC3339),
			UpdatedAt:         verificator.UpdatedAt.Format(time.RFC3339),
		}
	}

	return &BusinessTripResponse{
		ID:                 bt.GetID(),
		BusinessTripNumber: bt.GetBusinessTripNumber(),
		StartDate:          bt.GetStartDate().Format("2006-01-02"),
		EndDate:            bt.GetEndDate().Format("2006-01-02"),
		ActivityPurpose:    bt.GetActivityPurpose(),
		DestinationCity:    bt.GetDestinationCity(),
		SPDDate:            bt.GetSPDDate().Format("2006-01-02"),
		DepartureDate:      bt.GetDepartureDate().Format("2006-01-02"),
		ReturnDate:         bt.GetReturnDate().Format("2006-01-02"),
		Status:             string(bt.GetStatus()),
		DocumentLink:       bt.GetDocumentLink(),
		TotalCost:          bt.GetTotalCost(),
		Verificators:       verificators,
		Assignees:          assignees,
		CreatedAt:          bt.CreatedAt.Format(time.RFC3339),
		UpdatedAt:          bt.UpdatedAt.Format(time.RFC3339),
	}
}

func FromEntities(businessTrips []*entity.BusinessTrip, total int, page, limit int) *BusinessTripListResponse {
	btResponses := make([]BusinessTripResponse, len(businessTrips))

	for i, bt := range businessTrips {
		response := FromEntity(bt)
		btResponses[i] = *response
	}

	totalPages := (total + limit - 1) / limit

	return &BusinessTripListResponse{
		BusinessTrips: btResponses,
		Total:         total,
		Page:          page,
		Limit:         limit,
		TotalPages:    totalPages,
	}
}

// BusinessTripSummary represents the summary of a business trip
type BusinessTripSummary struct {
	BusinessTripID    string             `json:"business_trip_id"`
	TotalCost         float64            `json:"total_cost"`
	TotalAssignees    int                `json:"total_assignees"`
	TotalTransactions int                `json:"total_transactions"`
	CostByType        map[string]float64 `json:"cost_by_type"`
}

// AssigneeSummary represents the summary of an assignee
type AssigneeSummary struct {
	AssigneeID        string             `json:"assignee_id"`
	AssigneeName      string             `json:"assignee_name"`
	TotalCost         float64            `json:"total_cost"`
	TotalTransactions int                `json:"total_transactions"`
	CostByType        map[string]float64 `json:"cost_by_type"`
}
