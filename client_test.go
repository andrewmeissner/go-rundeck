package rundeck_test

import (
	"fmt"
	"testing"

	rundeck "github.com/andrewmeissner/go-rundeck"
)

func TestNewClient(t *testing.T) {
	client := rundeck.NewClient(&rundeck.Config{
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

	defaultClient := rundeck.NewClient(nil)

	if defaultClient == nil {
		t.Errorf("default client was not supposed to be nil")
	}

	if defaultClient.RundeckAddr != fmt.Sprintf("http://localhost:4440/api/%d", rundeck.DefaultAPIVersion24) {
		t.Errorf("default client's rundeck addr is malformed: got %s", defaultClient.RundeckAddr)
	}
}
