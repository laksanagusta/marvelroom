package excel

type RecapReport struct {
	StartDate            string     `json:"startDate"`
	EndDate              string     `json:"endDate"`
	ActivityPurpose      string     `json:"activityPurpose"` // This maps to Destination in current GenerateRecapExcelRequest
	DestinationCity      string     `json:"destinationCity"`
	SpdDate              string     `json:"spdDate"`
	DepartureDate        string     `json:"departureDate"`
	ReturnDate           string     `json:"returnDate"`
	ReceiptSignatureDate string     `json:"receiptSignatureDate"`
	Assignees            []Assignee `json:"assignees"`
}

type Assignee struct {
	Name           string        `json:"name"`
	SpdNumber      string        `json:"spd_number"`
	EmployeeID     string        `json:"employee_id"`
	EmployeeNumber string        `json:"employee_number"`
	Position       string        `json:"position"`
	Rank           string        `json:"rank"`
	Transactions   []Transaction `json:"transactions"`
}

type Transaction struct {
	Name            string `json:"name"`
	Type            string `json:"type"`
	Subtype         string `json:"subtype"`
	Amount          int32  `json:"amount"`
	TotalNight      *int32 `json:"total_night,omitempty"`
	Subtotal        int32  `json:"subtotal"`
	PaymentType     string `json:"payment_type"`
	Description     string `json:"description"`
	TransportDetail string `json:"transport_detail"`
}
