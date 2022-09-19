// Copyright (c) 2018, Randy Westlund. All rights reserved.
// This code is under the BSD-2-Clause license.

package quickbooks

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/guregu/null.v4"
)

// Customer represents a QuickBooks Customer object.
type Customer struct {
	Id                 string          `json:",omitempty"`
	SyncToken          string          `json:",omitempty"`
	MetaData           MetaData        `json:",omitempty"`
	Title              string          `json:",omitempty"`
	GivenName          string          `json:",omitempty"`
	MiddleName         string          `json:",omitempty"`
	FamilyName         string          `json:",omitempty"`
	Suffix             string          `json:",omitempty"`
	DisplayName        string          `json:",omitempty"`
	FullyQualifiedName string          `json:",omitempty"`
	CompanyName        string          `json:",omitempty"`
	PrintOnCheckName   string          `json:",omitempty"`
	Active             bool            `json:",omitempty"`
	PrimaryPhone       TelephoneNumber `json:",omitempty"`
	AlternatePhone     TelephoneNumber `json:",omitempty"`
	Mobile             TelephoneNumber `json:",omitempty"`
	Fax                TelephoneNumber `json:",omitempty"`
	CustomerTypeRef    ReferenceType   `json:",omitempty"`
	PrimaryEmailAddr   *EmailAddress   `json:",omitempty"`
	WebAddr            *WebSiteAddress `json:",omitempty"`
	// DefaultTaxCodeRef
	Taxable              *bool            `json:",omitempty"`
	TaxExemptionReasonId *string          `json:",omitempty"`
	BillAddr             *PhysicalAddress `json:",omitempty"`
	ShipAddr             *PhysicalAddress `json:",omitempty"`
	Notes                string           `json:",omitempty"`
	Job                  null.Bool        `json:",omitempty"`
	BillWithParent       bool             `json:",omitempty"`
	ParentRef            ReferenceType    `json:",omitempty"`
	Level                int              `json:",omitempty"`
	// SalesTermRef
	// PaymentMethodRef
	Balance         json.Number `json:",omitempty"`
	OpenBalanceDate Date        `json:",omitempty"`
	BalanceWithJobs json.Number `json:",omitempty"`
	// CurrencyRef
}

// GetAddress prioritizes the ship address, but falls back on bill address
func (c *Customer) GetAddress() PhysicalAddress {
	if c.ShipAddr != nil {
		return *c.ShipAddr
	}
	if c.BillAddr != nil {
		return *c.BillAddr
	}
	return PhysicalAddress{}
}

// GetWebsite de-nests the Website object
func (c *Customer) GetWebsite() string {
	if c.WebAddr != nil {
		return c.WebAddr.URI
	}
	return ""
}

// GetPrimaryEmail de-nests the PrimaryEmailAddr object
func (c *Customer) GetPrimaryEmail() string {
	if c.PrimaryEmailAddr != nil {
		return c.PrimaryEmailAddr.Address
	}
	return ""
}

// CreateCustomer creates the given Customer on the QuickBooks server,
// returning the resulting Customer object.
func (c *Client) CreateCustomer(customer *Customer) (*Customer, error) {
	var resp struct {
		Customer Customer
		Time     Date
	}

	if err := c.post("customer", customer, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.Customer, nil
}

// FindCustomers gets the full list of Customers in the QuickBooks account.
func (c *Client) FindCustomers() ([]Customer, error) {
	var resp struct {
		QueryResponse struct {
			Customers     []Customer `json:"Customer"`
			MaxResults    int
			StartPosition int
			TotalCount    int
		}
	}

	if err := c.query("SELECT COUNT(*) FROM Customer", &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.TotalCount == 0 {
		return nil, errors.New("no customers could be found")
	}

	customers := make([]Customer, 0, resp.QueryResponse.TotalCount)

	for i := 0; i < resp.QueryResponse.TotalCount; i += queryPageSize {
		query := "SELECT * FROM Customer ORDERBY Id STARTPOSITION " + strconv.Itoa(i+1) + " MAXRESULTS " + strconv.Itoa(queryPageSize)

		if err := c.query(query, &resp); err != nil {
			return nil, err
		}

		if resp.QueryResponse.Customers == nil {
			return nil, errors.New("no customers could be found")
		}

		customers = append(customers, resp.QueryResponse.Customers...)
	}

	return customers, nil
}

// FindCustomerById returns a customer with a given Id.
func (c *Client) FindCustomerById(id string) (*Customer, error) {
	var r struct {
		Customer Customer
		Time     Date
	}

	if err := c.get("customer/"+id, &r, nil); err != nil {
		return nil, err
	}

	return &r.Customer, nil
}

// FindCustomerByName gets a customer with a given name.
func (c *Client) FindCustomerByName(name string) (*Customer, error) {
	var resp struct {
		QueryResponse struct {
			Customer   []Customer
			TotalCount int
		}
	}

	query := "SELECT * FROM Customer WHERE DisplayName = '" + strings.Replace(name, "'", "''", -1) + "'"

	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	if len(resp.QueryResponse.Customer) == 0 {
		return nil, errors.New("no customers could be found")
	}

	return &resp.QueryResponse.Customer[0], nil
}

// QueryCustomers accepts an SQL query and returns all customers found using it
func (c *Client) QueryCustomers(query string) ([]Customer, error) {
	var resp struct {
		QueryResponse struct {
			Customers     []Customer `json:"Customer"`
			StartPosition int
			MaxResults    int
		}
	}

	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.Customers == nil {
		return nil, errors.New("could not find any customers")
	}

	return resp.QueryResponse.Customers, nil
}

// UpdateCustomer updates the given Customer on the QuickBooks server,
// returning the resulting Customer object. It's a sparse update, as not all QB
// fields are present in our Customer object.
func (c *Client) UpdateCustomer(customer *Customer) (*Customer, error) {
	if customer.Id == "" {
		return nil, errors.New("missing customer id")
	}

	existingCustomer, err := c.FindCustomerById(customer.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to find existing customer: %v", err)
	}

	customer.SyncToken = existingCustomer.SyncToken

	payload := struct {
		*Customer
		Sparse bool `json:"sparse"`
	}{
		Customer: customer,
		Sparse:   true,
	}

	var customerData struct {
		Customer Customer
		Time     Date
	}

	if err = c.post("customer", payload, &customerData, nil); err != nil {
		return nil, err
	}

	return &customerData.Customer, nil
}
