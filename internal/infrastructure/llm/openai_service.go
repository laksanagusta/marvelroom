package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"sandbox/internal/domain/service"
)

// OpenAIService implements LLMService interface using OpenAI API
type OpenAIService struct {
	apiKey     string
	httpClient *http.Client
	model      string
}

// NewOpenAIService creates a new OpenAI service instance
func NewOpenAIService(apiKey string) (service.LLMService, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required")
	}

	return &OpenAIService{
		apiKey:     apiKey,
		httpClient: DefaultHTTPClient(),
		model:      "gpt-4o-mini",
	}, nil
}

// CheckDocument performs document checking using OpenAI API
func (o *OpenAIService) CheckDocument(ctx context.Context, req *service.DocumentCheckRequest) (*service.DocumentCheckResponse, error) {
	if req == nil {
		return nil, fmt.Errorf("document check request is required")
	}

	log.Printf("[OpenAI] CheckDocument called with %d documents", len(req.Documents))
	for i, doc := range req.Documents {
		log.Printf("[OpenAI] Document %d: Name='%s', Type='%s', Size=%d bytes", i, doc.Name, doc.Type, len(doc.Data))
	}

	// Build the prompt
	prompt := o.buildPrompt(req.Number, req.Statement, req.Explanation, req.FillingGuide)

	// Prepare content for the request
	content := []OpenAIContent{
		{Type: "text", Text: prompt},
	}

	// Add documents as images (for image types) or describe them for other types
	maxDocuments := 5              // Limit number of documents
	maxDocSize := 10 * 1024 * 1024 // 10MB per document

	for i, doc := range req.Documents {
		log.Printf("[OpenAI] Processing document %d: %s (type: %s, size: %d bytes)", i, doc.Name, doc.Type, len(doc.Data))

		if i >= maxDocuments {
			log.Printf("[OpenAI] Skipping document %d: max document limit reached", i)
			break
		}

		if len(doc.Data) > maxDocSize {
			log.Printf("[OpenAI] Skipping document %d: file too large (%d bytes > %d bytes)", i, len(doc.Data), maxDocSize)
			continue
		}

		mimeType := GetMimeType(doc.Type)

		// Only include image types as vision content
		if strings.HasPrefix(mimeType, "image/") {
			base64Data := EncodeBase64(doc.Data)
			dataURL := fmt.Sprintf("data:%s;base64,%s", mimeType, base64Data)

			content = append(content, OpenAIContent{
				Type: "image_url",
				ImageURL: &OpenAIImageURL{
					URL: dataURL,
				},
			})
			log.Printf("[OpenAI] Added image document %d to request", i)
		} else {
			// For non-image files, add a note about the document
			content = append(content, OpenAIContent{
				Type: "text",
				Text: fmt.Sprintf("\n[Dokumen: %s (tipe: %s, ukuran: %d bytes)]", doc.Name, doc.Type, len(doc.Data)),
			})
			log.Printf("[OpenAI] Added document reference %d (non-image)", i)
		}
	}

	log.Printf("[OpenAI] Total content items being sent: %d", len(content))

	// Build the request
	openAIReq := OpenAIRequest{
		Model: o.model,
		Messages: []OpenAIMessage{
			{
				Role:    "user",
				Content: content,
			},
		},
		MaxTokens: 4096,
	}

	// Serialize request
	jsonBody, err := json.Marshal(openAIReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Make API call
	httpReq, err := http.NewRequestWithContext(ctx, "POST", OpenAIAPIURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", o.apiKey))

	resp, err := o.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to call OpenAI API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAI API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse response
	var openAIResp OpenAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI API response: %w", err)
	}

	if openAIResp.Error != nil {
		return nil, fmt.Errorf("OpenAI API error: %s", openAIResp.Error.Message)
	}

	rawText, err := GetTextFromOpenAIResponse(&openAIResp)
	if err != nil {
		return nil, err
	}

	// Parse the structured response
	result, err := o.parseLLMResponse(rawText)
	if err != nil {
		return nil, fmt.Errorf("failed to parse LLM response: %w", err)
	}

	result.Model = o.model

	return result, nil
}

// buildPrompt constructs the prompt for LAKIP document checking
func (o *OpenAIService) buildPrompt(number, statement, explanation, fillingGuide string) string {
	prompt := fmt.Sprintf(`Periksa dokumen yang diberikan sesuai dengan poin kertas kerja LAKIP berikut:

Nomor: %s
Pernyataan: %s
Penjelasan: %s
Petunjuk Pengisian: %s

Tugas Anda:
1. Analisis semua dokumen yang disediakan
2. Periksa apakah dokumen memenuhi persyaratan yang disebutkan dalam pernyataan dan penjelasan
3. Pertimbangkan petunjuk pengisian dalam evaluasi Anda
4. Berikan penilaian objektif tentang kelengkapan dan kepatuhan dokumen

Jawab dengan format JSON berikut:
{
  "isValid": true/false,
  "note": "Penjelasan tentang temuan, rekomendasi, atau alasan penilaian (maksimal 5-6 kalimat)"
}

Kriteria penilaian:
- isValid: true jika dokumen lengkap dan memenuhi persyaratan
- isValid: false jika dokumen tidak lengkap, tidak memenuhi persyaratan, atau ada masalah signifikan
- note: berikan penjelasan tentang temuan Anda

Dokumen yang akan dianalisis:`, number, statement, explanation, fillingGuide)

	return prompt
}

// parseLLMResponse parses the structured response from LLM
func (o *OpenAIService) parseLLMResponse(rawText string) (*service.DocumentCheckResponse, error) {
	jsonStr, err := ExtractJSONFromText(rawText)
	if err != nil {
		return o.parseFallbackResponse(rawText)
	}

	// Parse JSON
	var response struct {
		IsValid bool   `json:"isValid"`
		Note    string `json:"note"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
		// Fallback: try to extract boolean and note from text
		return o.parseFallbackResponse(rawText)
	}

	return &service.DocumentCheckResponse{
		IsValid: response.IsValid,
		Notes:   response.Note,
	}, nil
}

// parseFallbackResponse attempts to extract meaning when JSON parsing fails
func (o *OpenAIService) parseFallbackResponse(text string) (*service.DocumentCheckResponse, error) {
	lowerText := strings.ToLower(text)

	isValid := true
	if strings.Contains(lowerText, "tidak memenuhi") ||
		strings.Contains(lowerText, "tidak lengkap") ||
		strings.Contains(lowerText, "belum ada") ||
		strings.Contains(lowerText, "perlu perbaikan") ||
		strings.Contains(lowerText, "invalid") ||
		strings.Contains(lowerText, "false") {
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
