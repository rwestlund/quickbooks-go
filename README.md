# quickbooks-go
![Build](https://github.com/nsotgui/quickbooks-go/workflows/Build/badge.svg)
[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](http://godoc.org/github.com/rwestlund/quickbooks-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/nsotgui/quickbooks-go)](https://goreportcard.com/report/github.com/nsotgui/quickbooks-go)

quickbooks-go is a Go library that provides access to Intuit's QuickBooks
Online API.

**NOTE:** This library is very incomplete. I just implemented the minimum for my
use case. Pull requests welcome :)

# Example
See [_main.go_](./examples/main.go)
```go
// Call the discovery api to get latest endpoints (recommended to update 1 time per day)
discoveryApis := auth.CallDiscoveryAPI(quickbooks.DiscoverySandboxEndpoint)
authClient := auth.Client{
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
	Client:   auth.GetHttpClient(*bearerToken),
	Endpoint: quickbooks.SandboxEndpoint,
	RealmID:  realmId,
}

// Make a request!
info, _ := qb.FetchCompanyInfo()
fmt.Println(info)

// Revoke the token, this should be done only if a user unsubscribe from your app
authClient.RevokeToken(bearerToken.RefreshToken)
```

# License
BSD-2-Clause
