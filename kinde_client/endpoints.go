package kinde_client

import (
	"fmt"
)

func (c *Client) endpointApplicationsList() string {
	return fmt.Sprintf("%v/api/v1/applications", c.ApiUrl)
}

func (c *Client) endpointApplication(id string) string {
	return fmt.Sprintf("%v/api/v1/applications/%v", c.ApiUrl, id)
}
