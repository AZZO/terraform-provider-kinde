package kinde_client

func (c *Client) GetApplication(id string) (ApplicationResource, error) {

	url := c.endpointApplication(id)

	app, err := doGetRequest[responseApplicationGet](c, url)

	if err != nil {
		return ApplicationResource{}, err
	}

	return app.Application, nil
}
