package drive

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
}

type CreateFolderRequest struct {
	Name     string            `json:"name"`
	Parents  []string          `json:"parents,omitempty"`
	MimeType string            `json:"mimeType"`
	Properties map[string]string `json:"properties,omitempty"`
}

type FileResponse struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	MimeType    string            `json:"mimeType"`
	Parents     []string          `json:"parents"`
	CreatedTime string            `json:"createdTime"`
	WebViewLink string            `json:"webViewLink"`
}

func NewClient(apiKey string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiKey:  apiKey,
		baseURL: "https://www.googleapis.com/drive/v3",
	}
}

func (c *Client) CreateFolder(ctx context.Context, parentFolderID, folderName string) (string, error) {
	req := CreateFolderRequest{
		Name:     folderName,
		MimeType: "application/vnd.google-apps.folder",
	}

	if parentFolderID != "" {
		req.Parents = []string{parentFolderID}
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/files", c.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to create folder: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Drive API returned status %d: %s", resp.StatusCode, string(body))
	}

	var fileResp FileResponse
	if err := json.NewDecoder(resp.Body).Decode(&fileResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return fileResp.WebViewLink, nil
}

func (c *Client) DuplicateFile(ctx context.Context, templateID, targetFolderID, newFileName string) (string, error) {
	req := map[string]interface{}{
		"name": newFileName,
	}

	if targetFolderID != "" {
		req["parents"] = []string{targetFolderID}
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/files/%s/copy", c.baseURL, templateID)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("failed to duplicate file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Drive API returned status %d: %s", resp.StatusCode, string(body))
	}

	var fileResp FileResponse
	if err := json.NewDecoder(resp.Body).Decode(&fileResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return fileResp.WebViewLink, nil
}

