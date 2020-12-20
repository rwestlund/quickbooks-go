package auth

import (
	"context"
	"golang.org/x/oauth2"
	"net/http"
)

type Client struct {
	// The set of quickbooks APIs
	DiscoveryAPI DiscoveryAPI
	// The client ID
	ClientId string
	// The client Secret
	ClientSecret string
}

func GetHttpClient(bearerToken BearerToken) *http.Client {
	ctx := context.Background()
	token := oauth2.Token{
		AccessToken: bearerToken.AccessToken,
		TokenType:   "Bearer",
	}
	return oauth2.NewClient(ctx, oauth2.StaticTokenSource(&token))
}
