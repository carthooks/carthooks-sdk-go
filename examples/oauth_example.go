package main

import (
	"fmt"
	"log"
	"time"

	"github.com/carthooks/carthooks-sdk-go/carthooks"
)

// Example 1: Client Credentials Mode (Machine-to-Machine)
func clientCredentialsExample() {
	fmt.Println("=== Client Credentials Example ===")

	client := carthooks.NewClient(&carthooks.ClientConfig{
		BaseURL: "https://your-carthooks-instance.com",
		OAuth: &carthooks.OAuthConfig{
			ClientID:     "dvc-your-client-id",
			ClientSecret: "dvs-your-client-secret",
			AutoRefresh:  true,
		},
		Debug: true,
	})

	// Initialize OAuth with client credentials
	result := client.InitializeOAuth()
	if !result.Success {
		log.Fatalf("Failed to get token: %s", result.Error)
	}

	fmt.Printf("Got access token: %v\n", result.Data)

	// Now you can make API calls - tokens will be automatically refreshed
	itemsResult := client.GetItems(123, 456, 10, 0, nil)
	if itemsResult.Success {
		fmt.Printf("Items: %v\n", itemsResult.Data)
	} else {
		fmt.Printf("Failed to get items: %s\n", itemsResult.Error)
	}

	// Get current user info
	userResult := client.GetCurrentUser()
	if userResult.Success {
		fmt.Printf("Current user: %v\n", userResult.Data)
	}
}

// Example 2: Client Credentials with User Token Mode
func clientCredentialsWithUserExample() {
	fmt.Println("=== Client Credentials with User Token Example ===")

	client := carthooks.NewClient(&carthooks.ClientConfig{
		BaseURL: "https://your-carthooks-instance.com",
		OAuth: &carthooks.OAuthConfig{
			ClientID:     "dvc-your-client-id",
			ClientSecret: "dvs-your-client-secret",
			AutoRefresh:  true,
		},
	})

	// Initialize OAuth with user access token
	userAccessToken := "user-access-token-from-frontend"
	result := client.InitializeOAuth(userAccessToken)

	if result.Success {
		fmt.Printf("Got user-context token: %v\n", result.Data)

		// This token represents the user and has their permissions
		userInfo := client.GetCurrentUser()
		if userInfo.Success {
			fmt.Printf("User info: %v\n", userInfo.Data)
		}
	} else {
		fmt.Printf("Failed to get token: %s\n", result.Error)
	}
}

// Example 3: Authorization Code Flow
func authorizationCodeExample() {
	fmt.Println("=== Authorization Code Flow Example ===")

	client := carthooks.NewClient(&carthooks.ClientConfig{
		BaseURL: "https://your-carthooks-instance.com",
		OAuth: &carthooks.OAuthConfig{
			ClientID:     "dvc-your-client-id",
			ClientSecret: "dvs-your-client-secret",
			AutoRefresh:  true,
		},
	})

	// Step 1: Get authorization URL (this would be done in your web app)
	authRequest := &carthooks.OAuthAuthorizeCodeRequest{
		ClientID:       "dvc-your-client-id",
		RedirectURI:    "https://your-app.com/callback",
		State:          "random-state-string",
		TargetTenantID: 456, // For platform-level clients
	}

	authResult := client.GetOAuthAuthorizeCode(authRequest)
	if authResult.Success {
		if data, ok := authResult.Data.(map[string]interface{}); ok {
			if redirectURL, ok := data["redirect_url"].(string); ok {
				fmt.Printf("Redirect user to: %s\n", redirectURL)
			}
		}
	}

	// Step 2: Exchange authorization code for tokens (after user authorizes)
	authCode := "authorization-code-from-callback"
	tokenResult := client.ExchangeAuthorizationCode(authCode, "https://your-app.com/callback")

	if tokenResult.Success {
		fmt.Printf("Got tokens: %v\n", tokenResult.Data)

		// Now you can make API calls on behalf of the user
		userInfo := client.GetCurrentUser()
		if userInfo.Success {
			fmt.Printf("Authorized user: %v\n", userInfo.Data)
		}
	} else {
		fmt.Printf("Failed to exchange code: %s\n", tokenResult.Error)
	}
}

// Example 4: Manual Token Refresh
func manualRefreshExample() {
	fmt.Println("=== Manual Token Refresh Example ===")

	client := carthooks.NewClient(&carthooks.ClientConfig{
		OAuth: &carthooks.OAuthConfig{
			ClientID:     "dvc-your-client-id",
			ClientSecret: "dvs-your-client-secret",
			RefreshToken: "your-stored-refresh-token",
			AutoRefresh:  false, // Disable auto refresh
		},
	})

	// Manually refresh token
	refreshResult := client.RefreshOAuthToken()
	if refreshResult.Success {
		fmt.Printf("Refreshed token: %v\n", refreshResult.Data)

		// Get current tokens
		tokens := client.GetCurrentTokens()
		if tokens != nil {
			fmt.Printf("Current access token: %s\n", tokens.AccessToken)
			if tokens.RefreshToken != "" {
				fmt.Printf("New refresh token: %s\n", tokens.RefreshToken)
				// Store refresh_token for next time
			}
		}
	} else {
		fmt.Printf("Refresh failed: %s\n", refreshResult.Error)
	}
}

