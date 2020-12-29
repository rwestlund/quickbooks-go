package examples

import (
	"fmt"

	"github.com/rwestlund/quickbooks-go"
)

const (
	clientId     string = "<your-client-id>"
	clientSecret string = "<your-client-secret>"
)

func main() {

	// Call the discovery api to get latest endpoints (recommended to update 1 time per day)
	discoveryApis := quickbooks.CallDiscoveryAPI(quickbooks.DiscoverySandboxEndpoint)
	authClient := quickbooks.AuthClient{
		DiscoveryAPI: *discoveryApis,
		ClientId:     clientId,
		ClientSecret: clientSecret,
	}

	// To do first when you receive the authorization code from quickbooks callback
	authorizationCode := "<received-from-callback>"
	bearerToken, _ := authClient.RetrieveBearerToken(authorizationCode)
	// Save the bearer token inside a db

	// When the token expire, you can use the following function
	bearerToken, _ = authClient.RefreshToken(bearerToken.RefreshToken)

	// Initialize the quickbook client handle.
	realmId := "<realm-id>"
	var qb = quickbooks.Client{
		Client:   quickbooks.GetHttpClient(*bearerToken),
		Endpoint: quickbooks.SandboxEndpoint,
		RealmID:  realmId,
	}

	// Make a request!
	info, _ := qb.FetchCompanyInfo()
	fmt.Println(info)

	// Revoke the token, this should be done only if a user unsubscribe from your app
	authClient.RevokeToken(bearerToken.RefreshToken)
}
