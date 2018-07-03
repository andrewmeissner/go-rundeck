package rundeck

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Client is the basic client that interacts with the Rundeck API.
type Client struct {
	Config      *Config
	client      *http.Client
	RundeckAddr string
}

// NewClient returns a rundeck client
//
// If config is nil, then the configuration from DefaultConfig() will be used.
// DefaultConfig() assumes that the environment variable RUNDECK_TOKEN is set, and
// that its value is a valid Rundeck API token.
func NewClient(config *Config) *Client {
	if config == nil {
		config = DefaultConfig()
	}

	return &Client{
		Config: config,
		client: &http.Client{
			Jar: http.DefaultClient.Jar,
			Transport: &rundeckTransport{
				apiToken:            config.RundeckAuthToken,
				underlyingTransport: http.DefaultTransport,
			},
		},
		RundeckAddr: fmt.Sprintf("%s/api/%d", sanitizeAddr(config.ServerURL), config.APIVersion),
	}
}

// sanitizeAddr will remove all trailing slashes from the supplied ServerURL to ensure path correctness
func sanitizeAddr(addr string) string {
	for strings.HasSuffix(addr, "/") {
		addr = strings.TrimSuffix(addr, "/")
	}
	return addr
}

func (c *Client) get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return c.client.Do(req)
}

func (c *Client) post(url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	return c.client.Do(req)
}

func (c *Client) put(url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, err
	}

	return c.client.Do(req)
}

func (c *Client) delete(url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodDelete, url, body)
	if err != nil {
		return nil, err
	}

	return c.client.Do(req)
}
