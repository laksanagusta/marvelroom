package transaction

import (
	"context"
	"fmt"

	"sandbox/internal/infrastructure/excel"
)

type GenerateRecapExcelUseCase struct {
	excelGenerator *excel.Generator
}

func NewGenerateRecapExcelUseCase(excelGenerator *excel.Generator) *GenerateRecapExcelUseCase {
	return &GenerateRecapExcelUseCase{
		excelGenerator: excelGenerator,
	}
}

func (uc *GenerateRecapExcelUseCase) Execute(ctx context.Context, req RecapReportDTO) (*GenerateRecapExcelResponse, error) {
	recapReport := req.ToRecapReport()

	excelBuffer, err := uc.excelGenerator.GenerateRecapExcel(recapReport)
	if err != nil {
		return nil, fmt.Errorf("failed to generate excel file: %w", err)
	}

	return &GenerateRecapExcelResponse{
		FileContent: excelBuffer.Bytes(),
	}, nil
}
