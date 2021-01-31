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

func TestJournalEntryJSON(t *testing.T) {
	jsonFile, err := os.Open("data/testing/journal-entry.json")
	if err != nil {
		log.Fatal("When opening JSON file: ", err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal("When reading JSON file: ", err)
	}

	var r struct {
		JournalEntry JournalEntry
		Time         Date
	}
	err = json.Unmarshal(byteValue, &r)
	if err != nil {
		log.Fatal("When decoding JSON file: ", err)
	}
	assert.Equal(t, "0", r.JournalEntry.SyncToken)
	assert.Equal(t, "2015-07-27T00:00:00+00:00", r.JournalEntry.TxnDate.String())
	assert.False(t, r.JournalEntry.Adjustment)
	assert.Equal(t, "227", r.JournalEntry.ID)
	assert.Equal(t, "2015-06-29T12:33:57-07:00", r.JournalEntry.MetaData.CreateTime.String())
	assert.Equal(t, "2015-06-29T12:33:57-07:00", r.JournalEntry.MetaData.LastUpdatedTime.String())
}

func TestJournalEntry(t *testing.T) {
	qbClient, err := getClient()
	if err != nil {
		log.Println(err)
		t.Skip("Cannot instantiate the client")
	}

	jsonFile, err := os.Open("data/testing/journal-entry.json")
	if err != nil {
		t.Fatal("When opening JSON file: ", err)
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		t.Fatal("When reading JSON file: ", err)
	}

	var r struct {
		JournalEntry JournalEntry
		Time         Date
	}
	err = json.Unmarshal(byteValue, &r)
	if err != nil {
		t.Fatal("When decoding JSON file: ", err)
	}

	r.JournalEntry.ID = ""
	r.JournalEntry.TxnDate.Time = time.Now()
	createdJournalEntry, err := qbClient.CreateJournalEntry(&r.JournalEntry)
	assert.NoError(t, err)

	newJournalEntry, err := qbClient.GetJournalEntryByID(createdJournalEntry.ID)
	assert.NoError(t, err)
	assert.Equal(t, createdJournalEntry.ID, newJournalEntry.ID)

	createdJournalEntry.DocNumber = "NewNumber"
	newJournalEntry, err = qbClient.UpdateJournalEntry(createdJournalEntry)
	assert.NoError(t, err)
	assert.Equal(t, "NewNumber", newJournalEntry.DocNumber)
}
