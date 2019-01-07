// Copyright (c) 2018, Randy Westlund. All rights reserved.
// This code is under the BSD-2-Clause license.

package quickbooks

// EndpointURL specifies the endpoint to connect to.
type EndpointURL string

const (
	// ProductionEndpoint is for live apps.
	ProductionEndpoint EndpointURL = "https://quickbooks.api.intuit.com"
	// SandboxEndpoint is for testing.
	SandboxEndpoint EndpointURL = "https://sandbox-quickbooks.api.intuit.com"
)

const queryPageSize = 1000
