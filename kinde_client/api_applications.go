package kinde_client

import (
	"context"
	"fmt"
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

func (c *Client) UpdateApplication(ctx context.Context, id string, name string, logoutUris []string, redirectUris []string) (ApplicationResource, error) {
	url := c.endpointApplication(id)

	body := requestApplicationUpdate{
		Name:         name,
		LogoutUris:   logoutUris,
		RedirectUris: redirectUris,
	}

	app, err := doPostRequest[responseApplicationGet](c, ctx, url, body)
	if err != nil {
		return ApplicationResource{}, err
	}

	return app.Application, nil
}

func (c *Client) DeleteApplication(ctx context.Context, id string) error {
	url := c.endpointApplication(id)

	_, err := c.doRequest(ctx, "DELETE", url, nil, 3)
	return err
}

type Callbacks struct {
	LogoutUris   []string `json:"logout_uris"`
	RedirectUris []string `json:"redirect_uris"`
}

type responseCallbacks struct {
	Code      string    `json:"code"`
	Message   string    `json:"message"`
	Callbacks Callbacks `json:"callbacks"`
}

func (c *Client) GetCallbacks(ctx context.Context, id string) (Callbacks, error) {
	url := fmt.Sprintf("%v/api/v1/applications/%v/callbacks", c.IssuerUrl, id)

	callbacks, err := doGetRequest[responseCallbacks](c, ctx, url)
	if err != nil {
		return Callbacks{}, err
	}

	return callbacks.Callbacks, nil
}

func (c *Client) UpdateCallbacks(ctx context.Context, id string, callbacks Callbacks) (Callbacks, error) {
	url := fmt.Sprintf("%v/api/v1/applications/%v/callbacks", c.IssuerUrl, id)

	updated, err := doPostRequest[responseCallbacks](c, ctx, url, callbacks)
	if err != nil {
		return Callbacks{}, err
	}

	return updated.Callbacks, nil
}
