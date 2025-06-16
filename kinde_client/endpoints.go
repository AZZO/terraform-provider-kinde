package kinde_client

import (
	"fmt"
)

func (c *Client) endpointApplicationsList() string {
	return fmt.Sprintf("%v/api/v1/applications", c.IssuerUrl)
}

func (c *Client) endpointApplication(id string) string {
	return fmt.Sprintf("%v/api/v1/applications/%v", c.IssuerUrl, id)
}

func (c *Client) endpointApisList() string {
	return fmt.Sprintf("%v/api/v1/apis", c.IssuerUrl)
}

func (c *Client) endpointApi(id string) string {
	return fmt.Sprintf("%v/api/v1/apis/%v", c.IssuerUrl, id)
}
