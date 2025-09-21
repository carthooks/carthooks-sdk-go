package carthooks

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestOAuthConfig(t *testing.T) {
	config := &ClientConfig{
		OAuth: &OAuthConfig{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			AutoRefresh:  true,
		},
	}

	client := NewClient(config)

	if client.oauthConfig == nil {
		t.Fatal("OAuth config should be set")
	}

	if client.oauthConfig.ClientID != "test-client-id" {
		t.Errorf("Expected client ID 'test-client-id', got '%s'", client.oauthConfig.ClientID)
	}

	if !client.oauthConfig.AutoRefresh {
		t.Error("AutoRefresh should be true")
	}
}

func TestInitializeOAuth(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/open/api/oauth/token" {
			t.Errorf("Expected path '/open/api/oauth/token', got '%s'", r.URL.Path)
		}

		if r.Method != "POST" {
			t.Errorf("Expected POST method, got '%s'", r.Method)
		}

		// Check form data
		err := r.ParseForm()
		if err != nil {
			t.Fatalf("Failed to parse form: %v", err)
		}

		if r.Form.Get("grant_type") != "client_credentials" {
			t.Errorf("Expected grant_type 'client_credentials', got '%s'", r.Form.Get("grant_type"))
		}

		if r.Form.Get("client_id") != "test-client-id" {
			t.Errorf("Expected client_id 'test-client-id', got '%s'", r.Form.Get("client_id"))
		}

		// Return mock token response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"data": {
				"access_token": "test-access-token",
				"token_type": "Bearer",
				"expires_in": 3600,
				"scope": "api:full"
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(&ClientConfig{
		BaseURL: server.URL,
		OAuth: &OAuthConfig{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			AutoRefresh:  true,
		},
	})

	result := client.InitializeOAuth()

	if !result.Success {
		t.Fatalf("InitializeOAuth failed: %s", result.Error)
	}

	// Check if tokens were stored
	tokens := client.GetCurrentTokens()
	if tokens == nil {
		t.Fatal("Tokens should be stored")
	}

	if tokens.AccessToken != "test-access-token" {
		t.Errorf("Expected access token 'test-access-token', got '%s'", tokens.AccessToken)
	}

	if tokens.TokenType != "Bearer" {
		t.Errorf("Expected token type 'Bearer', got '%s'", tokens.TokenType)
	}

	if tokens.ExpiresIn != 3600 {
		t.Errorf("Expected expires_in 3600, got %d", tokens.ExpiresIn)
	}
}

func TestInitializeOAuthWithUserToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			t.Fatalf("Failed to parse form: %v", err)
		}

		if r.Form.Get("user_access_token") != "user-token-123" {
			t.Errorf("Expected user_access_token 'user-token-123', got '%s'", r.Form.Get("user_access_token"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"data": {
				"access_token": "user-access-token",
				"token_type": "Bearer",
				"expires_in": 3600,
				"refresh_token": "user-refresh-token",
				"scope": "api:user"
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(&ClientConfig{
		BaseURL: server.URL,
		OAuth: &OAuthConfig{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
		},
	})

	result := client.InitializeOAuth("user-token-123")

	if !result.Success {
		t.Fatalf("InitializeOAuth with user token failed: %s", result.Error)
	}

	tokens := client.GetCurrentTokens()
	if tokens.RefreshToken != "user-refresh-token" {
		t.Errorf("Expected refresh token 'user-refresh-token', got '%s'", tokens.RefreshToken)
	}
}

func TestRefreshOAuthToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			t.Fatalf("Failed to parse form: %v", err)
		}

		if r.Form.Get("grant_type") != "refresh_token" {
			t.Errorf("Expected grant_type 'refresh_token', got '%s'", r.Form.Get("grant_type"))
		}

		if r.Form.Get("refresh_token") != "test-refresh-token" {
			t.Errorf("Expected refresh_token 'test-refresh-token', got '%s'", r.Form.Get("refresh_token"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"data": {
				"access_token": "new-access-token",
				"token_type": "Bearer",
				"expires_in": 3600,
				"refresh_token": "new-refresh-token",
				"scope": "api:user"
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(&ClientConfig{
		BaseURL: server.URL,
		OAuth: &OAuthConfig{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			RefreshToken: "test-refresh-token",
		},
	})

	result := client.RefreshOAuthToken()

	if !result.Success {
		t.Fatalf("RefreshOAuthToken failed: %s", result.Error)
	}

	tokens := client.GetCurrentTokens()
	if tokens.AccessToken != "new-access-token" {
		t.Errorf("Expected access token 'new-access-token', got '%s'", tokens.AccessToken)
	}

	if tokens.RefreshToken != "new-refresh-token" {
		t.Errorf("Expected refresh token 'new-refresh-token', got '%s'", tokens.RefreshToken)
	}
}

func TestExchangeAuthorizationCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			t.Fatalf("Failed to parse form: %v", err)
		}

		if r.Form.Get("grant_type") != "authorization_code" {
			t.Errorf("Expected grant_type 'authorization_code', got '%s'", r.Form.Get("grant_type"))
		}

		if r.Form.Get("code") != "auth-code-123" {
			t.Errorf("Expected code 'auth-code-123', got '%s'", r.Form.Get("code"))
		}

		if r.Form.Get("redirect_uri") != "https://app.com/callback" {
			t.Errorf("Expected redirect_uri 'https://app.com/callback', got '%s'", r.Form.Get("redirect_uri"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"data": {
				"access_token": "auth-access-token",
				"token_type": "Bearer",
				"expires_in": 3600,
				"refresh_token": "auth-refresh-token",
				"scope": "api:user"
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(&ClientConfig{
		BaseURL: server.URL,
		OAuth: &OAuthConfig{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
		},
	})

	result := client.ExchangeAuthorizationCode("auth-code-123", "https://app.com/callback")

	if !result.Success {
		t.Fatalf("ExchangeAuthorizationCode failed: %s", result.Error)
	}

	tokens := client.GetCurrentTokens()
	if tokens.AccessToken != "auth-access-token" {
		t.Errorf("Expected access token 'auth-access-token', got '%s'", tokens.AccessToken)
	}
}

func TestEnsureValidToken(t *testing.T) {
	client := NewClient(&ClientConfig{
		OAuth: &OAuthConfig{
			ClientID:    "test-client-id",
			AutoRefresh: true,
		},
	})

	// Test with no token expiration set
	err := client.EnsureValidToken()
	if err != nil {
		t.Errorf("EnsureValidToken should not fail when no expiration is set: %v", err)
	}

	// Test with token that expires soon
	expiresAt := time.Now().Add(2 * time.Minute) // Expires in 2 minutes
	client.tokenExpiresAt = &expiresAt
	client.currentTokens = &OAuthTokens{
		AccessToken: "test-token",
	}

	// This should try to refresh (but will fail without a server)
	err = client.EnsureValidToken()
	if err == nil {
		t.Error("EnsureValidToken should fail when refresh is needed but no refresh token is available")
	}

	// Test with token that doesn't need refresh
	expiresAt = time.Now().Add(10 * time.Minute) // Expires in 10 minutes
	client.tokenExpiresAt = &expiresAt

	err = client.EnsureValidToken()
	if err != nil {
		t.Errorf("EnsureValidToken should not fail when token is still valid: %v", err)
	}
}

func TestOAuthConfigMethods(t *testing.T) {
	client := NewClient(&ClientConfig{})

	// Test setting OAuth config
	config := &OAuthConfig{
		ClientID:     "new-client-id",
		ClientSecret: "new-client-secret",
		AutoRefresh:  true,
	}

	client.SetOAuthConfig(config)

	retrievedConfig := client.GetOAuthConfig()
	if retrievedConfig == nil {
		t.Fatal("OAuth config should be set")
	}

	if retrievedConfig.ClientID != "new-client-id" {
		t.Errorf("Expected client ID 'new-client-id', got '%s'", retrievedConfig.ClientID)
	}
}

func TestErrorHandling(t *testing.T) {
	// Test client without OAuth config
	client := NewClient(&ClientConfig{})

	result := client.RefreshOAuthToken()
	if result.Success {
		t.Error("RefreshOAuthToken should fail without OAuth config")
	}

	if !strings.Contains(result.Error, "OAuth configuration not provided") {
		t.Errorf("Expected error about OAuth configuration, got: %s", result.Error)
	}

	// Test refresh without refresh token
	client.SetOAuthConfig(&OAuthConfig{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
	})

	result = client.RefreshOAuthToken()
	if result.Success {
		t.Error("RefreshOAuthToken should fail without refresh token")
	}

	if !strings.Contains(result.Error, "No refresh token available") {
		t.Errorf("Expected error about refresh token, got: %s", result.Error)
	}
}

func TestGetCurrentUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/open/api/v1/me" {
			t.Errorf("Expected path '/open/api/v1/me', got '%s'", r.URL.Path)
		}

		if r.Method != "GET" {
			t.Errorf("Expected GET method, got '%s'", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"data": {
				"user_id": 123,
				"username": "testuser",
				"email": "test@example.com",
				"tenant_id": 456,
				"tenant_name": "Test Tenant",
				"is_admin": true,
				"scope": ["api:user"]
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(&ClientConfig{
		BaseURL: server.URL,
	})

	result := client.GetCurrentUser()

	if !result.Success {
		t.Fatalf("GetCurrentUser failed: %s", result.Error)
	}

	if result.Data == nil {
		t.Fatal("GetCurrentUser should return user data")
	}
}

