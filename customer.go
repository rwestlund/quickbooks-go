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

// FetchCustomers gets the full list of Customers in the QuickBooks account.
func (c *Client) FetchCustomers() ([]Customer, error) {

	// See how many customers there are.
	var r struct {
		QueryResponse struct {
			TotalCount int
		}
	}
	err := c.query("SELECT COUNT(*) FROM Customer", &r)
	if err != nil {
		return nil, err
	}

	if r.QueryResponse.TotalCount == 0 {
		return make([]Customer, 0), nil
	}

	var customers = make([]Customer, 0, r.QueryResponse.TotalCount)
	for i := 0; i < r.QueryResponse.TotalCount; i += queryPageSize {
		var page, err = c.fetchCustomerPage(i + 1)
		if err != nil {
			return nil, err
		}
		customers = append(customers, page...)
	}
	return customers, nil
}

// Fetch one page of results, because we can't get them all in one query.
func (c *Client) fetchCustomerPage(startpos int) ([]Customer, error) {

	var r struct {
		QueryResponse struct {
			Customer      []Customer
			StartPosition int
			MaxResults    int
		}
	}
	q := "SELECT * FROM Customer ORDERBY Id STARTPOSITION " +
		strconv.Itoa(startpos) + " MAXRESULTS " + strconv.Itoa(queryPageSize)
	err := c.query(q, &r)
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
