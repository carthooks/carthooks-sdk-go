package carthooks

import (
	"testing"
)

func TestResult_GetData(t *testing.T) {
	tests := []struct {
		name    string
		result  *Result
		target  interface{}
		wantErr bool
	}{
		{
			name: "successful result with valid data",
			result: &Result{
				Success: true,
				Data: map[string]interface{}{
					"id":    1,
					"title": "Test Item",
				},
			},
			target:  &map[string]interface{}{},
			wantErr: false,
		},
		{
			name: "failed result",
			result: &Result{
				Success: false,
				Error:   "Test error",
			},
			target:  &map[string]interface{}{},
			wantErr: true,
		},
		{
			name: "successful result with nil data",
			result: &Result{
				Success: true,
				Data:    nil,
			},
			target:  &map[string]interface{}{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.result.GetData(tt.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("Result.GetData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestResult_GetRecords(t *testing.T) {
	tests := []struct {
		name    string
		result  *Result
		want    int // expected number of records
		wantErr bool
	}{
		{
			name: "valid records array",
			result: &Result{
				Success: true,
				Data: []interface{}{
					map[string]interface{}{
						"id":         1,
						"title":      "Item 1",
						"created_at": int64(1640995200),
						"updated_at": int64(1640995200),
						"creator":    uint(1),
						"fields":     map[string]interface{}{"f_1": "value1"},
					},
					map[string]interface{}{
						"id":         2,
						"title":      "Item 2",
						"created_at": int64(1640995300),
						"updated_at": int64(1640995300),
						"creator":    uint(1),
						"fields":     map[string]interface{}{"f_2": "value2"},
					},
				},
			},
			want:    2,
			wantErr: false,
		},
		{
			name: "empty records array",
			result: &Result{
				Success: true,
				Data:    []interface{}{},
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "failed result",
			result: &Result{
				Success: false,
				Error:   "Test error",
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			records, err := tt.result.GetRecords()
			if (err != nil) != tt.wantErr {
				t.Errorf("Result.GetRecords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(records) != tt.want {
				t.Errorf("Result.GetRecords() got %d records, want %d", len(records), tt.want)
			}
		})
	}
}

func TestResult_GetRecord(t *testing.T) {
	tests := []struct {
		name    string
		result  *Result
		wantID  uint
		wantErr bool
	}{
		{
			name: "valid single record",
			result: &Result{
				Success: true,
				Data: map[string]interface{}{
					"id":         uint(123),
					"title":      "Test Item",
					"created_at": int64(1640995200),
					"updated_at": int64(1640995200),
					"creator":    uint(1),
					"fields":     map[string]interface{}{"f_1": "value1"},
				},
			},
			wantID:  123,
			wantErr: false,
		},
		{
			name: "failed result",
			result: &Result{
				Success: false,
				Error:   "Test error",
			},
			wantID:  0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			record, err := tt.result.GetRecord()
			if (err != nil) != tt.wantErr {
				t.Errorf("Result.GetRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && record.ID != tt.wantID {
				t.Errorf("Result.GetRecord() got ID %d, want %d", record.ID, tt.wantID)
			}
		})
	}
}

func TestResult_GetString(t *testing.T) {
	tests := []struct {
		name    string
		result  *Result
		want    string
		wantErr bool
	}{
		{
			name: "valid string data",
			result: &Result{
				Success: true,
				Data:    "test string",
			},
			want:    "test string",
			wantErr: false,
		},
		{
			name: "non-string data",
			result: &Result{
				Success: true,
				Data:    123,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "failed result",
			result: &Result{
				Success: false,
				Error:   "Test error",
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.result.GetString()
			if (err != nil) != tt.wantErr {
				t.Errorf("Result.GetString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Result.GetString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResult_GetInt(t *testing.T) {
	tests := []struct {
		name    string
		result  *Result
		want    int
		wantErr bool
	}{
		{
			name: "valid int data",
			result: &Result{
				Success: true,
				Data:    123,
			},
			want:    123,
			wantErr: false,
		},
		{
			name: "valid float64 data",
			result: &Result{
				Success: true,
				Data:    float64(456),
			},
			want:    456,
			wantErr: false,
		},
		{
			name: "non-numeric data",
			result: &Result{
				Success: true,
				Data:    "not a number",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "failed result",
			result: &Result{
				Success: false,
				Error:   "Test error",
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.result.GetInt()
			if (err != nil) != tt.wantErr {
				t.Errorf("Result.GetInt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Result.GetInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResult_GetBool(t *testing.T) {
	tests := []struct {
		name    string
		result  *Result
		want    bool
		wantErr bool
	}{
		{
			name: "valid true bool data",
			result: &Result{
				Success: true,
				Data:    true,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "valid false bool data",
			result: &Result{
				Success: true,
				Data:    false,
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "non-bool data",
			result: &Result{
				Success: true,
				Data:    "not a bool",
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "failed result",
			result: &Result{
				Success: false,
				Error:   "Test error",
			},
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.result.GetBool()
			if (err != nil) != tt.wantErr {
				t.Errorf("Result.GetBool() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Result.GetBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResult_GetPagination(t *testing.T) {
	tests := []struct {
		name string
		result *Result
		want *PaginationMeta
	}{
		{
			name: "valid pagination metadata",
			result: &Result{
				Success: true,
				Meta: map[string]interface{}{
					"pagination": map[string]interface{}{
						"page":       1,
						"pageSize":   20,
						"total":      100,
						"totalPages": 5,
					},
				},
			},
			want: &PaginationMeta{
				Page:       1,
				PageSize:   20,
				Total:      100,
				TotalPages: 5,
			},
		},
		{
			name: "no pagination metadata",
			result: &Result{
				Success: true,
				Meta:    nil,
			},
			want: nil,
		},
		{
			name: "empty meta",
			result: &Result{
				Success: true,
				Meta:    map[string]interface{}{},
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.result.GetPagination()
			if tt.want == nil {
				if got != nil {
					t.Errorf("Result.GetPagination() = %v, want nil", got)
				}
			} else {
				if got == nil {
					t.Errorf("Result.GetPagination() = nil, want %v", tt.want)
					return
				}
				if got.Page != tt.want.Page || got.PageSize != tt.want.PageSize ||
					got.Total != tt.want.Total || got.TotalPages != tt.want.TotalPages {
					t.Errorf("Result.GetPagination() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
