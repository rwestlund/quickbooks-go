package quickbooks

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

const (
	BankAccountType                  = "Bank"
	OtherCurrentAssetAccountType     = "Other Current Asset"
	FixedAssetAccountType            = "Fixed Asset"
	OtherAssetAccountType            = "Other Asset"
	AccountsReceivableAccountType    = "Accounts Receivable"
	EquityAccountType                = "Equity"
	ExpenseAccountType               = "Expense"
	OtherExpenseAccountType          = "Other Expense"
	CostOfGoodsSoldAccountType       = "Cost of Goods Sold"
	AccountsPayableAccountType       = "Accounts Payable"
	CreditCardAccountType            = "Credit Card"
	LongTermLiabilityAccountType     = "Long Term Liability"
	OtherCurrentLiabilityAccountType = "Other Current Liability"
	IncomeAccountType                = "Income"
	OtherIncomeAccountType           = "Other Income"
)

type Account struct {
	ID                            string        `json:"Id,omitempty"`
	Name                          string        `json:",omitempty"`
	SyncToken                     string        `json:",omitempty"`
	AcctNum                       string        `json:",omitempty"`
	CurrencyRef                   ReferenceType `json:",omitempty"`
	ParentRef                     ReferenceType `json:",omitempty"`
	Description                   string        `json:",omitempty"`
	Active                        bool          `json:",omitempty"`
	MetaData                      MetaData      `json:",omitempty"`
	SubAccount                    bool          `json:",omitempty"`
	Classification                string        `json:",omitempty"`
	FullyQualifiedName            string        `json:",omitempty"`
	TxnLocationType               string        `json:",omitempty"`
	AccountType                   string        `json:",omitempty"`
	CurrentBalanceWithSubAccounts json.Number   `json:",omitempty"`
	AccountAlias                  string        `json:",omitempty"`
	TaxCodeRef                    ReferenceType `json:",omitempty"`
	AccountSubType                string        `json:",omitempty"`
	CurrentBalance                json.Number   `json:",omitempty"`
}

// CreateAccount creates the account
func (c *Client) CreateAccount(account *Account) (*Account, error) {
	var u, err = url.Parse(string(c.Endpoint))
	if err != nil {
		return nil, err
	}
	u.Path = "/v3/company/" + c.RealmID + "/account"
	var v = url.Values{}
	v.Add("minorversion", minorVersion)
	u.RawQuery = v.Encode()
	var j []byte
	j, err = json.Marshal(account)
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
		Account Account
		Time    Date
	}
	err = json.NewDecoder(res.Body).Decode(&r)
	return &r.Account, err
}

// QueryAccount gets the account
func (c *Client) QueryAccount(selectStatement string) ([]Account, error) {
	var r struct {
		QueryResponse struct {
			Account       []Account
			StartPosition int
			MaxResults    int
		}
	}
	err := c.query(selectStatement, &r)
	if err != nil {
		return nil, err
	}

	if r.QueryResponse.Account == nil {
		r.QueryResponse.Account = make([]Account, 0)
	}
	return r.QueryResponse.Account, nil
}

// GetAccounts gets the account
func (c *Client) GetAccounts(startpos int, pagesize int) ([]Account, error) {
	q := "SELECT * FROM Account ORDERBY Id STARTPOSITION " +
		strconv.Itoa(startpos) + " MAXRESULTS " + strconv.Itoa(pagesize)
	return c.QueryAccount(q)
}

// GetAccountByID returns an account with a given ID.
func (c *Client) GetAccountByID(id string) (*Account, error) {
	var u, err = url.Parse(string(c.Endpoint))
	if err != nil {
		return nil, err
	}
	u.Path = "/v3/company/" + c.RealmID + "/account/" + id
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
		Account Account
		Time    Date
	}
	err = json.NewDecoder(res.Body).Decode(&r)
	return &r.Account, err
}

// UpdateAccount updates the account
func (c *Client) UpdateAccount(account *Account) (*Account, error) {
	var u, err = url.Parse(string(c.Endpoint))
	if err != nil {
		return nil, err
	}
	u.Path = "/v3/company/" + c.RealmID + "/account"
	var v = url.Values{}
	v.Add("minorversion", minorVersion)
	u.RawQuery = v.Encode()
	var d = struct {
		*Account
		Sparse bool `json:"sparse"`
	}{
		Account: account,
		Sparse:  true,
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
		Account Account
		Time    Date
	}
	err = json.NewDecoder(res.Body).Decode(&r)
	return &r.Account, err
}
