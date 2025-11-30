package entity

import (
	"time"
)

type Meeting struct {
	ID          string          `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	StartTime   time.Time       `json:"start_time"`
	Timezone    string          `json:"timezone"`
	Duration    int             `json:"duration_minutes"`
	HostUserID  string          `json:"host_user_id"`
	JoinURL     string          `json:"join_url"`
	Password    string          `json:"password"`
	Options     MeetingOptions  `json:"options"`
	Metadata    MeetingMetadata `json:"metadata"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type MeetingOptions struct {
	CreateDriveFolder     bool             `json:"create_drive_folder"`
	DriveParentFolderID   string           `json:"drive_parent_folder_id"`
	DuplicateAbsenceForm  bool             `json:"duplicate_absence_form"`
	AbsenceFormTemplateID string           `json:"absence_form_template_id"`
	Notify                NotificationOpts `json:"notify"`
	Zoom                  ZoomOpts         `json:"zoom"`
}

type NotificationOpts struct {
	SendEmail bool     `json:"send_email"`
	Channels  []string `json:"channels"`
	Message   string   `json:"message"`
}

type ZoomOpts struct {
	WaitingRoom     bool   `json:"waiting_room"`
	RequirePassword bool   `json:"require_password"`
	AutoRecording   string `json:"auto_recording"`
	MuteUponEntry   bool   `json:"mute_upon_entry"`
}

type MeetingMetadata struct {
	ProjectID string   `json:"project_id"`
	Tags      []string `json:"tags"`
}

type MeetingResult struct {
	Meeting          Meeting `json:"meeting"`
	DriveFolderURL   string  `json:"drive_folder_url,omitempty"`
	AbsenceFormURL   string  `json:"absence_form_url,omitempty"`
	NotificationSent bool    `json:"notification_sent"`
}
