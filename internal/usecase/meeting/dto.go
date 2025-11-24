package meeting

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/invopop/validation"
	"sandbox/internal/domain/entity"
)

var datetimeRegex = regexp.MustCompile(`^\d{1,2}\s+(Januari|Februari|Maret|April|Mei|Juni|Juli|Agustus|September|Oktober|November|Desember)\s+\d{4}\s+\d{2}:\d{2}:\d{2}$`)

// Indonesian month name to month number mapping
var IndonesianMonths = map[string]time.Month{
	"Januari":   time.January,
	"Februari":  time.February,
	"Maret":     time.March,
	"April":     time.April,
	"Mei":       time.May,
	"Juni":      time.June,
	"Juli":      time.July,
	"Agustus":   time.August,
	"September": time.September,
	"Oktober":   time.October,
	"November":  time.November,
	"Desember":  time.December,
}

type CreateMeetingRequest struct {
	Title           string                  `json:"title"`
	Description     string                  `json:"description"`
	StartTime       string                  `json:"start_time"`
	Timezone        string                  `json:"timezone"`
	DurationMinutes int                     `json:"duration_minutes"`
	HostUserID      string                  `json:"host_user_id"`
	Options         MeetingOptionsDTO       `json:"options"`
	Metadata        MeetingMetadataDTO      `json:"metadata"`
}

// parseIndonesianDateTime parses Indonesian datetime format (e.g., "25 Oktober 2025 10:30:00")
func parseIndonesianDateTime(datetimeStr string) (time.Time, error) {
	// Split by space to separate date and time
	parts := strings.Fields(datetimeStr)
	if len(parts) != 5 {
		return time.Time{}, fmt.Errorf("invalid datetime format")
	}

	// Date parts: day, month, year
	dayStr := parts[0]
	monthStr := parts[1]
	yearStr := parts[2]

	// Time part
	timeStr := parts[3] + ":" + parts[4] // HH:MM:SS

	// Parse day
	day, err := strconv.Atoi(dayStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid day: %s", dayStr)
	}

	// Parse month
	month, exists := IndonesianMonths[monthStr]
	if !exists {
		return time.Time{}, fmt.Errorf("invalid month: %s", monthStr)
	}

	// Parse year
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid year: %s", yearStr)
	}

	// Parse time
	hourMinSec := strings.Split(timeStr, ":")
	if len(hourMinSec) != 3 {
		return time.Time{}, fmt.Errorf("invalid time format: %s", timeStr)
	}

	hour, err := strconv.Atoi(hourMinSec[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid hour: %s", hourMinSec[0])
	}

	minute, err := strconv.Atoi(hourMinSec[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid minute: %s", hourMinSec[1])
	}

	second, err := strconv.Atoi(hourMinSec[2])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid second: %s", hourMinSec[2])
	}

	return time.Date(year, month, day, hour, minute, second, 0, time.UTC), nil
}

func (req *CreateMeetingRequest) Validate() error {
	return validation.ValidateStruct(req,
		validation.Field(&req.Title, validation.Required, validation.Length(1, 200)),
		validation.Field(&req.Description, validation.Length(0, 1000)),
		validation.Field(&req.StartTime, validation.Required, validation.Match(datetimeRegex)),
		validation.Field(&req.Timezone, validation.Required),
		validation.Field(&req.DurationMinutes, validation.Required, validation.Min(1), validation.Max(480)),
		validation.Field(&req.HostUserID, validation.Required),
		validation.Field(&req.Options),
		validation.Field(&req.Metadata),
	)
}

type MeetingOptionsDTO struct {
	CreateDriveFolder     bool                  `json:"create_drive_folder"`
	DriveParentFolderID   string                `json:"drive_parent_folder_id"`
	DuplicateAbsenceForm  bool                  `json:"duplicate_absence_form"`
	AbsenceFormTemplateID string                `json:"absence_form_template_id"`
	Notify                NotificationOptsDTO   `json:"notify"`
	Zoom                  ZoomOptsDTO           `json:"zoom"`
}

