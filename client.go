// Copyright (c) 2018, Randy Westlund. All rights reserved.
// This code is under the BSD-2-Clause license.

/*
Package quickbooks provides access to Intuit's QuickBooks Online API.

NOTE: This library is very incomplete. I just implemented the minimum for my
use case. Pull requests welcome :)

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
*/
package quickbooks

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// EndpointURL specifies the endpoint to connect to.
type EndpointURL string

const (
	// ProductionEndpoint is for live apps.
	ProductionEndpoint EndpointURL = "https://quickbooks.api.intuit.com"
	// SandboxEndpoint is for testing.
	SandboxEndpoint EndpointURL = "https://sandbox-quickbooks.api.intuit.com"
)

// Client is your handle to the QuickBooks API.
type Client struct {
	// Get this from oauth2.NewClient().
	Client *http.Client
	// Set to ProductionEndpoint or SandboxEndpoint.
	Endpoint EndpointURL
	// The account ID you're connecting to.
	RealmID string
}

// FetchCompanyInfo returns the QuickBooks CompanyInfo object. This is a good
// test to check whether you're connected.
func (c *Client) FetchCompanyInfo() (*CompanyInfo, error) {
	var u, err = url.Parse(string(c.Endpoint))
	if err != nil {
		return nil, err
	}
	u.Path = "/v3/company/" + c.RealmID + "/companyinfo/" + c.RealmID
	var req *http.Request
	req, err = http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	var res *http.Response
	res, err = c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// TODO This could be better...
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(strconv.Itoa(res.StatusCode))
	}

	var r struct {
		CompanyInfo CompanyInfo
		Time        time.Time
	}
	err = json.NewDecoder(res.Body).Decode(&r)
	return &r.CompanyInfo, err
}

// FetchCustomers gets the full list of Customers in the QuickBooks account.
func (c *Client) FetchCustomers() ([]Customer, error) {
	var u, err = url.Parse(string(c.Endpoint))
	if err != nil {
		return nil, err
	}
	u.Path = "/v3/company/" + c.RealmID + "/query"

	var v = url.Values{}
	v.Add("query", "SELECT * FROM Customer")
	u.RawQuery = v.Encode()
	var req *http.Request
	req, err = http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	var res *http.Response
	res, err = c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// TODO This could be better...
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(strconv.Itoa(res.StatusCode))
	}

	var r struct {
		QueryResponse struct {
			Customer      []Customer
			StartPosition int
			MaxResults    int
		}
	}
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	// Make sure we don't return nil if there are no customers.
	if r.QueryResponse.Customer == nil {
		r.QueryResponse.Customer = make([]Customer, 0)
	}
	return r.QueryResponse.Customer, nil
}

// FetchItems returns the list of Items in the QuickBooks account. These are
// basically product types, and you need them to create invoices.
func (c *Client) FetchItems() ([]Item, error) {
	var u, err = url.Parse(string(c.Endpoint))
	if err != nil {
		return nil, err
	}
	u.Path = "/v3/company/" + c.RealmID + "/query"

	var v = url.Values{}
	v.Add("query", "SELECT * FROM Item")
	u.RawQuery = v.Encode()
	var req *http.Request
	req, err = http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	var res *http.Response
	res, err = c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// TODO This could be better...
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(strconv.Itoa(res.StatusCode))
	}

	var r struct {
		QueryResponse struct {
			Item          []Item
			StartPosition int
			MaxResults    int
		}
	}
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	// Make sure we don't return nil if there are no items.
	if r.QueryResponse.Item == nil {
		r.QueryResponse.Item = make([]Item, 0)
	}
	return r.QueryResponse.Item, nil
}

// FetchItem returns just one particular Item from QuickBooks, by ID.
func (c *Client) FetchItem(id string) (*Item, error) {
	var u, err = url.Parse(string(c.Endpoint))
	if err != nil {
		return nil, err
	}
	u.Path = "/v3/company/" + c.RealmID + "/item/" + id

	var req *http.Request
	req, err = http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	var res *http.Response
	res, err = c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// TODO This could be better...
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(strconv.Itoa(res.StatusCode))
	}

	var r struct {
		Item Item
		Time time.Time
	}
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r.Item, nil
}

// CreateInvoice creates the given Invoice on the QuickBooks server, returning
// the resulting Invoice object.
func (c *Client) CreateInvoice(inv *Invoice) (*Invoice, error) {
	var u, err = url.Parse(string(c.Endpoint))
	if err != nil {
		return nil, err
	}
	u.Path = "/v3/company/" + c.RealmID + "/invoice"
	var j []byte
	j, err = json.Marshal(inv)
	if err != nil {
		return nil, err
	}
	var req *http.Request
	req, err = http.NewRequest("POST", u.String(), bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	var res *http.Response
	res, err = c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// TODO This could be better...
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(strconv.Itoa(res.StatusCode))
	}

	var r struct {
		Invoice Invoice
		Time    time.Time
	}
	err = json.NewDecoder(res.Body).Decode(&r)
	return &r.Invoice, err
}
