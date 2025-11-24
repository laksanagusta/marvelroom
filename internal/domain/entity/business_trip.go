package entity

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

type BusinessTripType string

const (
	BusinessTripTypeDomestic      BusinessTripType = "domestic"
	BusinessTripTypeInternational BusinessTripType = "international"
)

type BusinessTripStatus string

const (
	BusinessTripStatusDraft     BusinessTripStatus = "draft"
	BusinessTripStatusOngoing   BusinessTripStatus = "ongoing"
	BusinessTripStatusCanceled  BusinessTripStatus = "canceled"
	BusinessTripStatusCompleted BusinessTripStatus = "completed"
)

// BusinessTrip represents a business trip entity
type BusinessTrip struct {
	ID                 string             `db:"id"`
	BusinessTripNumber sql.NullString     `db:"business_trip_number"`
	StartDate          time.Time          `db:"start_date"`
	EndDate            time.Time          `db:"end_date"`
	ActivityPurpose    string             `db:"activity_purpose"`
	DestinationCity    string             `db:"destination_city"`
	SPDDate            time.Time          `db:"spd_date"`
	DepartureDate      time.Time          `db:"departure_date"`
	ReturnDate         time.Time          `db:"return_date"`
	Status             BusinessTripStatus `db:"status"`
	DocumentLink       sql.NullString     `db:"document_link"`
	Assignees          []*Assignee        `db:"-"`
	CreatedAt          time.Time          `db:"created_at"`
	UpdatedAt          time.Time          `db:"updated_at"`
}

// Assignee represents an employee assigned to a business trip
type Assignee struct {
	ID             string         `db:"id"`
	BusinessTripID string         `db:"business_trip_id"`
	Name           string         `db:"name"`
	SPDNumber      string         `db:"spd_number"`
	EmployeeID     string         `db:"employee_id"`
	EmployeeName   string         `db:"employee_name"`
	EmployeeNumber string         `db:"employee_number"`
	Position       string         `db:"position"`
	Rank           string         `db:"rank"`
	Transactions   []*Transaction `db:"-"`
	CreatedAt      time.Time      `db:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at"`
}

type TransactionType string

const (
	TransactionTypeAccommodation TransactionType = "accommodation"
	TransactionTypeTransport     TransactionType = "transport"
	TransactionTypeOther         TransactionType = "other"
	TransactionTypeAllowance     TransactionType = "allowance"
)

type TransactionSubtype string

const (
	TransactionSubtypeHotel          TransactionSubtype = "hotel"
	TransactionSubtypeFlight         TransactionSubtype = "flight"
	TransactionSubtypeTrain          TransactionSubtype = "train"
	TransactionSubtypeTaxi           TransactionSubtype = "taxi"
	TransactionSubtypeDailyAllowance TransactionSubtype = "daily_allowance"
	TransactionSubtypeRentalCar      TransactionSubtype = "rental_car"
	TransactionSubtypeMeal           TransactionSubtype = "meal"
	TransactionSubtypeOther          TransactionSubtype = "other"
)

// Transaction represents a transaction for an assignee
type Transaction struct {
	ID              string             `db:"id"`
	AssigneeID      string             `db:"assignee_id"`
	Name            string             `db:"name"`
	Type            TransactionType    `db:"type"`
	Subtype         TransactionSubtype `db:"subtype"`
	Amount          float64            `db:"amount"`
	TotalNight      *int               `db:"total_night"`
	Subtotal        float64            `db:"subtotal"`
	Description     string             `db:"description"`
	TransportDetail string             `db:"transport_detail"`
	CreatedAt       time.Time          `db:"created_at"`
	UpdatedAt       time.Time          `db:"updated_at"`
}

// NewBusinessTrip creates a new business trip with validation
func NewBusinessTrip(startDate, endDate, spdDate, departureDate, returnDate time.Time, activityPurpose, destinationCity string) (*BusinessTrip, error) {
	// Validation
	if startDate.After(endDate) {
		return nil, errors.New("start date must be before or equal to end date")
	}

	if departureDate.After(returnDate) {
		return nil, errors.New("departure date must be before or equal to return date")
	}

	if spdDate.After(departureDate) {
		return nil, errors.New("SPD date must be before or equal to departure date")
	}

	if strings.TrimSpace(activityPurpose) == "" {
		return nil, errors.New("activity purpose is required")
	}

	if strings.TrimSpace(destinationCity) == "" {
		return nil, errors.New("destination city is required")
	}

	return &BusinessTrip{
		BusinessTripNumber: sql.NullString{},
		StartDate:          startDate,
		EndDate:            endDate,
		SPDDate:            spdDate,
		DepartureDate:      departureDate,
		ReturnDate:         returnDate,
		ActivityPurpose:    strings.TrimSpace(activityPurpose),
		DestinationCity:    strings.TrimSpace(destinationCity),
		Status:             BusinessTripStatusDraft,
		DocumentLink:       sql.NullString{},
		Assignees:          make([]*Assignee, 0),
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}, nil
}

