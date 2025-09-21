# Carthooks Go SDK OAuth Support

The Carthooks Go SDK now supports OAuth 2.0 authentication with automatic token refresh capabilities.

## Features

- ✅ **Client Credentials Flow** - Machine-to-machine authentication
- ✅ **Client Credentials + User Token** - Server-side apps acting on behalf of users  
- ✅ **Authorization Code Flow** - Standard OAuth 2.0 for web applications
- ✅ **Automatic Token Refresh** - Seamless token renewal before expiration
- ✅ **Token Management** - Store and retrieve tokens
- ✅ **Type Safety** - Full Go type safety for all OAuth operations

## Installation

```bash
go get github.com/carthooks/carthooks-sdk-go
```

## Supported Grant Types

### 1. Client Credentials (Machine-to-Machine)

For server-to-server communication without user context:

```go
package main

import (
    "fmt"
    "log"
    "github.com/carthooks/carthooks-sdk-go/carthooks"
)

func main() {
    client := carthooks.NewClient(&carthooks.ClientConfig{
        BaseURL: "https://your-instance.carthooks.com",
        OAuth: &carthooks.OAuthConfig{
            ClientID:     "dvc-your-client-id",
            ClientSecret: "dvs-your-client-secret",
            AutoRefresh:  true,
        },
    })

    // Initialize OAuth
    result := client.InitializeOAuth()
    if !result.Success {
        log.Fatalf("Failed to get token: %s", result.Error)
    }

    fmt.Printf("Access token: %v\n", result.Data)
    
    // Make API calls - tokens auto-refresh as needed
    items := client.GetItems(appID, collectionID, 10, 0, nil)
    if items.Success {
        fmt.Printf("Items: %v\n", items.Data)
    }
}
```

### 2. Client Credentials + User Token

For server-side applications acting on behalf of users:

```go
client := carthooks.NewClient(&carthooks.ClientConfig{
    OAuth: &carthooks.OAuthConfig{
        ClientID:     "dvc-your-client-id", 
        ClientSecret: "dvs-your-client-secret",
        AutoRefresh:  true,
    },
})

// Use user access token from your frontend
userToken := "user-access-token-from-frontend"
result := client.InitializeOAuth(userToken)

if result.Success {
    // This token represents the user and has their permissions
    userInfo := client.GetCurrentUser()
    if userInfo.Success {
        fmt.Printf("Acting as user: %v\n", userInfo.Data)
    }
}
```

### 3. Authorization Code Flow

For web applications with user authorization:

```go
client := carthooks.NewClient(&carthooks.ClientConfig{
    OAuth: &carthooks.OAuthConfig{
        ClientID:     "dvc-your-client-id",
        ClientSecret: "dvs-your-client-secret",
        AutoRefresh:  true,
    },
})

// Step 1: Get authorization URL (redirect user here)
authRequest := &carthooks.OAuthAuthorizeCodeRequest{
    ClientID:       "dvc-your-client-id",
    RedirectURI:    "https://your-app.com/oauth/callback",
    State:          "random-state-string",
    TargetTenantID: 123, // For platform-level clients
}

authResult := client.GetOAuthAuthorizeCode(authRequest)
if authResult.Success {
    // Redirect user to the authorization URL
    if data, ok := authResult.Data.(map[string]interface{}); ok {
        if redirectURL, ok := data["redirect_url"].(string); ok {
            fmt.Printf("Redirect user to: %s\n", redirectURL)
        }
    }
}

// Step 2: Exchange code for tokens (in your callback handler)
code := "authorization-code-from-callback"
tokenResult := client.ExchangeAuthorizationCode(code, "https://your-app.com/oauth/callback")

if tokenResult.Success {
    // Store tokens and make API calls
    tokens := client.GetCurrentTokens()
    fmt.Printf("Access token: %s\n", tokens.AccessToken)
    fmt.Printf("Refresh token: %s\n", tokens.RefreshToken)
}
```

## Token Management

### Automatic Refresh

Tokens are automatically refreshed 5 minutes before expiration:

```go
client := carthooks.NewClient(&carthooks.ClientConfig{
    OAuth: &carthooks.OAuthConfig{
        ClientID:     "dvc-your-client-id",
        ClientSecret: "dvs-your-client-secret",
        AutoRefresh:  true, // Default: false
    },
})

// All API calls will automatically refresh tokens as needed
result := client.QueryItems(appID, entityID, &carthooks.QueryOptions{
    Pagination: &carthooks.PaginationOptions{
        Page:     1,
        PageSize: 20,
    },
})
```

### Manual Refresh

You can also manually refresh tokens:

```go
// Refresh using stored refresh token
refreshResult := client.RefreshOAuthToken("stored-refresh-token")

if refreshResult.Success {
    fmt.Printf("New access token: %v\n", refreshResult.Data)
}

// Or refresh using configured refresh token
result := client.RefreshOAuthToken()
```

### Token Storage

Store and retrieve tokens for persistence:

