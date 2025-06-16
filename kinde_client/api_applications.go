package kinde_client

import (
	"context"
)

func (c *Client) GetApplication(ctx context.Context, id string) (ApplicationResource, error) {
	url := c.endpointApplication(id)

	app, err := doGetRequest[responseApplicationGet](c, ctx, url)
	if err != nil {
		return ApplicationResource{}, err
	}

	return app.Application, nil
}

func (c *Client) CreateApplication(ctx context.Context, name string, appType string) (ApplicationResource, error) {
	url := c.endpointApplicationsList()

	body := requestApplicationCreate{
		Name: name,
		Type: appType,
	}

	app, err := doPostRequest[responseApplicationCreate](c, ctx, url, body)
	if err != nil {
		return ApplicationResource{}, err
	}

	return ApplicationResource{
		Name:         name,
		Type:         appType,
		Id:           app.Application.Id,
		ClientId:     app.Application.ClientId,
		ClientSecret: app.Application.ClientSecret,
	}, nil
}
