package files

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PinataCloud/pinata-go-sdk/pinata/types"
)

// PrivateService provides operations for managing files on the private IPFS network
type PrivateService struct {
	config interface{}
}

// NewPrivateService creates a new PrivateService with the provided configuration
func NewPrivateService(config interface{}) *PrivateService {
	return &PrivateService{
		config: config,
	}
}

// Get retrieves a file by ID from the private IPFS network
func (s *PrivateService) Get(id string) (*types.File, error) {
	cfg := s.config.(*types.Config)
	url := fmt.Sprintf("%s/files/private/%s", cfg.APIUrl, id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+cfg.PinataJWT)
	req.Header.Set("Content-Type", "application/json")

	// Add custom headers if any
	for key, value := range cfg.CustomHeaders {
		req.Header.Set(key, value)
	}

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

	var response struct {
		Data *types.File `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Data, nil
}

// List retrieves a list of files from the private IPFS network
func (s *PrivateService) List(opts *ListOptions) (*types.FileListResponse, error) {
	cfg := s.config.(*types.Config)
	baseURL := fmt.Sprintf("%s/files/private", cfg.APIUrl)

	// Build query parameters
	params := url.Values{}

	if opts != nil {
		if opts.Name != "" {
			params.Add("name", opts.Name)
		}
		if opts.Group != "" {
			params.Add("group", opts.Group)
		}
		if opts.NoGroup {
			params.Add("group", "null")
		}
		if opts.CID != "" {
			params.Add("cid", opts.CID)
		}
		if opts.CIDPending {
			params.Add("cidPending", "true")
		}
		if opts.MimeType != "" {
			params.Add("mimeType", opts.MimeType)
		}
		if opts.Order != "" {
			params.Add("order", opts.Order)
		}
		if opts.Limit > 0 {
			params.Add("limit", strconv.Itoa(opts.Limit))
		}
		if opts.PageToken != "" {
			params.Add("pageToken", opts.PageToken)
		}

		// Add keyvalues if present
		if len(opts.KeyValues) > 0 {
			for key, value := range opts.KeyValues {
				params.Add(fmt.Sprintf("keyvalues[%s]", key), value)
			}
		}
	}

	// Append query parameters if any
	requestURL := baseURL
	if len(params) > 0 {
		requestURL = fmt.Sprintf("%s?%s", baseURL, params.Encode())
	}

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+cfg.PinataJWT)
	req.Header.Set("Content-Type", "application/json")

	// Add custom headers if any
	for key, value := range cfg.CustomHeaders {
		req.Header.Set(key, value)
	}

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

	var response struct {
		Data *types.FileListResponse `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Data, nil
}

// Update updates file metadata
func (s *PrivateService) Update(opts *UpdateOptions) (*types.File, error) {
	if opts == nil || opts.ID == "" {
		return nil, fmt.Errorf("file ID is required")
	}

	cfg := s.config.(*types.Config)
	url := fmt.Sprintf("%s/files/private/%s", cfg.APIUrl, opts.ID)

	payload, err := json.Marshal(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+cfg.PinataJWT)
	req.Header.Set("Content-Type", "application/json")

	// Add custom headers if any
	for key, value := range cfg.CustomHeaders {
		req.Header.Set(key, value)
	}

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

	var response struct {
		Data *types.File `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Data, nil
}

// Delete removes files by their IDs
func (s *PrivateService) Delete(ids []string) ([]types.DeleteResponse, error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("at least one file ID is required")
	}

	cfg := s.config.(*types.Config)

	var responses []types.DeleteResponse

	// Process each ID individually
	for _, id := range ids {
		url := fmt.Sprintf("%s/files/private/%s", cfg.APIUrl, id)

		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Authorization", "Bearer "+cfg.PinataJWT)
		req.Header.Set("Content-Type", "application/json")

		// Add custom headers if any
		for key, value := range cfg.CustomHeaders {
			req.Header.Set(key, value)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to send request: %w", err)
		}

		// Check response status
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
		}

		// Add to successful deletions
		responses = append(responses, types.DeleteResponse{
			ID:     id,
			Status: "deleted",
		})

		resp.Body.Close()
	}

	return responses, nil
}

// AddSwap creates a CID swap
func (s *PrivateService) AddSwap(opts *SwapOptions) (*types.SwapResponse, error) {
	if opts == nil || opts.CID == "" || opts.SwapCID == "" {
		return nil, fmt.Errorf("CID and swap CID are required")
	}

	cfg := s.config.(*types.Config)
	url := fmt.Sprintf("%s/files/private/swap/%s", cfg.APIUrl, opts.CID)

	payload := struct {
		SwapCID string `json:"swap_cid"`
	}{
		SwapCID: opts.SwapCID,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+cfg.PinataJWT)
	req.Header.Set("Content-Type", "application/json")

	// Add custom headers if any
	for key, value := range cfg.CustomHeaders {
		req.Header.Set(key, value)
	}

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

	var response struct {
		Data *types.SwapResponse `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Data, nil
}

// GetSwapHistory retrieves the swap history for a CID
func (s *PrivateService) GetSwapHistory(opts *SwapHistoryOptions) ([]types.SwapResponse, error) {
	if opts == nil || opts.CID == "" || opts.Domain == "" {
		return nil, fmt.Errorf("CID and domain are required")
	}

	cfg := s.config.(*types.Config)
	url := fmt.Sprintf("%s/files/private/swap/%s?domain=%s", cfg.APIUrl, opts.CID, url.QueryEscape(opts.Domain))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+cfg.PinataJWT)
	req.Header.Set("Content-Type", "application/json")

	// Add custom headers if any
	for key, value := range cfg.CustomHeaders {
		req.Header.Set(key, value)
	}

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

	var response struct {
		Data []types.SwapResponse `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Data, nil
}

// DeleteSwap removes a CID swap
func (s *PrivateService) DeleteSwap(cid string) error {
	if cid == "" {
		return fmt.Errorf("CID is required")
	}

	cfg := s.config.(*types.Config)
	url := fmt.Sprintf("%s/files/private/swap/%s", cfg.APIUrl, cid)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+cfg.PinataJWT)
	req.Header.Set("Content-Type", "application/json")

	// Add custom headers if any
	for key, value := range cfg.CustomHeaders {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// CreateAccessLink generates a temporary access link for a private IPFS file
func (s *PrivateService) CreateAccessLink(opts *types.AccessLinkOptions) (string, error) {
	if opts == nil || opts.CID == "" || opts.Expires <= 0 {
		return "", fmt.Errorf("CID and expiration time are required")
	}

	cfg := s.config.(*types.Config)
	url := fmt.Sprintf("%s/files/private/download_link", cfg.APIUrl)

	// Set default gateway if not provided
	gateway := opts.Gateway
	if gateway == "" {
		gateway = cfg.PinataGateway
	}

	// Set current time if not provided
	date := opts.Date
	if date == 0 {
		date = time.Now().Unix()
	}

	payload := struct {
		URL     string `json:"url"`
		Date    int64  `json:"date"`
		Expires int    `json:"expires"`
		Method  string `json:"method"`
	}{
		URL:     fmt.Sprintf("https://%s.mypinata.cloud/files/%s", gateway, opts.CID),
		Date:    date,
		Expires: opts.Expires,
		Method:  "GET",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+cfg.PinataJWT)
	req.Header.Set("Content-Type", "application/json")

	// Add custom headers if any
	for key, value := range cfg.CustomHeaders {
		req.Header.Set(key, value)
	}

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

	var response struct {
		Data string `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Clean up the URL (remove escaping)
	accessLink := strings.ReplaceAll(response.Data, "\\u0026", "&")
	accessLink = strings.Trim(accessLink, "\"")

	return accessLink, nil
}

// Vectorize adds vectors to a file for text search
func (s *PrivateService) Vectorize(fileID string) (*types.VectorizeResponse, error) {
	if fileID == "" {
		return nil, fmt.Errorf("file ID is required")
	}

	cfg := s.config.(*types.Config)
	url := fmt.Sprintf("%s/vectorize/files/%s", cfg.APIUrl, fileID)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+cfg.PinataJWT)
	req.Header.Set("Content-Type", "application/json")

	// Add custom headers if any
	for key, value := range cfg.CustomHeaders {
		req.Header.Set(key, value)
	}

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

	var response struct {
		Data *types.VectorizeResponse `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Data, nil
}

// DeleteVectors removes vectors from a file
func (s *PrivateService) DeleteVectors(fileID string) (*types.VectorizeResponse, error) {
	if fileID == "" {
		return nil, fmt.Errorf("file ID is required")
	}

	cfg := s.config.(*types.Config)
	url := fmt.Sprintf("%s/vectorize/files/%s", cfg.APIUrl, fileID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+cfg.PinataJWT)
	req.Header.Set("Content-Type", "application/json")

	// Add custom headers if any
	for key, value := range cfg.CustomHeaders {
		req.Header.Set(key, value)
	}

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

	var response struct {
		Data *types.VectorizeResponse `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Data, nil
}

// QueryVectors searches for files using vector similarity
func (s *PrivateService) QueryVectors(opts *types.VectorQueryOptions) (*types.VectorQueryResponse, error) {
	if opts == nil || opts.GroupID == "" || opts.Query == "" {
		return nil, fmt.Errorf("group ID and query text are required")
	}

	cfg := s.config.(*types.Config)
	url := fmt.Sprintf("%s/vectorize/groups/%s/query", cfg.APIUrl, opts.GroupID)

	payload := struct {
		Text string `json:"text"`
	}{
		Text: opts.Query,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+cfg.PinataJWT)
	req.Header.Set("Content-Type", "application/json")

	// Add custom headers if any
	for key, value := range cfg.CustomHeaders {
		req.Header.Set(key, value)
	}

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

	// Handle the response based on returnFile option
	if opts.ReturnFile {
		// Handle file response
		contentType := resp.Header.Get("Content-Type")
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		// Create response with file content
		return &types.VectorQueryResponse{
			ContentType: contentType,
			Data:        body,
		}, nil
	}

	// Standard query response
	var response struct {
		Data *types.VectorQueryResponse `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Data, nil
}
