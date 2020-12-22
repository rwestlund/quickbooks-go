package quickbooks

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

// Vendor describes a vendor.
type Vendor struct {
	ID               string       `json:"Id,omitempty"`
	SyncToken        string       `json:",omitempty"`
	Title            string       `json:",omitempty"`
	GivenName        string       `json:",omitempty"`
	MiddleName       string       `json:",omitempty"`
	Suffix           string       `json:",omitempty"`
	FamilyName       string       `json:",omitempty"`
	PrimaryEmailAddr EmailAddress `json:",omitempty"`
	DisplayName      string       `json:",omitempty"`
	// ContactInfo
	APAccountRef   ReferenceType   `json:",omitempty"`
	TermRef        ReferenceType   `json:",omitempty"`
	GSTIN          string          `json:",omitempty"`
	Fax            TelephoneNumber `json:",omitempty"`
	BusinessNumber string          `json:",omitempty"`
	// CurrencyRef
	HasTPAR           bool            `json:",omitempty"`
	TaxReportingBasis string          `json:",omitempty"`
	Mobile            TelephoneNumber `json:",omitempty"`
	PrimaryPhone      TelephoneNumber `json:",omitempty"`
	Active            bool            `json:",omitempty"`
	AlternatePhone    TelephoneNumber `json:",omitempty"`
	MetaData          MetaData        `json:",omitempty"`
	Vendor1099        bool            `json:",omitempty"`
	BillRate          json.Number     `json:",omitempty"`
	WebAddr           *WebSiteAddress `json:",omitempty"`
	CompanyName       string          `json:",omitempty"`
	// VendorPaymentBankDetail
	TaxIdentifier       string           `json:",omitempty"`
	AcctNum             string           `json:",omitempty"`
	GSTRegistrationType string           `json:",omitempty"`
	PrintOnCheckName    string           `json:",omitempty"`
	BillAddr            *PhysicalAddress `json:",omitempty"`
	Balance             json.Number      `json:",omitempty"`
}

// GetVendors gets the vendors
func (c *Client) GetVendors(startpos int) ([]Vendor, error) {

	var r struct {
		QueryResponse struct {
			Vendor        []Vendor
			StartPosition int
			MaxResults    int
		}
	}
	q := "SELECT * FROM Vendor ORDERBY Id STARTPOSITION " +
		strconv.Itoa(startpos) + " MAXRESULTS " + strconv.Itoa(queryPageSize)
	err := c.query(q, &r)
	if err != nil {
		return nil, err
	}

	if r.QueryResponse.Vendor == nil {
		r.QueryResponse.Vendor = make([]Vendor, 0)
	}
	return r.QueryResponse.Vendor, nil
}

// CreateVendor creates the vendor
func (c *Client) CreateVendor(vendor *Vendor) (*Vendor, error) {
	var u, err = url.Parse(string(c.Endpoint))
	if err != nil {
		return nil, err
	}
	u.Path = "/v3/company/" + c.RealmID + "/vendor"
	var v = url.Values{}
	v.Add("minorversion", minorVersion)
	u.RawQuery = v.Encode()
	var j []byte
	j, err = json.Marshal(vendor)
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
		Vendor Vendor
		Time   Date
	}
	err = json.NewDecoder(res.Body).Decode(&r)
	return &r.Vendor, err
}

// UpdateVendor updates the vendor
func (c *Client) UpdateVendor(vendor *Vendor) (*Vendor, error) {
	var u, err = url.Parse(string(c.Endpoint))
	if err != nil {
		return nil, err
	}
	u.Path = "/v3/company/" + c.RealmID + "/vendor"
	var v = url.Values{}
	v.Add("minorversion", minorVersion)
	u.RawQuery = v.Encode()
	var d = struct {
		*Vendor
		Sparse bool `json:"sparse"`
	}{
		Vendor: vendor,
		Sparse: true,
	}
	var j []byte
	j, err = json.Marshal(d)
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
		Vendor Vendor
		Time   Date
	}
	err = json.NewDecoder(res.Body).Decode(&r)
	return &r.Vendor, err
}
