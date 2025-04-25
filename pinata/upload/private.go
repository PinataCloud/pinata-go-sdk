package upload

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/PinataCloud/pinata-go-sdk/pinata/types"
)

// PublicService provides upload operations for the public IPFS network
type PrivateService struct {
	config interface{}
}

// NewPublicService creates a new PublicService with the provided configuration
func NewPrivateService(config interface{}) *PrivateService {
	return &PrivateService{
		config: config,
	}
}

// File uploads a file to the public IPFS network
func (s *PrivateService) File(file *os.File, opts *FileOptions) (*types.UploadResponse, error) {
	if file == nil {
		return nil, fmt.Errorf("file is required")
	}

	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}

	// Reset file position to start
	if _, err := file.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to reset file position: %w", err)
	}

	cfg := s.config.(*types.Config)
	url := fmt.Sprintf("%s/files", cfg.UploadUrl)

	// Create multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the network parameter
	if err := writer.WriteField("network", "private"); err != nil {
		return nil, fmt.Errorf("failed to add network field: %w", err)
	}

	// Add optional fields if provided
	if opts != nil {
		if opts.GroupID != "" {
			if err := writer.WriteField("group_id", opts.GroupID); err != nil {
				return nil, fmt.Errorf("failed to add group_id field: %w", err)
			}
		}

		// Use custom name or fallback to file name
		name := fileInfo.Name()
		if opts.FileName != "" {
			name = opts.FileName
		}

		if err := writer.WriteField("name", name); err != nil {
			return nil, fmt.Errorf("failed to add name field: %w", err)
		}

		// Add keyvalues if present
		if len(opts.KeyValues) > 0 {
			keyvaluesJSON, err := json.Marshal(opts.KeyValues)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal keyvalues: %w", err)
			}

			if err := writer.WriteField("keyvalues", string(keyvaluesJSON)); err != nil {
				return nil, fmt.Errorf("failed to add keyvalues field: %w", err)
			}
		}
	}

	// Add the file
	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("failed to copy file data: %w", err)
	}

	// Close the writer
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Create the request
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+cfg.PinataJWT)

	// Add custom headers if any
	for key, value := range cfg.CustomHeaders {
		req.Header.Set(key, value)
	}

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var response struct {
		Data *types.UploadResponse `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Data, nil
}

// FileArray uploads multiple files as a folder to the public IPFS network
func (s *PrivateService) FileArray(files []*os.File, opts *FileOptions) (*types.UploadResponse, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("at least one file is required")
	}

	cfg := s.config.(*types.Config)
	url := fmt.Sprintf("%s/files", cfg.UploadUrl)

	// Create multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add the network parameter
	if err := writer.WriteField("network", "private"); err != nil {
		return nil, fmt.Errorf("failed to add network field: %w", err)
	}

	// Add optional fields if provided
	if opts != nil {
		if opts.GroupID != "" {
			if err := writer.WriteField("group_id", opts.GroupID); err != nil {
				return nil, fmt.Errorf("failed to add group_id field: %w", err)
			}
		}

		// Use custom name for the folder if provided
		if opts.FileName != "" {
			if err := writer.WriteField("name", opts.FileName); err != nil {
				return nil, fmt.Errorf("failed to add name field: %w", err)
			}
		}

		// Add keyvalues if present
		if len(opts.KeyValues) > 0 {
			keyvaluesJSON, err := json.Marshal(opts.KeyValues)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal keyvalues: %w", err)
			}

			if err := writer.WriteField("keyvalues", string(keyvaluesJSON)); err != nil {
				return nil, fmt.Errorf("failed to add keyvalues field: %w", err)
			}
		}
	}

	// Add all files
	for _, file := range files {
		// Reset file position to start
		if _, err := file.Seek(0, 0); err != nil {
			return nil, fmt.Errorf("failed to reset file position: %w", err)
		}

		part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
		if err != nil {
			return nil, fmt.Errorf("failed to create form file: %w", err)
		}

		if _, err := io.Copy(part, file); err != nil {
			return nil, fmt.Errorf("failed to copy file data: %w", err)
		}
	}

	// Close the writer
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Create the request
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+cfg.PinataJWT)

	// Add custom headers if any
	for key, value := range cfg.CustomHeaders {
		req.Header.Set(key, value)
	}

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var response struct {
		Data *types.UploadResponse `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Data, nil
}

