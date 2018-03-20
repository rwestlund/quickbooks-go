// Copyright (c) 2018, Randy Westlund. All rights reserved.
// This code is under the BSD-2-Clause license.

package quickbooks

import (
	"time"

	null "gopkg.in/guregu/null.v3"
)

// Customer represents a QuickBooks Customer object.
type Customer struct {
	ID        string `json:"Id,omitempty"`
	SyncToken string `json:",omitempty"`
	//MetaData
	Title              null.String     `json:",omitempty"`
	GivenName          null.String     `json:",omitempty"`
	MiddleName         null.String     `json:",omitempty"`
	FamilyName         null.String     `json:",omitempty"`
	Suffix             null.String     `json:",omitempty"`
	DisplayName        string          `json:",omitempty"`
	FullyQualifiedName null.String     `json:",omitempty"`
	CompanyName        null.String     `json:",omitempty"`
	PrintOnCheckName   string          `json:",omitempty"`
	Active             bool            `json:",omitempty"`
	PrimaryPhone       TelephoneNumber `json:",omitempty"`
	AlternatePhone     TelephoneNumber `json:",omitempty"`
	Mobile             TelephoneNumber `json:",omitempty"`
	Fax                TelephoneNumber `json:",omitempty"`
	PrimaryEmailAddr   EmailAddress    `json:",omitempty"`
	//WebAddr
	//DefaultTaxCodeRef
	Taxable              bool            `json:",omitempty"`
	TaxExemptionReasonID string          `json:"TaxExemptionReasonId,omitempty"`
	BillAddr             PhysicalAddress `json:",omitempty"`
	ShipAddr             PhysicalAddress `json:",omitempty"`
	Notes                string          `json:",omitempty"`
	Job                  null.Bool       `json:",omitempty"`
	BillWithParent       bool            `json:",omitempty"`
	//ParentRef
	Level int `json:",omitempty"`
	//SalesTermRef
	//PaymentMethodRef
	Balance         float32   `json:",omitempty"`
	OpenBalanceDate time.Time `json:",omitempty"`
	BalanceWithJobs float32   `json:",omitempty"`
	//CurrencyRef
}
