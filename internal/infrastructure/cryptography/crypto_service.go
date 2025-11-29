package cryptography

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/skip2/go-qrcode"
)

// CryptoService defines the interface for cryptographic operations
type Service interface {
	// Hashing operations
	GenerateSHA256Hash(data []byte) string

	// Digital signature operations
	SignDocument(documentHash string) (string, error)
	VerifySignature(documentHash, signatureBase64 string) bool

	// QR payload operations
	GenerateQRPayload(docHash, userID, docID string) (*QRPayload, error)
	EncodeQRPayload(payload *QRPayload) (string, error)
	DecodeQRPayload(encodedPayload string) (*QRPayload, error)

	// QR code generation
	GenerateQRCode(text string) ([]byte, error)

	// Document verification
	VerifyDocument(encodedPayload string, documentData []byte) *VerificationResult

	// Key management
	GetPublicKey() string
}

// QRPayload represents the structure for QR code payload
type QRPayload struct {
	V       int    `json:"v"`       // Version
	DocHash string `json:"doc_hash"` // SHA-256 hash of document
	Sig     string `json:"sig"`      // Base64 encoded signature
	UID     string `json:"uid"`      // User ID who signed
	TS      int64  `json:"ts"`      // Unix timestamp in milliseconds
	DocID   string `json:"doc_id"`   // Document ID for audit/online validation
}

// SignatureStatus represents the verification status
type SignatureStatus string

const (
	StatusValid            SignatureStatus = "VALID"
	StatusInvalidSignature SignatureStatus = "INVALID_SIGNATURE"
	StatusHashMismatch     SignatureStatus = "HASH_MISMATCH"
	StatusPayloadCorrupted SignatureStatus = "PAYLOAD_CORRUPTED"
)

// VerificationResult contains the result of signature verification
type VerificationResult struct {
	Status       SignatureStatus `json:"status"`
	UID          string          `json:"uid,omitempty"`
	TS           int64           `json:"ts,omitempty"`
	DocID        string          `json:"doc_id,omitempty"`
	ErrorMessage string          `json:"error_message,omitempty"`
}

// CryptoService provides cryptographic operations for digital signatures
type CryptoService struct {
	privateKey ed25519.PrivateKey
	publicKey  ed25519.PublicKey
}

// NewCryptoService creates a new instance of CryptoService
// It will generate new keypair if none exists
func NewCryptoService() (*CryptoService, error) {
	// For now, generate a new keypair
	// In production, this should be loaded from secure storage
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate keypair: %w", err)
	}

	return &CryptoService{
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

// NewCryptoServiceWithKeys creates CryptoService with existing keys
func NewCryptoServiceWithKeys(publicKey, privateKey []byte) (*CryptoService, error) {
	if len(publicKey) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("invalid public key size")
	}
	if len(privateKey) != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("invalid private key size")
	}

	return &CryptoService{
		privateKey: ed25519.PrivateKey(privateKey),
		publicKey:  ed25519.PublicKey(publicKey),
	}, nil
}

// GenerateSHA256Hash generates SHA-256 hash of the given data
func (c *CryptoService) GenerateSHA256Hash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// SignDocument signs a document hash and returns base64 encoded signature
func (c *CryptoService) SignDocument(documentHash string) (string, error) {
	signature := ed25519.Sign(c.privateKey, []byte(documentHash))
	return base64.StdEncoding.EncodeToString(signature), nil
}

// VerifySignature verifies if the signature is valid for the given document hash
func (c *CryptoService) VerifySignature(documentHash, signatureBase64 string) bool {
	signature, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		return false
	}

	return ed25519.Verify(c.publicKey, []byte(documentHash), signature)
}

// GenerateQRPayload creates QR payload for signing
func (c *CryptoService) GenerateQRPayload(docHash, userID, docID string) (*QRPayload, error) {
	signature, err := c.SignDocument(docHash)
	if err != nil {
		return nil, fmt.Errorf("failed to sign document: %w", err)
	}

	payload := &QRPayload{
		V:       1,
		DocHash: docHash,
		Sig:     signature,
		UID:     userID,
		TS:      time.Now().UnixMilli(),
		DocID:   docID,
	}

	return payload, nil
}

// EncodeQRPayload encodes QR payload to base64 string
func (c *CryptoService) EncodeQRPayload(payload *QRPayload) (string, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	return base64.StdEncoding.EncodeToString(jsonData), nil
}

