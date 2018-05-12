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
	"io/ioutil"
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
		var msg []byte
		msg, err = ioutil.ReadAll(res.Body)
		return nil, errors.New(strconv.Itoa(res.StatusCode) + " " + string(msg))
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
		var msg []byte
		msg, err = ioutil.ReadAll(res.Body)
		return nil, errors.New(strconv.Itoa(res.StatusCode) + " " + string(msg))
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

// CreateCustomer creates the given Customer on the QuickBooks server,
// returning the resulting Customer object.
func (c *Client) CreateCustomer(customer *Customer) (*Customer, error) {
	var u, err = url.Parse(string(c.Endpoint))
	if err != nil {
		return nil, err
	}
	u.Path = "/v3/company/" + c.RealmID + "/customer"
	var j []byte
	j, err = json.Marshal(customer)
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
		var msg []byte
		msg, err = ioutil.ReadAll(res.Body)
		return nil, errors.New(strconv.Itoa(res.StatusCode) + " " + string(msg))
	}

	var r struct {
		Customer Customer
		Time     time.Time
	}
	err = json.NewDecoder(res.Body).Decode(&r)
	return &r.Customer, err
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
		var msg []byte
		msg, err = ioutil.ReadAll(res.Body)
		return nil, errors.New(strconv.Itoa(res.StatusCode) + " " + string(msg))
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
		var msg []byte
		msg, err = ioutil.ReadAll(res.Body)
		return nil, errors.New(strconv.Itoa(res.StatusCode) + " " + string(msg))
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
		var msg []byte
		msg, err = ioutil.ReadAll(res.Body)
		return nil, errors.New(strconv.Itoa(res.StatusCode) + " " + string(msg))
	}

	var r struct {
		Invoice Invoice
		Time    time.Time
	}
	err = json.NewDecoder(res.Body).Decode(&r)
	return &r.Invoice, err
}

// DeleteInvoice deletes the given Invoice by ID and sync token from the
// QuickBooks server.
func (c *Client) DeleteInvoice(id, syncToken string) error {
	var u, err = url.Parse(string(c.Endpoint))
	if err != nil {
		return err
	}
	u.Path = "/v3/company/" + c.RealmID + "/invoice?operation=delete"
	var j []byte
	j, err = json.Marshal(struct {
		ID        string `json:"Id"`
		SyncToken string
	}{
		ID:        id,
		SyncToken: syncToken,
	})
	if err != nil {
		return err
	}
	var req *http.Request
	req, err = http.NewRequest("POST", u.String(), bytes.NewBuffer(j))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	var res *http.Response
	res, err = c.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	//var b, _ = ioutil.ReadAll(res.Body)
	//log.Println(string(b))

	// If the invoice was already deleted, QuickBooks returns 400 :(
	// The response looks like this:
	// {"Fault":{"Error":[{"Message":"Object Not Found","Detail":"Object Not Found : Something you're trying to use has been made inactive. Check the fields with accounts, customers, items, vendors or employees.","code":"610","element":""}],"type":"ValidationFault"},"time":"2018-03-20T20:15:59.571-07:00"}

	// This is slightly horrifying and not documented in their API. When this
	// happens we just return success; the goal of deleting it has been
	// accomplished, just not by us.
	if res.StatusCode == http.StatusBadRequest {
		var r struct {
			Fault struct {
				Error []struct {
					Message string
					Detail  string
					Code    string `json:"code"`
					Element string `json:"element"`
				}
				Type string `json:"type"`
			}
			Time time.Time `json:"time"`
		}
		err = json.NewDecoder(res.Body).Decode(&r)
		if err != nil {
			return err
		}
		if r.Fault.Error[0].Message == "Object Not Found" {
			return nil
		}
	}
	// TODO This could be better...
	if res.StatusCode != http.StatusOK {
		var msg []byte
		msg, err = ioutil.ReadAll(res.Body)
		return errors.New(strconv.Itoa(res.StatusCode) + " " + string(msg))
	}

	// TODO they send something back, but is it useful?
	return nil
}
