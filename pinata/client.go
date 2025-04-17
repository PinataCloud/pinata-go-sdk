package pinata

import (
	"fmt"
	"io"
	"net/http"

	"github.com/PinataCloud/pinata-go-sdk/pinata/files"
	"github.com/PinataCloud/pinata-go-sdk/pinata/types"
)

// Client is the main Pinata SDK client
type Client struct {
	Config *types.Config
	Files  *files.Service
}

// DefaultAPIURL is the default API endpoint
const DefaultAPIURL = "https://api.pinata.cloud/v3"

// DefaultUploadURL is the default upload endpoint
const DefaultUploadURL = "https://uploads.pinata.cloud/v3"

// New creates a new Pinata SDK client with the provided JWT and gateway
func New(jwt string, gateway string) *Client {
	config := &types.Config{
		PinataJWT:     jwt,
		PinataGateway: gateway,
		APIUrl:        DefaultAPIURL,
		UploadUrl:     DefaultUploadURL,
		CustomHeaders: make(map[string]string),
	}

	return NewWithConfig(config)
}

// NewWithConfig creates a new Pinata SDK client with a custom configuration
func NewWithConfig(config *types.Config) *Client {
	client := &Client{
		Config: config,
	}

	// Initialize the services with the configuration
	client.Files = files.New(config)

	return client
}

// TestAuthentication tests if the JWT is valid
func (c *Client) TestAuthentication() (bool, error) {
	url := fmt.Sprintf("https://api.pinata.cloud/data/testAuthentication")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Config.PinataJWT)

	// Add custom headers if any
	for key, value := range c.Config.CustomHeaders {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("authentication failed (status %d): %s", resp.StatusCode, string(body))
	}

	return true, nil
}
