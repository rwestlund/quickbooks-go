// Copyright (c) 2018, Randy Westlund. All rights reserved.
// This code is under the BSD-2-Clause license.

package quickbooks

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
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
	Taxable             bool    `json:",omitempty"`
	SalesTaxIncluded    bool    `json:",omitempty"`
	UnitPrice           float32 `json:",omitempty"`
	Type                string
	IncomeAccountRef    ReferenceType
	ExpenseAccountRef   ReferenceType
	PurchaseDesc        string  `json:",omitempty"`
	PurchaseTaxIncluded bool    `json:",omitempty"`
	PurchaseCost        float32 `json:",omitempty"`
	AssetAccountRef     ReferenceType
	TrackQtyOnHand      bool `json:",omitempty"`
	//InvStartDate time.Time
	QtyOnHand          float32       `json:",omitempty"`
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
