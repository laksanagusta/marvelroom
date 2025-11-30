package entity

import (
	"strings"
)

type ExtractedTransactionType string

const (
	ExtractedTransactionTypeAccommodation ExtractedTransactionType = "accommodation"
	ExtractedTransactionTypeTransport     ExtractedTransactionType = "transport"
	ExtractedTransactionTypeOther         ExtractedTransactionType = "other"
)

type ExtractedTransaction struct {
	Name            string
	TxType          ExtractedTransactionType
	Subtype         string
	Amount          int32
	TotalNight      *int32
	Subtotal        int32
	Description     string
	TransportDetail string
	EmployeeID      string
	Position        string
	Rank            string
}

func NewExtractedTransaction(name, txType, subtype string, amount, subtotal int32, totalNight *int32, description string, transportDetail string, employeeID, position, rank string) (*ExtractedTransaction, error) {
	validType := ExtractedTransactionType(strings.ToLower(txType))
	if !isValidExtractedTransactionType(validType) {
		validType = ExtractedTransactionTypeOther
	}

	return &ExtractedTransaction{
		Name:            strings.TrimSpace(name),
		TxType:          validType,
		Subtype:         strings.TrimSpace(subtype),
		Amount:          amount,
		TotalNight:      totalNight,
		Subtotal:        subtotal,
		Description:     strings.TrimSpace(description),
		TransportDetail: strings.TrimSpace(transportDetail),
		EmployeeID:      strings.TrimSpace(employeeID),
		Position:        strings.TrimSpace(position),
		Rank:            strings.TrimSpace(rank),
	}, nil
}

// Getters
func (t *ExtractedTransaction) GetName() string {
	return t.Name
}

func (t *ExtractedTransaction) GetType() ExtractedTransactionType {
	return t.TxType
}

func (t *ExtractedTransaction) GetSubtype() string {
	return t.Subtype
}

func (t *ExtractedTransaction) GetAmount() int32 {
	return t.Amount
}

func (t *ExtractedTransaction) GetTotalNight() *int32 {
	return t.TotalNight
}

func (t *ExtractedTransaction) GetSubtotal() int32 {
	return t.Subtotal
}

func (t *ExtractedTransaction) GetDescription() string {
	return t.Description
}

func (t *ExtractedTransaction) GetTransportDetail() string {
	return t.TransportDetail
}

func (t *ExtractedTransaction) GetEmployeeID() string {
	return t.EmployeeID
}

func (t *ExtractedTransaction) GetPosition() string {
	return t.Position
}

func (t *ExtractedTransaction) GetRank() string {
	return t.Rank
}

func (t *ExtractedTransaction) IsAccommodation() bool {
	return t.TxType == ExtractedTransactionTypeAccommodation
}

func (t *ExtractedTransaction) IsTransport() bool {
	return t.TxType == ExtractedTransactionTypeTransport
}

func (t *ExtractedTransaction) CalculateTotal() int32 {
	if t.TotalNight != nil && *t.TotalNight > 0 {
		return t.Amount * *t.TotalNight
	}
	return t.Subtotal
}

func isValidExtractedTransactionType(t ExtractedTransactionType) bool {
	switch t {
	case ExtractedTransactionTypeAccommodation, ExtractedTransactionTypeTransport, ExtractedTransactionTypeOther:
		return true
	}
	return false
}
