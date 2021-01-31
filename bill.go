package quickbooks

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

type Bill struct {
	ID           string        `json:"Id,omitempty"`
	VendorRef    ReferenceType `json:",omitempty"`
	Line         []Line
	SyncToken    string        `json:",omitempty"`
	CurrencyRef  ReferenceType `json:",omitempty"`
	TxnDate      Date          `json:",omitempty"`
	APAccountRef ReferenceType `json:",omitempty"`
	SalesTermRef ReferenceType `json:",omitempty"`
	//LinkedTxn
	//GlobalTaxCalculation
	TotalAmt                json.Number `json:",omitempty"`
	TransactionLocationType string      `json:",omitempty"`
	DueDate                 Date        `json:",omitempty"`
	MetaData                MetaData    `json:",omitempty"`
	DocNumber               string
	PrivateNote             string        `json:",omitempty"`
	TxnTaxDetail            TxnTaxDetail  `json:",omitempty"`
	ExchangeRate            json.Number   `json:",omitempty"`
	DepartmentRef           ReferenceType `json:",omitempty"`
	IncludeInAnnualTPAR     bool          `json:",omitempty"`
	HomeBalance             json.Number   `json:",omitempty"`
	RecurDataRef            ReferenceType `json:",omitempty"`
	Balance                 json.Number   `json:",omitempty"`
}

func (c *Client) CreateBill(bill *Bill) (*Bill, error) {
	var u, err = url.Parse(string(c.Endpoint))
	if err != nil {
		return nil, err
	}
	u.Path = "/v3/company/" + c.RealmID + "/bill"
	var v = url.Values{}
	v.Add("minorversion", minorVersion)
	u.RawQuery = v.Encode()
	var j []byte
	j, err = json.Marshal(bill)
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
		Bill Bill
		Time Date
	}
	err = json.NewDecoder(res.Body).Decode(&r)
	return &r.Bill, err
}

// QueryBill gets the bill
func (c *Client) QueryBill(selectStatement string) ([]Bill, error) {
	var r struct {
		QueryResponse struct {
			Bill          []Bill
			StartPosition int
			MaxResults    int
		}
	}
	err := c.query(selectStatement, &r)
	if err != nil {
		return nil, err
	}

	if r.QueryResponse.Bill == nil {
		r.QueryResponse.Bill = make([]Bill, 0)
	}
	return r.QueryResponse.Bill, nil
}

// GetBills gets the bills
func (c *Client) GetBills(startpos int, pagesize int) ([]Bill, error) {
	q := "SELECT * FROM Bill ORDERBY Id STARTPOSITION " +
		strconv.Itoa(startpos) + " MAXRESULTS " + strconv.Itoa(pagesize)
	return c.QueryBill(q)
}

// GetBillByID returns a bill with a given ID.
func (c *Client) GetBillByID(id string) (*Bill, error) {
	var u, err = url.Parse(string(c.Endpoint))
	if err != nil {
		return nil, err
	}
	u.Path = "/v3/company/" + c.RealmID + "/bill/" + id
	var v = url.Values{}
	v.Add("minorversion", minorVersion)
	u.RawQuery = v.Encode()
	var req *http.Request
	req, err = http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
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
		Bill Bill
		Time Date
	}
	err = json.NewDecoder(res.Body).Decode(&r)
	return &r.Bill, err
}

// UpdateBill updates the bill
func (c *Client) UpdateBill(bill *Bill) (*Bill, error) {
	var u, err = url.Parse(string(c.Endpoint))
	if err != nil {
		return nil, err
	}
	u.Path = "/v3/company/" + c.RealmID + "/bill"
	var v = url.Values{}
	v.Add("minorversion", minorVersion)
	u.RawQuery = v.Encode()
	var d = struct {
		*Bill
		Sparse bool `json:"sparse"`
	}{
		Bill:   bill,
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
		Bill Bill
		Time Date
	}
	err = json.NewDecoder(res.Body).Decode(&r)
	return &r.Bill, err
}
