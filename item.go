// Copyright (c) 2018, Randy Westlund. All rights reserved.
// This code is under the BSD-2-Clause license.

package quickbooks

import (
	"encoding/json"
	"errors"
	"strconv"
)

// Item represents a QuickBooks Item object (a product type).
type Item struct {
	Id          string   `json:"Id,omitempty"`
	SyncToken   string   `json:",omitempty"`
	MetaData    MetaData `json:",omitempty"`
	Name        string
	SKU         string `json:"Sku,omitempty"`
	Description string `json:",omitempty"`
	Active      bool   `json:",omitempty"`
	// SubItem
	// ParentRef
	// Level
	// FullyQualifiedName
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
	// InvStartDate Date
	QtyOnHand          json.Number   `json:",omitempty"`
	SalesTaxCodeRef    ReferenceType `json:",omitempty"`
	PurchaseTaxCodeRef ReferenceType `json:",omitempty"`
}

func (c *Client) CreateItem(item *Item) (*Item, error) {
	var resp struct {
		Item Item
		Time Date
	}

	if err := c.post("item", item, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.Item, nil
}

// FindItems gets the full list of Items in the QuickBooks account.
func (c *Client) FindItems() ([]Item, error) {
	var resp struct {
		QueryResponse struct {
			Items         []Item `json:"Item"`
			MaxResults    int
			StartPosition int
			TotalCount    int
		}
	}

	if err := c.query("SELECT COUNT(*) FROM Item", &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.TotalCount == 0 {
		return nil, errors.New("no items could be found")
	}

	items := make([]Item, 0, resp.QueryResponse.TotalCount)

	for i := 0; i < resp.QueryResponse.TotalCount; i += queryPageSize {
		query := "SELECT * FROM Item ORDERBY Id STARTPOSITION " + strconv.Itoa(i+1) + " MAXRESULTS " + strconv.Itoa(queryPageSize)

		if err := c.query(query, &resp); err != nil {
			return nil, err
		}

		if resp.QueryResponse.Items == nil {
			return nil, errors.New("no items could be found")
		}

		items = append(items, resp.QueryResponse.Items...)
	}

	return items, nil
}

// FindItemById returns an item with a given Id.
func (c *Client) FindItemById(id string) (*Item, error) {
	var resp struct {
		Item Item
		Time Date
	}

	if err := c.get("item/"+id, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.Item, nil
}

// QueryItems accepts an SQL query and returns all items found using it
func (c *Client) QueryItems(query string) ([]Item, error) {
	var resp struct {
		QueryResponse struct {
			Items         []Item `json:"Item"`
			StartPosition int
			MaxResults    int
		}
	}

	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.Items == nil {
		return nil, errors.New("could not find any items")
	}

	return resp.QueryResponse.Items, nil
}

// UpdateItem updates the item
func (c *Client) UpdateItem(item *Item) (*Item, error) {
	if item.Id == "" {
		return nil, errors.New("missing item id")
	}

	existingItem, err := c.FindItemById(item.Id)
	if err != nil {
		return nil, err
	}

	item.SyncToken = existingItem.SyncToken

	payload := struct {
		*Item
		Sparse bool `json:"sparse"`
	}{
		Item:   item,
		Sparse: true,
	}

	var itemData struct {
		Item Item
		Time Date
	}

	if err = c.post("item", payload, &itemData, nil); err != nil {
		return nil, err
	}

	return &itemData.Item, err
}
