package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"sandbox/internal/domain/service"
)

// GeminiService implements LLMService interface using Google Gemini API
type GeminiService struct {
	apiKey     string
	httpClient *http.Client
	model      string
}

// NewGeminiService creates a new Gemini service instance
func NewGeminiService(apiKey string) (service.LLMService, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("Gemini API key is required")
	}

	return &GeminiService{
		apiKey:     apiKey,
		httpClient: LongTimeoutHTTPClient(),
		model:      GeminiModels[0], // Using first model (gemini-2.5-flash) as default
	}, nil
}

// CheckDocument performs document checking using Gemini API
func (g *GeminiService) CheckDocument(ctx context.Context, req *service.DocumentCheckRequest) (*service.DocumentCheckResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("document check request is required")
	}

	log.Printf("Gemini CheckDocument called with %d documents", len(req.Documents))
	for i, doc := range req.Documents {
		log.Printf("Document %d: Name='%s', Type='%s', Size=%d bytes", i, doc.Name, doc.Type, len(doc.Data))
	}

	// Build the prompt using shared function
	prompt := BuildDocumentCheckPrompt(req.Number, req.Classification, req.TopicDescription, req.DeskInstruction)

	// Prepare parts for the request
	parts := []GeminiPart{
		{Text: prompt},
	}

	// Add documents as inline data (limit to reasonable number/size)
	maxDocuments := 5              // Limit number of documents to avoid exceeding token limits
	maxDocSize := 10 * 1024 * 1024 // 10MB per document

	for i, doc := range req.Documents {
		log.Printf("Processing document %d: %s (type: %s, size: %d bytes)", i, doc.Name, doc.Type, len(doc.Data))

		if i >= maxDocuments {
			log.Printf("Skipping document %d: max document limit reached", i)
			break
		}

		if len(doc.Data) > maxDocSize {
			log.Printf("Skipping document %d: file too large (%d bytes > %d bytes)", i, len(doc.Data), maxDocSize)
			continue // Skip oversized documents
		}

		// Encode file data as base64
		mimeType := GetMimeType(doc.Type)
		data := EncodeBase64(doc.Data)

		log.Printf("Adding document %d to Gemini request: MIME type=%s, encoded size=%d chars", i, mimeType, len(data))

		parts = append(parts, GeminiPart{
			InlineData: &GeminiFile{
				MimeType: mimeType,
				Data:     data,
			},
		})
	}

	log.Printf("Total parts being sent to Gemini API: %d (1 text + %d documents)", len(parts), len(parts)-1)

	// Build the request
	geminiReq := GeminiAPIRequest{
		Contents: []GeminiContent{
			{
				Parts: parts,
			},
		},
	}

	// Serialize request
	jsonBody, err := json.Marshal(geminiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make API call
	apiURL := GetGeminiModelURL(g.model, g.apiKey)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := g.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call Gemini API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Gemini API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse response
	var geminiResp GeminiAPIResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return nil, fmt.Errorf("failed to parse Gemini API response: %w", err)
	}

	rawText, err := GetTextFromGeminiResponse(&geminiResp)
	if err != nil {
		return nil, err
	}

	// Parse the structured response
	result, err := g.parseLLMResponse(rawText)
	if err != nil {
		return nil, fmt.Errorf("failed to parse LLM response: %w", err)
	}

	result.Model = g.model

	return result, nil
}

// parseLLMResponse parses the structured response from LLM
func (g *GeminiService) parseLLMResponse(rawText string) (*service.DocumentCheckResponse, error) {
	jsonStr, err := ExtractJSONFromText(rawText)
	if err != nil {
		return g.parseFallbackResponse(rawText)
	}

	// Parse JSON
	var response struct {
		IsValid bool   `json:"isValid"`
		Note    string `json:"note"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
		// Fallback: try to extract boolean and note from text
		return g.parseFallbackResponse(rawText)
	}

	return &service.DocumentCheckResponse{
		IsValid: response.IsValid,
		Notes:   response.Note,
	}, nil
}

// parseFallbackResponse attempts to extract meaning when JSON parsing fails
func (g *GeminiService) parseFallbackResponse(text string) (*service.DocumentCheckResponse, error) {
	lowerText := CleanJSON(text)

	isValid := true
	if contains(lowerText, "tidak memenuhi", "tidak lengkap", "belum ada", "perlu perbaikan", "invalid", "false") {
		isValid = false
	}

	notes := text
	if len(notes) > 1000 {
		notes = notes[:1000] + "..."
	}

	return &service.DocumentCheckResponse{
		IsValid: isValid,
		Notes:   notes,
	}, nil
}

// contains checks if text contains any of the given substrings
func contains(text string, substrings ...string) bool {
	for _, s := range substrings {
		if len(s) > 0 && len(text) >= len(s) {
			for i := 0; i <= len(text)-len(s); i++ {
				if text[i:i+len(s)] == s {
					return true
				}
			}
		}
	}
	return false
}
