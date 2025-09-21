package main

import (
	"fmt"
	"time"

	"github.com/carthooks/carthooks-sdk-go/carthooks"
)

func main() {
	// Initialize client with configuration
	client := carthooks.NewClient(&carthooks.ClientConfig{
		BaseURL:     "https://api.carthooks.com",
		AccessToken: "your-access-token",
		Timeout:     30 * time.Second,
		Debug:       true,
	})

	// Example app and collection IDs
	appID := uint(123456)
	collectionID := uint(789012)

	fmt.Println("Carthooks Go SDK Examples")

	// Example 1: Basic CRUD Operations
	basicCRUDExamples(client, appID, collectionID)

	// Example 2: Advanced Querying
	advancedQueryExamples(client, appID, collectionID)

	// Example 3: Item Locking
	itemLockingExamples(client, appID, collectionID)

	// Example 4: File Upload
	fileUploadExamples(client)

	// Example 5: User Management
	userManagementExamples(client)

	// Example 6: Data Monitoring
	dataMonitoringExamples(client, appID, collectionID)

	// Example 7: Error Handling
	errorHandlingExamples(client, appID, collectionID)
}

func basicCRUDExamples(client *carthooks.Client, appID, collectionID uint) {
	fmt.Println("Basic CRUD Operations:")

	// Create an item
	fmt.Println("Creating new item...")
	createData := map[string]interface{}{
		"title":  "Test Item from Go SDK",
		"f_1001": "active",
		"f_1002": 100,
		"f_1003": []string{"tag1", "tag2"},
	}

	result := client.CreateItem(appID, collectionID, createData)
	if result.Success {
		var record carthooks.RecordFormat
		if err := result.GetData(&record); err == nil {
			fmt.Printf("Created item with ID: %d\n", record.ID)

			// Update the item
			fmt.Println("Updating item...")
			updateData := map[string]interface{}{
				"title":  "Updated Test Item",
				"f_1002": 150,
			}

			updateResult := client.UpdateItem(appID, collectionID, record.ID, updateData)
			if updateResult.Success {
				fmt.Println("Item updated successfully")
			}

			// Get the item
			fmt.Println("Retrieving item...")
			getResult := client.GetItemByID(appID, collectionID, record.ID, nil)
			if getResult.Success {
				var updatedRecord carthooks.RecordFormat
				if err := getResult.GetData(&updatedRecord); err == nil {
					fmt.Printf("Retrieved item: %s\n", updatedRecord.Title)
				}
			}

			// Delete the item
			fmt.Println("Deleting item...")
			deleteResult := client.DeleteItem(appID, collectionID, record.ID)
			if deleteResult.Success {
				fmt.Println("Item deleted successfully")
			}
		}
	} else {
		fmt.Printf("Failed to create item: %s\n", result.Error)
	}
}

func advancedQueryExamples(client *carthooks.Client, appID, collectionID uint) {
	fmt.Println("\nAdvanced Querying:")

	// Query with filters and pagination
	queryOptions := &carthooks.QueryOptions{
		Pagination: &carthooks.PaginationOptions{
			Page:      1,
			PageSize:  10,
			WithCount: true,
		},
		Filters: map[string]interface{}{
			"f_1001": map[string]interface{}{"$eq": "active"},
			"f_1002": map[string]interface{}{"$gte": 50},
		},
		Sort:   []string{"created_at:desc", "f_1002:asc"},
		Fields: []string{"title", "f_1001", "f_1002", "created_at"},
	}

	result := client.QueryItems(appID, collectionID, queryOptions)
	if result.Success {
		records, err := result.GetRecords()
		if err == nil {
			fmt.Printf("Found %d items matching criteria\n", len(records))

			// Show pagination info
			if pagination := result.GetPagination(); pagination != nil {
				fmt.Printf("Page %d of %d (Total: %d items)\n",
					pagination.Page, pagination.TotalPages, pagination.Total)
			}

			// Show first few records
			for i, record := range records {
				if i >= 3 {
					break
				}
				fmt.Printf("   - %s (ID: %d)\n", record.Title, record.ID)
			}
		}
	} else {
		fmt.Printf("Query failed: %s\n", result.Error)
	}
}

