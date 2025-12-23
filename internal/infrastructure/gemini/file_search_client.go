package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"time"
)

const (
	fileSearchBaseURL = "https://generativelanguage.googleapis.com/v1beta"
	fileUploadBaseURL = "https://generativelanguage.googleapis.com/upload/v1beta"
	defaultModel      = "gemini-2.5-flash"
)

// FileSearchClient handles Gemini File Search API operations
type FileSearchClient struct {
	apiKey     string
	httpClient *http.Client
}

// NewFileSearchClient creates a new FileSearchClient
func NewFileSearchClient(apiKey string) *FileSearchClient {
	return &FileSearchClient{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// FileSearchStore represents a Gemini File Search store
type FileSearchStore struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	CreateTime  string `json:"createTime"`
}

// FileSearchDocument represents a document in a store
type FileSearchDocument struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	State       string `json:"state"` // PROCESSING, ACTIVE, FAILED
	CreateTime  string `json:"createTime"`
}

// Citation represents a source citation from RAG
type Citation struct {
	DocumentName string `json:"document_name"`
	Content      string `json:"content"`
	StartIndex   int    `json:"start_index,omitempty"`
	EndIndex     int    `json:"end_index,omitempty"`
}

// ChatResponse represents a response from the chat API with RAG
type ChatResponse struct {
	Text      string     `json:"text"`
	Citations []Citation `json:"citations,omitempty"`
}

