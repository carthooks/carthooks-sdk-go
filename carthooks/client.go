package carthooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

// ClientConfig holds configuration options for the Carthooks client
type ClientConfig struct {
	BaseURL     string
	AccessToken string
	Timeout     time.Duration
	Headers     map[string]string
	Debug       bool
	OAuth       *OAuthConfig
}

// Client represents the Carthooks API client
type Client struct {
	baseURL        string
	accessToken    string
	httpClient     *http.Client
	headers        map[string]string
	debug          bool
	oauthConfig    *OAuthConfig
	currentTokens  *OAuthTokens
	tokenExpiresAt *time.Time
}

// NewClient creates a new Carthooks client with the given configuration
func NewClient(config *ClientConfig) *Client {
	if config == nil {
		config = &ClientConfig{}
	}

	// Set defaults from environment variables or fallback values
	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = os.Getenv("CARTHOOKS_API_URL")
		if baseURL == "" {
			baseURL = "https://api.carthooks.com"
		}
	}

	accessToken := config.AccessToken
	if accessToken == "" {
		accessToken = os.Getenv("CARTHOOKS_ACCESS_TOKEN")
	}

	timeout := config.Timeout
	if timeout == 0 {
		if timeoutStr := os.Getenv("CARTHOOKS_TIMEOUT"); timeoutStr != "" {
			if parsedTimeout, err := time.ParseDuration(timeoutStr + "s"); err == nil {
				timeout = parsedTimeout
			}
		}
		if timeout == 0 {
			timeout = 30 * time.Second
		}
	}

	debug := config.Debug
	if !debug {
		debug = os.Getenv("CARTHOOKS_SDK_DEBUG") == "true"
	}

	// Initialize headers
	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	// Add custom headers
	for k, v := range config.Headers {
		headers[k] = v
	}

	// Add authorization header if token is provided
	if accessToken != "" {
		headers["Authorization"] = "Bearer " + accessToken
	}

	client := &Client{
		baseURL:     baseURL,
		accessToken: accessToken,
		httpClient: &http.Client{
			Timeout: timeout,
		},
		headers: headers,
		debug:   debug,
	}

	// Set OAuth configuration if provided
	if config.OAuth != nil {
		client.oauthConfig = &OAuthConfig{
			ClientID:     config.OAuth.ClientID,
			ClientSecret: config.OAuth.ClientSecret,
			RefreshToken: config.OAuth.RefreshToken,
			AutoRefresh:  config.OAuth.AutoRefresh,
		}
		// Default auto refresh to true if not specified
		if client.oauthConfig.AutoRefresh == false && config.OAuth.RefreshToken != "" {
			client.oauthConfig.AutoRefresh = true
		}
	}

	return client
}

// SetAccessToken sets the access token for API authentication
func (c *Client) SetAccessToken(token string) {
	c.accessToken = token
	c.headers["Authorization"] = "Bearer " + token
}

// makeRequest performs an HTTP request and returns the response
func (c *Client) makeRequest(method, path string, body interface{}, params map[string]string) (*http.Response, error) {
	// Build URL
	fullURL := c.baseURL + path
	if len(params) > 0 {
		u, err := url.Parse(fullURL)
		if err != nil {
			return nil, fmt.Errorf("invalid URL: %w", err)
		}

		q := u.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
		fullURL = u.String()
	}

	// Prepare request body
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	// Create request
	req, err := http.NewRequest(method, fullURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	// Debug logging
	if c.debug {
		fmt.Printf("[DEBUG] %s %s\n", method, fullURL)
		if body != nil {
			if jsonData, err := json.Marshal(body); err == nil {
				fmt.Printf("[DEBUG] Request body: %s\n", string(jsonData))
			}
		}
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

// parseResponse parses the HTTP response into a Result
func (c *Client) parseResponse(resp *http.Response) *Result {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &Result{
			Success: false,
			Error:   fmt.Sprintf("failed to read response body: %v", err),
		}
	}

	if c.debug {
		fmt.Printf("[DEBUG] Response body: %s\n", string(body))
	}

	// Try to parse as JSON
	var apiResp struct {
		Data  interface{} `json:"data"`
		Error *struct {
			Message string `json:"message"`
			Code    string `json:"code"`
		} `json:"error"`
		TraceID string                 `json:"trace_id"`
		Meta    map[string]interface{} `json:"meta"`
	}

	if err := json.Unmarshal(body, &apiResp); err != nil {
		// If JSON parsing fails, treat as error
		return &Result{
			Success: false,
			Error:   string(body),
		}
	}

	result := &Result{
		TraceID: apiResp.TraceID,
		Meta:    apiResp.Meta,
	}

	if apiResp.Error != nil {
		result.Success = false
		result.Error = apiResp.Error.Message
	} else {
		result.Success = true
		result.Data = apiResp.Data
	}

	return result
}
