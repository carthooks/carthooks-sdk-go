package carthooks

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name   string
		config *ClientConfig
		want   string // expected base URL
	}{
		{
			name:   "default config",
			config: nil,
			want:   "https://api.carthooks.com",
		},
		{
			name: "custom config",
			config: &ClientConfig{
				BaseURL:     "https://custom.api.com",
				AccessToken: "test-token",
				Timeout:     10 * time.Second,
			},
			want: "https://custom.api.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.config)
			if client.baseURL != tt.want {
				t.Errorf("NewClient() baseURL = %v, want %v", client.baseURL, tt.want)
			}
		})
	}
}

func TestClient_SetAccessToken(t *testing.T) {
	client := NewClient(nil)
	token := "test-token-123"
	
	client.SetAccessToken(token)
	
	if client.accessToken != token {
		t.Errorf("SetAccessToken() accessToken = %v, want %v", client.accessToken, token)
	}
	
	expectedAuth := "Bearer " + token
	if client.headers["Authorization"] != expectedAuth {
		t.Errorf("SetAccessToken() Authorization header = %v, want %v", 
			client.headers["Authorization"], expectedAuth)
	}
}

func TestClient_GetItems(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		
		expectedPath := "/v1/apps/123/collections/456/items"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}
		
		// Check query parameters
		query := r.URL.Query()
		if query.Get("pagination[start]") != "0" {
			t.Errorf("Expected pagination[start]=0, got %s", query.Get("pagination[start]"))
		}
		if query.Get("pagination[limit]") != "20" {
			t.Errorf("Expected pagination[limit]=20, got %s", query.Get("pagination[limit]"))
		}
		
		// Return mock response
		response := map[string]interface{}{
			"data": []map[string]interface{}{
				{
					"id":         1,
					"title":      "Test Item 1",
					"created_at": 1640995200,
					"updated_at": 1640995200,
					"creator":    1,
					"fields": map[string]interface{}{
						"f_1001": "value1",
					},
				},
				{
					"id":         2,
					"title":      "Test Item 2",
					"created_at": 1640995300,
					"updated_at": 1640995300,
					"creator":    1,
					"fields": map[string]interface{}{
						"f_1002": "value2",
					},
				},
			},
			"meta": map[string]interface{}{
				"pagination": map[string]interface{}{
					"page":       1,
					"pageSize":   20,
					"total":      2,
					"totalPages": 1,
				},
			},
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	
	// Create client with mock server URL
	client := NewClient(&ClientConfig{
		BaseURL: server.URL,
	})
	
	// Test GetItems
	result := client.GetItems(123, 456, 20, 0, nil)
	
	if !result.Success {
		t.Errorf("GetItems() failed: %s", result.Error)
	}
	
	// Test data parsing
	records, err := result.GetRecords()
	if err != nil {
		t.Errorf("GetRecords() failed: %v", err)
	}
	
	if len(records) != 2 {
		t.Errorf("Expected 2 records, got %d", len(records))
	}
	
	if records[0].Title != "Test Item 1" {
		t.Errorf("Expected title 'Test Item 1', got '%s'", records[0].Title)
	}
	
	// Test pagination
	pagination := result.GetPagination()
	if pagination == nil {
		t.Error("Expected pagination metadata, got nil")
	} else {
		if pagination.Total != 2 {
			t.Errorf("Expected total 2, got %d", pagination.Total)
		}
	}
}

func TestClient_CreateItem(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		
		expectedPath := "/v1/apps/123/collections/456/items"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}
		
		// Verify content type
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}
		
		// Parse request body
		var requestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}
		
		// Verify request data
		data, ok := requestBody["data"].(map[string]interface{})
		if !ok {
			t.Error("Expected data field in request body")
		}
		
		if data["title"] != "Test Item" {
			t.Errorf("Expected title 'Test Item', got %v", data["title"])
		}
		
		// Return mock response
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"id":         123,
				"title":      data["title"],
				"created_at": 1640995200,
				"updated_at": 1640995200,
				"creator":    1,
				"fields":     data,
			},
		}
		
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	
	// Create client with mock server URL
	client := NewClient(&ClientConfig{
		BaseURL: server.URL,
	})
	
	// Test CreateItem
	data := map[string]interface{}{
		"title":  "Test Item",
		"f_1001": "test value",
	}
	
	result := client.CreateItem(123, 456, data)
	
	if !result.Success {
		t.Errorf("CreateItem() failed: %s", result.Error)
	}
	
	// Test data parsing
	record, err := result.GetRecord()
	if err != nil {
		t.Errorf("GetRecord() failed: %v", err)
	}
	
	if record.ID != 123 {
		t.Errorf("Expected ID 123, got %d", record.ID)
	}
	
	if record.Title != "Test Item" {
		t.Errorf("Expected title 'Test Item', got '%s'", record.Title)
	}
}

func TestClient_ErrorHandling(t *testing.T) {
	// Create mock server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"error": map[string]interface{}{
				"message": "Item not found",
				"code":    "NOT_FOUND",
			},
			"trace_id": "test-trace-123",
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()
	
	// Create client with mock server URL
	client := NewClient(&ClientConfig{
		BaseURL: server.URL,
	})
	
	// Test error handling
	result := client.GetItemByID(123, 456, 999, nil)
	
	if result.Success {
		t.Error("Expected failure, got success")
	}
	
	if result.Error != "Item not found" {
		t.Errorf("Expected error 'Item not found', got '%s'", result.Error)
	}
	
	if result.TraceID != "test-trace-123" {
		t.Errorf("Expected trace ID 'test-trace-123', got '%s'", result.TraceID)
	}
	
	// Test HasError method
	if !result.HasError() {
		t.Error("HasError() should return true for failed result")
	}
	
	// Test GetError method
	if result.GetError() != "Item not found" {
		t.Errorf("GetError() returned '%s', expected 'Item not found'", result.GetError())
	}
}
