package auth

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
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
func (c *Client) RetrieveBearerToken(authorizationCode string) (*BearerToken, error) {
	log.Println("Entering RetrieveBearerToken ")
	client := &http.Client{}
	data := url.Values{}
	//set parameters
	data.Set("grant_type", "authorization_code")
	data.Add("code", authorizationCode)
	data.Add("redirect_uri", "https://developer.intuit.com/v2/OAuth2Playground/RedirectUrl")

	request, err := http.NewRequest("POST", string(c.DiscoveryAPI.TokenEndpoint), bytes.NewBufferString(data.Encode()))
	if err != nil {
		log.Fatalln(err)
	}
	//set headers
	request.Header.Set("accept", "application/json")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	request.Header.Set("Authorization", "Basic "+basicAuth(c))

	resp, err := client.Do(request)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	bearerTokenResponse, err := getBearerTokenResponse([]byte(body))
	log.Println("Exiting RetrieveBearerToken")
	return bearerTokenResponse, err
}

//
// Call the refresh endpoint to generate new tokens
//
func (c *Client) RefreshToken(refreshToken string) (*BearerToken, error) {

	log.Println("Entering RefreshToken ")
	client := &http.Client{}
	data := url.Values{}

	//add parameters
	data.Set("grant_type", "refresh_token")
	data.Add("refresh_token", refreshToken)

	request, err := http.NewRequest("POST", string(c.DiscoveryAPI.TokenEndpoint), bytes.NewBufferString(data.Encode()))
	if err != nil {
		log.Fatalln(err)
	}
	//set the headers
	request.Header.Set("accept", "application/json")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	request.Header.Set("Authorization", "Basic "+basicAuth(c))

	resp, err := client.Do(request)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	bearerTokenResponse, err := getBearerTokenResponse([]byte(body))
	log.Println("Exiting RefreshToken ")
	return bearerTokenResponse, err
}

//
// Call the revoke endpoint to revoke tokens
//
func (c *Client) RevokeToken(refreshToken string) {
	log.Println("Entering RevokeToken ")
	client := &http.Client{}
	data := url.Values{}

	//add parameters
	data.Add("token", refreshToken)

	revokeEndpoint := c.DiscoveryAPI.RevocationEndpoint
	request, err := http.NewRequest("POST", revokeEndpoint, bytes.NewBufferString(data.Encode()))
	if err != nil {
		log.Fatalln(err)
	}
	//set headers
	request.Header.Set("accept", "application/json")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
	request.Header.Set("Authorization", "Basic "+basicAuth(c))

	resp, err := client.Do(request)
	defer resp.Body.Close()

	responseString := map[string]string{"response": "Revoke successful"}
	responseData, _ := json.Marshal(responseString)
	log.Println(responseData)
	log.Println("Exiting RevokeToken ")
}

func getBearerTokenResponse(body []byte) (*BearerToken, error) {
	var s = new(BearerToken)
	err := json.Unmarshal(body, &s)
	if err != nil {
		log.Fatalln("error getting BearerTokenResponse:", err)
	}
	return s, err
}

func basicAuth(c *Client) string {
	auth := c.ClientId + ":" + c.ClientSecret
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
