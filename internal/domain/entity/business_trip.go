package entity

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type BusinessTripType string

const (
	BusinessTripTypeDomestic      BusinessTripType = "domestic"
	BusinessTripTypeInternational BusinessTripType = "international"
)

type BusinessTripStatus string

const (
	BusinessTripStatusDraft         BusinessTripStatus = "draft"
	BusinessTripStatusReadyToVerify BusinessTripStatus = "ready_to_verify"
	BusinessTripStatusOngoing       BusinessTripStatus = "ongoing"
	BusinessTripStatusCanceled      BusinessTripStatus = "canceled"
	BusinessTripStatusCompleted     BusinessTripStatus = "completed"
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
	Verificators       []*Verificator     `db:"-"`
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
		ID:                 uuid.NewString(),
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

// AddVerificator adds a verificator to the business trip
func (bt *BusinessTrip) AddVerificator(userID, userName, employeeNumber, position string) (*Verificator, error) {
	// Validation
	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("verificator user ID is required")
	}

	if strings.TrimSpace(userName) == "" {
		return nil, errors.New("verificator name is required")
	}

	if strings.TrimSpace(employeeNumber) == "" {
		return nil, errors.New("verificator employee number is required")
	}

	if strings.TrimSpace(position) == "" {
		return nil, errors.New("verificator position is required")
	}

	// Check if user is already assigned as verificator for this business trip
	for _, verificator := range bt.Verificators {
		if verificator.UserID == userID {
			return nil, fmt.Errorf("user %s is already assigned as verificator for this business trip", userID)
		}
	}

	verificator := &Verificator{
		ID:                uuid.New().String(),
		BusinessTripID:    bt.ID,
		UserID:            strings.TrimSpace(userID),
		UserName:          strings.TrimSpace(userName),
		EmployeeNumber:    strings.TrimSpace(employeeNumber),
		Position:          strings.TrimSpace(position),
		Status:            VerificatorStatusPending,
		VerificationNotes: "",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	bt.Verificators = append(bt.Verificators, verificator)
	bt.UpdatedAt = time.Now()

	return verificator, nil
}

// RemoveVerificator removes a verificator from the business trip
func (bt *BusinessTrip) RemoveVerificator(userID string) error {
	if strings.TrimSpace(userID) == "" {
		return errors.New("verificator user ID is required")
	}

	for i, verificator := range bt.Verificators {
		if verificator.UserID == userID {
			// Remove from slice
			bt.Verificators = append(bt.Verificators[:i], bt.Verificators[i+1:]...)
			bt.UpdatedAt = time.Now()
			return nil
		}
	}

	return fmt.Errorf("verificator with user ID %s not found", userID)
}

// GetVerificatorByUserID returns a verificator by user ID
func (bt *BusinessTrip) GetVerificatorByUserID(userID string) *Verificator {
	for _, verificator := range bt.Verificators {
		if verificator.UserID == userID {
			return verificator
		}
	}
	return nil
}

// HasAllVerificatorsApproved returns true if all verificators have approved
func (bt *BusinessTrip) HasAllVerificatorsApproved() bool {
	if len(bt.Verificators) == 0 {
		return true // No verificators means no approval needed
	}

	for _, verificator := range bt.Verificators {
		if !verificator.IsApproved() {
			return false
		}
	}
	return true
}

// HasAnyVerificatorRejected returns true if any verificator has rejected
func (bt *BusinessTrip) HasAnyVerificatorRejected() bool {
	for _, verificator := range bt.Verificators {
		if verificator.IsRejected() {
			return true
		}
	}
	return false
}

// GetPendingVerificatorsCount returns the count of pending verificators
func (bt *BusinessTrip) GetPendingVerificatorsCount() int {
	count := 0
	for _, verificator := range bt.Verificators {
		if verificator.IsPending() {
			count++
		}
	}
	return count
}

// GetApprovedVerificatorsCount returns the count of approved verificators
func (bt *BusinessTrip) GetApprovedVerificatorsCount() int {
	count := 0
	for _, verificator := range bt.Verificators {
		if verificator.IsApproved() {
			count++
		}
	}
	return count
}

// GetRejectedVerificatorsCount returns the count of rejected verificators
func (bt *BusinessTrip) GetRejectedVerificatorsCount() int {
	count := 0
	for _, verificator := range bt.Verificators {
		if verificator.IsRejected() {
			count++
		}
	}
	return count
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
	case BusinessTripStatusDraft, BusinessTripStatusReadyToVerify, BusinessTripStatusOngoing, BusinessTripStatusCanceled, BusinessTripStatusCompleted:
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
		// Draft can transition to: ready_to_verify, ongoing, completed, canceled
		return targetStatus == BusinessTripStatusReadyToVerify || targetStatus == BusinessTripStatusOngoing || targetStatus == BusinessTripStatusCompleted || targetStatus == BusinessTripStatusCanceled
	case BusinessTripStatusReadyToVerify:
		// Ready to verify can transition to: ongoing, draft, canceled
		return targetStatus == BusinessTripStatusOngoing || targetStatus == BusinessTripStatusDraft || targetStatus == BusinessTripStatusCanceled
	case BusinessTripStatusOngoing:
		// Ongoing can transition to: completed, canceled, draft (going back), ready_to_verify
		return targetStatus == BusinessTripStatusCompleted || targetStatus == BusinessTripStatusCanceled || targetStatus == BusinessTripStatusDraft || targetStatus == BusinessTripStatusReadyToVerify
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

// GetVerificators returns the list of verificators
func (bt *BusinessTrip) GetVerificators() []*Verificator { return bt.Verificators }

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

// VerificatorStatus represents verification status
type VerificatorStatus string

const (
	VerificatorStatusPending  VerificatorStatus = "pending"
	VerificatorStatusApproved VerificatorStatus = "approved"
	VerificatorStatusRejected VerificatorStatus = "rejected"
)

// Verificator represents a user assigned to verify a business trip
type Verificator struct {
	ID                string            `db:"id"`
	BusinessTripID    string            `db:"business_trip_id"`
	UserID            string            `db:"user_id"`
	UserName          string            `db:"user_name"`
	EmployeeNumber    string            `db:"employee_number"`
	Position          string            `db:"position"`
	Status            VerificatorStatus `db:"status"`
	VerifiedAt        *time.Time        `db:"verified_at"`
	VerificationNotes string            `db:"verification_notes"`
	CreatedAt         time.Time         `db:"created_at"`
	UpdatedAt         time.Time         `db:"updated_at"`
}

// NewVerificator creates a new verificator with validation
func NewVerificator(businessTripID, userID, userName, employeeNumber, position string) (*Verificator, error) {
	// Validation
	if strings.TrimSpace(businessTripID) == "" {
		return nil, errors.New("business trip ID is required")
	}

	if strings.TrimSpace(userID) == "" {
		return nil, errors.New("user ID is required")
	}

	if strings.TrimSpace(userName) == "" {
		return nil, errors.New("user name is required")
	}

	if strings.TrimSpace(employeeNumber) == "" {
		return nil, errors.New("employee number is required")
	}

	if strings.TrimSpace(position) == "" {
		return nil, errors.New("position is required")
	}

	return &Verificator{
		ID:                uuid.New().String(),
		BusinessTripID:    strings.TrimSpace(businessTripID),
		UserID:            strings.TrimSpace(userID),
		UserName:          strings.TrimSpace(userName),
		EmployeeNumber:    strings.TrimSpace(employeeNumber),
		Position:          strings.TrimSpace(position),
		Status:            VerificatorStatusPending,
		VerificationNotes: "",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}, nil
}

// Approve marks the verificator as approved
func (v *Verificator) Approve(notes string) {
	v.Status = VerificatorStatusApproved
	now := time.Now()
	v.VerifiedAt = &now
	v.VerificationNotes = strings.TrimSpace(notes)
	v.UpdatedAt = time.Now()
}

// Reject marks the verificator as rejected
func (v *Verificator) Reject(notes string) {
	v.Status = VerificatorStatusRejected
	now := time.Now()
	v.VerifiedAt = &now
	v.VerificationNotes = strings.TrimSpace(notes)
	v.UpdatedAt = time.Now()
}

// IsPending returns true if verificator is pending
func (v *Verificator) IsPending() bool {
	return v.Status == VerificatorStatusPending
}

// IsApproved returns true if verificator is approved
func (v *Verificator) IsApproved() bool {
	return v.Status == VerificatorStatusApproved
}

// IsRejected returns true if verificator is rejected
func (v *Verificator) IsRejected() bool {
	return v.Status == VerificatorStatusRejected
}

// isValidVerificatorStatus checks if the verificator status is valid
func isValidVerificatorStatus(status VerificatorStatus) bool {
	switch status {
	case VerificatorStatusPending, VerificatorStatusApproved, VerificatorStatusRejected:
		return true
	default:
		return false
	}
}

// UpdateStatus updates the verificator status with validation
func (v *Verificator) UpdateStatus(newStatus VerificatorStatus, notes string) error {
	if !isValidVerificatorStatus(newStatus) {
		return fmt.Errorf("invalid verificator status: %s", newStatus)
	}

	v.Status = newStatus
	v.VerificationNotes = strings.TrimSpace(notes)

	// Set verified_at timestamp for approved/rejected status
	if newStatus == VerificatorStatusApproved || newStatus == VerificatorStatusRejected {
		now := time.Now()
		v.VerifiedAt = &now
	} else {
		v.VerifiedAt = nil
	}

	v.UpdatedAt = time.Now()
	return nil
}

// Getters
func (v *Verificator) GetID() string                { return v.ID }
func (v *Verificator) GetBusinessTripID() string    { return v.BusinessTripID }
func (v *Verificator) GetUserID() string            { return v.UserID }
func (v *Verificator) GetUserName() string          { return v.UserName }
func (v *Verificator) GetEmployeeNumber() string    { return v.EmployeeNumber }
func (v *Verificator) GetPosition() string          { return v.Position }
func (v *Verificator) GetStatus() VerificatorStatus { return v.Status }
func (v *Verificator) GetVerifiedAt() *time.Time    { return v.VerifiedAt }
func (v *Verificator) GetVerificationNotes() string { return v.VerificationNotes }

// VerificatorWithBusinessTrip represents a verificator with joined business trip data
type VerificatorWithBusinessTrip struct {
	ID                string            `db:"id"`
	BusinessTripID    string            `db:"business_trip_id"`
	UserID            string            `db:"user_id"`
	UserName          string            `db:"user_name"`
	EmployeeNumber    string            `db:"employee_number"`
	Position          string            `db:"position"`
	Status            VerificatorStatus `db:"status"`
	VerifiedAt        *time.Time        `db:"verified_at"`
	VerificationNotes string            `db:"verification_notes"`
	CreatedAt         time.Time         `db:"created_at"`
	UpdatedAt         time.Time         `db:"updated_at"`
	// Business trip fields
	BusinessTripNumber          sql.NullString     `db:"business_trip_number"`
	BusinessTripStartDate       time.Time          `db:"start_date"`
	BusinessTripEndDate         time.Time          `db:"end_date"`
	BusinessTripActivityPurpose string             `db:"activity_purpose"`
	BusinessTripDestinationCity string             `db:"destination_city"`
	BusinessTripSPDDate         time.Time          `db:"spd_date"`
	BusinessTripDepartureDate   time.Time          `db:"departure_date"`
	BusinessTripReturnDate      time.Time          `db:"return_date"`
	BusinessTripStatus          BusinessTripStatus `db:"business_trip_status"`
	BusinessTripDocumentLink    sql.NullString     `db:"document_link"`
}
