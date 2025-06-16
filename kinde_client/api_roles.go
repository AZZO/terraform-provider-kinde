package kinde_client

import (
	"context"
	"fmt"
)

// CreateRole creates a new role
func (c *Client) CreateRole(ctx context.Context, name string, description string) (RoleResource, error) {
	reqBody := requestRoleCreate{
		Name:        name,
		Description: description,
	}

	response, err := doPostRequest[responseRoleCreate](c, ctx, c.endpointRoleCreate(), reqBody)
	if err != nil {
		return RoleResource{}, err
	}

	return response.Role, nil
}

// GetRole retrieves a role by ID
func (c *Client) GetRole(ctx context.Context, roleId string) (RoleResource, error) {
	response, err := doGetRequest[responseRoleGet](c, ctx, c.endpointRoleGet(roleId))
	if err != nil {
		return RoleResource{}, err
	}

	return response.Role, nil
}

// UpdateRole updates an existing role
func (c *Client) UpdateRole(ctx context.Context, roleId string, name string, description string) (RoleResource, error) {
	url := c.endpointRoleUpdate(roleId)

	body := requestRoleCreate{
		Name:        name,
		Description: description,
	}

	response, err := doPostRequest[responseRoleGet](c, ctx, url, body)
	if err != nil {
		return RoleResource{}, err
	}

	return response.Role, nil
}

// DeleteRole deletes a role by ID
func (c *Client) DeleteRole(ctx context.Context, roleId string) error {
	_, err := doDeleteRequest[responseRoleDelete](c, ctx, c.endpointRoleDelete(roleId))
	return err
}

// GetRolePermissions retrieves permissions for a role
func (c *Client) GetRolePermissions(ctx context.Context, roleId string) ([]string, error) {
	url := c.endpointRolePermissionsGet(roleId)

	response, err := doGetRequest[responseRolePermissions](c, ctx, url)
	if err != nil {
		return nil, err
	}

	return response.Permissions, nil
}

// GetRoleScopes retrieves all scopes for a role
func (c *Client) GetRoleScopes(ctx context.Context, roleId string) ([]ApiScopeResource, error) {
	url := c.endpointRoleScopesGet(roleId)

	response, err := doGetRequest[responseRoleScopesGet](c, ctx, url)
	if err != nil {
		return nil, err
	}

	return response.Scopes, nil
}

// CreateRoleScope creates a new scope for a role
func (c *Client) CreateRoleScope(ctx context.Context, roleId string, name string, description string) (*ApiScopeResource, error) {
	url := c.endpointRoleScopeCreate(roleId)

	body := requestRoleScopeCreate{
		Name:        name,
		Description: description,
	}

	response, err := doPostRequest[responseRoleScopesGet](c, ctx, url, body)
	if err != nil {
		return nil, err
	}

	if len(response.Scopes) == 0 {
		return nil, fmt.Errorf("no scope returned after creation")
	}

	return &response.Scopes[0], nil
}

// DeleteRoleScope deletes a scope from a role
func (c *Client) DeleteRoleScope(ctx context.Context, roleId string, scopeId string) error {
	_, err := doDeleteRequest[responseRoleScopeDelete](c, ctx, c.endpointRoleScopeDelete(roleId, scopeId))
	return err
}

// UpdateRolePermissions updates the permissions for a role using PATCH
func (c *Client) UpdateRolePermissions(ctx context.Context, roleId string, newPermissions []string) error {
	// First, get current permissions
	currentPermissions, err := c.GetRolePermissions(ctx, roleId)
	if err != nil {
		return fmt.Errorf("error getting current permissions: %w", err)
	}

	// Find permissions to add and remove
	permissionsToAdd := make([]string, 0)
	permissionsToRemove := make([]string, 0)

	// Create a map of new permissions for quick lookup
	newPermsMap := make(map[string]bool)
	for _, p := range newPermissions {
		newPermsMap[p] = true
	}

	// Find permissions to add (in newPermissions but not in currentPermissions)
	for _, p := range newPermissions {
		found := false
		for _, cp := range currentPermissions {
			if p == cp {
				found = true
				break
			}
		}
		if !found {
			permissionsToAdd = append(permissionsToAdd, p)
		}
	}

	// Find permissions to remove (in currentPermissions but not in newPermissions)
	for _, p := range currentPermissions {
		if !newPermsMap[p] {
			permissionsToRemove = append(permissionsToRemove, p)
		}
	}

	// Add new permissions
	if len(permissionsToAdd) > 0 {
		url := c.endpointRolePermissionsUpdate(roleId)
		body := map[string]interface{}{
			"permissions": permissionsToAdd,
			"operation":   "add",
		}

		_, err := c.doRequest(ctx, "PATCH", url, body, 3)
		if err != nil {
			return fmt.Errorf("error adding permissions: %w", err)
		}
	}

	// Remove old permissions
	if len(permissionsToRemove) > 0 {
		url := c.endpointRolePermissionsUpdate(roleId)
		body := map[string]interface{}{
			"permissions": permissionsToRemove,
			"operation":   "delete",
		}

		_, err := c.doRequest(ctx, "PATCH", url, body, 3)
		if err != nil {
			return fmt.Errorf("error removing permissions: %w", err)
		}
	}

	return nil
}
