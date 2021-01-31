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

func TestPurchaseJSON(t *testing.T) {
	jsonFile, err := os.Open("data/testing/purchase.json")
	if err != nil {
		log.Fatal("When opening JSON file: ", err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal("When reading JSON file: ", err)
	}

	var r struct {
		Purchase Purchase
		Time     Date
	}
	err = json.Unmarshal(byteValue, &r)
	if err != nil {
		log.Fatal("When decoding JSON file: ", err)
	}
	assert.Equal(t, "0", r.Purchase.SyncToken)
	assert.Equal(t, "2015-07-27T00:00:00+00:00", r.Purchase.TxnDate.String())
	totalAmt, _ := r.Purchase.TotalAmt.Float64()
	assert.Equal(t, 10.0, totalAmt)
	assert.Equal(t, "Cash", r.Purchase.PaymentType)
	assert.Equal(t, "Checking", r.Purchase.AccountRef.Name)
	assert.Equal(t, "35", r.Purchase.AccountRef.Value)
	assert.Equal(t, "252", r.Purchase.ID)
	assert.Equal(t, "2015-07-27T10:37:26-07:00", r.Purchase.MetaData.CreateTime.String())
	assert.Equal(t, "2015-07-27T10:37:26-07:00", r.Purchase.MetaData.LastUpdatedTime.String())
}

func TestPurchase(t *testing.T) {
	qbClient, err := getClient()
	if err != nil {
		log.Println(err)
		t.Skip("Cannot instantiate the client")
	}

	jsonFile, err := os.Open("data/testing/purchase.json")
	if err != nil {
		t.Fatal("When opening JSON file: ", err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		t.Fatal("When reading JSON file: ", err)
	}

	var r struct {
		Purchase Purchase
		Time     Date
	}
	err = json.Unmarshal(byteValue, &r)
	if err != nil {
		t.Fatal("When decoding JSON file: ", err)
	}

	r.Purchase.ID = ""
	r.Purchase.TxnDate.Time = time.Now()
	createdPurchase, err := qbClient.CreatePurchase(&r.Purchase)
	assert.NoError(t, err)

	newPurchase, err := qbClient.GetPurchaseByID(createdPurchase.ID)
	assert.NoError(t, err)
	assert.Equal(t, createdPurchase.ID, newPurchase.ID)

	createdPurchase.DocNumber = "NewNumber"
	newPurchase, err = qbClient.UpdatePurchase(createdPurchase)
	assert.NoError(t, err)
	assert.Equal(t, "NewNumber", newPurchase.DocNumber)
}
