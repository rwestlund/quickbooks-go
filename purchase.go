package quickbooks

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

const (
	PaymentTypeCash       = "Cash"
	PaymentTypeCheck      = "Check"
	PaymentTypeCreditCard = "CreditCard"
)

type Purchase struct {
	ID          string `json:"Id,omitempty"`
	Line        []Line
	PaymentType string           `json:",omitempty"`
	AccountRef  ReferenceType    `json:",omitempty"`
	SyncToken   string           `json:",omitempty"`
	CurrencyRef ReferenceType    `json:",omitempty"`
	TxnDate     Date             `json:",omitempty"`
	PrintStatus string           `json:",omitempty"`
	RemitToAddr *PhysicalAddress `json:",omitempty"`
	TxnSource   string           `json:",omitempty"`
	//LinkedTxn
	//GlobalTaxCalculation
	TransactionLocationType string        `json:",omitempty"`
	MetaData                MetaData      `json:",omitempty"`
	DocNumber               string        `json:",omitempty"`
	PrivateNote             string        `json:",omitempty"`
	Credit                  bool          `json:",omitempty"`
	TxnTaxDetail            TxnTaxDetail  `json:",omitempty"`
	PaymentMethodRef        ReferenceType `json:",omitempty"`
	ExchangeRate            json.Number   `json:",omitempty"`
	DepartmentRef           ReferenceType `json:",omitempty"`
	EntityRef               ReferenceType `json:",omitempty"`
	IncludeInAnnualTPAR     bool          `json:",omitempty"`
	TotalAmt                json.Number   `json:",omitempty"`
	RecurDataRef            ReferenceType `json:",omitempty"`
}

// CreatePurchase creates the purchase
func (c *Client) CreatePurchase(purchase *Purchase) (*Purchase, error) {
	var u, err = url.Parse(string(c.Endpoint))
	if err != nil {
		return nil, err
	}
	u.Path = "/v3/company/" + c.RealmID + "/purchase"
	var v = url.Values{}
	v.Add("minorversion", minorVersion)
	u.RawQuery = v.Encode()
	var j []byte
	j, err = json.Marshal(purchase)
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
		Purchase Purchase
		Time     Date
	}
	err = json.NewDecoder(res.Body).Decode(&r)
	return &r.Purchase, err
}

// QueryPurchase gets the purchase
func (c *Client) QueryPurchase(selectStatement string) ([]Purchase, error) {
	var r struct {
		QueryResponse struct {
			Purchase      []Purchase
			StartPosition int
			MaxResults    int
		}
	}
	err := c.query(selectStatement, &r)
	if err != nil {
		return nil, err
	}

	if r.QueryResponse.Purchase == nil {
		r.QueryResponse.Purchase = make([]Purchase, 0)
	}
	return r.QueryResponse.Purchase, nil
}

// GetPurchases gets the purchase
func (c *Client) GetPurchases(startpos int, pagesize int) ([]Purchase, error) {
	q := "SELECT * FROM Purchase ORDERBY Id STARTPOSITION " +
		strconv.Itoa(startpos) + " MAXRESULTS " + strconv.Itoa(pagesize)
	return c.QueryPurchase(q)
}

// GetPurchaseByID returns an purchase with a given ID.
func (c *Client) GetPurchaseByID(id string) (*Purchase, error) {
	var u, err = url.Parse(string(c.Endpoint))
	if err != nil {
		return nil, err
	}
	u.Path = "/v3/company/" + c.RealmID + "/purchase/" + id
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
		Purchase Purchase
		Time     Date
	}
	err = json.NewDecoder(res.Body).Decode(&r)
	return &r.Purchase, err
}

// UpdatePurchase updates the purchase
func (c *Client) UpdatePurchase(purchase *Purchase) (*Purchase, error) {
	var u, err = url.Parse(string(c.Endpoint))
	if err != nil {
		return nil, err
	}
	u.Path = "/v3/company/" + c.RealmID + "/purchase"
	var v = url.Values{}
	v.Add("minorversion", minorVersion)
	u.RawQuery = v.Encode()
	var d = struct {
		*Purchase
		Sparse bool `json:"sparse"`
	}{
		Purchase: purchase,
		Sparse:   true,
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
		Purchase Purchase
		Time     Date
	}
	err = json.NewDecoder(res.Body).Decode(&r)
	return &r.Purchase, err
}
