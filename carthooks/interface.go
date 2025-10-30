package carthooks

// ClientInterface defines the interface for Carthooks SDK client
// This interface allows for easy mocking in tests
type ClientInterface interface {
	// Basic client methods
	SetAccessToken(token string)
	GetBaseURL() string
	
	// OAuth methods
	GetOAuthToken(request *OAuthTokenRequest) *Result
	RefreshOAuthToken(refreshToken ...string) *Result
	InitializeOAuth(userAccessToken ...string) *Result
	ExchangeAuthorizationCode(code, redirectURI string) *Result
	GetOAuthAuthorizeCode(request *OAuthAuthorizeCodeRequest) *Result
	GetCurrentUser() *Result
	GetUserTenants() *Result
	EnsureValidToken() error
	GetCurrentTokens() *OAuthTokens
	SetOAuthConfig(config *OAuthConfig)
	GetOAuthConfig() *OAuthConfig
	
	// Collection/Item methods
	GetItems(appID, collectionID uint, limit, start int, options map[string]string) *Result
	GetItemByID(appID, collectionID, itemID uint, fields []string) *Result
	QueryItems(appID, collectionID uint, options *QueryOptions) *Result
	CreateItem(appID, collectionID uint, data map[string]interface{}) *Result
	UpdateItem(appID, collectionID, itemID uint, data map[string]interface{}) *Result
	DeleteItem(appID, collectionID, itemID uint) *Result
	LockItem(appID, collectionID, itemID uint, options *LockOptions) *Result
	UnlockItem(appID, collectionID, itemID uint, lockID string) *Result
	
	// SubItem methods
	CreateSubItem(appID, collectionID, itemID, fieldID uint, data map[string]interface{}) *Result
	UpdateSubItem(appID, collectionID, itemID, fieldID, subItemID uint, data map[string]interface{}) *Result
	DeleteSubItem(appID, collectionID, itemID, fieldID, subItemID uint) *Result
	
	// Connection methods
	CreateConnection(appID uint, request *CreateConnectionRequest) *Result
	UpdateConnection(appID, connectionID uint, request *UpdateConnectionRequest) *Result
	GetConnection(appID, connectionID uint) *Result
	DeleteConnection(appID, connectionID uint) *Result
	CreateConnectionLog(appID, connectionID uint, request *CreateConnectionLogRequest) *Result
	CreateConnectionUsage(appID, connectionID uint, request *CreateConnectionUsageRequest) *Result
	
	// Advanced methods
	GetSubmissionToken(appID, collectionID uint, options *SubmissionTokenOptions) *Result
	UpdateSubmissionToken(appID, collectionID, itemID uint, options *UpdateTokenOptions) *Result
	GetUploadToken() *Result
	GetUser(userID uint) *Result
	GetUserByToken(token string) *Result
	StartWatchData(options *WatchDataOptions) *Result
	StopWatchData(options *WatchDataOptions) *Result
	GetCollections(appID uint) *Result
	GetCollection(appID, collectionID uint) *Result
	GetApps() *Result
	GetApp(appID uint) *Result
}

// Ensure Client implements ClientInterface
var _ ClientInterface = (*Client)(nil)

