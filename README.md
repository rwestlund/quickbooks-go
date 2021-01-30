# quickbooks-go
![Build](https://github.com/nsotgui/quickbooks-go/workflows/Build/badge.svg)
[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](http://godoc.org/github.com/nsotgui/quickbooks-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/nsotgui/quickbooks-go)](https://goreportcard.com/report/github.com/nsotgui/quickbooks-go)

quickbooks-go is a Go library that provides access to Intuit's QuickBooks
Online API.

**NOTE:** This library is very incomplete. I just implemented the minimum for my
use case. Pull requests welcome :)

# Example

## Authorization flow

See [_auth_flow_test.go_](./examples/auth_flow_test.go)
```go
clientId     := "<your-client-id>"
clientSecret := "<your-client-secret>"
realmId      := "<realm-id>"

qbClient, _ := quickbooks.NewQuickbooksClient(clientId, clientSecret, realmId, false, nil)

// To do first when you receive the authorization code from quickbooks callback
authorizationCode := "<received-from-callback>"
bearerToken, _ := qbClient.RetrieveBearerToken(authorizationCode)
// Save the bearer token inside a db

// When the token expire, you can use the following function
bearerToken, _ = qbClient.RefreshToken(bearerToken.RefreshToken)

// Make a request!
info, _ := qbClient.FetchCompanyInfo()
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

qbClient, _ := quickbooks.NewQuickbooksClient(clientId, clientSecret, realmId, false, &token)

// Make a request!
info, _ := qbClient.FetchCompanyInfo()
fmt.Println(info)
```

# License
BSD-2-Clause
