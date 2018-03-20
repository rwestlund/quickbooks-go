// Copyright (c) 2018, Randy Westlund. All rights reserved.
// This code is under the BSD-2-Clause license.

package quickbooks

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
