package examples

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/rwestlund/quickbooks-go"
	"github.com/stretchr/testify/require"
)

func TestAuthorizationFlow(t *testing.T) {
	_ = godotenv.Load("../.env") // Load .env file if it exists

	clientId := os.Getenv("QB_CLIENT_ID")
	clientSecret := os.Getenv("QB_CLIENT_SECRET")
	realmId := os.Getenv("QB_REALM_ID")
	authorizationCode := os.Getenv("QB_AUTH_CODE")

	if clientId == "" || clientSecret == "" || realmId == "" || authorizationCode == "" {
		t.Skip("Skipping TestAuthorizationFlow; QB_CLIENT_ID, QB_CLIENT_SECRET, QB_REALM_ID, or QB_AUTH_CODE not set")
	}

	qbClient, err := quickbooks.NewClient(clientId, clientSecret, realmId, false, "", nil)
	require.NoError(t, err)

	// To do first when you receive the authorization code from quickbooks callback
	redirectURI := "https://developer.intuit.com/v2/OAuth2Playground/RedirectUrl"
	bearerToken, err := qbClient.RetrieveBearerToken(authorizationCode, redirectURI)
	require.NoError(t, err)
	// Save the bearer token inside a db

	// When the token expire, you can use the following function
	bearerToken, err = qbClient.RefreshToken(bearerToken.RefreshToken)
	require.NoError(t, err)

	// Make a request!
	info, err := qbClient.FindCompanyInfo()
	require.NoError(t, err)
	fmt.Println(info)

	// Revoke the token, this should be done only if a user unsubscribe from your app
	require.NoError(t, qbClient.RevokeToken(bearerToken.RefreshToken))
}
