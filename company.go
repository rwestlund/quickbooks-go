// Copyright (c) 2018, Randy Westlund. All rights reserved.
// This code is under the BSD-2-Clause license.

package quickbooks

// CompanyInfo describes a company account.
type CompanyInfo struct {
	CompanyName string
	LegalName   string
	//CompanyAddr
	//CustomerCommunicationAddr
	//LegalAddr
	//PrimaryPhone
	//CompanyStartDate     time.Time
	CompanyStartDate     string
	FiscalYearStartMonth string
	Country              string
	//Email
	//WebAddr
	SupportedLanguages string
	//NameValue
	Domain    string
	ID        string `json:"Id"`
	SyncToken string
	Metadata  MetaData `json:",omitempty"`
}