// DecodeQRPayload decodes base64 QR payload string
func (c *CryptoService) DecodeQRPayload(encodedPayload string) (*QRPayload, error) {
	jsonData, err := base64.StdEncoding.DecodeString(encodedPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 payload: %w", err)
	}

	var payload QRPayload
	if err := json.Unmarshal(jsonData, &payload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	return &payload, nil
}

// GenerateQRCode generates QR code image as bytes
func (c *CryptoService) GenerateQRCode(text string) ([]byte, error) {
	qrCode, err := qrcode.Encode(text, qrcode.Medium, 256)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	return qrCode, nil
}

// VerifyDocument verifies document integrity and signature
func (c *CryptoService) VerifyDocument(encodedPayload string, documentData []byte) *VerificationResult {
	// Decode payload
	payload, err := c.DecodeQRPayload(encodedPayload)
	if err != nil {
		return &VerificationResult{
			Status:       StatusPayloadCorrupted,
			ErrorMessage: fmt.Sprintf("Failed to decode payload: %v", err),
		}
	}

	// Verify signature
	if !c.VerifySignature(payload.DocHash, payload.Sig) {
		return &VerificationResult{
			Status:       StatusInvalidSignature,
			ErrorMessage: "Signature is invalid",
		}
	}

	// Calculate current document hash
	currentHash := c.GenerateSHA256Hash(documentData)

	// Check if document has been modified
	if currentHash != payload.DocHash {
		return &VerificationResult{
			Status:       StatusHashMismatch,
			ErrorMessage: "Document has been modified",
			DocID:        payload.DocID,
		}
	}

	// Everything is valid
	return &VerificationResult{
		Status: StatusValid,
		UID:    payload.UID,
		TS:     payload.TS,
		DocID:  payload.DocID,
	}
}

// GetPublicKey returns the public key in base64 format for distribution
func (c *CryptoService) GetPublicKey() string {
	return base64.StdEncoding.EncodeToString(c.publicKey)
}

// GetPrivateKey returns the private key in base64 format (use with caution!)
func (c *CryptoService) GetPrivateKey() string {
	return base64.StdEncoding.EncodeToString(c.privateKey)
}

// LoadPublicKeyFromBase64 loads public key from base64 string
func LoadPublicKeyFromBase64(publicKeyBase64 string) (ed25519.PublicKey, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(publicKeyBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %w", err)
	}

	if len(keyBytes) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("invalid public key size: expected %d, got %d", ed25519.PublicKeySize, len(keyBytes))
	}

	return ed25519.PublicKey(keyBytes), nil
}

// LoadPrivateKeyFromBase64 loads private key from base64 string
func LoadPrivateKeyFromBase64(privateKeyBase64 string) (ed25519.PrivateKey, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(privateKeyBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %w", err)
	}

	if len(keyBytes) != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("invalid private key size: expected %d, got %d", ed25519.PrivateKeySize, len(keyBytes))
	}

	return ed25519.PrivateKey(keyBytes), nil
}

// VerifyDocumentOffline verifies document using provided public key (for offline verification)
func VerifyDocumentOffline(encodedPayload string, documentData []byte, publicKeyBase64 string) (*VerificationResult, error) {
	publicKey, err := LoadPublicKeyFromBase64(publicKeyBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to load public key: %w", err)
	}

	// Decode payload
	payload, err := base64.StdEncoding.DecodeString(encodedPayload)
	if err != nil {
		return &VerificationResult{
			Status:       StatusPayloadCorrupted,
			ErrorMessage: fmt.Sprintf("Failed to decode payload: %v", err),
		}, nil
	}

	var qrPayload QRPayload
	if err := json.Unmarshal(payload, &qrPayload); err != nil {
		return &VerificationResult{
			Status:       StatusPayloadCorrupted,
			ErrorMessage: fmt.Sprintf("Failed to unmarshal payload: %v", err),
		}, nil
	}

	// Verify signature
	signature, err := base64.StdEncoding.DecodeString(qrPayload.Sig)
	if err != nil {
		return &VerificationResult{
			Status:       StatusPayloadCorrupted,
			ErrorMessage: fmt.Sprintf("Failed to decode signature: %v", err),
		}, nil
	}

	if !ed25519.Verify(publicKey, []byte(qrPayload.DocHash), signature) {
		return &VerificationResult{
			Status:       StatusInvalidSignature,
			ErrorMessage: "Signature is invalid",
		}, nil
	}

	// Calculate current document hash
	crypto := &CryptoService{}
	currentHash := crypto.GenerateSHA256Hash(documentData)

	// Check if document has been modified
	if currentHash != qrPayload.DocHash {
		return &VerificationResult{
			Status:       StatusHashMismatch,
			ErrorMessage: "Document has been modified",
			DocID:        qrPayload.DocID,
		}, nil
	}

	// Everything is valid
	return &VerificationResult{
		Status: StatusValid,
		UID:    qrPayload.UID,
		TS:     qrPayload.TS,
		DocID:  qrPayload.DocID,
	}, nil
}