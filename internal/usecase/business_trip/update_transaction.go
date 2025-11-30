package business_trip

import (
	"context"
	"fmt"
	"strings"

	"github.com/invopop/validation"

	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/repository"
)

type UpdateTransactionUseCase struct {
	businessTripRepo repository.BusinessTripRepository
	assigneeRepo     repository.AssigneeRepository
}

func NewUpdateTransactionUseCase(businessTripRepo repository.BusinessTripRepository, assigneeRepo repository.AssigneeRepository) *UpdateTransactionUseCase {
	return &UpdateTransactionUseCase{
		businessTripRepo: businessTripRepo,
		assigneeRepo:     assigneeRepo,
	}
}

type UpdateTransactionRequest struct {
	BusinessTripID  string  `params:"businessTripId" json:"businessTripId"`
	AssigneeID      string  `params:"assigneeId" json:"assigneeId"`
	TransactionID   string  `params:"transactionId" json:"transactionId"`
	Name            string  `json:"name"`
	Type            string  `json:"type"`
	Subtype         string  `json:"subtype"`
	Amount          float64 `json:"amount"`
	TotalNight      *int    `json:"totalNight"`
	Description     string  `json:"description"`
	TransportDetail string  `json:"transportDetail"`
}

func (r UpdateTransactionRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.BusinessTripID, validation.Required),
		validation.Field(&r.AssigneeID, validation.Required),
		validation.Field(&r.TransactionID, validation.Required),
		validation.Field(&r.Name, validation.Required, validation.Length(1, 255)),
		validation.Field(&r.Type, validation.Required, validation.In("accommodation", "transport", "other", "allowance")),
		validation.Field(&r.Amount, validation.Required, validation.Min(0)),
	)
}

type UpdateTransactionResponse struct {
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

func (uc *UpdateTransactionUseCase) Execute(ctx context.Context, req UpdateTransactionRequest) (*UpdateTransactionResponse, error) {
	businessTrip, err := uc.businessTripRepo.GetByID(ctx, req.BusinessTripID)
	if err != nil {
		return nil, fmt.Errorf("failed to get business trip: %w", err)
	}
	if businessTrip == nil {
		return nil, fmt.Errorf("business trip not found")
	}

	assignee, err := uc.assigneeRepo.GetAssigneeByID(ctx, req.AssigneeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get assignee: %w", err)
	}
	if assignee == nil {
		return nil, fmt.Errorf("assignee not found")
	}

	if assignee.BusinessTripID != req.BusinessTripID {
		return nil, fmt.Errorf("assignee does not belong to the specified business trip")
	}

	transaction, err := uc.businessTripRepo.GetTransactionByID(ctx, req.TransactionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}
	if transaction == nil {
		return nil, fmt.Errorf("transaction not found")
	}

	if transaction.AssigneeID != req.AssigneeID {
		return nil, fmt.Errorf("transaction does not belong to the specified assignee")
	}

	txType := entity.TransactionType(req.Type)
	subtotal := req.Amount

	// For accommodation type, calculate subtotal based on total night
	if txType == entity.TransactionTypeAccommodation && req.TotalNight != nil && *req.TotalNight > 0 {
		subtotal = req.Amount * float64(*req.TotalNight)
	}

	// Update transaction details
	transaction.Name = strings.TrimSpace(req.Name)
	transaction.Type = txType
	transaction.Subtype = entity.TransactionSubtype(req.Subtype)
	transaction.Amount = req.Amount
	transaction.TotalNight = req.TotalNight
	transaction.Subtotal = subtotal
	transaction.Description = strings.TrimSpace(req.Description)
	transaction.TransportDetail = strings.TrimSpace(req.TransportDetail)

	// Save updated transaction
	updatedTransaction, err := uc.businessTripRepo.UpdateTransaction(ctx, transaction)
	if err != nil {
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}

	return &UpdateTransactionResponse{
		ID:              updatedTransaction.ID,
		AssigneeID:      updatedTransaction.AssigneeID,
		Name:            updatedTransaction.Name,
		Type:            string(updatedTransaction.Type),
		Subtype:         string(updatedTransaction.Subtype),
		Amount:          updatedTransaction.Amount,
		TotalNight:      updatedTransaction.TotalNight,
		Subtotal:        updatedTransaction.Subtotal,
		Description:     updatedTransaction.Description,
		TransportDetail: updatedTransaction.TransportDetail,
		CreatedAt:       updatedTransaction.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:       updatedTransaction.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}
