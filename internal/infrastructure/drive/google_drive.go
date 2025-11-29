package drive

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"sandbox/internal/domain/service"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// GoogleDriveService implements the DriveService interface for Google Drive
type GoogleDriveService struct {
	service    *drive.Service
	httpClient *http.Client
}

// NewGoogleDriveService creates a new Google Drive service instance using service account
func NewGoogleDriveService(credentialsFile string) (service.DriveService, error) {
	if credentialsFile == "" {
		credentialsFile = "dika-n8n-76c1c8c965e5.json"
	}

	ctx := context.Background()
	log.Printf("Creating Google Drive service with credentials file: %s", credentialsFile)
	srv, err := drive.NewService(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return nil, fmt.Errorf("Unable to create Drive client: %v", err)
	}

	log.Println("Successfully created Google Drive service")

	return &GoogleDriveService{
		service: srv,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// extractFolderID extracts the folder ID from a Google Drive URL or link
func (g *GoogleDriveService) extractFolderID(link string) (string, error) {
	// Handle different Google Drive URL formats
	patterns := []string{
		`drive\.google\.com/drive/folders/([a-zA-Z0-9_-]+)`,
		`drive\.google\.com/open\?id=([a-zA-Z0-9_-]+)`,
		`docs\.google\.com/spreadsheets/d/([a-zA-Z0-9_-]+)`,
		`docs\.google\.com/document/d/([a-zA-Z0-9_-]+)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(link)
		if len(matches) > 1 {
			return matches[1], nil
		}
	}

	// Check if the link itself is just an ID
	if strings.HasPrefix(link, "1") && len(link) > 20 {
		return link, nil
	}

	return "", fmt.Errorf("invalid Google Drive link format: %s", link)
}

// GetFilesFromFolder retrieves files from a Google Drive folder
func (g *GoogleDriveService) GetFilesFromFolder(ctx context.Context, folderLink string) ([]*service.DriveFile, error) {
	folderID, err := g.extractFolderID(folderLink)
	if err != nil {
		return nil, fmt.Errorf("failed to extract folder ID: %w", err)
	}

	log.Println("Extracted folder ID:", folderID)
	deadline, hasDeadline := ctx.Deadline()
	log.Printf("Context type: %T, Context deadline: %v (has: %v)", ctx, deadline, hasDeadline)

	// First, get folder details to ensure it exists
	folder, err := g.service.Files.Get(folderID).Fields("name", "id", "mimeType").Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to access folder: %w", err)
	}

	log.Printf("Found folder: %s (%s)", folder.Name, folder.Id)

	// Now get files in the folder
	query := fmt.Sprintf("'%s' in parents and trashed=false", folderID)
	files, err := g.service.Files.List().Q(query).
		Fields("files(id,name,mimeType,webContentLink,webViewLink)").
		Context(ctx).
		Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	var driveFiles []*service.DriveFile
	for _, file := range files.Files {
		// Filter to document types that are likely relevant for LAKIP
		if !g.isRelevantFileType(file.MimeType) {
			continue
		}

		driveFile := &service.DriveFile{
			ID:   file.Id,
			Name: file.Name,
			Type: g.getFileType(file.MimeType),
		}

		// Use webContentLink for downloading, fallback to webViewLink
		if file.WebContentLink != "" {
			driveFile.URL = file.WebContentLink
		} else if file.WebViewLink != "" {
			driveFile.URL = file.WebViewLink
		}

		driveFiles = append(driveFiles, driveFile)
	}

	log.Printf("Found %d relevant files in folder", len(driveFiles))
	return driveFiles, nil
}

// DownloadFile downloads file content from Google Drive
func (g *GoogleDriveService) DownloadFile(ctx context.Context, fileID string) ([]byte, error) {
	if fileID == "" {
		return nil, fmt.Errorf("file ID is required")
	}

	// Get file metadata first
	file, err := g.service.Files.Get(fileID).Fields("name", "mimeType").Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("file not found or access denied: %w", err)
	}

	log.Printf("Downloading file: %s (%s)", file.Name, file.Id)

	var resp *http.Response
	if strings.HasPrefix(file.MimeType, "application/vnd.google-apps") {
		// Google Apps files need to be exported
		exportMimeType := g.getExportMimeType(file.MimeType)
		resp, err = g.service.Files.Export(fileID, exportMimeType).Context(ctx).Download()
		if err != nil {
			return nil, fmt.Errorf("failed to export Google Apps file: %w", err)
		}
	} else {
		// Regular files can be downloaded directly
		resp, err = g.service.Files.Get(fileID).Context(ctx).Download()
		if err != nil {
			return nil, fmt.Errorf("failed to download file: %w", err)
		}
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	log.Printf("Successfully downloaded file: %s (%d bytes)", file.Name, len(content))
	return content, nil
}

// isRelevantFileType checks if the file type is relevant for LAKIP checking
func (g *GoogleDriveService) isRelevantFileType(mimeType string) bool {
	relevantTypes := []string{
		"application/pdf",
		"application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/vnd.ms-excel",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"application/vnd.ms-powerpoint",
		"application/vnd.openxmlformats-officedocument.presentationml.presentation",
		"text/plain",
		"image/jpeg",
		"image/png",
		"image/jpg",
		// Google Apps formats
		"application/vnd.google-apps.document",
		"application/vnd.google-apps.spreadsheet",
		"application/vnd.google-apps.presentation",
	}

	for _, relevantType := range relevantTypes {
		if strings.HasPrefix(mimeType, relevantType) {
			return true
		}
	}

	return false
}

// getFileType determines a simplified file type from MIME type
func (g *GoogleDriveService) getFileType(mimeType string) string {
	switch {
	case strings.Contains(mimeType, "pdf"):
		return "pdf"
	case strings.Contains(mimeType, "word") || strings.Contains(mimeType, "document"):
		return "document"
	case strings.Contains(mimeType, "excel") || strings.Contains(mimeType, "spreadsheet"):
		return "spreadsheet"
	case strings.Contains(mimeType, "powerpoint") || strings.Contains(mimeType, "presentation"):
		return "presentation"
	case strings.Contains(mimeType, "text"):
		return "text"
	case strings.Contains(mimeType, "image"):
		return "image"
	default:
		return "other"
	}
}

// getExportMimeType returns the appropriate export MIME type for Google Apps files
func (g *GoogleDriveService) getExportMimeType(googleAppsMimeType string) string {
	switch googleAppsMimeType {
	case "application/vnd.google-apps.document":
		return "application/pdf"
	case "application/vnd.google-apps.spreadsheet":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case "application/vnd.google-apps.presentation":
		return "application/pdf"
	default:
		return "application/pdf"
	}
}
