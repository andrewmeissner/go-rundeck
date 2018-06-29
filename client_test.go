package rundeck_test

import (
	"fmt"
	"os"
	"testing"

	rundeck "github.com/andrewmeissner/go-rundeck"
)

func TestMain(m *testing.M) {
	os.Setenv("RUNDECK_SERVER_URL", "http://localhost:4440//")
	os.Setenv("RUNDECK_TOKEN", "env-test-token")

	os.Exit(m.Run())
}

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

	defaultClient := rundeck.NewClient(rundeck.DefaultConfig())

	if defaultClient == nil {
		t.Errorf("default client was not supposed to be nil")
	}

	if defaultClient.RundeckAddr != fmt.Sprintf("http://localhost:4440/api/%d", rundeck.DefaultAPIVersion23) {
		t.Errorf("default config's rundeck addr is malformed")
	}

	os.Clearenv()

	superDefaultClient := rundeck.NewClient(nil)

	if superDefaultClient == nil {
		t.Errorf("super default client was not supposed to be nil")
	}

	if superDefaultClient.RundeckAddr != fmt.Sprintf("http://127.0.0.1:4440/api/%d", rundeck.DefaultAPIVersion23) {
		t.Errorf("super default client's rundeck addr is malformed")
	}
}
