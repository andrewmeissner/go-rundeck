package api_test

import (
	"testing"

	"github.com/andrewmeissner/go-rundeck/api"
)

func TestNewClient(t *testing.T) {
	client := api.NewClient(&api.Config{
		APIVersion:       21,
		RundeckAuthToken: "test-token",
		ServerURL:        "http://localhost:4440/",
	})

	if client == nil {
		t.Errorf("client was not supposed to be nil")
	}

	if client.RundeckAddr != "http://localhost:4440/api/21" {
		t.Errorf("client's RundeckAddr was malformed")
	}
}
