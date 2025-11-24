package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/repository"
)

type MeetingService struct {
	repo repository.MeetingRepository
}

func NewMeetingService(repo repository.MeetingRepository) *MeetingService {
	return &MeetingService{
		repo: repo,
	}
}

func (s *MeetingService) CreateMeeting(ctx context.Context, req entity.Meeting) (*entity.MeetingResult, error) {
	// Generate password if required
	if req.Options.Zoom.RequirePassword && req.Password == "" {
		req.Password = generatePassword()
	}

	// Create Zoom meeting
	meeting, err := s.repo.CreateZoomMeeting(ctx, &req)
	if err != nil {
		return nil, fmt.Errorf("failed to create Zoom meeting: %w", err)
	}

	result := &entity.MeetingResult{
		Meeting: *meeting,
	}

	// Create Drive folder if requested
	if req.Options.CreateDriveFolder {
		folderName := fmt.Sprintf("%s - %s", req.Title, time.Now().Format("2006-01-02"))
		driveURL, err := s.repo.CreateDriveFolder(ctx, req.Options.DriveParentFolderID, folderName)
		if err != nil {
			return nil, fmt.Errorf("failed to create Drive folder: %w", err)
		}
		result.DriveFolderURL = driveURL
	}

	// Duplicate absence form if requested
	if req.Options.DuplicateAbsenceForm && req.Options.AbsenceFormTemplateID != "" {
		// Extract folder ID from Drive URL if available
		folderID := ""
		if result.DriveFolderURL != "" {
			folderID = extractFolderIDFromURL(result.DriveFolderURL)
		}

		absenceFormURL, err := s.repo.DuplicateAbsenceForm(ctx, req.Options.AbsenceFormTemplateID, folderID)
		if err != nil {
			return nil, fmt.Errorf("failed to duplicate absence form: %w", err)
		}
		result.AbsenceFormURL = absenceFormURL
	}

	// Send notification if requested
	if req.Options.Notify.SendEmail {
		meetingURL := meeting.JoinURL
		if result.DriveFolderURL != "" {
			meetingURL += fmt.Sprintf("\nDrive Folder: %s", result.DriveFolderURL)
		}
		if result.AbsenceFormURL != "" {
			meetingURL += fmt.Sprintf("\nAbsence Form: %s", result.AbsenceFormURL)
		}

		err := s.repo.SendNotification(ctx, req.Options.Notify, meetingURL)
		if err != nil {
			return nil, fmt.Errorf("failed to send notification: %w", err)
		}
		result.NotificationSent = true
	}

	return result, nil
}

func generatePassword() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[i%len(charset)]
	}
	return string(b)
}

func extractFolderIDFromURL(url string) string {
	// Simple extraction - in real implementation, parse Google Drive URL
	parts := strings.Split(url, "/")
	for i, part := range parts {
		if part == "folders" && i+1 < len(parts) {
			return parts[i+1]
		}
	}
	return ""
}

