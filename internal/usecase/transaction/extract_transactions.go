package transaction

import (
	"context"
	"errors"

	"sandbox/internal/domain/repository"
	"sandbox/internal/domain/service"
)

type ExtractTransactionsUseCase struct {
	transactionService *service.TransactionService
}

func NewExtractTransactionsUseCase(transactionService *service.TransactionService) *ExtractTransactionsUseCase {
	return &ExtractTransactionsUseCase{
		transactionService: transactionService,
	}
}

func (uc *ExtractTransactionsUseCase) Execute(ctx context.Context, req ExtractTransactionsRequest) (*ExtractTransactionsResponse, error) {
	documents := make([]repository.Document, len(req.Files))
	for i, file := range req.Files {
		documents[i] = repository.Document{
			Content:  file.Content,
			MimeType: file.MimeType,
			Filename: file.Filename,
		}
	}

	result, err := uc.transactionService.ExtractTransactions(ctx, documents, "scanBusinessTripDocs")
	if err != nil {
		return nil, err
	}

	// Type assert to RecapReportDTO
	recapReport, ok := result.(*RecapReportDTO)
	if !ok {
		return nil, errors.New("extractor returned unexpected type")
	}

	return &ExtractTransactionsResponse{
		Report: *recapReport,
	}, nil
}
