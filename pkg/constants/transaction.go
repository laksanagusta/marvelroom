package constants

type TransactionType string

const (
	TransactionTypeAccommodation TransactionType = "accommodation"
	TransactionTypeTransport     TransactionType = "transport"
	TransactionTypeOther         TransactionType = "other"
)
