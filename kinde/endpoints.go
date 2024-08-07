package kinde

import (
	"fmt"
)

func endpointApplicationSpecific(config *Config, id string) string {
	return fmt.Sprintf("%v/api/v1/applications/%v", config.api_url, id)
}

func endpoitnApplicationGeneric(config *Config) string {
	return fmt.Sprintf("%v/api/v1/applications", config.api_url)
}

func EndpointApplicationGet(config *Config, id string) string {
	return endpointApplicationSpecific(config, id)
}

func EndpointApplicationCreate(config *Config) string {
	return endpoitnApplicationGeneric(config)
}

func EndpointApplicationUpdate(config *Config, id string) string {
	return endpointApplicationSpecific(config, id)
}

func EndpointApplicationDelete(config *Config, id string) string {
	return endpointApplicationSpecific(config, id)
}
