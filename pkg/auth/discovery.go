package auth

import (
	"encoding/json"
	"github.com/nsotgui/quickbooks-go/pkg"
	"io/ioutil"
	"log"
	"net/http"
)

// Call the discovery API.
// See https://developer.intuit.com/app/developer/qbo/docs/develop/authentication-and-authorization/openid-connect#discovery-document
//
func CallDiscoveryAPI(discoveryEndpoint quickbooks.EndpointURL) *DiscoveryAPI {
	log.Println("Entering CallDiscoveryAPI ")
	client := &http.Client{}
	request, err := http.NewRequest("GET", string(discoveryEndpoint), nil)
	if err != nil {
		log.Fatalln(err)
	}
	//set header
	request.Header.Set("accept", "application/json")

	resp, err := client.Do(request)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	discoveryAPIResponse, err := getDiscoveryAPIResponse(body)
	return discoveryAPIResponse
}

type DiscoveryAPI struct {
	Issuer                string `json:"issuer"`
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
	UserinfoEndpoint      string `json:"userinfo_endpoint"`
	RevocationEndpoint    string `json:"revocation_endpoint"`
	JwksUri               string `json:"jwks_uri"`
}

func getDiscoveryAPIResponse(body []byte) (*DiscoveryAPI, error) {
	var s = new(DiscoveryAPI)
	err := json.Unmarshal(body, &s)
	if err != nil {
		log.Fatalln("error getting DiscoveryAPIResponse:", err)
	}
	return s, err
}
