package kinde_client

import (
	"context"
)

// GetEnvironment retrieves information about a specific environment
func (c *Client) GetEnvironment(ctx context.Context, environmentId string) (EnvironmentResource, error) {
	response, err := doGetRequest[responseEnvironmentGet](c, ctx, c.endpointEnvironmentGet(environmentId))
	if err != nil {
		return EnvironmentResource{}, err
	}

	return response.Environment, nil
}
