package quickbooks

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// QueryPaged is like Query but adds the ability to specify a start position and page size ot the query.
func QueryPaged[T any](c *Client, query string, startpos int, pagesize int) ([]T, error) {
	selectStatement := fmt.Sprintf("%s STARTPOSITION %d MAXRESULTS %d", query, startpos, pagesize)
	return Query[T](c, selectStatement)

}

// Query allows you to query any QuickBooks Online entity and have the result unmarshalled into
// a slice of a type you specify.
//
// Example:
//
//	type JournalEntry struct {
//		 Id     string
//		 Amount float64
//	}
//	result, err := quickbooks.Query[JournalEntry](client, "SELECT * FROM JournalEntry ORDERBY Id")
func Query[T any](c *Client, query string) ([]T, error) {

	var resp struct {
		QueryResponse map[string]json.RawMessage
	}

	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	for key, value := range resp.QueryResponse {
		switch key {
		case "startPosition", "maxResults", "totalCount": // skip these...
		default:
			var data []T
			decoder := json.NewDecoder(bytes.NewReader(value))
			decoder.UseNumber()
			err := decoder.Decode(&data)
			if err != nil {
				return nil, err
			}
			return data, nil
		}
	}

	return []T{}, nil
}
