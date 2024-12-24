package kinde_client

import (
	"bytes"
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
		return output, fmt.Errorf("received status code %v for %v on url %v", res.StatusCode, req.Method, url)
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&output)

	if err != nil {
		return output, fmt.Errorf("failed to decode response: %w", err)
	}

	return output, nil
}

func doPostRequest[T any](c *Client, url string, body any) (T, error) {
	var output T

	b, err := json.Marshal(body)
	if err != nil {
		return output, fmt.Errorf("failed to marshal body: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		return output, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	req.Header.Add("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return output, fmt.Errorf("failed to make request: %w", err)
	} else if !isHttpStatusCodeSuccess(res.StatusCode) {
		var target map[string]any
		err = json.NewDecoder(res.Body).Decode(&target)
		if err != nil {
			return output, fmt.Errorf("received status code %v for %v on url %v, and failed to read response body: %w", res.StatusCode, req.Method, url, err)
		} else {
			return output, fmt.Errorf("received status code %v for %v on url %v: %+v - body was %s", res.StatusCode, req.Method, url, target, b)
		}
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
