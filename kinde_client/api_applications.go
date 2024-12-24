package kinde_client

func (c *Client) GetApplication(id string) (ApplicationResource, error) {

	url := c.endpointApplication(id)

	app, err := doGetRequest[responseApplicationGet](c, url)

	if err != nil {
		return ApplicationResource{}, err
	}

	return app.Application, nil
}

func (c *Client) CreateApplication(name string, appType string) (ApplicationResource, error) {
	url := c.endpointApplicationsList()

	body := requestApplicationCreate{
		Name: name,
		Type: appType,
	}

	app, err := doPostRequest[responseApplicationCreate](c, url, body)

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
