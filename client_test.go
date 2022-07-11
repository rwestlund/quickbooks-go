package quickbooks

import (
	"errors"
	"os"
)

// Client used for testing only
func getClient() (*Client, error) {
	clientId := os.Getenv("CLIENT_ID")
	if clientId == "" {
		return nil, errors.New("CLIENT_ID not defined")
	}
	clientSecret := os.Getenv("CLIENT_SECRET")
	if clientSecret == "" {
		return nil, errors.New("CLIENT_SECRET not defined")
	}
	realmID := os.Getenv("REALM_ID")
	if realmID == "" {
		return nil, errors.New("REALM_ID not defined")
	}
	accessToken := os.Getenv("ACCESS_TOKEN")
	if accessToken == "" {
		return nil, errors.New("ACCESS_TOKEN not defined")
	}
	token := BearerToken{
		AccessToken: accessToken,
		TokenType:   "Bearer",
	}
	return NewQuickbooksClient(clientId, clientSecret, realmID, false, &token)
}
