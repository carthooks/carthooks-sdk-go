package carthooks

import (
	"fmt"
)

// SubmissionTokenOptions represents options for creating submission tokens
type SubmissionTokenOptions struct {
	CallbackURL string   `json:"callback_url,omitempty"`
	RedirectURL string   `json:"redirect_url,omitempty"`
	TTL         int      `json:"ttl,omitempty"`
	Fields      []string `json:"fields,omitempty"`
}

// SubmissionToken represents a submission token response
type SubmissionToken struct {
	Token     string `json:"token"`
	URL       string `json:"url"`
	ExpiresAt string `json:"expires_at"`
}

// UpdateTokenOptions represents options for creating update tokens
type UpdateTokenOptions struct {
	TTL    int      `json:"ttl,omitempty"`
	Fields []string `json:"fields,omitempty"`
}

// UpdateToken represents an update token response
type UpdateToken struct {
	Token     string `json:"token"`
	URL       string `json:"url"`
	ExpiresAt string `json:"expires_at"`
}

// UploadToken represents an upload token response
type UploadToken struct {
	Token     string `json:"token"`
	URL       string `json:"url"`
	ExpiresAt string `json:"expires_at"`
}

// User represents user information
type User struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Avatar string `json:"avatar,omitempty"`
}

// WatchDataOptions represents options for watching data changes
type WatchDataOptions struct {
	EndpointURL      string                 `json:"endpoint_url"`
	EndpointType     string                 `json:"endpoint_type"`
	Name             string                 `json:"name"`
	AppID            uint                   `json:"app_id"`
	CollectionID     uint                   `json:"collection_id"`
	Filters          map[string]interface{} `json:"filters,omitempty"`
	Age              int                    `json:"age,omitempty"`
	WatchStartTime   int64                  `json:"watch_start_time,omitempty"`
}

// WatchDataResponse represents a watch data response
type WatchDataResponse struct {
	WatchID string `json:"watch_id"`
	Status  string `json:"status"`
}

// GetSubmissionToken gets a submission token for creating items
func (c *Client) GetSubmissionToken(appID, collectionID uint, options *SubmissionTokenOptions) *Result {
	path := fmt.Sprintf("/v1/apps/%d/collections/%d/submission-token", appID, collectionID)
	
	resp, err := c.makeRequest("POST", path, options, nil)
	if err != nil {
		return &Result{
			Success: false,
			Error:   err.Error(),
		}
	}
	
	return c.parseResponse(resp)
}

// UpdateSubmissionToken updates a submission token for an existing item
func (c *Client) UpdateSubmissionToken(appID, collectionID, itemID uint, options *UpdateTokenOptions) *Result {
	path := fmt.Sprintf("/v1/apps/%d/collections/%d/items/%d/update-token", appID, collectionID, itemID)
	
	resp, err := c.makeRequest("POST", path, options, nil)
	if err != nil {
		return &Result{
			Success: false,
			Error:   err.Error(),
		}
	}
	
	return c.parseResponse(resp)
}

// GetUploadToken gets a token for file uploads
func (c *Client) GetUploadToken() *Result {
	path := "/v1/uploads/token"
	
	resp, err := c.makeRequest("POST", path, nil, nil)
	if err != nil {
		return &Result{
			Success: false,
			Error:   err.Error(),
		}
	}
	
	return c.parseResponse(resp)
}

// GetUser gets user information by user ID
func (c *Client) GetUser(userID uint) *Result {
	path := fmt.Sprintf("/v1/users/%d", userID)
	
	resp, err := c.makeRequest("GET", path, nil, nil)
	if err != nil {
		return &Result{
			Success: false,
			Error:   err.Error(),
		}
	}
	
	return c.parseResponse(resp)
}

// GetUserByToken gets user information by token
func (c *Client) GetUserByToken(token string) *Result {
	path := fmt.Sprintf("/v1/user-token/%s", token)
	
	resp, err := c.makeRequest("GET", path, nil, nil)
	if err != nil {
		return &Result{
			Success: false,
			Error:   err.Error(),
		}
	}
	
	return c.parseResponse(resp)
}

// StartWatchData starts data monitoring
func (c *Client) StartWatchData(options *WatchDataOptions) *Result {
	path := "/v1/watch-data"
	
	resp, err := c.makeRequest("POST", path, options, nil)
	if err != nil {
		return &Result{
			Success: false,
			Error:   err.Error(),
		}
	}
	
	return c.parseResponse(resp)
}

// GetCollections gets collections for an app
func (c *Client) GetCollections(appID uint) *Result {
	path := fmt.Sprintf("/v1/apps/%d/collections", appID)
	
	resp, err := c.makeRequest("GET", path, nil, nil)
	if err != nil {
		return &Result{
			Success: false,
			Error:   err.Error(),
		}
	}
	
	return c.parseResponse(resp)
}

// GetCollection gets a specific collection
func (c *Client) GetCollection(appID, collectionID uint) *Result {
	path := fmt.Sprintf("/v1/apps/%d/collections/%d", appID, collectionID)
	
	resp, err := c.makeRequest("GET", path, nil, nil)
	if err != nil {
		return &Result{
			Success: false,
			Error:   err.Error(),
		}
	}
	
	return c.parseResponse(resp)
}

// GetApps gets available apps
func (c *Client) GetApps() *Result {
	path := "/v1/apps"
	
	resp, err := c.makeRequest("GET", path, nil, nil)
	if err != nil {
		return &Result{
			Success: false,
			Error:   err.Error(),
		}
	}
	
	return c.parseResponse(resp)
}

// GetApp gets a specific app
func (c *Client) GetApp(appID uint) *Result {
	path := fmt.Sprintf("/v1/apps/%d", appID)
	
	resp, err := c.makeRequest("GET", path, nil, nil)
	if err != nil {
		return &Result{
			Success: false,
			Error:   err.Error(),
		}
	}
	
	return c.parseResponse(resp)
}

// Collection represents a collection structure
type Collection struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// App represents an application structure
type App struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}
