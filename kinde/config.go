package kinde

import "net/http"

type Config struct {
	client  *http.Client
	api_url string
}