func (bt *BusinessTrip) AddAssignee(name, spdNumber, employeeID, employeeName, employeeNumber, position, rank string) (*Assignee, error) {
	// Validation
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("assignee name is required")
	}

	if strings.TrimSpace(spdNumber) == "" {
		return nil, errors.New("SPD number is required")
	}

	if strings.TrimSpace(employeeNumber) == "" {
		return nil, errors.New("employee number is required")
	}

	if strings.TrimSpace(position) == "" {
		return nil, errors.New("position is required")
	}

	if strings.TrimSpace(rank) == "" {
		return nil, errors.New("rank is required")
	}

	// Check if SPD number already exists for this business trip
	for _, assignee := range bt.Assignees {
		if assignee.SPDNumber == spdNumber {
			return nil, fmt.Errorf("SPD number %s already exists for this business trip", spdNumber)
		}
	}

	assignee := &Assignee{
		Name:           strings.TrimSpace(name),
		SPDNumber:      strings.TrimSpace(spdNumber),
		EmployeeID:     strings.TrimSpace(employeeID),
		EmployeeName:   strings.TrimSpace(employeeName),
		EmployeeNumber: strings.TrimSpace(employeeNumber),
		Position:       strings.TrimSpace(position),
		Rank:           strings.TrimSpace(rank),
		Transactions:   make([]*Transaction, 0),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	bt.Assignees = append(bt.Assignees, assignee)
	bt.UpdatedAt = time.Now()

	return assignee, nil
}

// NewTransaction creates a new transaction with validation
func NewTransaction(name string, txType TransactionType, subtype TransactionSubtype, amount, subtotal float64, totalNight *int, description, transportDetail string) (*Transaction, error) {
	// Validation
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("transaction name is required")
	}

	if !isValidTransactionType(txType) {
		return nil, fmt.Errorf("invalid transaction type: %s", txType)
	}

	if amount < 0 {
		return nil, errors.New("amount must be non-negative")
	}

	if subtotal < 0 {
		return nil, errors.New("subtotal must be non-negative")
	}

	if totalNight != nil && *totalNight < 0 {
		return nil, errors.New("total night must be non-negative")
	}

	// Calculate subtotal if not provided
	if txType == TransactionTypeAccommodation && totalNight != nil && *totalNight > 0 {
		subtotal = amount * float64(*totalNight)
	} else {
		subtotal = amount
	}

	transaction := &Transaction{
		Name:            strings.TrimSpace(name),
		Type:            txType,
		Subtype:         subtype,
		Amount:          amount,
		TotalNight:      totalNight,
		Subtotal:        subtotal,
		Description:     strings.TrimSpace(description),
		TransportDetail: strings.TrimSpace(transportDetail),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	return transaction, nil
}

// AddTransaction adds a transaction to an assignee
func (a *Assignee) AddTransaction(transaction *Transaction) error {
	if transaction == nil {
		return errors.New("transaction cannot be nil")
	}

	a.Transactions = append(a.Transactions, transaction)
	a.UpdatedAt = time.Now()

	return nil
}

// GetTotalCost calculates the total cost for the business trip
func (bt *BusinessTrip) GetTotalCost() float64 {
	var total float64
	for _, assignee := range bt.Assignees {
		total += assignee.GetTotalCost()
	}
	return total
}

// GetTotalCost calculates the total cost for an assignee
func (a *Assignee) GetTotalCost() float64 {
	var total float64
	for _, transaction := range a.Transactions {
		total += transaction.Subtotal
	}
	return total
}

// GetTransactionsByType returns transactions filtered by type
func (a *Assignee) GetTransactionsByType(txType TransactionType) []*Transaction {
	var transactions []*Transaction
	for _, transaction := range a.Transactions {
		if transaction.Type == txType {
			transactions = append(transactions, transaction)
		}
	}
	return transactions
}

// IsAccommodation returns true if transaction is accommodation type
func (t *Transaction) IsAccommodation() bool {
	return t.Type == TransactionTypeAccommodation
}

// IsTransport returns true if transaction is transport type
func (t *Transaction) IsTransport() bool {
	return t.Type == TransactionTypeTransport
}

// IsAllowance returns true if transaction is allowance type
func (t *Transaction) IsAllowance() bool {
	return t.Type == TransactionTypeAllowance
}

// CalculateTotal calculates the total for the transaction
func (t *Transaction) CalculateTotal() float64 {
	if t.TotalNight != nil && *t.TotalNight > 0 {
		return t.Amount * float64(*t.TotalNight)
	}
	return t.Subtotal
}

// isValidTransactionType checks if the transaction type is valid
func isValidTransactionType(txType TransactionType) bool {
	switch txType {
	case TransactionTypeAccommodation, TransactionTypeTransport, TransactionTypeOther, TransactionTypeAllowance:
		return true
	default:
		return false
	}
}

// isValidBusinessTripStatus checks if the business trip status is valid
func isValidBusinessTripStatus(status BusinessTripStatus) bool {
	switch status {
	case BusinessTripStatusDraft, BusinessTripStatusOngoing, BusinessTripStatusCanceled, BusinessTripStatusCompleted:
		return true
	default:
		return false
	}
}

