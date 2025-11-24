package handler

import (
	"log"

	transactionUC "sandbox/internal/usecase/transaction"
	"sandbox/internal/infrastructure/file"

	"github.com/gofiber/fiber/v2"
)

// TransactionHandler handles HTTP requests for transactions
type TransactionHandler struct {
	extractUseCase            *transactionUC.ExtractTransactionsUseCase
	fileProcessor             *file.Processor
	generateRecapExcelUseCase *transactionUC.GenerateRecapExcelUseCase
}

// NewTransactionHandler creates a new transaction handler
func NewTransactionHandler(
	extractUseCase *transactionUC.ExtractTransactionsUseCase,
	fileProcessor *file.Processor,
	generateRecapExcelUseCase *transactionUC.GenerateRecapExcelUseCase,
) *TransactionHandler {
	return &TransactionHandler{
		extractUseCase:            extractUseCase,
		fileProcessor:             fileProcessor,
		generateRecapExcelUseCase: generateRecapExcelUseCase,
	}
}

// UploadAndExtract handles the file upload and extraction endpoint
func (h *TransactionHandler) UploadAndExtract(c *fiber.Ctx) error {
	log.Println("Processing upload request")

	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse form data",
		})
	}

	fileHeaders := form.File["file"]
	if len(fileHeaders) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No files uploaded",
		})
	}

	// Process uploaded files
	processedFiles, err := h.fileProcessor.ProcessMultipleFiles(fileHeaders)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Convert to DTO
	fileUploads := make([]transactionUC.FileUpload, len(processedFiles))
	for i, pf := range processedFiles {
		fileUploads[i] = transactionUC.FileUpload{
			Content:  pf.Content,
			Filename: pf.Filename,
			MimeType: pf.MimeType,
		}
	}

	// Execute use case
	request := transactionUC.ExtractTransactionsRequest{
		Files: fileUploads,
	}

	response, err := h.extractUseCase.Execute(c.Context(), request)
	if err != nil {
		log.Printf("Error extracting transactions: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to extract transactions",
			"details": err.Error(),
		})
	}

	// Return the complete report structure as requested
	return c.JSON(response.Report)
}

func (h *TransactionHandler) GenerateRecapExcel(c *fiber.Ctx) error {
	log.Println("Generating Excel recap file")

	var reqBody transactionUC.RecapReportDTO
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	// Validate request using DTO validation method
	if err := reqBody.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	response, err := h.generateRecapExcelUseCase.Execute(c.Context(), reqBody)
	if err != nil {
		log.Printf("Error generating Excel recap: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to generate Excel recap file",
			"details": err.Error(),
		})
	}

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", "attachment; filename=\"kwitansi-perjadin.xlsx\"")
	return c.Send(response.FileContent)
}

// UploadAndExtractDetailed returns detailed response with count
func (h *TransactionHandler) UploadAndExtractDetailed(c *fiber.Ctx) error {
	log.Println("Processing upload request (detailed)")

	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse form data",
		})
	}

	fileHeaders := form.File["file"]
	if len(fileHeaders) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No files uploaded",
		})
	}

	// Process uploaded files
	processedFiles, err := h.fileProcessor.ProcessMultipleFiles(fileHeaders)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Convert to DTO
	fileUploads := make([]transactionUC.FileUpload, len(processedFiles))
	for i, pf := range processedFiles {
		fileUploads[i] = transactionUC.FileUpload{
			Content:  pf.Content,
			Filename: pf.Filename,
			MimeType: pf.MimeType,
		}
	}

	// Execute use case
	request := transactionUC.ExtractTransactionsRequest{
		Files: fileUploads,
	}

	response, err := h.extractUseCase.Execute(c.Context(), request)
	if err != nil {
		log.Printf("Error extracting transactions: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to extract transactions",
			"details": err.Error(),
		})
	}

	// Return full response
	return c.JSON(response)
}
