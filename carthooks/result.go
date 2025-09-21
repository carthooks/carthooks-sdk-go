package carthooks

import (
	"encoding/json"
	"fmt"
)

// Result represents the response from Carthooks API
type Result struct {
	Success bool                   `json:"success"`
	Data    interface{}            `json:"data,omitempty"`
	Error   string                 `json:"error,omitempty"`
	TraceID string                 `json:"trace_id,omitempty"`
	Meta    map[string]interface{} `json:"meta,omitempty"`
}

// String returns a string representation of the Result
func (r *Result) String() string {
	return fmt.Sprintf("CarthooksResult(success=%t, data=%v, error=%s)", r.Success, r.Data, r.Error)
}

// GetData attempts to unmarshal the result data into the provided interface
func (r *Result) GetData(v interface{}) error {
	if !r.Success {
		return fmt.Errorf("result is not successful: %s", r.Error)
	}

	if r.Data == nil {
		return fmt.Errorf("no data in result")
	}

	// Convert data to JSON and then unmarshal to target type
	jsonData, err := json.Marshal(r.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	if err := json.Unmarshal(jsonData, v); err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}

	return nil
}

// GetRecords is a convenience method to get a slice of RecordFormat
func (r *Result) GetRecords() ([]RecordFormat, error) {
	var records []RecordFormat
	if err := r.GetData(&records); err != nil {
		return nil, err
	}
	return records, nil
}

// GetRecord is a convenience method to get a single RecordFormat
func (r *Result) GetRecord() (*RecordFormat, error) {
	var record RecordFormat
	if err := r.GetData(&record); err != nil {
		return nil, err
	}
	return &record, nil
}

// GetString is a convenience method to get a string value from data
func (r *Result) GetString() (string, error) {
	if !r.Success {
		return "", fmt.Errorf("result is not successful: %s", r.Error)
	}

	if str, ok := r.Data.(string); ok {
		return str, nil
	}

	return "", fmt.Errorf("data is not a string")
}

// GetInt is a convenience method to get an integer value from data
func (r *Result) GetInt() (int, error) {
	if !r.Success {
		return 0, fmt.Errorf("result is not successful: %s", r.Error)
	}

	switch v := r.Data.(type) {
	case int:
		return v, nil
	case float64:
		return int(v), nil
	case string:
		// Try to parse string as int
		if i, err := fmt.Sscanf(v, "%d", new(int)); err == nil && i == 1 {
			var result int
			fmt.Sscanf(v, "%d", &result)
			return result, nil
		}
	}

	return 0, fmt.Errorf("data is not an integer")
}

// GetBool is a convenience method to get a boolean value from data
func (r *Result) GetBool() (bool, error) {
	if !r.Success {
		return false, fmt.Errorf("result is not successful: %s", r.Error)
	}

	if b, ok := r.Data.(bool); ok {
		return b, nil
	}

	return false, fmt.Errorf("data is not a boolean")
}

// HasError returns true if the result contains an error
func (r *Result) HasError() bool {
	return !r.Success || r.Error != ""
}

// GetError returns the error message if any
func (r *Result) GetError() string {
	if r.HasError() {
		return r.Error
	}
	return ""
}

// GetTraceID returns the trace ID for debugging
func (r *Result) GetTraceID() string {
	return r.TraceID
}

// GetMeta returns the metadata
func (r *Result) GetMeta() map[string]interface{} {
	return r.Meta
}

// GetPagination extracts pagination information from meta
func (r *Result) GetPagination() *PaginationMeta {
	if r.Meta == nil {
		return nil
	}

	if paginationData, ok := r.Meta["pagination"]; ok {
		var pagination PaginationMeta
		if jsonData, err := json.Marshal(paginationData); err == nil {
			if err := json.Unmarshal(jsonData, &pagination); err == nil {
				return &pagination
			}
		}
	}

	return nil
}

// PaginationMeta represents pagination metadata
type PaginationMeta struct {
	Page       int `json:"page"`
	PageSize   int `json:"pageSize"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}
