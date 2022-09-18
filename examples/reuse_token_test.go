package examples

import (
	"fmt"
	"testing"

	"github.com/rwestlund/quickbooks-go"
	"github.com/stretchr/testify/require"
)

func TestReuseToken(t *testing.T) {
	clientId := "<your-client-id>"
	clientSecret := "<your-client-secret>"
	realmId := "<realm-id>"

	token := quickbooks.BearerToken{
		RefreshToken: "<saved-refresh-token>",
		AccessToken:  "<saved-access-token>",
	}

	qbClient, err := quickbooks.NewClient(clientId, clientSecret, realmId, false, "", &token)
	require.NoError(t, err)

	// Make a request!
	info, err := qbClient.FindCompanyInfo()
	require.NoError(t, err)
	fmt.Println(info)
}
