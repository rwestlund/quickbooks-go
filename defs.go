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
	// DiscoverySandboxEndpoint is for testing.
	DiscoverySandboxEndpoint EndpointURL = "https://developer.api.intuit.com/.well-known/openid_sandbox_configuration"
	// DiscoveryProductionEndpoint is for live apps.
	DiscoveryProductionEndpoint EndpointURL = "https://developer.api.intuit.com/.well-known/openid_configuration"
)

const queryPageSize = 1000

const minorVersion = "52"
