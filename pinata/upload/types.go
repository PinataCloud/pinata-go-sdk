package upload

import (
	"io"
	"os"
)

// FileOptions represents options for file uploads
type FileOptions struct {
	FileName  string
	GroupID   string
	KeyValues map[string]string
	Vectorize bool
}

// Base64Options represents options for base64 uploads
type Base64Options struct {
	Name      string
	GroupID   string
	KeyValues map[string]string
	Vectorize bool
}

// JSONOptions represents options for JSON uploads
type JSONOptions struct {
	Name      string
	GroupID   string
	KeyValues map[string]string
	Vectorize bool
}

// URLOptions represents options for URL uploads
type URLOptions struct {
	Name      string
	GroupID   string
	KeyValues map[string]string
	Vectorize bool
}

// CIDOptions represents options for pinning an existing CID
type CIDOptions struct {
	CID       string
	Name      string
	GroupID   string
	KeyValues map[string]string
	HostNodes []string
}

// SignedUploadOptions represents options for creating a signed upload URL
type SignedUploadOptions struct {
	Date        int64
	Expires     int
	GroupID     string
	Name        string
	KeyValues   map[string]string
	Vectorize   bool
	MaxFileSize int64
	MimeTypes   []string
}

// FileData wraps either an os.File or io.Reader with additional metadata
type FileData struct {
	Reader      io.Reader
	Name        string
	Size        int64
	ContentType string
}

// NewFileData creates a FileData from an os.File
func NewFileData(file *os.File) (*FileData, error) {
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}

	return &FileData{
		Reader: file,
		Name:   info.Name(),
		Size:   info.Size(),
	}, nil
}

// NewCustomFileData creates a FileData with custom metadata
func NewCustomFileData(reader io.Reader, name string, size int64, contentType string) *FileData {
	return &FileData{
		Reader:      reader,
		Name:        name,
		Size:        size,
		ContentType: contentType,
	}
}
