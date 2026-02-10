package examples

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/rwestlund/quickbooks-go"
	"github.com/stretchr/testify/require"
)

func TestReuseToken(t *testing.T) {
	_ = godotenv.Load("../.env") // Load .env file if it exists

	clientId := os.Getenv("QB_CLIENT_ID")
	clientSecret := os.Getenv("QB_CLIENT_SECRET")
	realmId := os.Getenv("QB_REALM_ID")
	refreshToken := os.Getenv("QB_REFRESH_TOKEN")
	accessToken := os.Getenv("QB_ACCESS_TOKEN")

	if clientId == "" || clientSecret == "" || realmId == "" || refreshToken == "" || accessToken == "" {
		t.Skip("Skipping TestReuseToken; required environment variables not set")
	}

	token := quickbooks.BearerToken{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}

	qbClient, err := quickbooks.NewClient(clientId, clientSecret, realmId, false, "", &token)
	require.NoError(t, err)

	// Make a request!
	info, err := qbClient.FindCompanyInfo()
	require.NoError(t, err)
	fmt.Println(info)
}
