package cryptography

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"reflect"
	"time"
)

// DigitalSignatureService handles certificate-based digital signatures
type DigitalSignatureService struct {
	privateKeyBytes []byte // PEM bytes from base64 (required)
	publicKeyBytes  []byte // PEM bytes from base64 (required)
}

// NewDigitalSignatureService creates a new instance of DigitalSignatureService
// Accepts base64-decoded PEM bytes for private and public keys
func NewDigitalSignatureService(privateKeyBytes, publicKeyBytes []byte) *DigitalSignatureService {
	return &DigitalSignatureService{
		privateKeyBytes: privateKeyBytes,
		publicKeyBytes:  publicKeyBytes,
	}
}

// SignaturePayload represents the data to be signed
type SignaturePayload struct {
	UserID               string    `json:"user_id"`
	WorkPaperID          string    `json:"work_paper_id"`
	WorkPaperSignatureID string    `json:"work_paper_signature_id"`
	Timestamp            time.Time `json:"timestamp"`
}

// SignatureResult contains the signature data
type SignatureResult struct {
	Signature string    `json:"signature"`
	Payload   string    `json:"payload"`
	Timestamp time.Time `json:"timestamp"`
	Algorithm string    `json:"algorithm"`
}

// loadPrivateKey loads the RSA private key from bytes
func (s *DigitalSignatureService) loadPrivateKey() (*rsa.PrivateKey, error) {
	if len(s.privateKeyBytes) == 0 {
		return nil, fmt.Errorf("private key required: set PRIVATE_KEY_BASE64 environment variable")
	}

	privateBlock, _ := pem.Decode(s.privateKeyBytes)
	if privateBlock == nil || privateBlock.Type != "RSA PRIVATE KEY" && privateBlock.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	var privateKey interface{}
	var err error
	if privateBlock.Type == "RSA PRIVATE KEY" {
		privateKey, err = x509.ParsePKCS1PrivateKey(privateBlock.Bytes)
	} else {
		privateKey, err = x509.ParsePKCS8PrivateKey(privateBlock.Bytes)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("invalid key type: %s", reflect.TypeOf(privateKey))
	}

	return rsaPrivateKey, nil
}

// loadPublicKey loads the RSA public key from bytes
func (s *DigitalSignatureService) loadPublicKey() (*rsa.PublicKey, error) {
	if len(s.publicKeyBytes) == 0 {
		return nil, fmt.Errorf("public key required: set PUBLIC_KEY_BASE64 environment variable")
	}

	publicBlock, _ := pem.Decode(s.publicKeyBytes)
	if publicBlock == nil || publicBlock.Type != "PUBLIC KEY" && publicBlock.Type != "RSA PUBLIC KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}

	var publicKey interface{}
	var err error
	if publicBlock.Type == "RSA PUBLIC KEY" {
		publicKey, err = x509.ParsePKCS1PublicKey(publicBlock.Bytes)
	} else {
		publicKey, err = x509.ParsePKIXPublicKey(publicBlock.Bytes)
	}

	if err != nil {
		return nil, fmt.Errorf("parse public key: %w", err)
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("invalid key type: %s", reflect.TypeOf(publicKey))
	}

	return rsaPublicKey, nil
}

// SignPayload creates a digital signature for the given payload
func (s *DigitalSignatureService) SignPayload(payload *SignaturePayload) (*SignatureResult, error) {
	// Ensure timestamp is set
	if payload.Timestamp.IsZero() {
		payload.Timestamp = time.Now().UTC()
	}

	// Serialize payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize payload: %w", err)
	}

	// Hash the payload
	hash := sha256.Sum256(payloadBytes)

	// Load private key
	privateKey, err := s.loadPrivateKey()
	if err != nil {
		return nil, err
	}

	// Sign the hash with private key using PSS
	signatureRaw, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, hash[:], nil)
	if err != nil {
		return nil, fmt.Errorf("signing error: %w", err)
	}

	// Encode the signature to base64 for easy transport
	signature := base64.StdEncoding.EncodeToString(signatureRaw)

	return &SignatureResult{
		Signature: signature,
		Payload:   base64.StdEncoding.EncodeToString(payloadBytes),
		Timestamp: payload.Timestamp,
		Algorithm: "RSA-PSS-SHA256",
	}, nil
}

// VerifySignature verifies a digital signature
func (s *DigitalSignatureService) VerifySignature(signature string, payload *SignaturePayload) error {
	// Serialize payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to serialize payload: %w", err)
	}

	// Hash the payload
	hash := sha256.Sum256(payloadBytes)

	// Load public key
	publicKey, err := s.loadPublicKey()
	if err != nil {
		return err
	}

	// Decode the signature from base64
	signatureDecoded, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("decode signature: %w", err)
	}

	// Verify the signature with the public key using PSS
	err = rsa.VerifyPSS(publicKey, crypto.SHA256, hash[:], signatureDecoded, nil)
	if err != nil {
		return fmt.Errorf("invalid signature: %w", err)
	}

	return nil
}

// VerifySignatureFromBase64Payload verifies a signature using base64-encoded payload
func (s *DigitalSignatureService) VerifySignatureFromBase64Payload(signature, base64Payload string) error {
	// Decode the payload from base64
	payloadBytes, err := base64.StdEncoding.DecodeString(base64Payload)
	if err != nil {
		return fmt.Errorf("decode payload: %w", err)
	}

	// Parse payload
	var payload SignaturePayload
	err = json.Unmarshal(payloadBytes, &payload)
	if err != nil {
		return fmt.Errorf("parse payload: %w", err)
	}

	return s.VerifySignature(signature, &payload)
}

// CreatePayloadFromData creates a SignaturePayload from individual fields
func CreatePayloadFromData(userID, workPaperID, workPaperSignatureID string) *SignaturePayload {
	return &SignaturePayload{
		UserID:               userID,
		WorkPaperID:          workPaperID,
		WorkPaperSignatureID: workPaperSignatureID,
		Timestamp:            time.Now().UTC(),
	}
}
