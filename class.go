package quickbooks

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

type Class struct {
	ID                 string    `json:"Id,omitempty"`
	Name               string    `json:",omitempty"`
	SyncToken          string    `json:",omitempty"`
	ParentRef          ParentRef `json:",omitempty"`
	SubClass           bool      `json:",omitempty"`
	FullyQualifiedName string    `json:",omitempty"`
}

type ParentRef struct {
	Value string `json:"value"`
}

// GetClasses fetches classes based on a page size
func (c *Client) GetClasses(startpos int, pagesize int) ([]Class, error) {
	q := "SELECT * FROM Class ORDERBY Id STARTPOSITION " +
		strconv.Itoa(startpos) + " MAXRESULTS " + strconv.Itoa(pagesize)
	return c.QueryClasses(q)
}

// QueryClasses runs a select statement for classes
func (c *Client) QueryClasses(selectStatement string) ([]Class, error) {
	var r struct {
		QueryResponse struct {
			Class         []Class
			StartPosition int `json:"startPosition"`
			MaxResults    int `json:"totalCount"`
		}
	}
	err := c.query(selectStatement, &r)
	if err != nil {
		return nil, err
	}

	if r.QueryResponse.Class == nil {
		r.QueryResponse.Class = make([]Class, 0)
	}
	return r.QueryResponse.Class, nil
}

// GetClassByID returns an account with a given ID.
func (c *Client) GetClassByID(id string) (*Class, error) {
	var u, err = url.Parse(string(c.Endpoint))
	if err != nil {
		return nil, err
	}
	u.Path = "/v3/company/" + c.RealmID + "/class/" + id
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
		Class Class
		Time  Date
	}
	err = json.NewDecoder(res.Body).Decode(&r)
	return &r.Class, err
}
