package quickbooks

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

func TestBillJSON(t *testing.T) {
	jsonFile, err := os.Open("data/testing/bill.json")
	if err != nil {
		log.Fatal("When opening JSON file: ", err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal("When reading JSON file: ", err)
	}

	var r struct {
		Bill Bill
		Time Date
	}
	err = json.Unmarshal(byteValue, &r)
	if err != nil {
		log.Fatal("When decoding JSON file: ", err)
	}
	assert.Equal(t, "2", r.Bill.SyncToken)
	assert.Equal(t, "Accounts Payable (A/P)", r.Bill.APAccountRef.Name)
	assert.Equal(t, "33", r.Bill.APAccountRef.Value)
	assert.Equal(t, "Norton Lumber and Building Materials", r.Bill.VendorRef.Name)
	assert.Equal(t, "46", r.Bill.VendorRef.Value)
	assert.Equal(t, "2014-11-06T00:00:00+00:00", r.Bill.TxnDate.String())
	totalAmt, _ := r.Bill.TotalAmt.Float64()
	assert.Equal(t, 103.55, totalAmt)
	assert.Equal(t, "United States Dollar", r.Bill.CurrencyRef.Name)
	assert.Equal(t, "USD", r.Bill.CurrencyRef.Value)
	// LinkedTxn
	assert.Equal(t, "3", r.Bill.SalesTermRef.Value)
	assert.Equal(t, "2014-12-06T00:00:00+00:00", r.Bill.DueDate.String())
	assert.Equal(t, 1, len(r.Bill.Line))
	balance, _ := r.Bill.Balance.Int64()
	assert.Equal(t, int64(0), balance)
	assert.Equal(t, "25", r.Bill.ID)
	assert.Equal(t, "2014-11-06T15:37:25-08:00", r.Bill.MetaData.CreateTime.String())
	assert.Equal(t, "2015-02-09T10:11:11-08:00", r.Bill.MetaData.LastUpdatedTime.String())
}

func TestBill(t *testing.T) {
	qbClient, err := getClient()
	if err != nil {
		log.Println(err)
		t.Skip("Cannot instantiate the client")
	}

	jsonFile, err := os.Open("data/testing/bill.json")
	if err != nil {
		t.Fatal("When opening JSON file: ", err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		t.Fatal("When reading JSON file: ", err)
	}

	var r struct {
		Bill Bill
		Time Date
	}
	err = json.Unmarshal(byteValue, &r)
	if err != nil {
		t.Fatal("When decoding JSON file: ", err)
	}

	r.Bill.ID = ""
	r.Bill.TxnDate.Time = time.Now()
	createdBill, err := qbClient.CreateBill(&r.Bill)
	assert.NoError(t, err)

	newBill, err := qbClient.GetBillByID(createdBill.ID)
	assert.NoError(t, err)
	assert.Equal(t, createdBill.ID, newBill.ID)

	createdBill.DocNumber = "NewNumber"
	newBill, err = qbClient.UpdateBill(createdBill)
	assert.NoError(t, err)
	assert.Equal(t, "NewNumber", newBill.DocNumber)
}
