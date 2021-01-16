package examples

import (
	"fmt"
	"github.com/nsotgui/quickbooks-go"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAuthorizationFlow(t *testing.T) {
	clientId := "<your-client-id>"
	clientSecret := "<your-client-secret>"
	realmId := "<realm-id>"

	qbClient, err := quickbooks.NewQuickbooksClient(clientId, clientSecret, realmId, false, nil)
	require.NoError(t, err)

	// To do first when you receive the authorization code from quickbooks callback
	authorizationCode := "<received-from-callback>"
	bearerToken, err := qbClient.RetrieveBearerToken(authorizationCode)
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
