package meeting

import (
	"context"
	"fmt"

	"sandbox/internal/domain/service"
)

type CreateMeetingUseCase struct {
	meetingService *service.MeetingService
}

func NewCreateMeetingUseCase(meetingService *service.MeetingService) *CreateMeetingUseCase {
	return &CreateMeetingUseCase{
		meetingService: meetingService,
	}
}

func (uc *CreateMeetingUseCase) Execute(ctx context.Context, req CreateMeetingRequest) (*CreateMeetingResponse, error) {
	meetingEntity, err := req.ToDomain()
	if err != nil {
		return &CreateMeetingResponse{
			Success: false,
			Message: fmt.Sprintf("Invalid request format: %v", err),
		}, nil
	}

	result, err := uc.meetingService.CreateMeeting(ctx, *meetingEntity)
	if err != nil {
		return &CreateMeetingResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to create meeting: %v", err),
		}, nil
	}

	responseData := &MeetingResponseData{
		Meeting:          result.Meeting,
		DriveFolderURL:   result.DriveFolderURL,
		AbsenceFormURL:   result.AbsenceFormURL,
		NotificationSent: result.NotificationSent,
	}

	return &CreateMeetingResponse{
		Success: true,
		Message: "Meeting created successfully",
		Data:    responseData,
	}, nil
}
