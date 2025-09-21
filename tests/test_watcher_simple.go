package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/carthooks/carthooks-sdk-go/carthooks"
)

func processHandler(ctx interface{}, record map[string]interface{}) {
	fmt.Println("=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=")
	fmt.Println("New data item received:")
	fmt.Printf("Record: %+v\n", record)
	fmt.Println("=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=")
}

func main() {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	clientID := os.Getenv("CARTHOOKS_CLIENT_ID")
	clientSecret := os.Getenv("CARTHOOKS_CLIENT_SECRET")
	apiURL := os.Getenv("CARTHOOKS_API_URL")
	sqsQueueURL := os.Getenv("SQS_QUEUE_URL")

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

	collectionIDStr := os.Getenv("TEST_COLLECTION_ID")
	if collectionIDStr == "" {
		collectionIDStr = "3911747783"
	}
	collectionID, _ := strconv.ParseUint(collectionIDStr, 10, 32)

	// Get hostname for watcher ID
	hostname := os.Getenv("HOSTNAME")
	if hostname == "" {
		hostname = "test"
	}
	watcherID := fmt.Sprintf("subscribe-test-%s", hostname)

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
		log.Fatalf("Failed to initialize OAuth: %s", result.Error)
	}

	log.Printf("watcher_id: %s", watcherID)

	// Create watcher
	watcher, err := carthooks.NewWatcherBuilder(client, watcherID).
		WithApp(uint(appID), uint(collectionID)).
		WithSQS(sqsQueueURL, "ap-southeast-1").
		WithFilters(map[string]interface{}{
			"f_1009": map[string]interface{}{
				"$eq": 1,
			},
		}).
		WithHandler(processHandler).
		Build()

	if err != nil {
		log.Fatalf("Failed to create watcher: %v", err)
	}

	fmt.Println("Starting to listen for data...")

	// Run watcher (this will block)
	if err := watcher.Run(); err != nil {
		log.Fatalf("Watcher error: %v", err)
	}
}
