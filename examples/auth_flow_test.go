package examples

import (
	"fmt"
	"testing"

	"github.com/neonmoose/quickbooks-go"
	"github.com/stretchr/testify/require"
)

func TestAuthorizationFlow(t *testing.T) {
	clientId := "<your-client-id>"
	clientSecret := "<your-client-secret>"
	realmId := "<realm-id>"

	qbClient, err := quickbooks.NewQuickbooksClient(clientId, clientSecret, realmId, false, nil)
	require.NoError(t, err)

	// Get the authorization url
	qbClient.GetAuthorizationUrl(quickbooks.AccountingScope, "random-string", "https://localhost/redirect_uri")

	// To do first when you receive the authorization code from quickbooks callback
	authorizationCode := "<received-from-callback>"
	redirectURI := "https://developer.intuit.com/v2/OAuth2Playground/RedirectUrl"
	bearerToken, err := qbClient.RetrieveBearerToken(authorizationCode, redirectURI)
	require.NoError(t, err)
	// Save the bearer token inside a db

	// When the token expire, you can use the following function
	bearerToken, err = qbClient.RefreshToken(bearerToken.RefreshToken)
	require.NoError(t, err)

	// Make a request!
	info, err := qbClient.FetchCompanyInfo()
	require.NoError(t, err)
	fmt.Println(info)

	// Revoke the token, this should be done only if a user unsubscribe from your app
	qbClient.RevokeToken(bearerToken.RefreshToken)
}
