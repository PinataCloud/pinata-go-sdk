package files

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	types "github.com/PinataCloud/pinata-go-sdk/pinata/types"
)

// PublicService provides operations for managing files on the public IPFS network
type PublicService struct {
	config interface{}
}

// NewPublicService creates a new PublicService with the provided configuration
func NewPublicService(config interface{}) *PublicService {
	return &PublicService{
		config: config,
	}
}

// Get retrieves a file by ID from the public IPFS network
func (s *PublicService) Get(id string) (*types.File, error) {
	cfg := s.config.(*types.Config)
	url := fmt.Sprintf("%s/files/public/%s", cfg.APIUrl, id)

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

// List retrieves a list of files from the public IPFS network
func (s *PublicService) List(opts *ListOptions) (*types.FileListResponse, error) {
	cfg := s.config.(*types.Config)
	baseURL := fmt.Sprintf("%s/files/public", cfg.APIUrl)

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
func (s *PublicService) Update(opts *UpdateOptions) (*types.File, error) {
	if opts == nil || opts.ID == "" {
		return nil, fmt.Errorf("file ID is required")
	}

	cfg := s.config.(*types.Config)
	url := fmt.Sprintf("%s/files/public/%s", cfg.APIUrl, opts.ID)

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
func (s *PublicService) Delete(ids []string) ([]types.DeleteResponse, error) {
	if len(ids) == 0 {
		return nil, fmt.Errorf("at least one file ID is required")
	}

	cfg := s.config.(*types.Config)

	var responses []types.DeleteResponse

	// Process each ID individually
	for _, id := range ids {
		url := fmt.Sprintf("%s/files/public/%s", cfg.APIUrl, id)

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
func (s *PublicService) AddSwap(opts *SwapOptions) (*types.SwapResponse, error) {
	if opts == nil || opts.CID == "" || opts.SwapCID == "" {
		return nil, fmt.Errorf("CID and swap CID are required")
	}

	cfg := s.config.(*types.Config)
	url := fmt.Sprintf("%s/files/public/swap/%s", cfg.APIUrl, opts.CID)

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
func (s *PublicService) GetSwapHistory(opts *SwapHistoryOptions) ([]types.SwapResponse, error) {
	if opts == nil || opts.CID == "" || opts.Domain == "" {
		return nil, fmt.Errorf("CID and domain are required")
	}

	cfg := s.config.(*types.Config)
	url := fmt.Sprintf("%s/files/public/swap/%s?domain=%s", cfg.APIUrl, opts.CID, url.QueryEscape(opts.Domain))

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
func (s *PublicService) DeleteSwap(cid string) error {
	if cid == "" {
		return fmt.Errorf("CID is required")
	}

	cfg := s.config.(*types.Config)
	url := fmt.Sprintf("%s/files/public/swap/%s", cfg.APIUrl, cid)

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

// PinByHash pins a CID that already exists on IPFS
func (s *PublicService) PinByHash(opts *PinByHashOptions) (*types.PinByHashResponse, error) {
	if opts == nil || opts.CID == "" {
		return nil, fmt.Errorf("CID is required")
	}

	cfg := s.config.(*types.Config)
	url := fmt.Sprintf("%s/files/public/pin_by_cid", cfg.APIUrl)

	jsonData, err := json.Marshal(opts)
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

	var response struct {
		Data *types.PinByHashResponse `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Data, nil
}

// Queue returns a list of pin by hash requests
func (s *PublicService) Queue(opts *PinQueueOptions) (*types.PinQueueResponse, error) {
	cfg := s.config.(*types.Config)
	baseURL := fmt.Sprintf("%s/files/public/pin_by_cid", cfg.APIUrl)

	// Build query parameters
	params := url.Values{}

	if opts != nil {
		if opts.Sort != "" {
			params.Add("order", opts.Sort)
		}
		if opts.Status != "" {
			params.Add("status", opts.Status)
		}
		if opts.CID != "" {
			params.Add("cid", opts.CID)
		}
		if opts.Limit > 0 {
			params.Add("limit", strconv.Itoa(opts.Limit))
		}
		if opts.PageToken != "" {
			params.Add("pageToken", opts.PageToken)
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
		Data *types.PinQueueResponse `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Data, nil
}

// CancelPinRequest cancels a pin by hash request
func (s *PublicService) CancelPinRequest(id string) error {
	if id == "" {
		return fmt.Errorf("request ID is required")
	}

	cfg := s.config.(*types.Config)
	url := fmt.Sprintf("%s/files/public/pin_by_cid/%s", cfg.APIUrl, id)

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
