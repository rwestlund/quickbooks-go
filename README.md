# quickbooks-go
![Build](https://github.com/rwestlund/quickbooks-go/workflows/Build/badge.svg)
[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](http://godoc.org/github.com/rwestlund/quickbooks-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/rwestlund/quickbooks-go)](https://goreportcard.com/report/github.com/rwestlund/quickbooks-go)

quickbooks-go is a Go library that provides access to Intuit's QuickBooks
Online API.

**NOTE:** This library is incomplete. I implemented the minimum for my
use case. Pull requests welcome :)

# Example

## Authorization flow

Before you can initialize the client, you'll need to obtain an authorization code. You can see an example of this from QuickBooks' [OAuth Playground](https://developer.intuit.com/app/developer/playground).

See [_auth_flow_test.go_](./examples/auth_flow_test.go)
```go
clientId     := "<your-client-id>"
clientSecret := "<your-client-secret>"
realmId      := "<realm-id>"

qbClient, err := quickbooks.NewClient(clientId, clientSecret, realmId, false, "", nil)
if err != nil {
	log.Fatalln(err)
}

// To do first when you receive the authorization code from quickbooks callback
authorizationCode := "<received-from-callback>"
redirectURI := "https://developer.intuit.com/v2/OAuth2Playground/RedirectUrl"

bearerToken, err := qbClient.RetrieveBearerToken(authorizationCode, redirectURI)
if err != nil {
	log.Fatalln(err)
}
// Save the bearer token inside a db

// When the token expire, you can use the following function
bearerToken, err = qbClient.RefreshToken(bearerToken.RefreshToken)
if err != nil {
	log.Fatalln(err)
}

// Make a request!
info, err := qbClient.FindCompanyInfo()
if err != nil {
	log.Fatalln(err)
}

fmt.Println(info)

// Revoke the token, this should be done only if a user unsubscribe from your app
qbClient.RevokeToken(bearerToken.RefreshToken)
```

## Re-using tokens

See [_reuse_token_test.go_](./examples/reuse_token_test.go)
```go
clientId     := "<your-client-id>"
clientSecret := "<your-client-secret>"
realmId      := "<realm-id>"

token := quickbooks.BearerToken{
	RefreshToken:           "<saved-refresh-token>",
	AccessToken:            "<saved-access-token>",
}

qbClient, err := quickbooks.NewClient(clientId, clientSecret, realmId, false, "", &token)
if err != nil {
	log.Fatalln(err)
}

// Make a request!
info, err := qbClient.FindCompanyInfo()
if err != nil {
	log.Fatalln(err)
}

fmt.Println(info)
```

# License
BSD-2-Clause
