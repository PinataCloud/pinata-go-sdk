package files

// ListOptions represents options for the List method
type ListOptions struct {
	Name       string
	Group      string
	NoGroup    bool
	CID        string
	CIDPending bool
	MimeType   string
	KeyValues  map[string]string
	Order      string
	Limit      int
	PageToken  string
}

// UpdateOptions represents options for the Update method
type UpdateOptions struct {
	ID        string            `json:"-"`
	Name      string            `json:"name,omitempty"`
	KeyValues map[string]string `json:"keyvalues,omitempty"`
}

// SwapOptions represents options for AddSwap method
type SwapOptions struct {
	CID     string `json:"-"`
	SwapCID string `json:"swap_cid"`
}

// SwapHistoryOptions represents options for GetSwapHistory method
type SwapHistoryOptions struct {
	CID    string `json:"-"`
	Domain string `json:"-"`
}

// DeleteOptions represents options for the Delete method
type DeleteOptions struct {
	IDs []string
}

// PinByHashOptions represents options for the PinByHash method
type PinByHashOptions struct {
	CID       string            `json:"cid"`
	Name      string            `json:"name,omitempty"`
	GroupID   string            `json:"group_id,omitempty"`
	KeyValues map[string]string `json:"keyvalues,omitempty"`
	HostNodes []string          `json:"host_nodes,omitempty"`
}

// PinQueueOptions represents options for querying the pin queue
type PinQueueOptions struct {
	Sort      string
	Status    string
	CID       string
	Limit     int
	PageToken string
}
