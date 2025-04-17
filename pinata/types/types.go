package types

// Config holds the configuration for the Pinata SDK client
type Config struct {
	PinataJWT        string
	PinataGateway    string
	PinataGatewayKey string
	CustomHeaders    map[string]string
	APIUrl           string
	UploadUrl        string
}

// File represents a file stored on Pinata
type File struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	CID           string            `json:"cid"`
	Size          int64             `json:"size"`
	CreatedAt     string            `json:"created_at"`
	NumberOfFiles int               `json:"number_of_files"`
	MimeType      string            `json:"mime_type"`
	GroupID       *string           `json:"group_id"`
	KeyValues     map[string]string `json:"keyvalues"`
	Vectorized    bool              `json:"vectorized"`
	Network       string            `json:"network,omitempty"`
	IsDuplicate   bool              `json:"is_duplicate,omitempty"`
}

// FileListResponse represents the response for listing files
type FileListResponse struct {
	Files         []File `json:"files"`
	NextPageToken string `json:"next_page_token"`
}

// DeleteResponse represents the response for deleting a file
type DeleteResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

// SwapResponse represents a CID swap record
type SwapResponse struct {
	MappedCID string `json:"mapped_cid"`
	CreatedAt string `json:"created_at"`
}

// PinByHashResponse represents the response for pinning by hash
type PinByHashResponse struct {
	ID         string            `json:"id"`
	CID        string            `json:"cid"`
	Status     string            `json:"status"`
	Name       string            `json:"name"`
	DateQueued string            `json:"date_queued"`
	KeyValues  map[string]string `json:"keyvalues"`
	HostNodes  []string          `json:"host_nodes"`
	GroupID    *string           `json:"group_id"`
}

// PinQueueItem represents an item in the pin queue
type PinQueueItem struct {
	ID         string            `json:"id"`
	CID        string            `json:"cid"`
	Status     string            `json:"status"`
	Name       string            `json:"name"`
	DateQueued string            `json:"date_queued"`
	KeyValues  map[string]string `json:"keyvalues"`
	HostNodes  []string          `json:"host_nodes"`
	GroupID    *string           `json:"group_id"`
}

// PinQueueResponse represents the response for pin queue listing
type PinQueueResponse struct {
	Items         []PinQueueItem `json:"jobs"`
	NextPageToken string         `json:"next_page_token"`
}

// AccessLinkOptions represents options for creating an access link
type AccessLinkOptions struct {
	CID     string
	Date    int64
	Expires int
	Gateway string
}

// VectorizeResponse represents the response for vectorizing a file
type VectorizeResponse struct {
	Status bool `json:"status"`
}

// VectorQueryOptions represents options for querying vectors
type VectorQueryOptions struct {
	GroupID    string
	Query      string
	ReturnFile bool
}

// VectorMatch represents a match from a vector query
type VectorMatch struct {
	FileID string  `json:"file_id"`
	CID    string  `json:"cid"`
	Score  float64 `json:"score"`
}

// VectorQueryResponse represents the response for a vector query
type VectorQueryResponse struct {
	Count       int           `json:"count,omitempty"`
	Matches     []VectorMatch `json:"matches,omitempty"`
	ContentType string        `json:"-"`
	Data        []byte        `json:"-"`
}

// Group represents a file group on Pinata
type Group struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	IsPublic  bool   `json:"is_public"`
	CreatedAt string `json:"created_at"`
}

// GroupListResponse represents the response for listing groups
type GroupListResponse struct {
	Groups        []Group `json:"groups"`
	NextPageToken string  `json:"next_page_token"`
}

// UploadResponse represents the response from an upload
type UploadResponse struct {
	ID            string            `json:"id"`
	Name          string            `json:"name"`
	CID           string            `json:"cid"`
	Size          int64             `json:"size"`
	CreatedAt     string            `json:"created_at"`
	NumberOfFiles int               `json:"number_of_files"`
	MimeType      string            `json:"mime_type"`
	GroupID       *string           `json:"group_id"`
	KeyValues     map[string]string `json:"keyvalues"`
	Vectorized    bool              `json:"vectorized"`
	Network       string            `json:"network,omitempty"`
	IsDuplicate   bool              `json:"is_duplicate,omitempty"`
}

// Key represents an API key
type Key struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Key       string    `json:"key"`
	Secret    string    `json:"secret"`
	MaxUses   int       `json:"max_uses"`
	Uses      int       `json:"uses"`
	UserID    string    `json:"user_id"`
	Scopes    KeyScopes `json:"scopes"`
	Revoked   bool      `json:"revoked"`
	CreatedAt string    `json:"createdAt"`
	UpdatedAt string    `json:"updatedAt"`
}

// KeyScopes represents the scopes for an API key
type KeyScopes struct {
	Admin     bool `json:"admin"`
	Endpoints struct {
		Pinning struct {
			PinFileToIPFS bool `json:"pinFileToIPFS"`
			PinJSONToIPFS bool `json:"pinJSONToIPFS"`
		} `json:"pinning"`
	} `json:"endpoints"`
}

// KeyListResponse represents the response for listing API keys
type KeyListResponse struct {
	Keys  []Key `json:"keys"`
	Count int   `json:"count"`
}

// SignatureResponse represents a signature for a CID
type SignatureResponse struct {
	CID       string `json:"cid"`
	Signature string `json:"signature"`
}

// AnalyticsItem represents an analytics data point
type AnalyticsItem struct {
	Value     string `json:"value"`
	Requests  int    `json:"requests"`
	Bandwidth int64  `json:"bandwidth"`
}

// AnalyticsResponse represents the response for analytics queries
type AnalyticsResponse struct {
	Data []AnalyticsItem `json:"data"`
}

// TimePeriod represents a time period for time series analytics
type TimePeriod struct {
	PeriodStartTime string `json:"period_start_time"`
	Requests        int    `json:"requests"`
	Bandwidth       int64  `json:"bandwidth"`
}

// TimeSeriesResponse represents the response for time series analytics
type TimeSeriesResponse struct {
	TotalRequests  int          `json:"total_requests"`
	TotalBandwidth int64        `json:"total_bandwidth"`
	TimePeriods    []TimePeriod `json:"time_periods"`
}
