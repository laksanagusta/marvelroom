package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"sandbox/internal/domain/entity"
)

type Client struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
}

type EmailRequest struct {
	To       []string `json:"to"`
	Subject  string   `json:"subject"`
	Body     string   `json:"body"`
	HTML     bool     `json:"html"`
}

type EmailResponse struct {
	MessageID string `json:"message_id"`
	Status    string `json:"status"`
}

func NewClient(apiKey string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiKey:  apiKey,
		baseURL: "https://api.notification-service.com/v1", // Example service
	}
}

func (c *Client) SendEmail(ctx context.Context, to []string, subject, body string) error {
	req := EmailRequest{
		To:      to,
		Subject: subject,
		Body:    body,
		HTML:    true,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/send", c.baseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Notification API returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func (c *Client) SendNotification(ctx context.Context, opts entity.NotificationOpts, meetingURL string) error {
	if !opts.SendEmail {
		return nil
	}

	subject := "Meeting Created: New Meeting Scheduled"

	body := fmt.Sprintf(`
		<h2>Meeting Details</h2>
		<p>%s</p>
		<p><strong>Meeting Link:</strong> <a href="%s">Join Meeting</a></p>
	`, opts.Message, meetingURL)

	// In a real implementation, you would get the recipient email from the host_user_id
	// For now, we'll use a placeholder
	recipients := []string{"host@example.com"} // This should be resolved from host_user_id

	err := c.SendEmail(ctx, recipients, subject, body)
	if err != nil {
		return fmt.Errorf("failed to send notification: %w", err)
	}

	return nil
}


