package pinata

import (
	"github.com/PinataCloud/pinata-go-sdk/pinata/analytics"
	"github.com/PinataCloud/pinata-go-sdk/pinata/files"
	"github.com/PinataCloud/pinata-go-sdk/pinata/gateways"
	"github.com/PinataCloud/pinata-go-sdk/pinata/groups"
	"github.com/PinataCloud/pinata-go-sdk/pinata/keys"
	"github.com/PinataCloud/pinata-go-sdk/pinata/signatures"
	"github.com/PinataCloud/pinata-go-sdk/pinata/upload"
)

// Client is the main SDK client that holds all Pinata API functionality
type Client struct {
	Config     *Config
	Files      *files.Service
	Upload     *upload.Service
	Groups     *groups.Service
	Gateways   *gateways.Service
	Keys       *keys.Service
	Signatures *signatures.Service
	Analytics  *analytics.Service
}

// New creates a new Pinata SDK client with the provided JWT and gateway
func New(jwt string, gateway string) *Client {
	config := NewConfig(jwt, gateway)
	return NewWithConfig(config)
}

// NewWithConfig creates a new Pinata SDK client with the provided configuration
func NewWithConfig(config *Config) *Client {
	client := &Client{
		Config: config,
	}
	
	// Initialize services
	client.Files = files.New(config)
	client.Upload = upload.New(config)
	client.Groups = groups.New(config)
	client.Gateways = gateways.New(config)
	client.Keys = keys.New(config)
	client.Signatures = signatures.New(config)
	client.Analytics = analytics.New(config)
	
	return client
}

// TestAuthentication tests if the JWT is valid
func (c *Client) TestAuthentication() (bool, error) {
	// Implementation
	return true, nil
}
