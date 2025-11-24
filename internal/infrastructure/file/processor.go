package file

import (
	"errors"
	"io"
	"mime/multipart"
	"strings"
)

// Processor handles file processing operations
type Processor struct {
	allowedMimeTypes map[string]bool
	maxFileSize      int64
}

// NewProcessor creates a new file processor with default settings
func NewProcessor() *Processor {
	return &Processor{
		allowedMimeTypes: map[string]bool{
			"image/png":       true,
			"image/jpeg":      true,
			"application/pdf": true,
		},
		maxFileSize: 10 * 1024 * 1024, // 10MB
	}
}

// ProcessedFile represents a processed file
type ProcessedFile struct {
	Content  []byte
	Filename string
	MimeType string
}

func (p *Processor) ProcessUploadedFile(fileHeader *multipart.FileHeader) (*ProcessedFile, error) {
	if fileHeader == nil {
		return nil, errors.New("file header is nil")
	}

	if fileHeader.Size > p.maxFileSize {
		return nil, errors.New("file size exceeds maximum allowed size")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	mimeType := p.detectMimeType(fileHeader.Filename)

	if !p.allowedMimeTypes[mimeType] {
		return nil, errors.New("file type not allowed")
	}

	return &ProcessedFile{
		Content:  content,
		Filename: fileHeader.Filename,
		MimeType: mimeType,
	}, nil
}

func (p *Processor) ProcessMultipleFiles(fileHeaders []*multipart.FileHeader) ([]*ProcessedFile, error) {
	if len(fileHeaders) == 0 {
		return nil, errors.New("no files provided")
	}

	processedFiles := make([]*ProcessedFile, 0, len(fileHeaders))

	for _, fileHeader := range fileHeaders {
		processed, err := p.ProcessUploadedFile(fileHeader)
		if err != nil {
			return nil, err
		}
		processedFiles = append(processedFiles, processed)
	}

	return processedFiles, nil
}

func (p *Processor) detectMimeType(filename string) string {
	lowerFilename := strings.ToLower(filename)

	if strings.HasSuffix(lowerFilename, ".png") {
		return "image/png"
	} else if strings.HasSuffix(lowerFilename, ".jpg") || strings.HasSuffix(lowerFilename, ".jpeg") {
		return "image/jpeg"
	} else if strings.HasSuffix(lowerFilename, ".pdf") {
		return "application/pdf"
	}

	return "application/octet-stream"
}