func itemLockingExamples(client *carthooks.Client, appID, collectionID uint) {
	fmt.Println("\nItem Locking:")

	// Get first item to demonstrate locking
	result := client.GetItems(appID, collectionID, 1, 0, nil)
	if result.Success {
		records, err := result.GetRecords()
		if err == nil && len(records) > 0 {
			itemID := records[0].ID

			// Lock the item
			lockOptions := &carthooks.LockOptions{
				LockTimeout: 300, // 5 minutes
				LockID:      "go-sdk-example",
				Subject:     "Processing in Go SDK example",
			}

			lockResult := client.LockItem(appID, collectionID, itemID, lockOptions)
			if lockResult.Success {
				fmt.Printf("Locked item %d\n", itemID)

				// Simulate some processing
				time.Sleep(2 * time.Second)

				// Unlock the item
				unlockResult := client.UnlockItem(appID, collectionID, itemID, "go-sdk-example")
				if unlockResult.Success {
					fmt.Printf("Unlocked item %d\n", itemID)
				}
			} else {
				fmt.Printf("Failed to lock item: %s\n", lockResult.Error)
			}
		}
	}
}

func fileUploadExamples(client *carthooks.Client) {
	fmt.Println("\nFile Upload:")

	result := client.GetUploadToken()
	if result.Success {
		var token carthooks.UploadToken
		if err := result.GetData(&token); err == nil {
			fmt.Printf("Got upload token: %s\n", token.Token)
			fmt.Printf("Upload URL: %s\n", token.URL)
			fmt.Printf("Expires at: %s\n", token.ExpiresAt)
		}
	} else {
		fmt.Printf("Failed to get upload token: %s\n", result.Error)
	}
}

func userManagementExamples(client *carthooks.Client) {
	fmt.Println("User Management:")

	// Get user by token (example)
	result := client.GetUserByToken("example-token")
	if result.Success {
		var user carthooks.User
		if err := result.GetData(&user); err == nil {
			fmt.Printf("User: %s (%s)\n", user.Name, user.Email)
		}
	} else {
		fmt.Printf("User token example (expected to fail): %s\n", result.Error)
	}
}

func dataMonitoringExamples(client *carthooks.Client, appID, collectionID uint) {
	fmt.Println("Data Monitoring:")

	watchOptions := &carthooks.WatchDataOptions{
		EndpointURL:  "https://sqs.us-east-1.amazonaws.com/123456789/my-queue",
		EndpointType: "sqs",
		Name:         "go-sdk-watch",
		AppID:        appID,
		CollectionID: collectionID,
		Filters: map[string]interface{}{
			"f_1001": map[string]interface{}{"$eq": "active"},
		},
		Age: 86400, // 24 hours
	}

	result := client.StartWatchData(watchOptions)
	if result.Success {
		var watchResponse carthooks.WatchDataResponse
		if err := result.GetData(&watchResponse); err == nil {
			fmt.Printf("Started watching data: %s\n", watchResponse.WatchID)
		}
	} else {
		fmt.Printf("Watch example (may fail without proper SQS setup): %s\n", result.Error)
	}
}

func errorHandlingExamples(client *carthooks.Client, appID, collectionID uint) {
	fmt.Println("\nError Handling:")

	// Intentionally cause an error by using invalid ID
	result := client.GetItemByID(appID, collectionID, 999999999, nil)

	// Method 1: Check Success field
	if !result.Success {
		fmt.Printf("Method 1 - Error detected: %s\n", result.Error)
		if result.TraceID != "" {
			fmt.Printf("Trace ID: %s\n", result.TraceID)
		}
	}

	// Method 2: Use HasError method
	if result.HasError() {
		fmt.Printf("Method 2 - Error: %s\n", result.GetError())
	}

	// Method 3: Use GetData with error handling
	var record carthooks.RecordFormat
	if err := result.GetData(&record); err != nil {
		fmt.Printf("Method 3 - Data parsing failed: %v\n", err)
	}
}
