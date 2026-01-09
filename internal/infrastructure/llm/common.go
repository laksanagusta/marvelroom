package llm

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// API URLs
const (
	GeminiAPIBaseURL = "https://generativelanguage.googleapis.com/v1beta/models"
	OpenAIAPIURL     = "https://api.openai.com/v1/chat/completions"
)

// GeminiModels is the list of Gemini models to try in order of preference
var GeminiModels = []string{
	"gemini-2.5-flash",
	"gemini-2.5-flash-lite",
}

// DefaultHTTPClient creates a new HTTP client with default timeout
func DefaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 300 * time.Second,
	}
}

// LongTimeoutHTTPClient creates an HTTP client with longer timeout for document processing
func LongTimeoutHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 2000 * time.Second,
	}
}

// GetGeminiModelURL returns the API URL for a specific Gemini model
func GetGeminiModelURL(model, apiKey string) string {
	return fmt.Sprintf("%s/%s:generateContent?key=%s", GeminiAPIBaseURL, model, apiKey)
}

// EncodeBase64 encodes data to base64 string
func EncodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// MimeTypeMap contains mappings from file type to MIME type
var MimeTypeMap = map[string]string{
	"pdf":          "application/pdf",
	"document":     "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	"spreadsheet":  "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
	"presentation": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
	"text":         "text/plain",
	"image":        "image/jpeg",
	"jpeg":         "image/jpeg",
	"jpg":          "image/jpeg",
	"png":          "image/png",
}

// GetMimeType returns the MIME type for a given file type
func GetMimeType(fileType string) string {
	if mimeType, exists := MimeTypeMap[fileType]; exists {
		return mimeType
	}
	return "application/octet-stream"
}

// GeminiAPIRequest represents the common request structure for Gemini API
type GeminiAPIRequest struct {
	Contents []GeminiContent `json:"contents"`
}

// GeminiContent represents content structure for Gemini API
type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

// GeminiPart represents a part of content for Gemini API
type GeminiPart struct {
	Text       string      `json:"text,omitempty"`
	InlineData *GeminiFile `json:"inline_data,omitempty"`
}

// GeminiFile represents file data for Gemini API
type GeminiFile struct {
	MimeType string `json:"mime_type"`
	Data     string `json:"data"`
}

// GeminiAPIResponse represents the common response structure from Gemini API
type GeminiAPIResponse struct {
	Candidates []GeminiCandidate `json:"candidates"`
}

// GeminiCandidate represents a candidate response from Gemini API
type GeminiCandidate struct {
	Content GeminiContent `json:"content"`
}

// OpenAIRequest represents the request structure for OpenAI API
type OpenAIRequest struct {
	Model     string          `json:"model"`
	Messages  []OpenAIMessage `json:"messages"`
	MaxTokens int             `json:"max_tokens,omitempty"`
}

// OpenAIMessage represents a message in OpenAI API
type OpenAIMessage struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"` // Can be string or []OpenAIContent
}

// OpenAIContent represents content in OpenAI message for vision API
type OpenAIContent struct {
	Type     string          `json:"type"`
	Text     string          `json:"text,omitempty"`
	ImageURL *OpenAIImageURL `json:"image_url,omitempty"`
}

// OpenAIImageURL represents image URL in OpenAI content
type OpenAIImageURL struct {
	URL string `json:"url"`
}

// OpenAIResponse represents the response structure from OpenAI API
type OpenAIResponse struct {
	Choices []OpenAIChoice `json:"choices"`
	Error   *OpenAIError   `json:"error,omitempty"`
}

// OpenAIChoice represents a choice in OpenAI response
type OpenAIChoice struct {
	Message OpenAIMessageResponse `json:"message"`
}

// OpenAIMessageResponse represents message in OpenAI response
type OpenAIMessageResponse struct {
	Content string `json:"content"`
}

// OpenAIError represents error in OpenAI response
type OpenAIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

// CleanJSON removes markdown code blocks and whitespace from JSON string
func CleanJSON(s string) string {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "```json")
	s = strings.TrimPrefix(s, "```JSON")
	s = strings.TrimSuffix(s, "```")
	return strings.TrimSpace(s)
}

