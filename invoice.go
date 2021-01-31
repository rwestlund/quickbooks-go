// Copyright (c) 2018, Randy Westlund. All rights reserved.
// This code is under the BSD-2-Clause license.

package quickbooks

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

// Invoice represents a QuickBooks Invoice object.
type Invoice struct {
	ID        string   `json:"Id,omitempty"`
	SyncToken string   `json:",omitempty"`
	MetaData  MetaData `json:",omitempty"`
	//CustomField
	DocNumber string `json:",omitempty"`
	TxnDate   Date   `json:",omitempty"`
	//DepartmentRef
	PrivateNote string `json:",omitempty"`
	//LinkedTxn
	Line         []Line
	TxnTaxDetail TxnTaxDetail `json:",omitempty"`
	CustomerRef  ReferenceType
	CustomerMemo MemoRef         `json:",omitempty"`
	BillAddr     PhysicalAddress `json:",omitempty"`
	ShipAddr     PhysicalAddress `json:",omitempty"`
	ClassRef     ReferenceType   `json:",omitempty"`
	SalesTermRef ReferenceType   `json:",omitempty"`
	DueDate      Date            `json:",omitempty"`
	//GlobalTaxCalculation
	ShipMethodRef ReferenceType `json:",omitempty"`
	ShipDate      Date          `json:",omitempty"`
	TrackingNum   string        `json:",omitempty"`
	TotalAmt      json.Number   `json:",omitempty"`
	//CurrencyRef
	ExchangeRate          json.Number  `json:",omitempty"`
	HomeAmtTotal          json.Number  `json:",omitempty"`
	HomeBalance           json.Number  `json:",omitempty"`
	ApplyTaxAfterDiscount bool         `json:",omitempty"`
	PrintStatus           string       `json:",omitempty"`
	EmailStatus           string       `json:",omitempty"`
	BillEmail             EmailAddress `json:",omitempty"`
	BillEmailCC           EmailAddress `json:"BillEmailCc,omitempty"`
	BillEmailBCC          EmailAddress `json:"BillEmailBcc,omitempty"`
	//DeliveryInfo
	Balance                      json.Number   `json:",omitempty"`
	TxnSource                    string        `json:",omitempty"`
	AllowOnlineCreditCardPayment bool          `json:",omitempty"`
	AllowOnlineACHPayment        bool          `json:",omitempty"`
	Deposit                      json.Number   `json:",omitempty"`
	DepositToAccountRef          ReferenceType `json:",omitempty"`
}

// TxnTaxDetail ...
type TxnTaxDetail struct {
	TxnTaxCodeRef ReferenceType `json:",omitempty"`
	TotalTax      json.Number   `json:",omitempty"`
	TaxLine       []Line        `json:",omitempty"`
}

const (
	PostingTypeCredit         = "Credit"
	PostingTypeDebit          = "Debit"
	BillableStatusBillable    = "Billable"
	BillableStatusNotBillable = "NotBillable"
	BillableStatus            = "HasBeenBilled"
)

type JournalEntryLineDetail struct {
	JournalCodeRef  ReferenceType `json:",omitempty"`
	PostingType     string        `json:",omitempty"`
	AccountRef      ReferenceType `json:",omitempty"`
	TaxApplicableOn string        `json:",omitempty"`
	TaxInclusiveAmt json.Number   `json:",omitempty"`
	ClassRef        ReferenceType `json:",omitempty"`
	DepartmentRef   ReferenceType `json:",omitempty"`
	TaxCodeRef      ReferenceType `json:",omitempty"`
	BillableStatus  string        `json:",omitempty"`
	Entity          ReferenceType `json:",omitempty"`
}

// AccountBasedExpenseLineDetail
type AccountBasedExpenseLineDetail struct {
	AccountRef ReferenceType
	TaxAmount  json.Number `json:",omitempty"`
	//TaxInclusiveAmt json.Number              `json:",omitempty"`
	//ClassRef        ReferenceType `json:",omitempty"`
	//TaxCodeRef      ReferenceType `json:",omitempty"`
	// MarkupInfo MarkupInfo `json:",omitempty"`
	//BillableStatus BillableStatusEnum       `json:",omitempty"`
	//CustomerRef    ReferenceType `json:",omitempty"`
}

// Line ...
type Line struct {
	ID                            string `json:"Id,omitempty"`
	LineNum                       int    `json:",omitempty"`
	Description                   string `json:",omitempty"`
	Amount                        json.Number
	DetailType                    string
	AccountBasedExpenseLineDetail AccountBasedExpenseLineDetail `json:",omitempty"`
	JournalEntryLineDetail        JournalEntryLineDetail        `json:",omitempty"`
	SalesItemLineDetail           SalesItemLineDetail           `json:",omitempty"`
	DiscountLineDetail            DiscountLineDetail            `json:",omitempty"`
	TaxLineDetail                 TaxLineDetail                 `json:",omitempty"`
}

