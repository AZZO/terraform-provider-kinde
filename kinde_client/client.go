package kinde_client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	HTTPClient *http.Client
	Token      string
	Auth       authStruct
	IssuerUrl  string
	ApiUrl     string
}

type authStruct struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func NewClient(ctx context.Context, issuerUrl, clientId, clientSecret string) (*Client, error) {
	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		IssuerUrl:  issuerUrl,
	}

	c.Auth = authStruct{
		ClientId:     clientId,
		ClientSecret: clientSecret,
	}

	tok, err := c.GetToken(ctx)

	if err != nil {
		return nil, err
	}

	c.Token = *tok

	return &c, nil
}

func doGetRequest[T any](c *Client, url string) (T, error) {
	var output T

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return output, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return output, fmt.Errorf("failed to make request: %w", err)
	} else if !isHttpStatusCodeSuccess(res.StatusCode) {
		return output, fmt.Errorf("received status code %v for url %v", res.StatusCode, url)
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&output)

	if err != nil {
		return output, fmt.Errorf("failed to decode response: %w", err)
	}

	return output, nil
}

func isHttpStatusCodeSuccess(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}
