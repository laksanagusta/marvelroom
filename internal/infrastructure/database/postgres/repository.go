package meeting

import (
	"context"
	"fmt"
	"time"

	"sandbox/internal/domain/entity"
	"sandbox/internal/domain/repository"
	"sandbox/internal/infrastructure/drive"
	"sandbox/internal/infrastructure/notification"
	"sandbox/internal/infrastructure/zoom"
)

type Repository struct {
	zoomClient         *zoom.Client
	driveClient        *drive.Client
	notificationClient *notification.Client
}

func NewRepository(
	zoomClient *zoom.Client,
	driveClient *drive.Client,
	notificationClient *notification.Client,
) repository.MeetingRepository {
	return &Repository{
		zoomClient:         zoomClient,
		driveClient:        driveClient,
		notificationClient: notificationClient,
	}
}

func (r *Repository) CreateZoomMeeting(ctx context.Context, meeting *entity.Meeting) (*entity.Meeting, error) {
	return r.zoomClient.CreateZoomMeeting(ctx, *meeting)
}

func (r *Repository) CreateDriveFolder(ctx context.Context, parentFolderID, folderName string) (string, error) {
	return r.driveClient.CreateFolder(ctx, parentFolderID, folderName)
}

func (r *Repository) DuplicateAbsenceForm(ctx context.Context, templateID, folderID string) (string, error) {
	if templateID == "" {
		return "", fmt.Errorf("template ID is required")
	}

	newFileName := fmt.Sprintf("Absence Form - %s", time.Now().Format("2006-01-02-15-04"))

	return r.driveClient.DuplicateFile(ctx, templateID, folderID, newFileName)
}

func (r *Repository) SendNotification(ctx context.Context, opts interface{}, meetingURL string) error {
	notifyOpts, ok := opts.(entity.NotificationOpts)
	if !ok {
		return fmt.Errorf("invalid notification options type")
	}
	return r.notificationClient.SendNotification(ctx, notifyOpts, meetingURL)
}
