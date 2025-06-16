package kinde_client

import (
	"context"
	"fmt"
)

func (c *Client) GetApi(ctx context.Context, id string) (ApiResource, error) {
	url := c.endpointApi(id)

	api, err := doGetRequest[responseApiGet](c, ctx, url)
	if err != nil {
		return ApiResource{}, err
	}

	return api.Api, nil
}

func (c *Client) CreateApi(ctx context.Context, name string, audience string) (ApiResource, error) {
	url := c.endpointApisList()

	body := requestApiCreate{
		Name:     name,
		Audience: audience,
	}

	api, err := doPostRequest[responseApiCreate](c, ctx, url, body)
	if err != nil {
		return ApiResource{}, err
	}

	return api.Api, nil
}

func (c *Client) UpdateApi(ctx context.Context, id string, name string, audience string) (ApiResource, error) {
	url := c.endpointApi(id)

	body := requestApiCreate{
		Name:     name,
		Audience: audience,
	}

	api, err := doPostRequest[responseApiGet](c, ctx, url, body)
	if err != nil {
		return ApiResource{}, err
	}

	return api.Api, nil
}

// DeleteApi deletes an API by ID.
func (c *Client) DeleteApi(ctx context.Context, id string) error {
	url := c.endpointApi(id)
	_, err := doDeleteRequest[responseApiDelete](c, ctx, url)
	return err
}

// GetApiScopes retrieves all scopes for an API
func (c *Client) GetApiScopes(ctx context.Context, apiId string) ([]ApiScopeResource, error) {
	url := c.endpointApiScopesGet(apiId)

	response, err := doGetRequest[responseApiScopesGet](c, ctx, url)
	if err != nil {
		return nil, err
	}

	return response.Scopes, nil
}

// CreateApiScope creates a new scope for an API
func (c *Client) CreateApiScope(ctx context.Context, apiId string, name string, description string) (*ApiScopeResource, error) {
	url := c.endpointApiScopeCreate(apiId)

	body := requestApiScopeCreate{
		Name:        name,
		Description: description,
	}

	response, err := doPostRequest[responseApiScopesGet](c, ctx, url, body)
	if err != nil {
		return nil, err
	}

	if len(response.Scopes) == 0 {
		return nil, fmt.Errorf("no scope returned after creation")
	}

	return &response.Scopes[0], nil
}

// DeleteApiScope deletes an API scope by ID.
func (c *Client) DeleteApiScope(ctx context.Context, apiId string, scopeId string) error {
	url := c.endpointApiScopeDelete(apiId, scopeId)
	_, err := doDeleteRequest[responseApiScopeDelete](c, ctx, url)
	return err
}

// GetApis retrieves all APIs
func (c *Client) GetApis(ctx context.Context) ([]ApiResource, error) {
	url := c.endpointApisList()

	response, err := doGetRequest[responseApisGet](c, ctx, url)
	if err != nil {
		return nil, err
	}

	return response.Apis, nil
}

// GetApiScope retrieves a single scope by ID
func (c *Client) GetApiScope(ctx context.Context, scopeId string) (*ApiScopeResource, error) {
	// First, get all APIs
	apis, err := c.GetApis(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting APIs: %w", err)
	}

	// Search through each API's scopes to find the one we want
	for _, api := range apis {
		scopes, err := c.GetApiScopes(ctx, api.Id)
		if err != nil {
			continue // Skip this API if we can't get its scopes
		}

		for _, scope := range scopes {
			if scope.Id == scopeId {
				return &scope, nil
			}
		}
	}

	return nil, fmt.Errorf("scope not found: %s", scopeId)
}