// Example 5: Token Management with Callbacks
func tokenManagementExample() {
	fmt.Println("=== Token Management Example ===")

	// Simulate loading refresh token from storage
	storedRefreshToken := loadRefreshTokenFromStorage()

	client := carthooks.NewClient(&carthooks.ClientConfig{
		OAuth: &carthooks.OAuthConfig{
			ClientID:     "dvc-your-client-id",
			ClientSecret: "dvs-your-client-secret",
			RefreshToken: storedRefreshToken,
			AutoRefresh:  true,
		},
	})

	// Try to refresh token on startup if we have one
	if storedRefreshToken != "" {
		refreshResult := client.RefreshOAuthToken()
		if refreshResult.Success {
			fmt.Println("Token refreshed on startup")
			saveTokensToStorage(client.GetCurrentTokens())
		}
	} else {
		// Initialize with client credentials if no refresh token
		initResult := client.InitializeOAuth()
		if initResult.Success {
			fmt.Println("Initialized with client credentials")
			saveTokensToStorage(client.GetCurrentTokens())
		}
	}

	// Check current token status
	currentTokens := client.GetCurrentTokens()
	if currentTokens != nil {
		fmt.Printf("Current access token: %s\n", currentTokens.AccessToken)
		fmt.Printf("Token scope: %s\n", currentTokens.Scope)
	}
}

// Example 6: Error Handling
func errorHandlingExample() {
	fmt.Println("=== Error Handling Example ===")

	// Client without OAuth config
	clientWithoutOAuth := carthooks.NewClient(&carthooks.ClientConfig{
		BaseURL: "https://test.carthooks.com",
	})

	result := clientWithoutOAuth.RefreshOAuthToken()
	if !result.Success {
		fmt.Printf("Expected error: %s\n", result.Error)
	}

	// Client with invalid credentials
	clientWithBadCreds := carthooks.NewClient(&carthooks.ClientConfig{
		OAuth: &carthooks.OAuthConfig{
			ClientID:     "invalid-client-id",
			ClientSecret: "invalid-secret",
		},
	})

	badResult := clientWithBadCreds.InitializeOAuth()
	if !badResult.Success {
		fmt.Printf("Authentication failed: %s\n", badResult.Error)
	}
}

// Helper functions for token storage (implement based on your needs)
func loadRefreshTokenFromStorage() string {
	// In a real application, load from database, file, etc.
	return ""
}

func saveTokensToStorage(tokens *carthooks.OAuthTokens) {
	if tokens == nil {
		return
	}
	// In a real application, save to database, file, etc.
	fmt.Printf("Saving tokens to storage: %s\n", tokens.AccessToken)
}

// Example 7: Advanced Usage with Custom Configuration
func advancedUsageExample() {
	fmt.Println("=== Advanced Usage Example ===")

	client := carthooks.NewClient(&carthooks.ClientConfig{
		BaseURL: "https://your-carthooks-instance.com",
		Timeout: 60 * time.Second,
		Headers: map[string]string{
			"X-Custom-Header": "custom-value",
		},
		OAuth: &carthooks.OAuthConfig{
			ClientID:     "dvc-your-client-id",
			ClientSecret: "dvs-your-client-secret",
			AutoRefresh:  true,
		},
		Debug: true,
	})

	// Initialize OAuth
	if result := client.InitializeOAuth(); result.Success {
		fmt.Println("OAuth initialized successfully")

		// Make API calls with automatic token refresh
		queryOptions := &carthooks.QueryOptions{
			Pagination: &carthooks.PaginationOptions{
				Page:      1,
				PageSize:  20,
				WithCount: true,
			},
			Filters: map[string]interface{}{
				"status": "active",
			},
			Sort: []string{"created_at:desc"},
		}

		result := client.QueryItems(123, 456, queryOptions)
		if result.Success {
			fmt.Printf("Query result: %v\n", result.Data)
		}

		// Update OAuth config dynamically
		client.SetOAuthConfig(&carthooks.OAuthConfig{
			ClientID:     "dvc-new-client-id",
			ClientSecret: "dvs-new-secret",
			AutoRefresh:  true,
		})

		// Get current config
		config := client.GetOAuthConfig()
		if config != nil {
			fmt.Printf("Current OAuth config: %+v\n", config)
		}
	}
}

func main() {
	fmt.Println("Carthooks Go SDK OAuth Examples")
	fmt.Println("================================")

	// Run examples (uncomment the ones you want to test)

	// clientCredentialsExample()
	// clientCredentialsWithUserExample()
	// authorizationCodeExample()
	// manualRefreshExample()
	// tokenManagementExample()
	// errorHandlingExample()
	// advancedUsageExample()

	fmt.Println("\nNote: Update the client credentials and URLs before running these examples")
}
