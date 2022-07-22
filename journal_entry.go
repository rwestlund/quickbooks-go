package quickbooks

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

type JournalEntry struct {
	ID                      string         `json:"Id,omitempty"`
	Line                    []Line         `json:",omitempty"`
	SyncToken               string         `json:",omitempty"`
	CurrencyRef             ReferenceType  `json:",omitempty"`
	DocNumber               string         `json:",omitempty"`
	PrivateNote             string         `json:",omitempty"`
	TxnDate                 Date           `json:",omitempty"`
	ExchangeRate            json.Number    `json:",omitempty"`
	TaxRateRef              *ReferenceType `json:",omitempty"`
	TransactionLocationType string         `json:",omitempty"`
	TxnTaxDetail            TxnTaxDetail   `json:",omitempty"`
	//GlobalTaxCalculation
	Adjustment   bool          `json:",omitempty"`
	MetaData     MetaData      `json:",omitempty"`
	RecurDataRef ReferenceType `json:",omitempty"`
	TotalAmt     json.Number   `json:",omitempty"`
}

// CreateJournalEntry creates the journalEntry
func (c *Client) CreateJournalEntry(journalEntry *JournalEntry, opts ...ClientOpt) (*JournalEntry, error) {
	var u, err = url.Parse(string(c.Endpoint))
	if err != nil {
		return nil, err
	}
	u.Path = "/v3/company/" + c.RealmID + "/journalentry"
	var v = url.Values{}
	v.Add("minorversion", minorVersion)

	for _, o := range opts {
		if o.Type == ClientOptTypeQueryParameter {
			v.Add(o.Name, o.Value)
		}
	}

	u.RawQuery = v.Encode()
	var j []byte
	j, err = json.Marshal(journalEntry)
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
		JournalEntry JournalEntry
		Time         Date
	}
	err = json.NewDecoder(res.Body).Decode(&r)
	return &r.JournalEntry, err
}

// QueryJournalEntry gets the journalEntry
func (c *Client) QueryJournalEntry(selectStatement string) ([]JournalEntry, error) {
	var r struct {
		QueryResponse struct {
			JournalEntry  []JournalEntry
			StartPosition int
			MaxResults    int
		}
	}
	err := c.query(selectStatement, &r)
	if err != nil {
		return nil, err
	}

	if r.QueryResponse.JournalEntry == nil {
		r.QueryResponse.JournalEntry = make([]JournalEntry, 0)
	}
	return r.QueryResponse.JournalEntry, nil
}

// GetJournalEntrys gets the journalEntry
func (c *Client) GetJournalEntrys(startpos int, pagesize int) ([]JournalEntry, error) {
	q := "SELECT * FROM JournalEntry ORDERBY Id STARTPOSITION " +
		strconv.Itoa(startpos) + " MAXRESULTS " + strconv.Itoa(pagesize)
	return c.QueryJournalEntry(q)
}

// GetJournalEntryByID returns an journalEntry with a given ID.
func (c *Client) GetJournalEntryByID(id string) (*JournalEntry, error) {
	var u, err = url.Parse(string(c.Endpoint))
	if err != nil {
		return nil, err
	}
	u.Path = "/v3/company/" + c.RealmID + "/journalentry/" + id
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
		JournalEntry JournalEntry
		Time         Date
	}
	err = json.NewDecoder(res.Body).Decode(&r)
	return &r.JournalEntry, err
}

// UpdateJournalEntry updates the journalEntry
func (c *Client) UpdateJournalEntry(journalEntry *JournalEntry, opts ...ClientOpt) (*JournalEntry, error) {
	var u, err = url.Parse(string(c.Endpoint))
	if err != nil {
		return nil, err
	}
	u.Path = "/v3/company/" + c.RealmID + "/journalentry"
	var v = url.Values{}
	v.Add("minorversion", minorVersion)

	for _, o := range opts {
		if o.Type == ClientOptTypeQueryParameter {
			v.Add(o.Name, o.Value)
		}
	}

	u.RawQuery = v.Encode()
	var d = struct {
		*JournalEntry
		Sparse bool `json:"sparse"`
	}{
		JournalEntry: journalEntry,
		Sparse:       true,
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
		JournalEntry JournalEntry
		Time         Date
	}
	err = json.NewDecoder(res.Body).Decode(&r)
	return &r.JournalEntry, err
}
