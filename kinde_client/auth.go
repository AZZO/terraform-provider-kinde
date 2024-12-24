package kinde_client

import (
	"context"
	"fmt"
	"time"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func (c *Client) GetToken(ctx context.Context) (*string, error) {
	apiBaseUrl := c.IssuerUrl
	tokenIssuer := apiBaseUrl
	apiAudience := tokenIssuer + "/api"
	tokenUrl := tokenIssuer + "/oauth2/token"

	config := &clientcredentials.Config{
		ClientID:     c.Auth.ClientId,
		ClientSecret: c.Auth.ClientSecret,
		TokenURL:     tokenUrl,
		Scopes:       []string{},
		EndpointParams: map[string][]string{
			"audience": {apiAudience},
		},
	}

	context, cancelTimer := context.WithTimeout(ctx, 5*time.Second)
	defer cancelTimer()
	m2mToken, err := config.Token(context)

	if err != nil {
		return nil, err
	}

	err = verifyToken(ctx, m2mToken, tokenIssuer, apiAudience)

	if err != nil {
		return nil, err
	}

	return &m2mToken.AccessToken, nil
}

func verifyToken(ctx context.Context, m2mToken *oauth2.Token, tokenIssuer string, audience string) error {
	jwksURL := fmt.Sprintf("%v/.well-known/jwks", tokenIssuer)

	jwks, err := keyfunc.NewDefaultCtx(ctx, []string{jwksURL}) // Context is used to end the refresh goroutine.
	if err != nil {
		return err
	}

	parsedToken, err := jwt.Parse(m2mToken.AccessToken, jwks.Keyfunc,
		jwt.WithValidMethods([]string{"RS256"}), // verifying the signing algorithm
		jwt.WithIssuer(tokenIssuer),             // verifying the token issuer
		jwt.WithAudience(audience))              // verifying that the token is for correct audience

	if err != nil {
		return fmt.Errorf("error verifying token %v", err)
	}

	tflog.Trace(ctx, "m2m token", map[string]interface{}{"claims": parsedToken.Claims})

	return nil
}