// ExtractJSONFromText extracts JSON object from text
func ExtractJSONFromText(text string) (string, error) {
	cleanText := CleanJSON(text)

	jsonStart := strings.Index(cleanText, "{")
	jsonEnd := strings.LastIndex(cleanText, "}")

	if jsonStart == -1 || jsonEnd == -1 || jsonStart >= jsonEnd {
		return "", errors.New("no valid JSON found in response")
	}

	return cleanText[jsonStart : jsonEnd+1], nil
}

// CallGeminiAPI makes an API call to Gemini with the given request
func CallGeminiAPI(httpClient *http.Client, apiURL string, request *GeminiAPIRequest) (*GeminiAPIResponse, error) {
	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call Gemini API: %w", err)
	}
	defer resp.Body.Close()

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Gemini API error (status %d): %s", resp.StatusCode, string(bodyResp))
	}

	var geminiResp GeminiAPIResponse
	if err := json.Unmarshal(bodyResp, &geminiResp); err != nil {
		return nil, fmt.Errorf("failed to parse Gemini API response: %w", err)
	}

	return &geminiResp, nil
}

// CallOpenAIAPI makes an API call to OpenAI with the given request
func CallOpenAIAPI(httpClient *http.Client, apiKey string, request *OpenAIRequest) (*OpenAIResponse, error) {
	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", OpenAIAPIURL, strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call OpenAI API: %w", err)
	}
	defer resp.Body.Close()

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAI API error (status %d): %s", resp.StatusCode, string(bodyResp))
	}

	var openAIResp OpenAIResponse
	if err := json.Unmarshal(bodyResp, &openAIResp); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI API response: %w", err)
	}

	if openAIResp.Error != nil {
		return nil, fmt.Errorf("OpenAI API error: %s", openAIResp.Error.Message)
	}

	return &openAIResp, nil
}

// GetTextFromGeminiResponse extracts text from Gemini API response
func GetTextFromGeminiResponse(resp *GeminiAPIResponse) (string, error) {
	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", errors.New("empty response from Gemini API")
	}
	return resp.Candidates[0].Content.Parts[0].Text, nil
}

// GetTextFromOpenAIResponse extracts text from OpenAI API response
func GetTextFromOpenAIResponse(resp *OpenAIResponse) (string, error) {
	if len(resp.Choices) == 0 {
		return "", errors.New("empty response from OpenAI API")
	}
	return resp.Choices[0].Message.Content, nil
}

// BuildDocumentCheckPrompt constructs the prompt for LAKIP document checking
// This is a shared function used by both Gemini and OpenAI services
func BuildDocumentCheckPrompt(number, classification, topicDescription, deskInstruction string) string {
	// Build context section if topic description is available
	contextSection := ""
	if topicDescription != "" {
		contextSection = fmt.Sprintf(`
Konteks Topic:
%s
`, topicDescription)
	}

	prompt := fmt.Sprintf(`Periksa dokumen yang diberikan sesuai dengan poin kertas kerja LAKIP berikut:

Nomor: %s
Klasifikasi: %s
%s
Instruksi Desk:
%s

Tugas Anda:
1. Analisis semua dokumen yang disediakan
2. Periksa apakah isi dokumen telah memenuhi persyaratan yang disebutkan dalam instruksi desk
3. Berikan penilaian objektif dan berbasis bukti tentang kelengkapan dan kepatuhan dokumen
4. Jika dokumen tidak memenuhi persyaratan, berikan rekomendasi untuk perbaikan
5. Jika menurutmu sudah memenuhi persyaratan, berikan bukti pernyataan mana yang membuat poin ini telah memenuhi persyaratan

Jawab dengan format JSON berikut:
{
  "isValid": true/false,
  "note": "Penjelasan tentang temuan, rekomendasi, atau alasan penilaian (maksimal 5-6 kalimat)"
}

Kriteria penilaian:
- isValid: true jika dokumen lengkap dan memenuhi persyaratan
- isValid: false jika dokumen tidak lengkap, tidak memenuhi persyaratan, atau ada masalah signifikan
- note: berikan penjelasan tentang temuan Anda

Dokumen yang akan dianalisis:`, number, classification, contextSection, deskInstruction)

	return prompt
}
