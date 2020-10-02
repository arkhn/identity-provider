package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/ory/hydra-client-go/client"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
)

func main() {
	clientConfigFile := os.Getenv("CLIENTS_FILE")
	hydraURL := os.Getenv("HYDRA_URL")

	clientURL, err := url.Parse(hydraURL)
	if err != nil {
		panic(err)
	}
	hydra := client.NewHTTPClientWithConfig(nil, &client.TransportConfig{
		Schemes: []string{clientURL.Scheme}, Host: clientURL.Host, BasePath: clientURL.Path,
	})

	byteValue, err := ioutil.ReadFile(clientConfigFile)
	if err != nil {
		panic(err)
	}

	userConfigs := []models.OAuth2Client{}
	err = json.Unmarshal(byteValue, &userConfigs)
	if err != nil {
		panic(err)
	}

	for _, userConfig := range userConfigs {
		created, err := registerClient(hydra, userConfig)
		if err != nil {
			log.Println(err)
		} else {
			fmt.Printf("Created client: %v\n", created)
		}
	}
}

type transporter struct {
	*http.Transport
}

func (t *transporter) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("X-Forwarded-Proto", "https")
	return t.Transport.RoundTrip(req)
}

func registerClient(hydra *client.OryHydra, cc models.OAuth2Client) (*admin.CreateOAuth2ClientCreated, error) {

	// Delete previously existing client
	httpClient := &http.Client{
		Transport: &transporter{
			Transport: &http.Transport{},
		},
	}
	deleteParams := admin.NewDeleteOAuth2ClientParams().
		WithID(cc.ClientID).
		WithContext(context.Background()).
		WithHTTPClient(httpClient)
	_, err := hydra.Admin.DeleteOAuth2Client(deleteParams)
	if err != nil {
		return nil, err
	}

	// Create a new client
	createParams := admin.NewCreateOAuth2ClientParams().
		WithBody(&cc).
		WithHTTPClient(httpClient)
	created, err := hydra.Admin.CreateOAuth2Client(createParams)
	if err != nil {
		switch e := err.(type) {
		case *admin.CreateOAuth2ClientConflict:
			return nil, fmt.Errorf("client %s already exists: %s", cc.ClientID, e.GetPayload().ErrorDescription)
		default:
			return nil, err
		}
	}
	return created, nil
}
