package kinde_client

import (
	"context"
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

func (c *Client) DeleteApi(ctx context.Context, id string) error {
	url := c.endpointApi(id)

	_, err := c.doRequest(ctx, "DELETE", url, nil, 3)
	return err
}
