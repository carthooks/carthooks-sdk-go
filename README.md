# Carthooks Go SDK

A comprehensive Go SDK for the Carthooks API with full type safety and extensive functionality.

## Installation

```bash
go get github.com/carthooks/carthooks-sdk-go
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/carthooks/carthooks-sdk-go/carthooks"
)

func main() {
    // Initialize client
    client := carthooks.NewClient(&carthooks.ClientConfig{
        BaseURL:     "https://api.carthooks.com",
        AccessToken: "your-access-token",
        Debug:       true,
    })
    
    // Get items from a collection
    result := client.GetItems(123456, 789012, 20, 0, nil)
    if result.Success {
        var records []carthooks.RecordFormat
        if err := result.GetData(&records); err == nil {
            fmt.Printf("Found %d items\n", len(records))
        }
    } else {
        log.Printf("Error: %s", result.Error)
    }
}
```

## Configuration

### Environment Variables

```bash
export CARTHOOKS_API_URL="https://api.carthooks.com"
export CARTHOOKS_ACCESS_TOKEN="your-access-token"
export CARTHOOKS_TIMEOUT="30s"
export CARTHOOKS_SDK_DEBUG="true"
```

### Programmatic Configuration

```go
config := &carthooks.ClientConfig{
    BaseURL:     "https://api.carthooks.com",
    AccessToken: "your-access-token",
    Timeout:     30 * time.Second,
    Headers: map[string]string{
        "Custom-Header": "value",
    },
    Debug: true,
}

client := carthooks.NewClient(config)
```

## Basic Operations

### Get Items

```go
// Simple get with pagination
result := client.GetItems(appID, collectionID, 20, 0, nil)

// With additional options
options := map[string]string{
    "sort": "created_at:desc",
}
result := client.GetItems(appID, collectionID, 20, 0, options)
```

### Get Single Item

```go
// Get item by ID
result := client.GetItemByID(appID, collectionID, itemID, nil)

// Get specific fields only
fields := []string{"title", "f_1001", "f_1002"}
result := client.GetItemByID(appID, collectionID, itemID, fields)
```

### Query Items with Advanced Filtering

```go
queryOptions := &carthooks.QueryOptions{
    Pagination: &carthooks.PaginationOptions{
        Page:      1,
        PageSize:  20,
        WithCount: true,
    },
    Filters: map[string]interface{}{
        "f_1001": map[string]interface{}{"$eq": "active"},
        "f_1002": map[string]interface{}{"$gt": 100},
    },
    Sort: []string{"f_1001:asc", "updated_at:desc"},
    Fields: []string{"title", "f_1001", "f_1002"},
}

result := client.QueryItems(appID, collectionID, queryOptions)
```

### Create Item

```go
data := map[string]interface{}{
    "title": "New Item",
    "f_1001": "active",
    "f_1002": 150,
}

result := client.CreateItem(appID, collectionID, data)
if result.Success {
    var record carthooks.RecordFormat
    if err := result.GetData(&record); err == nil {
        fmt.Printf("Created item with ID: %d\n", record.ID)
    }
}
```

### Update Item

```go
data := map[string]interface{}{
    "title": "Updated Item",
    "f_1002": 200,
}

result := client.UpdateItem(appID, collectionID, itemID, data)
```

### Delete Item

```go
result := client.DeleteItem(appID, collectionID, itemID)
if result.Success {
    fmt.Println("Item deleted successfully")
}
```

## Advanced Features

### Item Locking

```go
// Lock an item
lockOptions := &carthooks.LockOptions{
    LockTimeout: 600,  // 10 minutes
    LockID:      "my-process-id",
    Subject:     "Processing item",
}

result := client.LockItem(appID, collectionID, itemID, lockOptions)

// Unlock the item
result = client.UnlockItem(appID, collectionID, itemID, "my-process-id")
```

### Subform Operations

```go
// Create sub-item
subData := map[string]interface{}{
    "sub_field_1": "value1",
    "sub_field_2": "value2",
}
result := client.CreateSubItem(appID, collectionID, itemID, fieldID, subData)

// Update sub-item
result = client.UpdateSubItem(appID, collectionID, itemID, fieldID, subItemID, subData)

// Delete sub-item
result = client.DeleteSubItem(appID, collectionID, itemID, fieldID, subItemID)
```

### File Upload

```go
// Get upload token
result := client.GetUploadToken()
if result.Success {
    var token carthooks.UploadToken
    if err := result.GetData(&token); err == nil {
        fmt.Printf("Upload URL: %s\n", token.URL)
    }
}
```