// CanTransitionTo checks if the business trip can transition to the target status
func (bt *BusinessTrip) CanTransitionTo(targetStatus BusinessTripStatus) bool {
	if !isValidBusinessTripStatus(targetStatus) {
		return false
	}

	currentStatus := bt.Status

	// Cannot transition from completed status
	if currentStatus == BusinessTripStatusCompleted {
		return false
	}

	// Allow transitions based on current status
	switch currentStatus {
	case BusinessTripStatusDraft:
		// Draft can transition to: ongoing, completed, canceled
		return targetStatus == BusinessTripStatusOngoing || targetStatus == BusinessTripStatusCompleted || targetStatus == BusinessTripStatusCanceled
	case BusinessTripStatusOngoing:
		// Ongoing can transition to: completed, canceled, draft (going back)
		return targetStatus == BusinessTripStatusCompleted || targetStatus == BusinessTripStatusCanceled || targetStatus == BusinessTripStatusDraft
	case BusinessTripStatusCanceled:
		// Canceled can transition to: draft (reactivation)
		return targetStatus == BusinessTripStatusDraft
	default:
		return false
	}
}

// UpdateStatus updates the business trip status with validation
func (bt *BusinessTrip) UpdateStatus(newStatus BusinessTripStatus) error {
	if !bt.CanTransitionTo(newStatus) {
		return fmt.Errorf("cannot transition from %s to %s", bt.Status, newStatus)
	}

	// Additional validation: if transitioning to completed, document link is required
	if newStatus == BusinessTripStatusCompleted {
		if !bt.DocumentLink.Valid || bt.DocumentLink.String == "" {
			return fmt.Errorf("document link is required when marking business trip as completed")
		}
	}

	bt.Status = newStatus
	bt.UpdatedAt = time.Now()
	return nil
}

// UpdateDocumentLink updates the document link
func (bt *BusinessTrip) UpdateDocumentLink(documentLink string) {
	trimmed := strings.TrimSpace(documentLink)
	if trimmed == "" {
		bt.DocumentLink = sql.NullString{}
	} else {
		bt.DocumentLink = sql.NullString{String: trimmed, Valid: true}
	}
	bt.UpdatedAt = time.Now()
}

// SetBusinessTripNumber sets the business trip number
func (bt *BusinessTrip) SetBusinessTripNumber(number string) {
	bt.BusinessTripNumber = sql.NullString{String: number, Valid: true}
	bt.UpdatedAt = time.Now()
}

// Getters
func (bt *BusinessTrip) GetID() string { return bt.ID }

func (bt *BusinessTrip) GetBusinessTripNumber() string {
	if bt.BusinessTripNumber.Valid {
		return bt.BusinessTripNumber.String
	}
	return ""
}
func (bt *BusinessTrip) GetStartDate() time.Time       { return bt.StartDate }
func (bt *BusinessTrip) GetEndDate() time.Time         { return bt.EndDate }
func (bt *BusinessTrip) GetActivityPurpose() string    { return bt.ActivityPurpose }
func (bt *BusinessTrip) GetDestinationCity() string    { return bt.DestinationCity }
func (bt *BusinessTrip) GetSPDDate() time.Time         { return bt.SPDDate }
func (bt *BusinessTrip) GetDepartureDate() time.Time   { return bt.DepartureDate }
func (bt *BusinessTrip) GetReturnDate() time.Time      { return bt.ReturnDate }
func (bt *BusinessTrip) GetStatus() BusinessTripStatus { return bt.Status }
func (bt *BusinessTrip) GetDocumentLink() string {
	if bt.DocumentLink.Valid {
		return bt.DocumentLink.String
	}
	return ""
}
func (bt *BusinessTrip) GetAssignees() []*Assignee { return bt.Assignees }

func (a *Assignee) GetID() string                   { return a.ID }
func (a *Assignee) GetBusinessTripID() string       { return a.BusinessTripID }
func (a *Assignee) GetName() string                 { return a.Name }
func (a *Assignee) GetSPDNumber() string            { return a.SPDNumber }
func (a *Assignee) GetEmployeeID() string           { return a.EmployeeID }
func (a *Assignee) GetEmployeeName() string         { return a.EmployeeName }
func (a *Assignee) GetEmployeeNumber() string       { return a.EmployeeNumber }
func (a *Assignee) GetPosition() string             { return a.Position }
func (a *Assignee) GetRank() string                 { return a.Rank }
func (a *Assignee) GetTransactions() []*Transaction { return a.Transactions }

func (t *Transaction) GetID() string                  { return t.ID }
func (t *Transaction) GetAssigneeID() string          { return t.AssigneeID }
func (t *Transaction) GetName() string                { return t.Name }
func (t *Transaction) GetType() TransactionType       { return t.Type }
func (t *Transaction) GetSubtype() TransactionSubtype { return t.Subtype }
func (t *Transaction) GetAmount() float64             { return t.Amount }
func (t *Transaction) GetTotalNight() *int            { return t.TotalNight }
func (t *Transaction) GetSubtotal() float64           { return t.Subtotal }
func (t *Transaction) GetDescription() string         { return t.Description }
func (t *Transaction) GetTransportDetail() string     { return t.TransportDetail }
