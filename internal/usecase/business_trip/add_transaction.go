package business_trip

import (
	"context"
	"time"

	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/repository"
)

type AddTransactionUseCase struct {
	businessTripRepo repository.BusinessTripRepository
	assigneeRepo     repository.AssigneeRepository
}

func NewAddTransactionUseCase(businessTripRepo repository.BusinessTripRepository, assigneeRepo repository.AssigneeRepository) *AddTransactionUseCase {
	return &AddTransactionUseCase{
		businessTripRepo: businessTripRepo,
		assigneeRepo:     assigneeRepo,
	}
}

func (uc *AddTransactionUseCase) Execute(ctx context.Context, assigneeID string, req TransactionRequest) (*TransactionResponse, error) {
	assignee, err := uc.assigneeRepo.GetAssigneeByID(ctx, assigneeID)
	if err != nil {
		return nil, err
	}
	if assignee == nil {
		return nil, entity.ErrAssigneeNotFound
	}

	transaction := &entity.Transaction{
		Name:            req.Name,
		Type:            entity.TransactionType(req.Type),
		Subtype:         entity.TransactionSubtype(req.Subtype),
		Amount:          req.Amount,
		TotalNight:      req.TotalNight,
		Description:     req.Description,
		TransportDetail: req.TransportDetail,
		AssigneeID:      assigneeID,
	}

	createdTransaction, err := uc.businessTripRepo.CreateTransaction(ctx, transaction)
	if err != nil {
		return nil, err
	}

	return &TransactionResponse{
		ID:              createdTransaction.GetID(),
		Name:            createdTransaction.GetName(),
		Type:            string(createdTransaction.GetType()),
		Subtype:         string(createdTransaction.GetSubtype()),
		Amount:          createdTransaction.GetAmount(),
		TotalNight:      createdTransaction.GetTotalNight(),
		Subtotal:        createdTransaction.GetSubtotal(),
		Description:     createdTransaction.GetDescription(),
		TransportDetail: createdTransaction.GetTransportDetail(),
		CreatedAt:       createdTransaction.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       createdTransaction.UpdatedAt.Format(time.RFC3339),
	}, nil
}