### User Management

```go
// Get user by ID
result := client.GetUser(userID)

// Get user by token
result = client.GetUserByToken("user-token")
```

### Data Monitoring

```go
watchOptions := &carthooks.WatchDataOptions{
    EndpointURL:  "https://sqs.region.amazonaws.com/account/queue-name",
    EndpointType: "sqs",
    Name:         "my-watch",
    AppID:        appID,
    CollectionID: collectionID,
    Filters: map[string]interface{}{
        "f_1001": map[string]interface{}{"$eq": "active"},
    },
    Age: 86400, // 24 hours
}

result := client.StartWatchData(watchOptions)
```

### Connection Management

The SDK provides comprehensive support for managing hooklet connections:

```go
// Create a new connection
connectionReq := &carthooks.CreateConnectionRequest{
    HookletID:   123,
    DevClientID: 456,
    Title:       "My Integration Connection",
    IconUrl:     "https://example.com/icon.png",
    Description: "This connection integrates with external service",
}

result := client.CreateConnection(appID, connectionReq)
if result.Success {
    var connection carthooks.Connection
    if err := result.ParseData(&connection); err == nil {
        fmt.Printf("Connection created: ID=%d, Title=%s\n", connection.ID, connection.Title)
    }
}

// Create connection log entries
logReq := &carthooks.CreateConnectionLogRequest{
    Status:  uint8(carthooks.ConnectionLogStatusCreated),
    Message: "Connection established successfully",
}

logResult := client.CreateConnectionLog(appID, connectionID, logReq)

// Record connection usage
usageReq := &carthooks.CreateConnectionUsageRequest{
    Usage: 100, // Record 100 usage units
}

usageResult := client.CreateConnectionUsage(appID, connectionID, usageReq)
```

#### Connection Status Constants

```go
const (
    ConnectionStatusPending  ConnectionStatus = 0
    ConnectionStatusActive   ConnectionStatus = 1
    ConnectionStatusInactive ConnectionStatus = 2
)
```

#### Connection Log Status Constants

```go
const (
    ConnectionLogStatusCreated ConnectionLogStatus = 1
    ConnectionLogStatusUpdated ConnectionLogStatus = 2
    ConnectionLogStatusWarn    ConnectionLogStatus = 3
    ConnectionLogStatusError   ConnectionLogStatus = 4
)
```

## Error Handling

```go
result := client.GetItems(appID, collectionID, 20, 0, nil)

// Method 1: Check Success field
if !result.Success {
    log.Printf("Error: %s", result.Error)
    if result.TraceID != "" {
        log.Printf("Trace ID: %s", result.TraceID)
    }
}

// Method 2: Use HasError method
if result.HasError() {
    log.Printf("Error: %s", result.GetError())
}

// Method 3: Use GetData with error handling
var records []carthooks.RecordFormat
if err := result.GetData(&records); err != nil {
    log.Printf("Failed to parse data: %v", err)
}
```

## Working with Results

```go
result := client.GetItems(appID, collectionID, 20, 0, nil)

// Get records
records, err := result.GetRecords()
if err != nil {
    log.Printf("Error: %v", err)
    return
}

// Get pagination info
if pagination := result.GetPagination(); pagination != nil {
    fmt.Printf("Page %d of %d (Total: %d)\n", 
        pagination.Page, pagination.TotalPages, pagination.Total)
}

// Get metadata
meta := result.GetMeta()
if meta != nil {
    fmt.Printf("Metadata: %+v\n", meta)
}
```

## Type Safety

The SDK provides full type safety with predefined structures:

```go
// RecordFormat represents a data record
type RecordFormat struct {
    ID        uint                   `json:"id"`
    Title     string                 `json:"title"`
    CreatedAt int64                  `json:"created_at"`
    UpdatedAt int64                  `json:"updated_at"`
    Creator   uint                   `json:"creator"`
    Fields    map[string]interface{} `json:"fields"`
}

// EventMessage for webhook/SQS events
type EventMessage struct {
    Version string           `json:"version"`
    Meta    EventMessageMeta `json:"meta"`
    Payload interface{}      `json:"payload"`
}
```

## Debug Mode

Enable debug mode to see detailed request/response information:

```go
client := carthooks.NewClient(&carthooks.ClientConfig{
    Debug: true,
})

// Or via environment variable
// export CARTHOOKS_SDK_DEBUG=true
```

## License

MIT License
