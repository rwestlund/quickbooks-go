// Copyright (c) 2018, Randy Westlund. All rights reserved.
// This code is under the BSD-2-Clause license.

package quickbooks

import null "gopkg.in/guregu/null.v3"

// Invoice represents a QuickBooks Invoice object.
type Invoice struct {
	ID        string `json:"Id,omitempty"`
	SyncToken string `json:",omitempty"`
	//MetaData
	//CustomField
	DocNumber string `json:",omitempty"`
	//TxnDate   time.Time `json:",omitempty"`
	//DepartmentRef
	PrivateNote string `json:",omitempty"`
	//LinkedTxn
	Line []SalesItemLine
	//TxnTaxDetail
	CustomerRef  ReferenceType
	CustomerMemo MemoRef         `json:",omitempty"`
	BillAddr     PhysicalAddress `json:",omitempty"`
	ShipAddr     PhysicalAddress `json:",omitempty"`
	ClassRef     ReferenceType   `json:",omitempty"`
	SalesTermRef ReferenceType   `json:",omitempty"`
	//DueDate      time.Time       `json:",omitempty"`
	//GlobalTaxCalculation
	ShipMethodRef ReferenceType `json:",omitempty"`
	//ShipDate      time.Time     `json:",omitempty"`
	TrackingNum string  `json:",omitempty"`
	TotalAmt    float32 `json:",omitempty"`
	//CurrencyRef
	ExchangeRate          float32      `json:",omitempty"`
	HomeAmtTotal          float32      `json:",omitempty"`
	HomeBalance           float32      `json:",omitempty"`
	ApplyTaxAfterDiscount bool         `json:",omitempty"`
	PrintStatus           string       `json:",omitempty"`
	EmailStatus           string       `json:",omitempty"`
	BillEmail             EmailAddress `json:",omitempty"`
	BillEmailCC           EmailAddress `json:"BillEmailCc,omitempty"`
	BillEmailBCC          EmailAddress `json:"BillEmailBcc,omitempty"`
	//DeliveryInfo
	Balance                      float32       `json:",omitempty"`
	TxnSource                    string        `json:",omitempty"`
	AllowOnlineCreditCardPayment bool          `json:",omitempty"`
	AllowOnlineACHPayment        bool          `json:",omitempty"`
	Deposit                      float32       `json:",omitempty"`
	DepositToAccountRef          ReferenceType `json:",omitempty"`
}

// SalesItemLine ...
type SalesItemLine struct {
	ID                  string `json:"Id,omitempty"`
	LineNum             int    `json:",omitempty"`
	Description         string `json:",omitempty"`
	Amount              float32
	DetailType          string
	SalesItemLineDetail SalesItemLineDetail
}

// SalesItemLineDetail ...
type SalesItemLineDetail struct {
	ItemRef   ReferenceType `json:",omitempty"`
	ClassRef  ReferenceType `json:",omitempty"`
	UnitPrice float32       `json:",omitempty"`
	//MarkupInfo
	Qty             int           `json:",omitempty"`
	ItemAccountRef  ReferenceType `json:",omitempty"`
	TaxCodeRef      ReferenceType `json:",omitempty"`
	ServiceDate     null.Time     `json:",omitempty"`
	TaxInclusiveAmt float32       `json:",omitempty"`
	DiscountRate    float32       `json:",omitempty"`
	DiscountAmt     float32       `json:",omitempty"`
}
