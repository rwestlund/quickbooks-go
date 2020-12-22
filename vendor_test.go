package quickbooks

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVendor(t *testing.T) {
	jsonFile, err := os.Open("data/testing/vendor.json")
	require.NoError(t, err)
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	require.NoError(t, err)

	var r struct {
		Vendor Vendor
		Time   Date
	}
	err = json.Unmarshal(byteValue, &r)
	require.NoError(t, err)
	assert.NotNil(t, r.Vendor.PrimaryEmailAddr)
	assert.False(t, r.Vendor.Vendor1099)
	assert.Equal(t, "Bessie", r.Vendor.GivenName)
	assert.Equal(t, "Books by Bessie", r.Vendor.DisplayName)
	assert.NotNil(t, r.Vendor.BillAddr)
	assert.Equal(t, "0", r.Vendor.SyncToken)
	assert.Equal(t, "Books by Bessie", r.Vendor.PrintOnCheckName)
	assert.Equal(t, "Williams", r.Vendor.FamilyName)
	assert.NotNil(t, r.Vendor.PrimaryPhone)
	assert.Equal(t, "1345", r.Vendor.AcctNum)
	assert.Equal(t, "Books by Bessie", r.Vendor.CompanyName)
	assert.NotNil(t, r.Vendor.WebAddr)
	assert.True(t, r.Vendor.Active)
	assert.Equal(t, "0", r.Vendor.Balance.String())
	assert.Equal(t, "30", r.Vendor.ID)
	assert.Equal(t, "2014-09-12T10:07:56-07:00", r.Vendor.MetaData.CreateTime.String())
	assert.Equal(t, "2014-09-17T11:13:46-07:00", r.Vendor.MetaData.LastUpdatedTime.String())
}
