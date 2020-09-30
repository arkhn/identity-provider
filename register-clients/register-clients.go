package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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

func registerClient(hydra *client.OryHydra, cc models.OAuth2Client) (*admin.CreateOAuth2ClientCreated, error) {

	// Delete previously existing client
	_, err := hydra.Admin.DeleteOAuth2Client(
		&admin.DeleteOAuth2ClientParams{ID: cc.ClientID, Context: context.Background()},
	)
	if err != nil {
		return nil, err
	}

	// Create a new client
	created, err := hydra.Admin.CreateOAuth2Client(admin.NewCreateOAuth2ClientParams().WithBody(&cc))
	if err != nil {
		switch e := err.(type) {
		case *admin.CreateOAuth2ClientConflict:
			return nil, fmt.Errorf("Client %s already exists: %s\n", cc.ClientID, e.GetPayload().ErrorDescription)
		default:
			return nil, err
		}
	}
	return created, nil
}
