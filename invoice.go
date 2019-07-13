// Copyright (c) 2018, Randy Westlund. All rights reserved.
// This code is under the BSD-2-Clause license.

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

	null "gopkg.in/guregu/null.v3"
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
	Line         []SalesItemLine
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
	TotalAmt      float32       `json:",omitempty"`
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

// TxnTaxDetail ...
type TxnTaxDetail struct {
	TxnTaxCodeRef ReferenceType `json:",omitempty"`
	TotalTax      float32       `json:",omitempty"`
	TaxLine       []Line        `json:",omitempty"`
}

// Line ...
type Line struct {
	Amount float32 `json:",omitempty"`
	// Must be set to "TaxLineDetail".
	DetailType    string
	TaxLineDetail TaxLineDetail
}

// TaxLineDetail ...
type TaxLineDetail struct {
	PercentBased     bool    `json:",omitempty"`
	NetAmountTaxable float32 `json:",omitempty"`
	//TaxInclusiveAmount float32 `json:",omitempty"`
	//OverrideDeltaAmount
	TaxPercent float32 `json:',omitempty"`
	TaxRateRef ReferenceType
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
	u.Path = "/v3/company/" + c.RealmID + "/invoice"
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
	req, err = http.NewRequest("POST", u.String()+"?operation=delete", bytes.NewBuffer(j))
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
