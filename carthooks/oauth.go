package carthooks

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// GetOAuthToken gets OAuth token using various grant types
func (c *Client) GetOAuthToken(request *OAuthTokenRequest) *Result {
	// Use form-encoded data for OAuth token requests (OAuth 2.0 standard)
	formData := url.Values{}
	formData.Set("grant_type", request.GrantType)
	formData.Set("client_id", request.ClientID)
	formData.Set("client_secret", request.ClientSecret)

	if request.UserAccessToken != "" {
		formData.Set("user_access_token", request.UserAccessToken)
	}
	if request.Code != "" {
		formData.Set("code", request.Code)
	}
	if request.RedirectURI != "" {
		formData.Set("redirect_uri", request.RedirectURI)
	}
	if request.RefreshToken != "" {
		formData.Set("refresh_token", request.RefreshToken)
	}

	// Create a custom request for form data
	resp, err := c.makeFormRequest("POST", "/oauth/token", formData)
	if err != nil {
		return &Result{
			Success: false,
			Error:   err.Error(),
		}
	}

	result := c.parseResponse(resp)

	// Store tokens if this is our client and request was successful
	if result.Success && c.oauthConfig != nil && request.ClientID == c.oauthConfig.ClientID {
		if tokenData, ok := result.Data.(map[string]interface{}); ok {
			tokens := &OAuthTokens{}
			if accessToken, ok := tokenData["access_token"].(string); ok {
				tokens.AccessToken = accessToken
			}
			if tokenType, ok := tokenData["token_type"].(string); ok {
				tokens.TokenType = tokenType
			}
			if expiresIn, ok := tokenData["expires_in"].(float64); ok {
				tokens.ExpiresIn = int(expiresIn)
			}
			if refreshToken, ok := tokenData["refresh_token"].(string); ok {
				tokens.RefreshToken = refreshToken
			}
			if scope, ok := tokenData["scope"].(string); ok {
				tokens.Scope = scope
			}

			// Store tokens and expiration time
			c.currentTokens = tokens
			if tokens.ExpiresIn > 0 {
				expiresAt := time.Now().Add(time.Duration(tokens.ExpiresIn) * time.Second)
				c.tokenExpiresAt = &expiresAt
			}

			// Update authorization header
			c.SetAccessToken(tokens.AccessToken)
		}
	}

	return result
}

// RefreshOAuthToken refreshes the OAuth token using refresh token
func (c *Client) RefreshOAuthToken(refreshToken ...string) *Result {
	if c.oauthConfig == nil {
		return &Result{
			Success: false,
			Error:   "OAuth configuration not provided",
		}
	}

	var tokenToUse string
	if len(refreshToken) > 0 && refreshToken[0] != "" {
		tokenToUse = refreshToken[0]
	} else if c.oauthConfig.RefreshToken != "" {
		tokenToUse = c.oauthConfig.RefreshToken
	} else if c.currentTokens != nil && c.currentTokens.RefreshToken != "" {
		tokenToUse = c.currentTokens.RefreshToken
	}

	if tokenToUse == "" {
		return &Result{
			Success: false,
			Error:   "No refresh token available",
		}
	}

	request := &OAuthTokenRequest{
		GrantType:    "refresh_token",
		ClientID:     c.oauthConfig.ClientID,
		ClientSecret: c.oauthConfig.ClientSecret,
		RefreshToken: tokenToUse,
	}

	return c.GetOAuthToken(request)
}

// InitializeOAuth initializes OAuth with client credentials
func (c *Client) InitializeOAuth(userAccessToken ...string) *Result {
	if c.oauthConfig == nil {
		return &Result{
			Success: false,
			Error:   "OAuth configuration not provided",
		}
	}

	request := &OAuthTokenRequest{
		GrantType:    "client_credentials",
		ClientID:     c.oauthConfig.ClientID,
		ClientSecret: c.oauthConfig.ClientSecret,
	}

	if len(userAccessToken) > 0 && userAccessToken[0] != "" {
		request.UserAccessToken = userAccessToken[0]
	}

	return c.GetOAuthToken(request)
}

// ExchangeAuthorizationCode exchanges authorization code for tokens
func (c *Client) ExchangeAuthorizationCode(code, redirectURI string) *Result {
	if c.oauthConfig == nil {
		return &Result{
			Success: false,
			Error:   "OAuth configuration not provided",
		}
	}

	request := &OAuthTokenRequest{
		GrantType:    "authorization_code",
		ClientID:     c.oauthConfig.ClientID,
		ClientSecret: c.oauthConfig.ClientSecret,
		Code:         code,
		RedirectURI:  redirectURI,
	}

	return c.GetOAuthToken(request)
}

// GetOAuthAuthorizeCode gets OAuth authorization code (requires authentication)
func (c *Client) GetOAuthAuthorizeCode(request *OAuthAuthorizeCodeRequest) *Result {
	resp, err := c.makeRequest("POST", "/oauth/get-authorize-code", request, nil)
	if err != nil {
		return &Result{
			Success: false,
			Error:   err.Error(),
		}
	}

	return c.parseResponse(resp)
}

// GetCurrentUser gets current user information (requires OAuth token)
func (c *Client) GetCurrentUser() *Result {
	resp, err := c.makeRequest("GET", "/v1/me", nil, nil)
	if err != nil {
		return &Result{
			Success: false,
			Error:   err.Error(),
		}
	}

	return c.parseResponse(resp)
}

// EnsureValidToken checks if token needs refresh and refreshes if necessary
func (c *Client) EnsureValidToken() error {
	if c.oauthConfig == nil || !c.oauthConfig.AutoRefresh || c.tokenExpiresAt == nil {
		return nil
	}

	// Check if token expires within 5 minutes
	fiveMinutesFromNow := time.Now().Add(5 * time.Minute)
	if c.tokenExpiresAt.After(fiveMinutesFromNow) {
		return nil
	}

	// Try to refresh token
	result := c.RefreshOAuthToken()
	if !result.Success {
		return fmt.Errorf("failed to refresh token: %s", result.Error)
	}

	return nil
}

// GetCurrentTokens returns the current OAuth tokens
func (c *Client) GetCurrentTokens() *OAuthTokens {
	return c.currentTokens
}

// SetOAuthConfig sets the OAuth configuration
func (c *Client) SetOAuthConfig(config *OAuthConfig) {
	c.oauthConfig = &OAuthConfig{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RefreshToken: config.RefreshToken,
		AutoRefresh:  config.AutoRefresh,
	}
}

// GetOAuthConfig returns the current OAuth configuration
func (c *Client) GetOAuthConfig() *OAuthConfig {
	return c.oauthConfig
}

// makeFormRequest makes an HTTP request with form-encoded data
func (c *Client) makeFormRequest(method, path string, formData url.Values) (*http.Response, error) {
	// Build URL
	fullURL := c.baseURL + path

	// Create request with form data
	req, err := http.NewRequest(method, fullURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set form content type
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	// Add other headers (except Authorization for OAuth token requests)
	for k, v := range c.headers {
		if k != "Authorization" && k != "Content-Type" {
			req.Header.Set(k, v)
		}
	}

	// Debug logging
	if c.debug {
		fmt.Printf("[DEBUG] %s %s\n", method, fullURL)
		fmt.Printf("[DEBUG] Form data: %s\n", formData.Encode())
	}

	// Make request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Debug response
	if c.debug {
		fmt.Printf("[DEBUG] Response status: %s\n", resp.Status)
	}

	return resp, nil
}
