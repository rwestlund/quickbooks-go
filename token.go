package quickbooks

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"
	"net/url"
)

type BearerToken struct {
	RefreshToken           string `json:"refresh_token"`
	AccessToken            string `json:"access_token"`
	TokenType              string `json:"token_type"`
	IdToken                string `json:"id_token"`
	ExpiresIn              int64  `json:"expires_in"`
	XRefreshTokenExpiresIn int64  `json:"x_refresh_token_expires_in"`
}

//
// Method to retrieve access token (bearer token)
// This method can only be called once
//
func (c *Client) RetrieveBearerToken(authorizationCode, redirectURI string) (*BearerToken, error) {
	client := &http.Client{}
	data := url.Values{}
	//set parameters
	data.Set("grant_type", "authorization_code")
	data.Add("code", authorizationCode)
	data.Add("redirect_uri", redirectURI)

	request, err := http.NewRequest("POST", string(c.discoveryAPI.TokenEndpoint), bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}
	//set headers
	request.Header.Set("accept", "application/json")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	request.Header.Set("Authorization", "Basic "+basicAuth(c))

	resp, err := client.Do(request)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(body))
	}

	bearerTokenResponse, err := getBearerTokenResponse([]byte(body))
	return bearerTokenResponse, err
}

//
// Call the refresh endpoint to generate new tokens
//
func (c *Client) RefreshToken(refreshToken string) (*BearerToken, error) {
	client := &http.Client{}
	data := url.Values{}

	//add parameters
	data.Set("grant_type", "refresh_token")
	data.Add("refresh_token", refreshToken)

	request, err := http.NewRequest("POST", string(c.discoveryAPI.TokenEndpoint), bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}
	//set the headers
	request.Header.Set("accept", "application/json")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	request.Header.Set("Authorization", "Basic "+basicAuth(c))

	resp, err := client.Do(request)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(body))
	}

	bearerTokenResponse, err := getBearerTokenResponse([]byte(body))
	c.Client = getHttpClient(bearerTokenResponse)
	return bearerTokenResponse, err
}

//
// Call the revoke endpoint to revoke tokens
//
func (c *Client) RevokeToken(refreshToken string) error {
	client := &http.Client{}
	data := url.Values{}

	//add parameters
	data.Add("token", refreshToken)

	revokeEndpoint := c.discoveryAPI.RevocationEndpoint
	request, err := http.NewRequest("POST", revokeEndpoint, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return err
	}
	//set headers
	request.Header.Set("accept", "application/json")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	request.Header.Set("Authorization", "Basic "+basicAuth(c))

	resp, err := client.Do(request)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(string(body))
	}

	c.Client = nil
	return nil
}

func getBearerTokenResponse(body []byte) (*BearerToken, error) {
	var s = new(BearerToken)
	err := json.Unmarshal(body, &s)
	if err != nil {
		return nil, errors.New(string(body))
	}
	return s, err
}

func basicAuth(c *Client) string {
	auth := c.clientId + ":" + c.clientSecret
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func getHttpClient(bearerToken *BearerToken) *http.Client {
	ctx := context.Background()
	token := oauth2.Token{
		AccessToken: bearerToken.AccessToken,
		TokenType:   "Bearer",
	}
	return oauth2.NewClient(ctx, oauth2.StaticTokenSource(&token))
}
