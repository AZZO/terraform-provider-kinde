package kinde_client

import (
	"context"
	"fmt"
)

// GetEnvironment retrieves information about a specific environment
func (c *Client) GetEnvironment(ctx context.Context, environmentId string) (*EnvironmentResource, error) {
	url := fmt.Sprintf("%s/api/v1/environment/%s", c.IssuerUrl, environmentId)

	response, err := doGetRequest[responseEnvironmentGet](c, ctx, url)
	if err != nil {
		return nil, fmt.Errorf("error getting environment: %w", err)
	}

	return &response.Environment, nil
}