// TaxLineDetail ...
type TaxLineDetail struct {
	PercentBased     bool        `json:",omitempty"`
	NetAmountTaxable json.Number `json:",omitempty"`
	//TaxInclusiveAmount json.Number `json:",omitempty"`
	//OverrideDeltaAmount
	TaxPercent json.Number `json:",omitempty"`
	TaxRateRef ReferenceType
}

// SalesItemLineDetail ...
type SalesItemLineDetail struct {
	ItemRef   ReferenceType `json:",omitempty"`
	ClassRef  ReferenceType `json:",omitempty"`
	UnitPrice json.Number   `json:",omitempty"`
	//MarkupInfo
	Qty             float32       `json:",omitempty"`
	ItemAccountRef  ReferenceType `json:",omitempty"`
	TaxCodeRef      ReferenceType `json:",omitempty"`
	ServiceDate     Date          `json:",omitempty"`
	TaxInclusiveAmt json.Number   `json:",omitempty"`
	DiscountRate    json.Number   `json:",omitempty"`
	DiscountAmt     json.Number   `json:",omitempty"`
}

// DiscountLineDetail ...
type DiscountLineDetail struct {
	PercentBased    bool
	DiscountPercent float32 `json:",omitempty"`
}

// FetchInvoices gets the full list of Invoices in the QuickBooks account.
func (c *Client) FetchInvoices() ([]Invoice, error) {

	// See how many invoices there are.
	var r struct {
		QueryResponse struct {
			TotalCount int
		}
	}
	err := c.query("SELECT COUNT(*) FROM Invoice", &r)
	if err != nil {
		return nil, err
	}

	if r.QueryResponse.TotalCount == 0 {
		return make([]Invoice, 0), nil
	}

	var invoices = make([]Invoice, 0, r.QueryResponse.TotalCount)
	for i := 0; i < r.QueryResponse.TotalCount; i += queryPageSize {
		var page, err = c.fetchInvoicePage(i + 1)
		if err != nil {
			return nil, err
		}
		invoices = append(invoices, page...)
	}
	return invoices, nil
}

// Fetch one page of results, because we can't get them all in one query.
func (c *Client) fetchInvoicePage(startpos int) ([]Invoice, error) {

	var r struct {
		QueryResponse struct {
			Invoice       []Invoice
			StartPosition int
			MaxResults    int
		}
	}
	q := "SELECT * FROM Invoice ORDERBY Id STARTPOSITION " +
		strconv.Itoa(startpos) + " MAXRESULTS " + strconv.Itoa(queryPageSize)
	err := c.query(q, &r)
	if err != nil {
		return nil, err
	}

	// Make sure we don't return nil if there are no invoices.
	if r.QueryResponse.Invoice == nil {
		r.QueryResponse.Invoice = make([]Invoice, 0)
	}
	return r.QueryResponse.Invoice, nil
}

// CreateInvoice creates the given Invoice on the QuickBooks server, returning
// the resulting Invoice object.
func (c *Client) CreateInvoice(inv *Invoice) (*Invoice, error) {
	var u, err = url.Parse(string(c.Endpoint))
	if err != nil {
		return nil, err
	}
	u.Path = "/v3/company/" + c.RealmID + "/invoice"
	var v = url.Values{}
	v.Add("minorversion", minorVersion)
	u.RawQuery = v.Encode()
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

	if res.StatusCode != http.StatusOK {
		return nil, parseFailure(res)
	}

	var r struct {
		Invoice Invoice
		Time    Date
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
	u.Path = "/v3/company/" + c.RealmID + "/invoice"
	var v = url.Values{}
	v.Add("minorversion", minorVersion)
	v.Add("operation", "delete")
	u.RawQuery = v.Encode()
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
	// {"Fault":{"Error":[{"Message":"Object Not Found","Detail":"Object Not Found : Something you're trying to use has been made inactive. Check the fields with accounts, invoices, items, vendors or employees.","code":"610","element":""}],"type":"ValidationFault"},"time":"2018-03-20T20:15:59.571-07:00"}

	// This is slightly horrifying and not documented in their API. When this
	// happens we just return success; the goal of deleting it has been
	// accomplished, just not by us.
	if res.StatusCode == http.StatusBadRequest {
		var r Failure
		err = json.NewDecoder(res.Body).Decode(&r)
		if err != nil {
			return err
		}
		if r.Fault.Error[0].Message == "Object Not Found" {
			return nil
		}
	}
	if res.StatusCode != http.StatusOK {
		return parseFailure(res)
	}

	// TODO they send something back, but is it useful?
	return nil
}
