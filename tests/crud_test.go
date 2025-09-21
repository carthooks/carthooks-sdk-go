package tests

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/carthooks/carthooks-sdk-go/carthooks"
)

var (
	testClient       *carthooks.Client
	testAppID        uint
	testCollectionID uint
	testItemID       uint
)

func TestMain(m *testing.M) {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	// Setup test client
	setupTestClient()

	// Run tests
	code := m.Run()

	// Cleanup
	cleanup()

	os.Exit(code)
}

func setupTestClient() {
	clientID := os.Getenv("CARTHOOKS_CLIENT_ID")
	clientSecret := os.Getenv("CARTHOOKS_CLIENT_SECRET")
	apiURL := os.Getenv("CARTHOOKS_API_URL")

	// Remove quotes if present
	if len(clientID) > 2 && clientID[0] == '"' && clientID[len(clientID)-1] == '"' {
		clientID = clientID[1 : len(clientID)-1]
	}
	if len(clientSecret) > 2 && clientSecret[0] == '"' && clientSecret[len(clientSecret)-1] == '"' {
		clientSecret = clientSecret[1 : len(clientSecret)-1]
	}

	// Parse test IDs
	appIDStr := os.Getenv("TEST_APP_ID")
	if appIDStr == "" {
		appIDStr = "3883548539"
	}
	appID, _ := strconv.ParseUint(appIDStr, 10, 32)
	testAppID = uint(appID)

	collectionIDStr := os.Getenv("TEST_COLLECTION_ID")
	if collectionIDStr == "" {
		collectionIDStr = "3911747783"
	}
	collectionID, _ := strconv.ParseUint(collectionIDStr, 10, 32)
	testCollectionID = uint(collectionID)

	// Create OAuth config
	oauthConfig := &carthooks.OAuthConfig{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		AutoRefresh:  true,
	}

	// Create client
	config := &carthooks.ClientConfig{
		BaseURL: apiURL,
		OAuth:   oauthConfig,
		Debug:   true,
	}

	testClient = carthooks.NewClient(config)

	// Initialize OAuth
	result := testClient.InitializeOAuth()
	if !result.Success {
		log.Fatalf("Failed to initialize OAuth: %s", result.Error)
	}

	log.Printf("✅ Test client initialized successfully")
	log.Printf("📋 App ID: %d, Collection ID: %d", testAppID, testCollectionID)
}

func cleanup() {
	// Clean up test data if needed
	if testItemID != 0 {
		log.Printf("🧹 Cleaning up test item: %d", testItemID)
		result := testClient.DeleteItem(testAppID, testCollectionID, testItemID)
		if result.Success {
			log.Printf("✅ Test item deleted successfully")
		} else {
			log.Printf("⚠️ Failed to delete test item: %s", result.Error)
		}
	}
}

func TestOAuthInitialization(t *testing.T) {
	t.Log("🧪 Testing OAuth initialization...")

	// OAuth should already be initialized in setupTestClient
	tokens := testClient.GetCurrentTokens()
	require.NotNil(t, tokens, "OAuth tokens should not be nil")
	assert.NotEmpty(t, tokens.AccessToken, "Access token should not be empty")
	assert.Equal(t, "Bearer", tokens.TokenType, "Token type should be Bearer")
	assert.Greater(t, tokens.ExpiresIn, 0, "Expires in should be greater than 0")

	t.Logf("✅ OAuth tokens: %s..., expires in: %d seconds",
		tokens.AccessToken[:30], tokens.ExpiresIn)
}

func TestGetCurrentUser(t *testing.T) {
	t.Log("🧪 Testing get current user...")

	result := testClient.GetCurrentUser()
	require.True(t, result.Success, "Get current user should succeed: %s", result.Error)

	// Parse user data
	userData, ok := result.Data.(map[string]interface{})
	require.True(t, ok, "User data should be a map")

	userID, exists := userData["user_id"]
	assert.True(t, exists, "User ID should exist")

	username, exists := userData["username"]
	assert.True(t, exists, "Username should exist")

	tenantID, exists := userData["tenant_id"]
	assert.True(t, exists, "Tenant ID should exist")

	t.Logf("✅ Current user: ID=%v, Username=%v, TenantID=%v", userID, username, tenantID)
}

