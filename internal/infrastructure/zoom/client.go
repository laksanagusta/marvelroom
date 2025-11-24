package zoom

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
	apiSecret  string
	baseURL    string
}

type CreateMeetingRequest struct {
	Topic        string                 `json:"topic"`
	Type         int                    `json:"type"`
	StartTime    string                 `json:"start_time"`
	Duration     int                    `json:"duration"`
	Timezone     string                 `json:"timezone"`
	Password     string                 `json:"password,omitempty"`
	Settings     ZoomMeetingSettings    `json:"settings"`
}

type ZoomMeetingSettings struct {
	WaitingRoom   bool   `json:"waiting_room"`
	RequirePassword bool `json:"require_password"`
	AutoRecording string `json:"auto_recording"`
	MuteUponEntry bool   `json:"mute_upon_entry"`
	HostVideo     bool   `json:"host_video"`
	ParticipantVideo bool `json:"participant_video"`
}

type CreateMeetingResponse struct {
	UUID        string              `json:"uuid"`
	ID          int64               `json:"id"`
	HostID      string              `json:"host_id"`
	Topic       string              `json:"topic"`
	Type        int                 `json:"type"`
	Status      string              `json:"status"`
	StartTime   string              `json:"start_time"`
	Duration    int                 `json:"duration"`
	Timezone    string              `json:"timezone"`
	Password    string              `json:"password"`
	H323Password string             `json:"h323_password"`
	PstnPassword string             `json:"pstn_password"`
	Settings    ZoomMeetingSettings  `json:"settings"`
	JoinURL     string              `json:"join_url"`
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func NewClient(apiKey, apiSecret string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiKey:    apiKey,
		apiSecret: apiSecret,
		baseURL:   "https://api.zoom.us/v2",
	}
}

func (c *Client) getAccessToken(ctx context.Context) (string, error) {
	authURL := "https://zoom.us/oauth/token"

	req, err := http.NewRequestWithContext(ctx, "POST", authURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create auth request: %w", err)
	}

	req.SetBasicAuth(c.apiKey, c.apiSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	q := req.URL.Query()
	q.Add("grant_type", "account_credentials")
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get access token: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("auth failed with status %d: %s", resp.StatusCode, string(body))
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return "", fmt.Errorf("failed to decode auth response: %w", err)
	}

	return authResp.AccessToken, nil
}

func (c *Client) CreateZoomMeeting(ctx context.Context, meeting entity.Meeting) (*entity.Meeting, error) {
	accessToken, err := c.getAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	zoomReq := CreateMeetingRequest{
		Topic:     meeting.Title,
		Type:      2, // Scheduled meeting
		StartTime: meeting.StartTime.Format("2006-01-02T15:04:05Z"),
		Duration:  meeting.Duration,
		Timezone:  meeting.Timezone,
		Password:  meeting.Password,
		Settings: ZoomMeetingSettings{
			WaitingRoom:      meeting.Options.Zoom.WaitingRoom,
			AutoRecording:    meeting.Options.Zoom.AutoRecording,
			MuteUponEntry:    meeting.Options.Zoom.MuteUponEntry,
			HostVideo:        true,
			ParticipantVideo: false,
		},
	}

	reqBody, err := json.Marshal(zoomReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/users/%s/meetings", c.baseURL, meeting.HostUserID)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create meeting: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Zoom API returned status %d: %s", resp.StatusCode, string(body))
	}

	var zoomResp CreateMeetingResponse
	if err := json.NewDecoder(resp.Body).Decode(&zoomResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Update meeting with Zoom response
	result := meeting
	result.ID = fmt.Sprintf("%d", zoomResp.ID)
	result.JoinURL = zoomResp.JoinURL
	result.Password = zoomResp.Password
	result.CreatedAt = time.Now()
	result.UpdatedAt = time.Now()

	return &result, nil
}