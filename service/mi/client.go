package mi

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

const (
	serverEndpoint = ""
)

type client struct {
	serverURL *url.URL

	httpClient http.Client
}

func NewRouter(serverAddr string) (*client, error) {
	serverURL, err := url.Parse(serverAddr)
	if err != nil {
		return nil, fmt.Errorf("invalid server address: %v", err)
	}

	return &client{
		serverURL: serverURL,
	}, nil
}

func (c *client) Connect() error {
	addr, err := c.apiUrl(serverEndpoint)
	if err != nil {
		return fmt.Errorf("invalid connection address: %v", err)
	}

	response, err := c.httpClient.Get(addr)
	if err != nil {
		return err
	}

	log.Printf("Request '%s' result: %v", addr, response)
	return nil
}

func (c *client) apiUrl(path string) (string, error) {
	apiUrl, err := c.serverURL.Parse(path)
	if err != nil {
		return "", fmt.Errorf("failed to get url for server URL '%s' and path '%s': %v", c.serverURL.String(), path, err)
	}

	return apiUrl.String(), nil
}
