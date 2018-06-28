package api

import (
	"fmt"
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
func NewClient(config *Config) *Client {
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
