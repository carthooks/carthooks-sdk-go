package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"

	"github.com/carthooks/carthooks-sdk-go/carthooks"
)

func main() {
	fmt.Println("🧪 Carthooks Go SDK Manual Test Runner")
	fmt.Println("======================================")

	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	// Setup client
	client := setupClient()
	if client == nil {
		log.Fatal("❌ Failed to setup client")
	}

	// Parse test IDs
	appID := parseUintEnv("TEST_APP_ID", 3883548539)
	collectionID := parseUintEnv("TEST_COLLECTION_ID", 3911747783)

	fmt.Printf("📋 Test Configuration:\n")
	fmt.Printf("   App ID: %d\n", appID)
	fmt.Printf("   Collection ID: %d\n", collectionID)
	fmt.Printf("   API URL: %s\n", os.Getenv("CARTHOOKS_API_URL"))

	// Run tests
	runCRUDTests(client, appID, collectionID)
	runWatcherTests(client, appID, collectionID)

	fmt.Println("\n✅ All manual tests completed!")
}

func setupClient() *carthooks.Client {
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

	client := carthooks.NewClient(config)

	// Initialize OAuth
	result := client.InitializeOAuth()
	if !result.Success {
		log.Printf("❌ Failed to initialize OAuth: %s", result.Error)
		return nil
	}

	fmt.Println("✅ Client initialized successfully")
	return client
}

func runCRUDTests(client *carthooks.Client, appID, collectionID uint) {
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println(" CRUD TESTS")
	fmt.Println(strings.Repeat("=", 50))

	var testItemID uint

	// Test 1: Get current user
	fmt.Println("\n🧪 Test 1: Get Current User")
	userResult := client.GetCurrentUser()
	if userResult.Success {
		fmt.Println("✅ Get current user: SUCCESS")
		if userData, ok := userResult.Data.(map[string]interface{}); ok {
			fmt.Printf("   User ID: %v\n", userData["user_id"])
			fmt.Printf("   Username: %v\n", userData["username"])
			fmt.Printf("   Tenant ID: %v\n", userData["tenant_id"])
		}
	} else {
		fmt.Printf("❌ Get current user: FAILED - %s\n", userResult.Error)
	}

	// Test 2: Get items
	fmt.Println("\n🧪 Test 2: Get Items")
	itemsResult := client.GetItems(appID, collectionID, 5, 0, nil)
	if itemsResult.Success {
		fmt.Println("✅ Get items: SUCCESS")
		if itemsData, ok := itemsResult.Data.(map[string]interface{}); ok {
			fmt.Printf("   Data keys: %v\n", getMapKeys(itemsData))
		}
	} else {
		fmt.Printf("❌ Get items: FAILED - %s\n", itemsResult.Error)
	}

	// Test 3: Create item
	fmt.Println("\n🧪 Test 3: Create Item")
	testData := map[string]interface{}{
		"title":       fmt.Sprintf("Go SDK Test Item - %d", time.Now().Unix()),
		"f_1009":      1,
		"description": "Created by Go SDK manual test",
		"status":      "active",
	}

	createResult := client.CreateItem(appID, collectionID, testData)
	if createResult.Success {
		fmt.Println("✅ Create item: SUCCESS")
		if itemData, ok := createResult.Data.(map[string]interface{}); ok {
			if id, exists := itemData["id"]; exists {
				if idFloat, ok := id.(float64); ok {
					testItemID = uint(idFloat)
					fmt.Printf("   Created item ID: %d\n", testItemID)
				} else if idStr, ok := id.(string); ok {
					if parsed, err := strconv.ParseUint(idStr, 10, 32); err == nil {
						testItemID = uint(parsed)
						fmt.Printf("   Created item ID: %d\n", testItemID)
					}
				}
			}
		}
	} else {
		fmt.Printf("❌ Create item: FAILED - %s\n", createResult.Error)
		return // Can't continue without an item
	}

	// Test 4: Get item by ID
	fmt.Println("\n🧪 Test 4: Get Item by ID")
	getResult := client.GetItemByID(appID, collectionID, testItemID, nil)
	if getResult.Success {
		fmt.Println("✅ Get item by ID: SUCCESS")
	} else {
		fmt.Printf("❌ Get item by ID: FAILED - %s\n", getResult.Error)
	}

	// Test 5: Update item
	fmt.Println("\n🧪 Test 5: Update Item")
	updateData := map[string]interface{}{
		"title":       fmt.Sprintf("Updated Go SDK Test Item - %d", time.Now().Unix()),
		"description": "Updated by Go SDK manual test",
		"status":      "updated",
	}

	updateResult := client.UpdateItem(appID, collectionID, testItemID, updateData)
	if updateResult.Success {
		fmt.Println("✅ Update item: SUCCESS")
	} else {
		fmt.Printf("❌ Update item: FAILED - %s\n", updateResult.Error)
	}

	// Test 6: Delete item
	fmt.Println("\n🧪 Test 6: Delete Item")
	deleteResult := client.DeleteItem(appID, collectionID, testItemID)
	if deleteResult.Success {
		fmt.Println("✅ Delete item: SUCCESS")
	} else {
		fmt.Printf("❌ Delete item: FAILED - %s\n", deleteResult.Error)
	}
}

func runWatcherTests(client *carthooks.Client, appID, collectionID uint) {
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println(" WATCHER TESTS")
	fmt.Println(strings.Repeat("=", 50))

	sqsQueueURL := os.Getenv("SQS_QUEUE_URL")
	if sqsQueueURL == "" {
		fmt.Println("⚠️ SQS_QUEUE_URL not set, skipping watcher tests")
		return
	}

	// Test 1: Start watch data
	fmt.Println("\n🧪 Test 1: Start Watch Data")
	watchName := fmt.Sprintf("go-sdk-test-watch-%d", time.Now().Unix())

	options := &carthooks.WatchDataOptions{
		EndpointURL:  sqsQueueURL,
		EndpointType: "sqs",
		Name:         watchName,
		AppID:        appID,
		CollectionID: collectionID,
		Filters: map[string]interface{}{
			"f_1009": map[string]interface{}{
				"$eq": 1,
			},
		},
		Age:            432000, // 5 days
		WatchStartTime: 0,
	}

	fmt.Printf("📋 Watch Configuration:\n")
	fmt.Printf("   Name: %s\n", watchName)
	fmt.Printf("   SQS Queue: %s\n", sqsQueueURL)
	fmt.Printf("   Filters: %+v\n", options.Filters)

	watchResult := client.StartWatchData(options)
	if watchResult.Success {
		fmt.Println("✅ Start watch data: SUCCESS")
		if responseData, ok := watchResult.Data.(map[string]interface{}); ok {
			fmt.Printf("   Response: %+v\n", responseData)
		}
	} else {
		fmt.Printf("❌ Start watch data: FAILED - %s\n", watchResult.Error)
	}
}

func parseUintEnv(key string, defaultValue uint) uint {
	str := os.Getenv(key)
	if str == "" {
		return defaultValue
	}
	val, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return defaultValue
	}
	return uint(val)
}

func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