// JSON uploads a JSON object to the public IPFS network
func (s *PrivateService) JSON(data interface{}, opts *JSONOptions) (*types.UploadResponse, error) {
	if data == nil {
		return nil, fmt.Errorf("JSON data is required")
	}

	// Marshal the JSON data
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal JSON data: %w", err)
	}

	// Create a temporary file to hold the JSON data
	tmpFile, err := os.CreateTemp("", "pinata-json-*.json")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up

	if _, err := tmpFile.Write(jsonData); err != nil {
		return nil, fmt.Errorf("failed to write to temporary file: %w", err)
	}

	// Reset file position to start
	if _, err := tmpFile.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to reset file position: %w", err)
	}

	// Create file options
	fileOpts := &FileOptions{
		GroupID:   opts.GroupID,
		KeyValues: opts.KeyValues,
	}

	// Use custom name or default
	if opts.Name != "" {
		fileOpts.FileName = opts.Name
	} else {
		fileOpts.FileName = "data.json"
	}

	// Use the File method to upload
	return s.File(tmpFile, fileOpts)
}

// Base64 uploads base64-encoded data to the public IPFS network
func (s *PrivateService) Base64(data string, opts *Base64Options) (*types.UploadResponse, error) {
	if data == "" {
		return nil, fmt.Errorf("base64 data is required")
	}

	// Decode the base64 data
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 data: %w", err)
	}

	// Create a temporary file to hold the decoded data
	tmpFile, err := os.CreateTemp("", "pinata-base64-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up

	if _, err := tmpFile.Write(decoded); err != nil {
		return nil, fmt.Errorf("failed to write to temporary file: %w", err)
	}

	// Reset file position to start
	if _, err := tmpFile.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to reset file position: %w", err)
	}

	// Create file options
	fileOpts := &FileOptions{
		GroupID:   opts.GroupID,
		KeyValues: opts.KeyValues,
	}

	// Use custom name or default
	if opts.Name != "" {
		fileOpts.FileName = opts.Name
	} else {
		fileOpts.FileName = "file"
	}

	// Use the File method to upload
	return s.File(tmpFile, fileOpts)
}

// URL uploads the content of a URL to the public IPFS network
func (s *PrivateService) URL(targetURL string, opts *URLOptions) (*types.UploadResponse, error) {
	if targetURL == "" {
		return nil, fmt.Errorf("URL is required")
	}

	// Fetch the content from the URL
	client := &http.Client{}
	resp, err := client.Get(targetURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL content: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("URL returned non-OK status: %d", resp.StatusCode)
	}

	// Create a temporary file to hold the content
	tmpFile, err := os.CreateTemp("", "pinata-url-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up

	// Copy the content to the temporary file
	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		return nil, fmt.Errorf("failed to copy URL content: %w", err)
	}

	// Reset file position to start
	if _, err := tmpFile.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to reset file position: %w", err)
	}

	// Create file options
	fileOpts := &FileOptions{
		GroupID:   opts.GroupID,
		KeyValues: opts.KeyValues,
	}

	// Use custom name or extract from URL
	if opts.Name != "" {
		fileOpts.FileName = opts.Name
	} else {
		// Extract filename from URL or use default
		urlPath := strings.Split(targetURL, "/")
		if len(urlPath) > 0 && urlPath[len(urlPath)-1] != "" {
			fileOpts.FileName = urlPath[len(urlPath)-1]
		} else {
			fileOpts.FileName = "file"
		}
	}

	// Use the File method to upload
	return s.File(tmpFile, fileOpts)
}

// CreateSignedURL generates a signed URL for client-side uploads
func (s *PrivateService) CreateSignedURL(opts *SignedUploadOptions) (string, error) {
	if opts == nil || opts.Expires <= 0 {
		return "", fmt.Errorf("expiration time is required")
	}

	cfg := s.config.(*types.Config)
	url := fmt.Sprintf("%s/files/sign", cfg.UploadUrl)

	// Set current time if not provided
	date := opts.Date
	if date == 0 {
		date = time.Now().Unix()
	}

	// Build request payload
	payload := map[string]interface{}{
		"date":    date,
		"expires": opts.Expires,
		"network": "private",
	}

	// Add optional fields
	if opts.GroupID != "" {
		payload["group_id"] = opts.GroupID
	}

	if opts.Name != "" {
		payload["filename"] = opts.Name
	}

	if len(opts.KeyValues) > 0 {
		payload["keyvalues"] = opts.KeyValues
	}

	if opts.Vectorize {
		payload["vectorize"] = true
	}

	if opts.MaxFileSize > 0 {
		payload["max_file_size"] = opts.MaxFileSize
	}

	if len(opts.MimeTypes) > 0 {
		payload["allow_mime_types"] = opts.MimeTypes
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Create the request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.PinataJWT)

	// Add custom headers if any
	for key, value := range cfg.CustomHeaders {
		req.Header.Set(key, value)
	}

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var response struct {
		Data string `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Data, nil
}
