package service

import (
	"context"
	"errors"

	"sandbox/internal/domain/repository"
)

// TransactionService provides domain business logic for transactions
type TransactionService struct {
	extractor repository.ExtractorRepository
}

// NewTransactionService creates a new transaction service
func NewTransactionService(extractor repository.ExtractorRepository) *TransactionService {
	return &TransactionService{
		extractor: extractor,
	}
}

// ExtractTransactions extracts transactions from multiple documents
func (s *TransactionService) ExtractTransactions(ctx context.Context, documents []repository.Document, promptType string) (interface{}, error) {
	if len(documents) == 0 {
		return nil, errors.New("no documents provided")
	}

	result, err := s.extractor.ExtractFromDocuments(ctx, documents, promptType)
	if err != nil {
		return nil, err
	}

	return result, nil
}
