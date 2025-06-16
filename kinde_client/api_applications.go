package kinde_client

import (
	"context"
)

func (c *Client) GetApplication(ctx context.Context, id string) (ApplicationResource, error) {
	url := c.endpointApplication(id)

	response, err := doGetRequest[responseApplicationGet](c, ctx, url)
	if err != nil {
		return ApplicationResource{}, err
	}

	return response.Application, nil
}

func (c *Client) GetApplications(ctx context.Context) ([]ApplicationResource, error) {
	url := c.endpointApplicationsList()

	response, err := doGetRequest[responseApplicationGet](c, ctx, url)
	if err != nil {
		return nil, err
	}

	return []ApplicationResource{response.Application}, nil
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

	response, err := doPostRequest[responseApplicationGet](c, ctx, url, body)
	if err != nil {
		return ApplicationResource{}, err
	}

	return response.Application, nil
}

func (c *Client) DeleteApplication(ctx context.Context, id string) error {
	_, err := doDeleteRequest[responseApplicationDelete](c, ctx, c.endpointApplication(id))
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

func (c *Client) GetApplicationCallbacks(ctx context.Context, id string) (Callbacks, error) {
	callbacks, err := doGetRequest[responseCallbacks](c, ctx, c.endpointApplicationCallbacksGet(id))
	if err != nil {
		return Callbacks{}, err
	}

	return callbacks.Callbacks, nil
}

func (c *Client) UpdateApplicationCallbacks(ctx context.Context, id string, callbacks Callbacks) (Callbacks, error) {
	updated, err := doPostRequest[responseCallbacks](c, ctx, c.endpointApplicationCallbacksUpdate(id), callbacks)
	if err != nil {
		return Callbacks{}, err
	}

	return updated.Callbacks, nil
}
