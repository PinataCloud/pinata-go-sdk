package pinata

// Config holds the configuration for the Pinata SDK client
type Config struct {
	PinataJWT        string
	PinataGateway    string
	PinataGatewayKey string
	CustomHeaders    map[string]string
	APIUrl           string
	UploadUrl        string
}

// DefaultAPIUrl is the default API endpoint
const DefaultAPIUrl = "https://api.pinata.cloud/v3"

// DefaultUploadUrl is the default upload endpoint
const DefaultUploadUrl = "https://uploads.pinata.cloud/v3"

// NewConfig creates a default configuration with provided JWT and gateway
func NewConfig(jwt string, gateway string) *Config {
	return &Config{
		PinataJWT:     jwt,
		PinataGateway: gateway,
		APIUrl:        DefaultAPIUrl,
		UploadUrl:     DefaultUploadUrl,
		CustomHeaders: make(map[string]string),
	}
}
