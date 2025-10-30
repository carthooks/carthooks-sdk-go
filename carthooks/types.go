package carthooks

import "strconv"

type UrlSets struct {
	// Three sizes: original size, 128x128px, 26x26px
	FullSizeUrl string `json:"full_size_url"` // Original size
	ThumbUrl    string `json:"thumb_url"`     // 128x128px
	IconUrl     string `json:"icon_url"`      // 26x26px
}

type ApiImageResult struct {
	Url      *UrlSets `json:"url"`
	Meta     any      `json:"meta"`
	Expired  int      `json:"expired"`
	FileSize int64    `json:"file_size"`
	Created  int      `json:"created"`
}

type RecordFormat struct {
	ID        uint                   `json:"id"`
	Title     string                 `json:"title"`
	CreatedAt int64                  `json:"created_at"`
	UpdatedAt int64                  `json:"updated_at"`
	Creator   uint                   `json:"creator"`
	Fields    map[string]interface{} `json:"fields"`
}

type EventMessage struct {
	Version string           `json:"version"`
	Meta    EventMessageMeta `json:"meta"`
	Payload any              `json:"payload"`
}

type EventCode string

const (
	EventCodeRecordCreated EventCode = "collection.item.created"
	EventCodeRecordUpdated EventCode = "collection.item.updated"
)

type EventMessageMeta struct {
	TenantID     uint      `json:"tenant_id"`
	CollectionID uint      `json:"collection_id"`
	Event        EventCode `json:"event"`
	TriggerType  string    `json:"trigger_type"`
	TriggerName  string    `json:"trigger_name,omitempty"`
}

func (e *EventMessageMeta) ToMap() map[string]string {
	return map[string]string{
		"tenant_id":     strconv.FormatUint(uint64(e.TenantID), 10),
		"collection_id": strconv.FormatUint(uint64(e.CollectionID), 10),
		"trigger_type":  e.TriggerType,
		"trigger_name":  e.TriggerName,
	}
}

// Connection represents a hooklet connection
type Connection struct {
	ID          uint   `json:"id"`
	TenantID    uint   `json:"tenant_id"`
	AppID       uint   `json:"app_id"`
	HookletID   uint   `json:"hooklet_id"`
	DevClientID uint   `json:"dev_client_id"`
	Title       string `json:"title"`
	Status      uint8  `json:"status"`
	IconUrl     string `json:"icon_url"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ConnectionStatus represents connection status constants
type ConnectionStatus uint8

const (
	ConnectionStatusPending  ConnectionStatus = 0
	ConnectionStatusActive   ConnectionStatus = 1
	ConnectionStatusInactive ConnectionStatus = 2
)

// ConnectionLog represents a connection log entry
type ConnectionLog struct {
	ID           uint   `json:"id"`
	TenantID     uint   `json:"tenant_id"`
	ConnectionID uint   `json:"connection_id"`
	Status       uint8  `json:"status"`
	Message      string `json:"message"`
	CreatedAt    string `json:"created_at"`
}

// ConnectionLogStatus represents connection log status constants
type ConnectionLogStatus uint8

const (
	ConnectionLogStatusCreated ConnectionLogStatus = 1
	ConnectionLogStatusUpdated ConnectionLogStatus = 2
	ConnectionLogStatusWarn    ConnectionLogStatus = 3
	ConnectionLogStatusError   ConnectionLogStatus = 4
)

// ConnectionUsage represents connection usage data
type ConnectionUsage struct {
	ID           uint   `json:"id"`
	TenantID     uint   `json:"tenant_id"`
	ConnectionID uint   `json:"connection_id"`
	Usage        int64  `json:"usage"`
	CreatedAt    string `json:"created_at"`
}

// CreateConnectionRequest represents the request body for creating a connection
type CreateConnectionRequest struct {
	HookletID    string `json:"hooklet_id"`
	Title        string `json:"title"`
	IconUrl      string `json:"icon_url,omitempty"`
	Description  string `json:"description,omitempty"`
	VendorTaskID string `json:"vendor_task_id"`
}

// CreateConnectionLogRequest represents the request body for creating a connection log
type CreateConnectionLogRequest struct {
	Status  uint8  `json:"status"`
	Message string `json:"message"`
}

// CreateConnectionUsageRequest represents the request body for creating connection usage
type CreateConnectionUsageRequest struct {
	Usage int64 `json:"usage"`
}

// UpdateConnectionRequest represents the request body for updating a connection
type UpdateConnectionRequest struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"` // "active", "inactive"
	IconUrl     string `json:"icon_url,omitempty"`
}

// OAuth related types

// OAuthConfig holds OAuth configuration
type OAuthConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RefreshToken string `json:"refresh_token,omitempty"`
	AutoRefresh  bool   `json:"auto_refresh"`
}

// OAuthTokens represents OAuth token response
type OAuthTokens struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope"`
}

// OAuthTokenRequest represents OAuth token request
type OAuthTokenRequest struct {
	GrantType       string `json:"grant_type"`
	ClientID        string `json:"client_id"`
	ClientSecret    string `json:"client_secret"`
	UserAccessToken string `json:"user_access_token,omitempty"`
	Code            string `json:"code,omitempty"`
	RedirectURI     string `json:"redirect_uri,omitempty"`
	RefreshToken    string `json:"refresh_token,omitempty"`
}

// OAuthAuthorizeCodeRequest represents OAuth authorization code request
type OAuthAuthorizeCodeRequest struct {
	ClientID       string `json:"client_id"`
	RedirectURI    string `json:"redirect_uri"`
	State          string `json:"state"`
	TargetTenantID uint   `json:"target_tenant_id,omitempty"`
}

// OAuthAuthorizeCodeResponse represents OAuth authorization code response
type OAuthAuthorizeCodeResponse struct {
	RedirectURL string `json:"redirect_url"`
}

// UserInfo represents current user information
type UserInfo struct {
	UserID     uint     `json:"user_id"`
	Username   string   `json:"username"`
	Email      string   `json:"email"`
	TenantID   uint     `json:"tenant_id"`
	TenantName string   `json:"tenant_name"`
	IsAdmin    bool     `json:"is_admin"`
	Scope      []string `json:"scope"`
}
