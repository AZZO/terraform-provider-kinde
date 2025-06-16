package kinde_client

import (
	"context"
	"encoding/json"
)

// CreatePermission creates a new permission
func (c *Client) CreatePermission(ctx context.Context, name string, key string, description string) (PermissionResource, error) {
	body := requestPermissionCreate{
		Name:        name,
		Key:         key,
		Description: description,
	}

	respBody, err := c.doRequest(ctx, "POST", c.endpointPermissionCreate(), body, 3)
	if err != nil {
		return PermissionResource{}, err
	}

	var response responsePermissionCreate
	if err := json.Unmarshal(respBody, &response); err != nil {
		return PermissionResource{}, err
	}

	return response.Permission, nil
}

// GetPermission retrieves a permission by ID
func (c *Client) GetPermission(ctx context.Context, permissionId string) (PermissionResource, error) {
	respBody, err := c.doRequest(ctx, "GET", c.endpointPermissionGet(permissionId), nil, 3)
	if err != nil {
		return PermissionResource{}, err
	}

	var response responsePermissionGet
	if err := json.Unmarshal(respBody, &response); err != nil {
		return PermissionResource{}, err
	}

	return response.Permission, nil
}

// UpdatePermission updates an existing permission
func (c *Client) UpdatePermission(ctx context.Context, permissionId string, name string, description string) (PermissionResource, error) {
	body := requestPermissionUpdate{
		Name:        name,
		Description: description,
	}

	respBody, err := c.doRequest(ctx, "POST", c.endpointPermissionUpdate(permissionId), body, 3)
	if err != nil {
		return PermissionResource{}, err
	}

	var response responsePermissionUpdate
	if err := json.Unmarshal(respBody, &response); err != nil {
		return PermissionResource{}, err
	}

	return response.Permission, nil
}

// DeletePermission deletes a permission by ID.
func (c *Client) DeletePermission(ctx context.Context, permissionId string) error {
	_, err := doDeleteRequest[responsePermissionDelete](c, ctx, c.endpointPermissionDelete(permissionId))
	return err
}
