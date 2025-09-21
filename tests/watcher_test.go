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
	watcherClient       *carthooks.Client
	watcherAppID        uint
	watcherCollectionID uint
	sqsQueueURL         string
)

func setupWatcherTest() {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	clientID := os.Getenv("CARTHOOKS_CLIENT_ID")
	clientSecret := os.Getenv("CARTHOOKS_CLIENT_SECRET")
	apiURL := os.Getenv("CARTHOOKS_API_URL")
	sqsQueueURL = os.Getenv("SQS_QUEUE_URL")

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
	watcherAppID = uint(appID)

	collectionIDStr := os.Getenv("TEST_COLLECTION_ID")
	if collectionIDStr == "" {
		collectionIDStr = "3911747783"
	}
	collectionID, _ := strconv.ParseUint(collectionIDStr, 10, 32)
	watcherCollectionID = uint(collectionID)

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

	watcherClient = carthooks.NewClient(config)

	// Initialize OAuth
	result := watcherClient.InitializeOAuth()
	if !result.Success {
		log.Fatalf("Failed to initialize OAuth for watcher test: %s", result.Error)
	}

	log.Printf("✅ Watcher test client initialized successfully")
}

func TestStartWatchData(t *testing.T) {
	setupWatcherTest()

	t.Log("🧪 Testing start watch data...")

	if sqsQueueURL == "" {
		t.Skip("⚠️ SQS_QUEUE_URL not set, skipping watcher test")
		return
	}

	// Create watch data options
	watchName := fmt.Sprintf("test-watch-go-%d", time.Now().Unix())

	options := &carthooks.WatchDataOptions{
		EndpointURL:  sqsQueueURL,
		EndpointType: "sqs",
		Name:         watchName,
		AppID:        watcherAppID,
		CollectionID: watcherCollectionID,
		Filters: map[string]interface{}{
			"f_1009": map[string]interface{}{
				"$eq": 1,
			},
		},
		Age:            432000, // 5 days
		WatchStartTime: 0,
	}

	t.Logf("📋 Watch configuration:")
	t.Logf("   Name: %s", watchName)
	t.Logf("   App ID: %d", watcherAppID)
	t.Logf("   Collection ID: %d", watcherCollectionID)
	t.Logf("   SQS Queue: %s", sqsQueueURL)
	t.Logf("   Filters: %+v", options.Filters)

	result := watcherClient.StartWatchData(options)
	require.True(t, result.Success, "Start watch data should succeed: %s", result.Error)

	// Parse response data
	responseData, ok := result.Data.(map[string]interface{})
	require.True(t, ok, "Response data should be a map")

	t.Logf("✅ Watch data started successfully")
	t.Logf("📋 Response data: %+v", responseData)

	// Verify response contains expected fields
	assert.Contains(t, responseData, "endpoint_url", "Response should contain endpoint_url")
	assert.Contains(t, responseData, "endpoint_type", "Response should contain endpoint_type")
	assert.Contains(t, responseData, "name", "Response should contain name")
	assert.Contains(t, responseData, "app_id", "Response should contain app_id")
	assert.Contains(t, responseData, "collection_id", "Response should contain collection_id")

	// Verify values
	if endpointURL, exists := responseData["endpoint_url"]; exists {
		assert.Equal(t, sqsQueueURL, endpointURL, "Endpoint URL should match")
	}

	if endpointType, exists := responseData["endpoint_type"]; exists {
		assert.Equal(t, "sqs", endpointType, "Endpoint type should be sqs")
	}

	if name, exists := responseData["name"]; exists {
		assert.Equal(t, watchName, name, "Watch name should match")
	}
}

func TestWatchDataWithDifferentFilters(t *testing.T) {
	setupWatcherTest()

	t.Log("🧪 Testing watch data with different filters...")

	if sqsQueueURL == "" {
		t.Skip("⚠️ SQS_QUEUE_URL not set, skipping watcher test")
		return
	}

	testCases := []struct {
		name    string
		filters map[string]interface{}
	}{
		{
			name: "Simple equality filter",
			filters: map[string]interface{}{
				"status": map[string]interface{}{
					"$eq": "active",
				},
			},
		},
		{
			name: "Multiple filters",
			filters: map[string]interface{}{
				"f_1009": map[string]interface{}{
					"$eq": 1,
				},
				"status": map[string]interface{}{
					"$eq": "active",
				},
			},
		},
		{
			name: "Range filter",
			filters: map[string]interface{}{
				"created_at": map[string]interface{}{
					"$gte": time.Now().AddDate(0, 0, -7).Unix(), // Last 7 days
				},
			},
		},
	}

	for i, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			watchName := fmt.Sprintf("test-watch-filter-%d-%d", i, time.Now().Unix())

			options := &carthooks.WatchDataOptions{
				EndpointURL:    sqsQueueURL,
				EndpointType:   "sqs",
				Name:           watchName,
				AppID:          watcherAppID,
				CollectionID:   watcherCollectionID,
				Filters:        tc.filters,
				Age:            432000,
				WatchStartTime: 0,
			}

			result := watcherClient.StartWatchData(options)
			require.True(t, result.Success, "Start watch data should succeed: %s", result.Error)

			t.Logf("✅ Watch data started with %s", tc.name)
		})
	}
}