func (opts *MeetingOptionsDTO) Validate() error {
	if opts.DuplicateAbsenceForm && opts.AbsenceFormTemplateID == "" {
		return validation.NewError("absence_form_template_id", "absence_form_template_id is required when duplicate_absence_form is true")
	}

	return validation.ValidateStruct(opts,
		validation.Field(&opts.DriveParentFolderID),
		validation.Field(&opts.AbsenceFormTemplateID),
		validation.Field(&opts.Notify),
		validation.Field(&opts.Zoom),
	)
}

type NotificationOptsDTO struct {
	SendEmail bool     `json:"send_email"`
	Channels  []string `json:"channels"`
	Message   string   `json:"message"`
}

func (opts *NotificationOptsDTO) Validate() error {
	if opts.SendEmail {
		if len(opts.Channels) == 0 {
			return validation.NewError("channels", "channels are required when send_email is true")
		}
		if opts.Message == "" {
			return validation.NewError("message", "message is required when send_email is true")
		}
	}

	return validation.ValidateStruct(opts,
		validation.Field(&opts.Channels),
		validation.Field(&opts.Message, validation.Length(0, 1000)),
	)
}

type ZoomOptsDTO struct {
	WaitingRoom     bool   `json:"waiting_room"`
	RequirePassword bool   `json:"require_password"`
	AutoRecording   string `json:"auto_recording"`
}

func (opts *ZoomOptsDTO) Validate() error {
	if opts.AutoRecording != "" && opts.AutoRecording != "none" && opts.AutoRecording != "local" && opts.AutoRecording != "cloud" {
		return validation.NewError("auto_recording", "auto_recording must be one of: none, local, cloud")
	}

	return validation.ValidateStruct(opts)
}

type MeetingMetadataDTO struct {
	ProjectID string   `json:"project_id"`
	Tags      []string `json:"tags"`
}

func (meta *MeetingMetadataDTO) Validate() error {
	return validation.ValidateStruct(meta,
		validation.Field(&meta.ProjectID),
		validation.Field(&meta.Tags, validation.Each(validation.Length(1, 50))),
	)
}

type CreateMeetingResponse struct {
	Success           bool   `json:"success"`
	Message           string `json:"message"`
	Data              *MeetingResponseData `json:"data,omitempty"`
}

type MeetingResponseData struct {
	Meeting           entity.Meeting `json:"meeting"`
	DriveFolderURL    string          `json:"drive_folder_url,omitempty"`
	AbsenceFormURL    string          `json:"absence_form_url,omitempty"`
	NotificationSent  bool            `json:"notification_sent"`
}

func (req *CreateMeetingRequest) ToDomain() (*entity.Meeting, error) {
	startTime, err := parseIndonesianDateTime(req.StartTime)
	if err != nil {
		return nil, fmt.Errorf("invalid datetime format: %v (expected format: '25 Oktober 2025 10:30:00')", err)
	}

	return &entity.Meeting{
		Title:       req.Title,
		Description: req.Description,
		StartTime:   startTime,
		Timezone:    req.Timezone,
		Duration:    req.DurationMinutes,
		HostUserID:  req.HostUserID,
		Options: entity.MeetingOptions{
			CreateDriveFolder:     req.Options.CreateDriveFolder,
			DriveParentFolderID:   req.Options.DriveParentFolderID,
			DuplicateAbsenceForm:  req.Options.DuplicateAbsenceForm,
			AbsenceFormTemplateID: req.Options.AbsenceFormTemplateID,
			Notify: entity.NotificationOpts{
				SendEmail: req.Options.Notify.SendEmail,
				Channels:  req.Options.Notify.Channels,
				Message:   req.Options.Notify.Message,
			},
			Zoom: entity.ZoomOpts{
				WaitingRoom:     req.Options.Zoom.WaitingRoom,
				RequirePassword: req.Options.Zoom.RequirePassword,
				AutoRecording:   req.Options.Zoom.AutoRecording,
				MuteUponEntry:   true, // Default as mentioned
			},
		},
		Metadata: entity.MeetingMetadata{
			ProjectID: req.Metadata.ProjectID,
			Tags:      req.Metadata.Tags,
		},
	}, nil
}

