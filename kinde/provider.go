package kinde

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"golang.org/x/oauth2/clientcredentials"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"issuer_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("KINDE_ISSUER_URL", nil),
				Description: "Root URL for the Kinde environment",
			},
			"client_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("KINDE_CLIENT_ID", nil),
				Description: "Client ID for the management API application",
			},
			"client_secret": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("KINDE_CLIENT_SECRET", nil),
				Description: "Client secret for the management API application",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"kinde_application": ResourceApplication(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"kinde_application": DataSourceApplication(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	apiBaseUrl := d.Get("issuer_url").(string)
	tokenIssuer := apiBaseUrl
	apiAudience := tokenIssuer + "/api"
	tokenUrl := tokenIssuer + "/oauth2/token"

	config := &clientcredentials.Config{
		ClientID:     d.Get("client_id").(string),
		ClientSecret: d.Get("client_secret").(string),
		TokenURL:     tokenUrl,
		Scopes:       []string{},
		EndpointParams: map[string][]string{
			"audience": {apiAudience},
		},
	}

	context, _ := context.WithTimeout(context.Background(), 5*time.Second)
	m2mToken, err := config.Token(context)
	if err != nil {
		return nil, err
	}

	VerifyToken(m2mToken, tokenIssuer, apiAudience)
	client := config.Client(context)

	return &Config{
		client:  client,
		api_url: apiBaseUrl,
	}, nil
}