func TestWatchDataConfiguration(t *testing.T) {
	setupWatcherTest()

	t.Log("🧪 Testing watch data configuration options...")

	if sqsQueueURL == "" {
		t.Skip("⚠️ SQS_QUEUE_URL not set, skipping watcher test")
		return
	}

	// Test different age values
	ageTestCases := []struct {
		name string
		age  int
	}{
		{"1 hour", 3600},
		{"1 day", 86400},
		{"1 week", 604800},
		{"1 month", 2592000},
	}

	for i, tc := range ageTestCases {
		t.Run(fmt.Sprintf("Age_%s", tc.name), func(t *testing.T) {
			watchName := fmt.Sprintf("test-watch-age-%d-%d", i, time.Now().Unix())

			options := &carthooks.WatchDataOptions{
				EndpointURL:    sqsQueueURL,
				EndpointType:   "sqs",
				Name:           watchName,
				AppID:          watcherAppID,
				CollectionID:   watcherCollectionID,
				Age:            tc.age,
				WatchStartTime: 0,
			}

			result := watcherClient.StartWatchData(options)
			require.True(t, result.Success, "Start watch data should succeed: %s", result.Error)

			// Verify age in response
			responseData, ok := result.Data.(map[string]interface{})
			require.True(t, ok, "Response data should be a map")

			if age, exists := responseData["age"]; exists {
				// Age might be returned as float64 from JSON
				if ageFloat, ok := age.(float64); ok {
					assert.Equal(t, float64(tc.age), ageFloat, "Age should match")
				}
			}

			t.Logf("✅ Watch data configured with age: %s (%d seconds)", tc.name, tc.age)
		})
	}
}

func TestWatchDataErrorHandling(t *testing.T) {
	setupWatcherTest()

	t.Log("🧪 Testing watch data error handling...")

	// Test with invalid endpoint URL
	t.Run("Invalid_Endpoint_URL", func(t *testing.T) {
		options := &carthooks.WatchDataOptions{
			EndpointURL:  "invalid-url",
			EndpointType: "sqs",
			Name:         "test-invalid-endpoint",
			AppID:        watcherAppID,
			CollectionID: watcherCollectionID,
			Age:          3600,
		}

		result := watcherClient.StartWatchData(options)
		// This might succeed or fail depending on server validation
		// Just log the result
		t.Logf("📋 Invalid endpoint URL result: Success=%v, Error=%s", result.Success, result.Error)
	})

	// Test with missing required fields
	t.Run("Missing_Required_Fields", func(t *testing.T) {
		options := &carthooks.WatchDataOptions{
			// Missing EndpointURL
			EndpointType: "sqs",
			Name:         "test-missing-fields",
			AppID:        watcherAppID,
			CollectionID: watcherCollectionID,
		}

		result := watcherClient.StartWatchData(options)
		// This should fail
		assert.False(t, result.Success, "Should fail with missing endpoint URL")
		t.Logf("📋 Missing fields result: Error=%s", result.Error)
	})

	// Test with invalid app/collection IDs
	t.Run("Invalid_IDs", func(t *testing.T) {
		if sqsQueueURL == "" {
			t.Skip("⚠️ SQS_QUEUE_URL not set, skipping test")
			return
		}

		options := &carthooks.WatchDataOptions{
			EndpointURL:  sqsQueueURL,
			EndpointType: "sqs",
			Name:         "test-invalid-ids",
			AppID:        99999999, // Invalid app ID
			CollectionID: 99999999, // Invalid collection ID
			Age:          3600,
		}

		result := watcherClient.StartWatchData(options)
		// This might succeed or fail depending on server validation
		t.Logf("📋 Invalid IDs result: Success=%v, Error=%s", result.Success, result.Error)
	})
}
