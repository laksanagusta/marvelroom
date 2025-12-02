package llm

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"sandbox/internal/domain/service"
)

// GeminiService implements LLMService interface using Google Gemini API
type GeminiService struct {
	apiKey     string
	httpClient *http.Client
	model      string
}

// geminiRequest represents the request structure for Gemini API
type geminiRequest struct {
	Contents []geminiContent `json:"contents"`
}

// geminiContent represents content structure for Gemini API
type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

// geminiPart represents a part of content for Gemini API
type geminiPart struct {
	Text       string      `json:"text,omitempty"`
	InlineData *geminiFile `json:"inline_data,omitempty"`
}

// geminiFile represents file data for Gemini API
type geminiFile struct {
	MimeType string `json:"mime_type"`
	Data     string `json:"data"`
}

// geminiResponse represents the response structure from Gemini API
type geminiResponse struct {
	Candidates []geminiCandidate `json:"candidates"`
}

// geminiCandidate represents a candidate response from Gemini API
type geminiCandidate struct {
	Content geminiContent `json:"content"`
}

// NewGeminiService creates a new Gemini service instance
func NewGeminiService(apiKey string) (service.LLMService, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("Gemini API key is required")
	}

	return &GeminiService{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 2000 * time.Second, // Longer timeout for document processing
		},
		model: "gemini-2.5-flash", // Using flash model for better multimodal capabilities
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

	// Build the prompt
	prompt := g.buildPrompt(req.Number, req.Statement, req.Explanation, req.FillingGuide)

	// Prepare parts for the request
	parts := []geminiPart{
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
		mimeType := g.getMimeType(doc.Type)
		data := g.encodeBase64(doc.Data)

		log.Printf("Adding document %d to Gemini request: MIME type=%s, encoded size=%d chars", i, mimeType, len(data))

		parts = append(parts, geminiPart{
			InlineData: &geminiFile{
				MimeType: mimeType,
				Data:     data,
			},
		})
	}

	log.Printf("Total parts being sent to Gemini API: %d (1 text + %d documents)", len(parts), len(parts)-1)

	// Build the request
	geminiReq := geminiRequest{
		Contents: []geminiContent{
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
	apiURL := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", g.model, g.apiKey)

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
	var geminiResp geminiResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return nil, fmt.Errorf("failed to parse Gemini API response: %w", err)
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("empty response from Gemini API")
	}

	// Extract the text response
	rawText := geminiResp.Candidates[0].Content.Parts[0].Text

	// Parse the structured response
	result, err := g.parseLLMResponse(rawText)
	if err != nil {
		return nil, fmt.Errorf("failed to parse LLM response: %w", err)
	}

	result.Model = g.model

	return result, nil
}

// buildPrompt constructs the prompt for LAKIP document checking
func (g *GeminiService) buildPrompt(number, statement, explanation, fillingGuide string) string {
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
func (g *GeminiService) parseLLMResponse(rawText string) (*service.DocumentCheckResponse, error) {
	// Clean the response text
	cleanText := strings.TrimSpace(rawText)

	// Extract JSON from response (handle markdown code blocks)
	if strings.Contains(cleanText, "```json") {
		start := strings.Index(cleanText, "```json") + 7
		end := strings.LastIndex(cleanText, "```")
		if end > start {
			cleanText = cleanText[start:end]
		}
	} else if strings.Contains(cleanText, "```") {
		start := strings.Index(cleanText, "```") + 3
		end := strings.LastIndex(cleanText, "```")
		if end > start {
			cleanText = cleanText[start:end]
		}
	}

	// Find JSON object boundaries
	jsonStart := strings.Index(cleanText, "{")
	jsonEnd := strings.LastIndex(cleanText, "}")

	if jsonStart == -1 || jsonEnd == -1 || jsonStart >= jsonEnd {
		return nil, fmt.Errorf("no valid JSON found in response: %s", cleanText)
	}

	jsonStr := cleanText[jsonStart : jsonEnd+1]

	// Parse JSON
	var response struct {
		IsValid bool   `json:"isValid"`
		Note    string `json:"note"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &response); err != nil {
		// Fallback: try to extract boolean and note from text
		return g.parseFallbackResponse(cleanText)
	}

	return &service.DocumentCheckResponse{
		IsValid: response.IsValid,
		Notes:   response.Note,
	}, nil
}

// parseFallbackResponse attempts to extract meaning when JSON parsing fails
func (g *GeminiService) parseFallbackResponse(text string) (*service.DocumentCheckResponse, error) {
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

// getMimeType maps file type to MIME type
func (g *GeminiService) getMimeType(fileType string) string {
	mimeTypeMap := map[string]string{
		"pdf":          "application/pdf",
		"document":     "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"spreadsheet":  "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"presentation": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
		"text":         "text/plain",
		"image":        "image/jpeg", // Default image type
		"jpeg":         "image/jpeg",
		"jpg":          "image/jpeg",
		"png":          "image/png",
	}

	if mimeType, exists := mimeTypeMap[fileType]; exists {
		return mimeType
	}

	return "application/octet-stream" // Default binary type
}

// encodeBase64 encodes data to base64
func (g *GeminiService) encodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}
