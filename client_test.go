package rundeck_test

import (
	"fmt"
	"os"
	"testing"

	rundeck "github.com/andrewmeissner/go-rundeck"
)

func TestMain(m *testing.M) {
	if os.Getenv("RUNDECK_TOKEN") == "" {
		fmt.Println("Please make sure RUNDECK_TOKEN is set before running tests")
		os.Exit(1)
	}
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

	defaultClient := rundeck.NewClient(nil)

	if defaultClient == nil {
		t.Errorf("default client was not supposed to be nil")
	}

	if defaultClient.RundeckAddr != fmt.Sprintf("http://localhost:4440/api/%d", rundeck.APIVersion24) {
		t.Errorf("default client's rundeck addr is malformed: got %s", defaultClient.RundeckAddr)
	}
}

func TestSetToken(t *testing.T) {
	testTokenOne := "testTokenOne"
	testTokenTwo := "testTokenTwo"
	testURL := "http://localhost:4440"

	configOne := rundeck.Config{
		APIVersion:       rundeck.APIVersion24,
		RundeckAuthToken: testTokenOne,
		ServerURL:        testURL,
	}

	client := rundeck.NewClient(&configOne)

	client.SetAPIToken(testTokenTwo)

	if client.Config.RundeckAuthToken == testTokenOne {
		t.Errorf("token should have been changed")
	}

	if client.Config.RundeckAuthToken != testTokenTwo {
		t.Errorf("token was %s but should have been %s", client.Config.RundeckAuthToken, testTokenTwo)
	}
}
