package kinde_client

import (
	"fmt"
)

// Role endpoints
func (c *Client) endpointRoleCreate() string {
	return fmt.Sprintf("%s/api/v1/roles", c.IssuerUrl)
}

func (c *Client) endpointRoleGet(roleId string) string {
	return fmt.Sprintf("%s/api/v1/roles/%s", c.IssuerUrl, roleId)
}

func (c *Client) endpointRoleUpdate(roleId string) string {
	return fmt.Sprintf("%s/api/v1/roles/%s", c.IssuerUrl, roleId)
}

func (c *Client) endpointRoleDelete(roleId string) string {
	return fmt.Sprintf("%s/api/v1/roles/%s", c.IssuerUrl, roleId)
}

func (c *Client) endpointRolePermissionsGet(roleId string) string {
	return fmt.Sprintf("%s/api/v1/roles/%s/permissions", c.IssuerUrl, roleId)
}

func (c *Client) endpointRolePermissionsUpdate(roleId string) string {
	return fmt.Sprintf("%s/api/v1/roles/%s/permissions", c.IssuerUrl, roleId)
}

func (c *Client) endpointRoleScopesGet(roleId string) string {
	return fmt.Sprintf("%s/api/v1/roles/%s/scopes", c.IssuerUrl, roleId)
}

func (c *Client) endpointRoleScopeCreate(roleId string) string {
	return fmt.Sprintf("%s/api/v1/roles/%s/scopes", c.IssuerUrl, roleId)
}

func (c *Client) endpointRoleScopeDelete(roleId string, scopeId string) string {
	return fmt.Sprintf("%s/api/v1/roles/%s/scopes/%s", c.IssuerUrl, roleId, scopeId)
}

func (c *Client) endpointApplicationsList() string {
	return fmt.Sprintf("%s/api/v1/applications", c.IssuerUrl)
}

func (c *Client) endpointApplication(id string) string {
	return fmt.Sprintf("%s/api/v1/applications/%v", c.IssuerUrl, id)
}

func (c *Client) endpointApisList() string {
	return fmt.Sprintf("%s/api/v1/apis", c.IssuerUrl)
}

func (c *Client) endpointApi(id string) string {
	return fmt.Sprintf("%s/api/v1/apis/%v", c.IssuerUrl, id)
}

func (c *Client) endpointApiScopesGet(apiId string) string {
	return fmt.Sprintf("%s/api/v1/apis/%v/scopes", c.IssuerUrl, apiId)
}

func (c *Client) endpointApiScopeCreate(apiId string) string {
	return fmt.Sprintf("%s/api/v1/apis/%v/scopes", c.IssuerUrl, apiId)
}

func (c *Client) endpointApiScopeDelete(apiId string, scopeId string) string {
	return fmt.Sprintf("%s/api/v1/apis/%v/scopes/%v", c.IssuerUrl, apiId, scopeId)
}

// Permission endpoints
func (c *Client) endpointPermissionCreate() string {
	return fmt.Sprintf("%s/api/v1/permissions", c.IssuerUrl)
}

func (c *Client) endpointPermissionGet(permissionId string) string {
	return fmt.Sprintf("%s/api/v1/permissions/%s", c.IssuerUrl, permissionId)
}

func (c *Client) endpointPermissionUpdate(permissionId string) string {
	return fmt.Sprintf("%s/api/v1/permissions/%s", c.IssuerUrl, permissionId)
}

func (c *Client) endpointPermissionDelete(permissionId string) string {
	return fmt.Sprintf("%s/api/v1/permissions/%s", c.IssuerUrl, permissionId)
}

// Environment endpoints
func (c *Client) endpointEnvironmentGet(environmentId string) string {
	return fmt.Sprintf("%s/api/v1/environment/%s", c.IssuerUrl, environmentId)
}

// Application callback endpoints
func (c *Client) endpointApplicationCallbacksGet(id string) string {
	return fmt.Sprintf("%s/api/v1/applications/%s/callbacks", c.IssuerUrl, id)
}

func (c *Client) endpointApplicationCallbacksUpdate(id string) string {
	return fmt.Sprintf("%s/api/v1/applications/%s/callbacks", c.IssuerUrl, id)
}
