package carthooks

import (
	"fmt"
	"strconv"
)

// QueryOptions represents options for querying items
type QueryOptions struct {
	Pagination *PaginationOptions         `json:"pagination,omitempty"`
	Filters    map[string]interface{}     `json:"filters,omitempty"`
	Sort       []string                   `json:"sort,omitempty"`
	Fields     []string                   `json:"fields,omitempty"`
}

// PaginationOptions represents pagination parameters
type PaginationOptions struct {
	Page      int  `json:"page,omitempty"`
	PageSize  int  `json:"pageSize,omitempty"`
	WithCount bool `json:"withCount,omitempty"`
}

// LockOptions represents options for locking items
type LockOptions struct {
	LockTimeout int    `json:"lockTimeout,omitempty"`
	LockID      string `json:"lockId,omitempty"`
	Subject     string `json:"lockSubject,omitempty"`
}

// GetItems retrieves items from a collection with pagination
func (c *Client) GetItems(appID, collectionID uint, limit, start int, options map[string]string) *Result {
	path := fmt.Sprintf("/v1/apps/%d/collections/%d/items", appID, collectionID)
	
	params := map[string]string{
		"pagination[start]": strconv.Itoa(start),
		"pagination[limit]": strconv.Itoa(limit),
	}
	
	// Add additional options
	for k, v := range options {
		params[k] = v
	}
	
	resp, err := c.makeRequest("GET", path, nil, params)
	if err != nil {
		return &Result{
			Success: false,
			Error:   err.Error(),
		}
	}
	
	return c.parseResponse(resp)
}

// GetItemByID retrieves a specific item by ID
func (c *Client) GetItemByID(appID, collectionID, itemID uint, fields []string) *Result {
	path := fmt.Sprintf("/v1/apps/%d/collections/%d/items/%d", appID, collectionID, itemID)
	
	params := map[string]string{}
	if len(fields) > 0 {
		// Convert fields slice to comma-separated string
		fieldsStr := ""
		for i, field := range fields {
			if i > 0 {
				fieldsStr += ","
			}
			fieldsStr += field
		}
		params["fields"] = fieldsStr
	}
	
	resp, err := c.makeRequest("GET", path, nil, params)
	if err != nil {
		return &Result{
			Success: false,
			Error:   err.Error(),
		}
	}
	
	return c.parseResponse(resp)
}

// QueryItems queries items with advanced filtering and sorting
func (c *Client) QueryItems(appID, collectionID uint, options *QueryOptions) *Result {
	path := fmt.Sprintf("/v1/apps/%d/collections/%d/items/query", appID, collectionID)
	
	resp, err := c.makeRequest("POST", path, options, nil)
	if err != nil {
		return &Result{
			Success: false,
			Error:   err.Error(),
		}
	}
	
	return c.parseResponse(resp)
}

// CreateItem creates a new item in a collection
func (c *Client) CreateItem(appID, collectionID uint, data map[string]interface{}) *Result {
	path := fmt.Sprintf("/v1/apps/%d/collections/%d/items", appID, collectionID)
	
	body := map[string]interface{}{
		"data": data,
	}
	
	resp, err := c.makeRequest("POST", path, body, nil)
	if err != nil {
		return &Result{
			Success: false,
			Error:   err.Error(),
		}
	}
	
	return c.parseResponse(resp)
}

// UpdateItem updates an existing item
func (c *Client) UpdateItem(appID, collectionID, itemID uint, data map[string]interface{}) *Result {
	path := fmt.Sprintf("/v1/apps/%d/collections/%d/items/%d", appID, collectionID, itemID)
	
	body := map[string]interface{}{
		"data": data,
	}
	
	resp, err := c.makeRequest("PUT", path, body, nil)
	if err != nil {
		return &Result{
			Success: false,
			Error:   err.Error(),
		}
	}
	
	return c.parseResponse(resp)
}

// DeleteItem deletes an item from a collection
func (c *Client) DeleteItem(appID, collectionID, itemID uint) *Result {
	path := fmt.Sprintf("/v1/apps/%d/collections/%d/items/%d", appID, collectionID, itemID)
	
	resp, err := c.makeRequest("DELETE", path, nil, nil)
	if err != nil {
		return &Result{
			Success: false,
			Error:   err.Error(),
		}
	}
	
	return c.parseResponse(resp)
}

// LockItem locks an item to prevent concurrent modifications
func (c *Client) LockItem(appID, collectionID, itemID uint, options *LockOptions) *Result {
	path := fmt.Sprintf("/v1/apps/%d/collections/%d/items/%d/lock", appID, collectionID, itemID)
	
	body := map[string]interface{}{}
	if options != nil {
		if options.LockTimeout > 0 {
			body["lockTimeout"] = options.LockTimeout
		}
		if options.LockID != "" {
			body["lockId"] = options.LockID
		}
		if options.Subject != "" {
			body["lockSubject"] = options.Subject
		}
	}
	
	resp, err := c.makeRequest("POST", path, body, nil)
	if err != nil {
		return &Result{
			Success: false,
			Error:   err.Error(),
		}
	}
	
	return c.parseResponse(resp)
}

// UnlockItem unlocks a previously locked item
func (c *Client) UnlockItem(appID, collectionID, itemID uint, lockID string) *Result {
	path := fmt.Sprintf("/v1/apps/%d/collections/%d/items/%d/unlock", appID, collectionID, itemID)
	
	body := map[string]interface{}{}
	if lockID != "" {
		body["lockId"] = lockID
	}
	
	resp, err := c.makeRequest("POST", path, body, nil)
	if err != nil {
		return &Result{
			Success: false,
			Error:   err.Error(),
		}
	}
	
	return c.parseResponse(resp)
}

// CreateSubItem creates a sub-item in a subform field
func (c *Client) CreateSubItem(appID, collectionID, itemID, fieldID uint, data map[string]interface{}) *Result {
	path := fmt.Sprintf("/v1/apps/%d/collections/%d/items/%d/subform/%d", appID, collectionID, itemID, fieldID)
	
	body := map[string]interface{}{
		"data": data,
	}
	
	resp, err := c.makeRequest("POST", path, body, nil)
	if err != nil {
		return &Result{
			Success: false,
			Error:   err.Error(),
		}
	}
	
	return c.parseResponse(resp)
}

// UpdateSubItem updates a sub-item in a subform field
func (c *Client) UpdateSubItem(appID, collectionID, itemID, fieldID, subItemID uint, data map[string]interface{}) *Result {
	path := fmt.Sprintf("/v1/apps/%d/collections/%d/items/%d/subform/%d/items/%d", appID, collectionID, itemID, fieldID, subItemID)
	
	body := map[string]interface{}{
		"data": data,
	}
	
	resp, err := c.makeRequest("PUT", path, body, nil)
	if err != nil {
		return &Result{
			Success: false,
			Error:   err.Error(),
		}
	}
	
	return c.parseResponse(resp)
}

// DeleteSubItem deletes a sub-item from a subform field
func (c *Client) DeleteSubItem(appID, collectionID, itemID, fieldID, subItemID uint) *Result {
	path := fmt.Sprintf("/v1/apps/%d/collections/%d/items/%d/subform/%d/items/%d", appID, collectionID, itemID, fieldID, subItemID)
	
	resp, err := c.makeRequest("DELETE", path, nil, nil)
	if err != nil {
		return &Result{
			Success: false,
			Error:   err.Error(),
		}
	}
	
	return c.parseResponse(resp)
}
