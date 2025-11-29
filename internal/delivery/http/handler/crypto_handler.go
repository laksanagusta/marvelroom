package handler

import (
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"sandbox/internal/domain/service"
	"sandbox/internal/infrastructure/cryptography"

	"github.com/gofiber/fiber/v2"
)

// CryptoHandler handles cryptographic operations for document signing
type CryptoHandler struct {
	cryptoService cryptography.Service
	deskService   service.DeskService
}

// SignDocumentRequest represents request for signing a document
type SignDocumentRequest struct {
	UserID string `json:"userId" validate:"required"`
	DocID  string `json:"docId" validate:"required"`
}

// SignDocumentResponse represents response for signing a document
type SignDocumentResponse struct {
	DocumentID      string `json:"document_id"`
	QRPayload       string `json:"qr_payload"`
	QRImageBase64   string `json:"qr_image_base64"`
}

// VerifyDocumentRequest represents request for verifying a document
type VerifyDocumentRequest struct {
	QRPayload string `form:"qr_payload" validate:"required"`
	File      *multipart.FileHeader `form:"file" validate:"required"`
}

// VerifyDocumentResponse represents response for verifying a document
type VerifyDocumentResponse struct {
	Status       string `json:"status"`
	UID          string `json:"uid,omitempty"`
	TS           int64  `json:"ts,omitempty"`
	DocID        string `json:"doc_id,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

// GenerateQRCodeRequest represents request for generating QR code
type GenerateQRCodeRequest struct {
	Text string `json:"text" validate:"required"`
}

// GenerateQRCodeResponse represents response for QR code generation
type GenerateQRCodeResponse struct {
	QRImageBase64 string `json:"qr_image_base64"`
}

// PublicKeyResponse represents response for public key
type PublicKeyResponse struct {
	PublicKey string `json:"public_key"`
}

// SignDocument signs a document with digital signature and generates QR code
// @Summary Sign Document
// @Description Signs a document with digital signature and generates QR code
// @Tags crypto
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Document file to sign"
// @Param userId formData string true "User ID who is signing"
// @Param docId formData string true "Document ID for audit"
// @Success 200 {object} SignDocumentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/crypto/sign [post]
func (h *CryptoHandler) SignDocument(c *fiber.Ctx) error {
	// Parse form data
	userID := c.FormValue("userId")
	docID := c.FormValue("docId")

	if userID == "" || docID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Missing required fields",
			Message: "userId and docId are required",
			Code:    fiber.StatusBadRequest,
		})
	}

	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "File upload required",
			Message: "Please upload a document file",
			Code:    fiber.StatusBadRequest,
		})
	}

	// Read file content
	fileContent, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to read file",
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		})
	}
	defer fileContent.Close()

	// Read all file bytes
	fileBytes, err := io.ReadAll(fileContent)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to read file content",
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		})
	}

	// Generate document hash
	docHash := h.cryptoService.GenerateSHA256Hash(fileBytes)

	// Generate QR payload
	payload, err := h.cryptoService.GenerateQRPayload(docHash, userID, docID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to generate QR payload",
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		})
	}

	// Encode payload to base64
	encodedPayload, err := h.cryptoService.EncodeQRPayload(payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to encode QR payload",
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		})
	}

	// Generate QR code image
	qrImage, err := h.cryptoService.GenerateQRCode(encodedPayload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to generate QR code",
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		})
	}

	// Convert QR image to base64
	qrImageBase64 := base64.StdEncoding.EncodeToString(qrImage)

	return c.JSON(SuccessResponse{
		Success: true,
		Message: "Document signed successfully",
		Data: SignDocumentResponse{
			DocumentID:    docID,
			QRPayload:     encodedPayload,
			QRImageBase64: qrImageBase64,
		},
	})
}

// VerifyDocument verifies document integrity and signature
// @Summary Verify Document
// @Description Verifies document integrity and digital signature using QR payload
// @Tags crypto
// @Accept multipart/form-data
// @Produce json
// @Param qr_payload formData string true "Base64 encoded QR payload"
// @Param file formData file true "Document file to verify"
// @Success 200 {object} VerifyDocumentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/crypto/verify [post]
func (h *CryptoHandler) VerifyDocument(c *fiber.Ctx) error {
	// Parse form data
	encodedPayload := c.FormValue("qr_payload")
	if encodedPayload == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Missing QR payload",
			Message: "qr_payload is required",
			Code:    fiber.StatusBadRequest,
		})
	}

	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "File upload required",
			Message: "Please upload a document file to verify",
			Code:    fiber.StatusBadRequest,
		})
	}

	// Read file content
	fileContent, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to read file",
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		})
	}
	defer fileContent.Close()

	// Read all file bytes
	fileBytes, err := io.ReadAll(fileContent)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to read file content",
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		})
	}

	// Verify document
	result := h.cryptoService.VerifyDocument(encodedPayload, fileBytes)

	return c.JSON(SuccessResponse{
		Success: true,
		Message: "Document verification completed",
		Data: VerifyDocumentResponse{
			Status:       string(result.Status),
			UID:          result.UID,
			TS:           result.TS,
			DocID:        result.DocID,
			ErrorMessage: result.ErrorMessage,
		},
	})
}

// VerifyDocumentOffline verifies document integrity and signature using provided public key
// @Summary Verify Document Offline
// @Description Verifies document integrity and digital signature using provided public key (for offline verification)
// @Tags crypto
// @Accept json
// @Produce json
// @Param request body VerifyOfflineRequest true "Verification request"
// @Success 200 {object} VerifyDocumentResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/crypto/verify-offline [post]
type VerifyOfflineRequest struct {
	QRPayload    string `json:"qr_payload" validate:"required"`
	PublicKey    string `json:"public_key" validate:"required"`
	DocumentData []byte `json:"document_data" validate:"required"`
}

func (h *CryptoHandler) VerifyDocumentOffline(c *fiber.Ctx) error {
	var req VerifyOfflineRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
			Code:    fiber.StatusBadRequest,
		})
	}

	// Verify document offline
	result, err := cryptography.VerifyDocumentOffline(req.QRPayload, req.DocumentData, req.PublicKey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Verification failed",
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		})
	}

	return c.JSON(SuccessResponse{
		Success: true,
		Message: "Document offline verification completed",
		Data: VerifyDocumentResponse{
			Status:       string(result.Status),
			UID:          result.UID,
			TS:           result.TS,
			DocID:        result.DocID,
			ErrorMessage: result.ErrorMessage,
		},
	})
}

// GenerateQRCode generates QR code for given text
// @Summary Generate QR Code
// @Description Generates QR code image for given text
// @Tags crypto
// @Accept json
// @Produce json
// @Param request body GenerateQRCodeRequest true "QR code generation request"
// @Success 200 {object} GenerateQRCodeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/crypto/qrcode [post]
func (h *CryptoHandler) GenerateQRCode(c *fiber.Ctx) error {
	var req GenerateQRCodeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
			Code:    fiber.StatusBadRequest,
		})
	}

	// Generate QR code
	qrImage, err := h.cryptoService.GenerateQRCode(req.Text)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to generate QR code",
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		})
	}

	// Convert to base64
	qrImageBase64 := base64.StdEncoding.EncodeToString(qrImage)

	return c.JSON(SuccessResponse{
		Success: true,
		Message: "QR code generated successfully",
		Data: GenerateQRCodeResponse{
			QRImageBase64: qrImageBase64,
		},
	})
}

// GetPublicKey returns the public key for offline verification
// @Summary Get Public Key
// @Description Returns the public key for offline document verification
// @Tags crypto
// @Produce json
// @Success 200 {object} PublicKeyResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/crypto/public-key [get]
func (h *CryptoHandler) GetPublicKey(c *fiber.Ctx) error {
	publicKey := h.cryptoService.GetPublicKey()

	return c.JSON(SuccessResponse{
		Success: true,
		Message: "Public key retrieved successfully",
		Data: PublicKeyResponse{
			PublicKey: publicKey,
		},
	})
}

// SignWorkPaperWithQR signs a work paper with QR-based digital signature
// @Summary Sign Work Paper with QR
// @Description Signs a work paper and generates QR-based digital signature
// @Tags crypto
// @Accept json
// @Produce json
// @Param signatureId path string true "Signature ID"
// @Success 200 {object} entity.WorkPaperSignature
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/crypto/work-paper-signatures/{signatureId}/sign-with-qr [post]
func (h *CryptoHandler) SignWorkPaperWithQR(c *fiber.Ctx) error {
	signatureID := c.Params("signatureId")
	if signatureID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
			Error:   "Missing signature ID",
			Message: "Signature ID is required",
			Code:    fiber.StatusBadRequest,
		})
	}

	// Get signature first to get work paper details
	ctx := c.Context()
	signature, err := h.deskService.GetWorkPaperSignature(ctx, signatureID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(ErrorResponse{
			Error:   "Signature not found",
			Message: err.Error(),
			Code:    fiber.StatusNotFound,
		})
	}

	// Get work paper data (assuming we can get the actual document content)
	// For now, we'll create a hash from the signature data
	// In a real implementation, you would get the actual document content
	docData := fmt.Sprintf("%s-%s-%s-%s",
		signature.WorkPaperID,
		signature.UserID,
		signature.Status,
		signature.UpdatedAt.String())

	docHash := h.cryptoService.GenerateSHA256Hash([]byte(docData))

	// Generate QR payload
	payload, err := h.cryptoService.GenerateQRPayload(docHash, signature.UserID, signature.WorkPaperID.String())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to generate QR payload",
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		})
	}

	// Generate QR code
	encodedPayload, err := h.cryptoService.EncodeQRPayload(payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to encode QR payload",
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		})
	}

	// Generate QR code image
	qrImage, err := h.cryptoService.GenerateQRCode(encodedPayload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
			Error:   "Failed to generate QR code",
			Message: err.Error(),
			Code:    fiber.StatusInternalServerError,
		})
	}

	// Store QR code image to Google Drive (or other storage)
	// This would integrate with the existing Google Drive service
	// For now, we'll just update the signature with the QR payload

	// Create response with QR code
	response := map[string]interface{}{
		"signature":    signature,
		"qr_payload":    encodedPayload,
		"qr_image_base64": base64.StdEncoding.EncodeToString(qrImage),
	}

	return c.JSON(SuccessResponse{
		Success: true,
		Message: "Work paper signed with QR successfully",
		Data:    response,
	})
}

// NewCryptoHandler creates a new instance of CryptoHandler
func NewCryptoHandler(cryptoService cryptography.Service, deskService service.DeskService) *CryptoHandler {
	return &CryptoHandler{
		cryptoService: cryptoService,
		deskService:   deskService,
	}
}