// CreateFileSearchStore creates a new file search store
func (c *FileSearchClient) CreateFileSearchStore(ctx context.Context, displayName string) (*FileSearchStore, error) {
	url := fmt.Sprintf("%s/fileSearchStores?key=%s", fileSearchBaseURL, c.apiKey)

	payload := map[string]string{
		"displayName": displayName,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var store FileSearchStore
	if err := json.Unmarshal(body, &store); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &store, nil
}

// UploadToFileSearchStore uploads a file to the store
func (c *FileSearchClient) UploadToFileSearchStore(ctx context.Context, storeName string, fileName string, content []byte, mimeType string) (*FileSearchDocument, error) {
	url := fmt.Sprintf("%s/%s:uploadToFileSearchStore?key=%s", fileUploadBaseURL, storeName, c.apiKey)

	// Log upload attempt
	fmt.Printf("[Gemini] Upload URL: %s/%s:uploadToFileSearchStore\n", fileUploadBaseURL, storeName)
	fmt.Printf("[Gemini] Uploading file: %s (size: %d bytes, mimeType: %s)\n", fileName, len(content), mimeType)

	// Create multipart form with custom MIME header
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Create form file with custom headers (matches what curl --form sends)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, fileName))
	if mimeType != "" {
		h.Set("Content-Type", mimeType)
	} else {
		h.Set("Content-Type", "application/octet-stream")
	}

	part, err := writer.CreatePart(h)
	if err != nil {
		return nil, fmt.Errorf("failed to create form part: %w", err)
	}
	if _, err := part.Write(content); err != nil {
		return nil, fmt.Errorf("failed to write file content: %w", err)
	}
	writer.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Log response
	fmt.Printf("[Gemini] Upload Response (status %d): %s\n", resp.StatusCode, string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Response contains operation details, parse to get document info
	var result struct {
		Name     string `json:"name"`
		Done     bool   `json:"done"`
		Response struct {
			Document FileSearchDocument `json:"document"`
		} `json:"response"`
		Error struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Check for error in response
	if result.Error.Message != "" {
		fmt.Printf("[Gemini] Upload Error: %s\n", result.Error.Message)
		return nil, fmt.Errorf("Gemini error: %s", result.Error.Message)
	}

	fmt.Printf("[Gemini] Upload result - Name: %s, Done: %v\n", result.Name, result.Done)

	if result.Done {
		fmt.Printf("[Gemini] Document created: %s (State: %s)\n", result.Response.Document.Name, result.Response.Document.State)
		return &result.Response.Document, nil
	}

	// If not done, return operation name and let caller poll
	return &FileSearchDocument{
		Name:  result.Name,
		State: "PROCESSING",
	}, nil
}

// QueryWithFileSearch performs a RAG query using the file search store
func (c *FileSearchClient) QueryWithFileSearch(ctx context.Context, storeName string, query string, chatHistory []ChatMessage) (*ChatResponse, error) {
	url := fmt.Sprintf("%s/models/%s:generateContent?key=%s", fileSearchBaseURL, defaultModel, c.apiKey)

	// Build contents array with chat history
	contents := make([]map[string]interface{}, 0)
	for _, msg := range chatHistory {
		// Convert role: Gemini expects "user" or "model", not "assistant"
		role := msg.Role
		if role == "assistant" {
			role = "model"
		}
		contents = append(contents, map[string]interface{}{
			"role": role,
			"parts": []map[string]string{
				{"text": msg.Content},
			},
		})
	}

	// Add current query
	contents = append(contents, map[string]interface{}{
		"role": "user",
		"parts": []map[string]string{
			{"text": query},
		},
	})

	payload := map[string]interface{}{
		"contents": contents,
		"tools": []map[string]interface{}{
			{
				"fileSearch": map[string]interface{}{
					"fileSearchStoreNames": []string{storeName},
				},
			},
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse response
	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
			GroundingMetadata struct {
				GroundingChunks []struct {
					RetrievedContext struct {
						Uri   string `json:"uri"`
						Title string `json:"title"`
						Text  string `json:"text"`
					} `json:"retrievedContext"`
				} `json:"groundingChunks"`
			} `json:"groundingMetadata"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	chatResp := &ChatResponse{}

	if len(result.Candidates) > 0 && len(result.Candidates[0].Content.Parts) > 0 {
		chatResp.Text = result.Candidates[0].Content.Parts[0].Text

		// Extract citations from grounding metadata
		for _, chunk := range result.Candidates[0].GroundingMetadata.GroundingChunks {
			chatResp.Citations = append(chatResp.Citations, Citation{
				DocumentName: chunk.RetrievedContext.Title,
				Content:      chunk.RetrievedContext.Text,
			})
		}
	}

	return chatResp, nil
}

// ChatMessage for building conversation history
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// DeleteFileSearchStore deletes a file search store
func (c *FileSearchClient) DeleteFileSearchStore(ctx context.Context, storeName string) error {
	url := fmt.Sprintf("%s/%s?force=true&key=%s", fileSearchBaseURL, storeName, c.apiKey)

	fmt.Printf("[Gemini] Deleting FileSearchStore: %s\n", storeName)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("[Gemini] Delete store response (status %d): %s\n", resp.StatusCode, string(body))

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	fmt.Printf("[Gemini] FileSearchStore deleted successfully: %s\n", storeName)
	return nil
}

// DeleteDocument deletes a document from a file search store
func (c *FileSearchClient) DeleteDocument(ctx context.Context, storeName, documentName string) error {
	// documentName from Gemini is already a full path like "fileSearchStores/xxx/documents/yyy"
	// So we use it directly. Add force=true to force delete even if document has chunks
	url := fmt.Sprintf("%s/%s?force=true&key=%s", fileSearchBaseURL, documentName, c.apiKey)

	fmt.Printf("[Gemini] Deleting document: %s\n", documentName)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("[Gemini] Delete document response (status %d): %s\n", resp.StatusCode, string(body))

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	fmt.Printf("[Gemini] Document deleted successfully: %s\n", documentName)
	return nil
}

// ListDocuments lists all documents in a file search store
func (c *FileSearchClient) ListDocuments(ctx context.Context, storeName string) ([]FileSearchDocument, error) {
	url := fmt.Sprintf("%s/%s/documents?key=%s", fileSearchBaseURL, storeName, c.apiKey)

	// Log the request URL (without API key for security)
	fmt.Printf("[Gemini] ListDocuments URL: %s/%s/documents\n", fileSearchBaseURL, storeName)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Log the raw response
	fmt.Printf("[Gemini] ListDocuments Response (status %d): %s\n", resp.StatusCode, string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Documents []FileSearchDocument `json:"documents"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Log parsed documents
	for i, doc := range result.Documents {
		fmt.Printf("[Gemini] Document %d: Name=%s, DisplayName=%s, State=%s\n", i, doc.Name, doc.DisplayName, doc.State)
	}

	return result.Documents, nil
}

// GetOperationStatus checks the status of an async operation
func (c *FileSearchClient) GetOperationStatus(ctx context.Context, operationName string) (bool, error) {
	url := fmt.Sprintf("%s/%s?key=%s", fileSearchBaseURL, operationName, c.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Done bool `json:"done"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return false, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return result.Done, nil
}