```go
// Load refresh token from storage on startup
storedRefreshToken := loadFromDatabase() // Your implementation

client := carthooks.NewClient(&carthooks.ClientConfig{
    OAuth: &carthooks.OAuthConfig{
        ClientID:     "dvc-your-client-id",
        ClientSecret: "dvs-your-client-secret",
        RefreshToken: storedRefreshToken,
        AutoRefresh:  true,
    },
})

// Try to refresh on startup
if storedRefreshToken != "" {
    client.RefreshOAuthToken()
}

// Get current tokens for storage
currentTokens := client.GetCurrentTokens()
if currentTokens != nil {
    saveToDatabase(currentTokens.RefreshToken) // Your implementation
}
```

## Configuration Options

### OAuthConfig Struct

```go
type OAuthConfig struct {
    ClientID     string `json:"client_id"`      // Your OAuth client ID
    ClientSecret string `json:"client_secret"`  // Your OAuth client secret  
    RefreshToken string `json:"refresh_token"`  // Stored refresh token
    AutoRefresh  bool   `json:"auto_refresh"`   // Auto-refresh tokens
}
```

### ClientConfig with OAuth

```go
type ClientConfig struct {
    BaseURL     string                 // API base URL
    AccessToken string                 // Direct access token (alternative to OAuth)
    Timeout     time.Duration          // Request timeout
    Headers     map[string]string      // Custom headers
    Debug       bool                   // Enable debug logging
    OAuth       *OAuthConfig           // OAuth configuration
}
```

## Error Handling

The SDK provides detailed error information:

```go
result := client.InitializeOAuth()

if !result.Success {
    fmt.Printf("OAuth failed: %s\n", result.Error)
    if result.TraceID != "" {
        fmt.Printf("Trace ID: %s\n", result.TraceID)
    }
}

// Handle specific OAuth errors
refreshResult := client.RefreshOAuthToken()
if !refreshResult.Success {
    if strings.Contains(refreshResult.Error, "invalid_refresh_token") {
        // Refresh token expired, need to re-authorize
        redirectToAuthorizationFlow()
    }
}
```

## Token Expiration

- **Access Tokens**: 24 hours
- **Refresh Tokens**: 12 months
- **Auto-refresh**: Triggered 5 minutes before expiration

## Security Best Practices

1. **Store secrets securely**: Never expose client secrets in logs or client-side code
2. **Use HTTPS**: Always use secure connections for token exchange
3. **Rotate tokens**: Implement proper token rotation and storage
4. **Handle expiration**: Gracefully handle token expiration scenarios
5. **Scope limitation**: Request only necessary scopes for your application

## API Endpoints

The SDK automatically handles these OAuth endpoints:

- `POST /open/api/oauth/token` - Token exchange and refresh
- `POST /api/oauth/get-authorize-code` - Get authorization code  
- `GET /open/api/v1/me` - Get current user information

## Complete Example

```go
package main

import (
    "fmt"
    "log"
    "time"
    "github.com/carthooks/carthooks-sdk-go/carthooks"
)

func main() {
    // Create client with OAuth configuration
    client := carthooks.NewClient(&carthooks.ClientConfig{
        BaseURL: "https://your-instance.carthooks.com",
        Timeout: 30 * time.Second,
        OAuth: &carthooks.OAuthConfig{
            ClientID:     "dvc-your-client-id",
            ClientSecret: "dvs-your-client-secret",
            AutoRefresh:  true,
        },
        Debug: true,
    })

    // Initialize OAuth
    result := client.InitializeOAuth()
    if !result.Success {
        log.Fatalf("OAuth initialization failed: %s", result.Error)
    }

    fmt.Println("OAuth initialized successfully")

    // Get current user
    userResult := client.GetCurrentUser()
    if userResult.Success {
        fmt.Printf("Current user: %v\n", userResult.Data)
    }

    // Query items with auto token refresh
    queryOptions := &carthooks.QueryOptions{
        Pagination: &carthooks.PaginationOptions{
            Page:      1,
            PageSize:  10,
            WithCount: true,
        },
        Filters: map[string]interface{}{
            "status": "active",
        },
    }

    itemsResult := client.QueryItems(123, 456, queryOptions)
    if itemsResult.Success {
        fmt.Printf("Items: %v\n", itemsResult.Data)
    }

    // Create new item
    newItem := map[string]interface{}{
        "title": "New Item",
        "status": "active",
    }

    createResult := client.CreateItem(123, 456, newItem)
    if createResult.Success {
        fmt.Printf("Created item: %v\n", createResult.Data)
    }
}
```

## Migration from Access Token

If you're currently using direct access tokens:

```go
// Old way
client := carthooks.NewClient(&carthooks.ClientConfig{
    AccessToken: "your-access-token",
})

// New way with OAuth
client := carthooks.NewClient(&carthooks.ClientConfig{
    OAuth: &carthooks.OAuthConfig{
        ClientID:     "dvc-your-client-id",
        ClientSecret: "dvs-your-client-secret",
        AutoRefresh:  true,
    },
})

client.InitializeOAuth()
```

The OAuth approach provides better security and automatic token management.

## Examples

See `examples/oauth_example.go` for complete working examples of all OAuth flows.

## Testing

```bash
# Run tests
go test ./...

# Run with verbose output
go test -v ./...
```
