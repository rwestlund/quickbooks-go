package quickbooks

import (
	"errors"
	"strconv"
)

type Payment struct {
	SyncToken           string        `json:",omitempty"`
	Domain              string        `json:"domain,omitempty"`
	DepositToAccountRef ReferenceType `json:",omitempty"`
	UnappliedAmt        float64       `json:",omitempty"`
	TxnDate             string        `json:",omitempty"`
	TotalAmt            float64       `json:",omitempty"`
	ProcessPayment      bool          `json:",omitempty"`
	Sparse              bool          `json:"sparse,omitempty"`
	Line                []PaymentLine `json:",omitempty"`
	CustomerRef         ReferenceType `json:",omitempty"`
	Id                  string        `json:",omitempty"`
	MetaData            MetaData      `json:",omitempty"`
}

type PaymentLine struct {
	Amount    float64 `json:",omitempty"`
	LinkedTxn []struct {
		TxnLineId string `json:",omitempty"`
		TxnId     string `json:",omitempty"`
		TxnType   string `json:",omitempty"`
	} `json:",omitempty"`
}

// CreatePayment creates the given payment within QuickBooks
func (c *Client) CreatePayment(payment *Payment) (*Payment, error) {
	var resp struct {
		Payment Payment
		Time    Date
	}

	if err := c.post("payment", payment, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.Payment, nil
}

func (c *Client) DeletePayment(payment *Payment) error {
	if payment.Id == "" || payment.SyncToken == "" {
		return errors.New("missing id/sync token")
	}

	return c.post("payment", payment, nil, map[string]string{"operation": "delete"})
}

// FindPayments gets the full list of Payments in the QuickBooks account.
func (c *Client) FindPayments() ([]Payment, error) {
	var resp struct {
		QueryResponse struct {
			Payments      []Payment `json:"Payment"`
			MaxResults    int
			StartPosition int
			TotalCount    int
		}
	}

	if err := c.query("SELECT COUNT(*) FROM Payment", &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.TotalCount == 0 {
		return nil, errors.New("no payments could be found")
	}

	payments := make([]Payment, 0, resp.QueryResponse.TotalCount)

	for i := 0; i < resp.QueryResponse.TotalCount; i += queryPageSize {
		query := "SELECT * FROM Payment ORDERBY Id STARTPOSITION " + strconv.Itoa(i+1) + " MAXRESULTS " + strconv.Itoa(queryPageSize)

		if err := c.query(query, &resp); err != nil {
			return nil, err
		}

		if resp.QueryResponse.Payments == nil {
			return nil, errors.New("no payments could be found")
		}

		payments = append(payments, resp.QueryResponse.Payments...)
	}

	return payments, nil
}

// FindPaymentById returns an payment with a given Id.
func (c *Client) FindPaymentById(id string) (*Payment, error) {
	var resp struct {
		Payment Payment
		Time    Date
	}

	if err := c.get("payment/"+id, &resp, nil); err != nil {
		return nil, err
	}

	return &resp.Payment, nil
}

// QueryPayments accepts an SQL query and returns all payments found using it
func (c *Client) QueryPayments(query string) ([]Payment, error) {
	var resp struct {
		QueryResponse struct {
			Payments      []Payment `json:"Payment"`
			StartPosition int
			MaxResults    int
		}
	}

	if err := c.query(query, &resp); err != nil {
		return nil, err
	}

	if resp.QueryResponse.Payments == nil {
		return nil, errors.New("could not find any payments")
	}

	return resp.QueryResponse.Payments, nil
}

// UpdatePayment updates the payment
func (c *Client) UpdatePayment(payment *Payment) (*Payment, error) {
	if payment.Id == "" {
		return nil, errors.New("missing payment id")
	}

	existingPayment, err := c.FindPaymentById(payment.Id)
	if err != nil {
		return nil, err
	}

	payment.SyncToken = existingPayment.SyncToken

	payload := struct {
		*Payment
		Sparse bool `json:"sparse"`
	}{
		Payment: payment,
		Sparse:  true,
	}

	var paymentData struct {
		Payment Payment
		Time    Date
	}

	if err = c.post("payment", payload, &paymentData, nil); err != nil {
		return nil, err
	}

	return &paymentData.Payment, err
}

// VoidPayment voids the given payment
func (c *Client) VoidPayment(payment Payment) error {
	if payment.Id == "" {
		return errors.New("missing payment id")
	}

	existingPayment, err := c.FindPaymentById(payment.Id)
	if err != nil {
		return err
	}

	payment.SyncToken = existingPayment.SyncToken

	return c.post("payment", payment, nil, map[string]string{"operation": "update", "include": "void"})
}
