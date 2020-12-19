package quickbooks

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
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

func CreateBill(c Client, bill *Bill) (*Bill, error) {
	var u, err = url.Parse(string(c.Endpoint))
	if err != nil {
		return nil, err
	}
	u.Path = "/v3/company/" + c.RealmID + "/bill"
	var v = url.Values{}
	v.Add("minorversion", "55")
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
