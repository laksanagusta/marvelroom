package business_trip

import (
	"context"
	"fmt"

	"sandbox/internal/domain/repository"
)

type GetTransactionUseCase struct {
	businessTripRepo repository.BusinessTripRepository
}

func NewGetTransactionUseCase(businessTripRepo repository.BusinessTripRepository) *GetTransactionUseCase {
	return &GetTransactionUseCase{
		businessTripRepo: businessTripRepo,
	}
}

type GetTransactionResponse struct {
	ID              string  `json:"id"`
	AssigneeID      string  `json:"assigneeId"`
	Name            string  `json:"name"`
	Type            string  `json:"type"`
	Subtype         string  `json:"subtype"`
	Amount          float64 `json:"amount"`
	TotalNight      *int    `json:"totalNight,omitempty"`
	Subtotal        float64 `json:"subtotal"`
	Description     string  `json:"description,omitempty"`
	TransportDetail string  `json:"transportDetail,omitempty"`
	CreatedAt       string  `json:"createdAt"`
	UpdatedAt       string  `json:"updatedAt"`
}

func (uc *GetTransactionUseCase) Execute(ctx context.Context, transactionID string) (*GetTransactionResponse, error) {
	transaction, err := uc.businessTripRepo.GetTransactionByID(ctx, transactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}
	if transaction == nil {
		return nil, fmt.Errorf("transaction not found")
	}

	return &GetTransactionResponse{
		ID:              transaction.ID,
		AssigneeID:      transaction.AssigneeID,
		Name:            transaction.Name,
		Type:            string(transaction.Type),
		Subtype:         string(transaction.Subtype),
		Amount:          transaction.Amount,
		TotalNight:      transaction.TotalNight,
		Subtotal:        transaction.Subtotal,
		Description:     transaction.Description,
		TransportDetail: transaction.TransportDetail,
		CreatedAt:       transaction.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:       transaction.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}