func TestGetItems(t *testing.T) {
	t.Log("🧪 Testing get items...")

	result := testClient.GetItems(testAppID, testCollectionID, 5, 0, nil)
	require.True(t, result.Success, "Get items should succeed: %s", result.Error)

	// Parse items data
	itemsData, ok := result.Data.(map[string]interface{})
	require.True(t, ok, "Items data should be a map")

	items, exists := itemsData["items"]
	assert.True(t, exists, "Items array should exist")

	t.Logf("✅ Retrieved items successfully, data keys: %v", getMapKeys(itemsData))
}

func TestCreateItem(t *testing.T) {
	t.Log("🧪 Testing create item...")

	// Create test data
	testData := map[string]interface{}{
		"title":       fmt.Sprintf("Test Item - %d", time.Now().Unix()),
		"f_1009":      1, // This matches our filter condition
		"description": "Created by Go SDK test",
		"status":      "active",
	}

	result := testClient.CreateItem(testAppID, testCollectionID, testData)
	require.True(t, result.Success, "Create item should succeed: %s", result.Error)

	// Parse created item data
	itemData, ok := result.Data.(map[string]interface{})
	require.True(t, ok, "Item data should be a map")

	// Extract item ID for cleanup
	if id, exists := itemData["id"]; exists {
		if idFloat, ok := id.(float64); ok {
			testItemID = uint(idFloat)
		} else if idStr, ok := id.(string); ok {
			if parsed, err := strconv.ParseUint(idStr, 10, 32); err == nil {
				testItemID = uint(parsed)
			}
		}
	}

	t.Logf("✅ Item created successfully: %d", testItemID)
	t.Logf("📋 Item data keys: %v", getMapKeys(itemData))
}

func TestGetItemByID(t *testing.T) {
	// First create an item if we don't have one
	if testItemID == 0 {
		TestCreateItem(t)
	}

	t.Log("🧪 Testing get item by ID...")

	result := testClient.GetItemByID(testAppID, testCollectionID, testItemID, nil)
	require.True(t, result.Success, "Get item by ID should succeed: %s", result.Error)

	// Parse item data
	itemData, ok := result.Data.(map[string]interface{})
	require.True(t, ok, "Item data should be a map")

	// Verify item ID matches
	if id, exists := itemData["id"]; exists {
		assert.Equal(t, testItemID, id, "Item ID should match")
	}

	t.Logf("✅ Retrieved item by ID: %s", testItemID)
}

func TestUpdateItem(t *testing.T) {
	// First create an item if we don't have one
	if testItemID == 0 {
		TestCreateItem(t)
	}

	t.Log("🧪 Testing update item...")

	// Update data
	updateData := map[string]interface{}{
		"title":       fmt.Sprintf("Updated Test Item - %d", time.Now().Unix()),
		"description": "Updated by Go SDK test",
		"status":      "updated",
	}

	result := testClient.UpdateItem(testAppID, testCollectionID, testItemID, updateData)
	require.True(t, result.Success, "Update item should succeed: %s", result.Error)

	t.Logf("✅ Item updated successfully: %s", testItemID)
}

func TestLockAndUnlockItem(t *testing.T) {
	// First create an item if we don't have one
	if testItemID == 0 {
		TestCreateItem(t)
	}

	t.Log("🧪 Testing lock and unlock item...")

	// Lock item
	lockOptions := &carthooks.LockOptions{
		LockTimeout: 600,
		LockID:      "test-lock-id",
		Subject:     "Go SDK Test",
	}

	lockResult := testClient.LockItem(testAppID, testCollectionID, testItemID, lockOptions)
	require.True(t, lockResult.Success, "Lock item should succeed: %s", lockResult.Error)

	t.Logf("✅ Item locked successfully: %s", testItemID)

	// Unlock item
	unlockResult := testClient.UnlockItem(testAppID, testCollectionID, testItemID, "test-lock-id")
	require.True(t, unlockResult.Success, "Unlock item should succeed: %s", unlockResult.Error)

	t.Logf("✅ Item unlocked successfully: %s", testItemID)
}

func TestDeleteItem(t *testing.T) {
	// First create an item if we don't have one
	if testItemID == 0 {
		TestCreateItem(t)
	}

	t.Log("🧪 Testing delete item...")

	result := testClient.DeleteItem(testAppID, testCollectionID, testItemID)
	require.True(t, result.Success, "Delete item should succeed: %s", result.Error)

	t.Logf("✅ Item deleted successfully: %d", testItemID)

	// Clear testItemID so cleanup doesn't try to delete it again
	testItemID = 0
}

// Helper function to get map keys
func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
