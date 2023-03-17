// Copyright (c) 2020, Randy Westlund. All rights reserved.
// This code is under the BSD-2-Clause license.

package quickbooks

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// Error implements the error interface.
func (f Failure) Error() string {
	var text, err = json.Marshal(f)
	if err != nil {
		return "When marshalling error:" + err.Error()
	}
	return string(text)
}

// Failure is the outermost struct that holds an error response.
type Failure struct {
	Fault Fault `json:"fault"`
}

type Fault struct {
	Error []struct {
		Message string
		Detail  string
		Code    string `json:"code"`
		Element string `json:"element"`
	}
	Type string `json:"type"`
}

// parseFailure takes a response reader and tries to parse a Failure.
func parseFailure(res *http.Response) error {
	var msg, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.New("When reading response body:" + err.Error())
	}
	var errStruct Failure
	err = json.Unmarshal(msg, &errStruct)
	if err != nil {
		return err
	}
	return errStruct
}
