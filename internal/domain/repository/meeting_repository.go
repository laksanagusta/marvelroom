package repository

import (
	"context"
	"sandbox/internal/domain/entity"
)

type MeetingRepository interface {
	CreateZoomMeeting(ctx context.Context, meeting *entity.Meeting) (*entity.Meeting, error)
	CreateDriveFolder(ctx context.Context, parentFolderID, folderName string) (string, error)
	DuplicateAbsenceForm(ctx context.Context, templateID, folderID string) (string, error)
	SendNotification(ctx context.Context, opts interface{}, meetingURL string) error
}