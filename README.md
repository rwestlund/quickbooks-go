# quickbooks-go
quickbooks-go is a Go library that provides access to Intuit's QuickBooks
Online API.

**NOTE:** This library is very incomplete. I just implemented the minimum for my
use case. Pull requests welcome :)

# Example
```
// Do this after you go through the normal OAuth process.
var client = oauth2.NewClient(ctx, tokenSource)

// Initialize the client handle.
var qb = quickbooks.Client{
    Client: client,
    Endpoint: quickbooks.SandboxEndpoint,
    RealmID: "some company account ID"'
}

// Make a request!
var companyInfo, err = qb.FetchCompanyInfo()
```

# License
BSD-2-Clause
