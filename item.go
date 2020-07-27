// Copyright (c) 2018, Randy Westlund. All rights reserved.
// This code is under the BSD-2-Clause license.

package quickbooks

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

// Item represents a QuickBooks Item object (a product type).
type Item struct {
	ID        string `json:"Id,omitempty"`
	SyncToken string `json:",omitempty"`
	//MetaData
	Name        string
	SKU         string `json:"Sku,omitempty"`
	Description string `json:",omitempty"`
	Active      bool   `json:",omitempty"`
	//SubItem
	//ParentRef
	//Level
	//FullyQualifiedName
	Taxable             bool        `json:",omitempty"`
	SalesTaxIncluded    bool        `json:",omitempty"`
	UnitPrice           json.Number `json:",omitempty"`
	Type                string
	IncomeAccountRef    ReferenceType
	ExpenseAccountRef   ReferenceType
	PurchaseDesc        string      `json:",omitempty"`
	PurchaseTaxIncluded bool        `json:",omitempty"`
	PurchaseCost        json.Number `json:",omitempty"`
	AssetAccountRef     ReferenceType
	TrackQtyOnHand      bool `json:",omitempty"`
	//InvStartDate Date
	QtyOnHand          json.Number   `json:",omitempty"`
	SalesTaxCodeRef    ReferenceType `json:",omitempty"`
	PurchaseTaxCodeRef ReferenceType `json:",omitempty"`
}

// FetchItems returns the list of Items in the QuickBooks account. These are
// basically product types, and you need them to create invoices.
func (c *Client) FetchItems() ([]Item, error) {
	var r struct {
		QueryResponse struct {
			Item          []Item
			StartPosition int
			MaxResults    int
		}
	}
	err := c.query("SELECT * FROM Item MAXRESULTS "+strconv.Itoa(queryPageSize), &r)
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
	var v = url.Values{}
	v.Add("minorversion", minorVersion)
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

	if res.StatusCode != http.StatusOK {
		return nil, parseFailure(res)
	}

	var r struct {
		Item Item
		Time Date
	}
	err = json.NewDecoder(res.Body).Decode(&r)
	if err != nil {
		return nil, err
	}
	return &r.Item, nil
}
