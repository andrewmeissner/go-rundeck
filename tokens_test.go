package rundeck_test

import (
	"encoding/json"
	"fmt"
	"testing"

	rundeck "github.com/andrewmeissner/go-rundeck"
)

func TestListTokens(t *testing.T) {
	client := rundeck.NewClient(rundeck.DefaultConfig())
	tokens, err := client.Tokens().List()
	if err != nil {
		t.Error(err)
	}

	if len(tokens) < 1 {
		t.Errorf("there should be at least one token")
	}

	bs, err := json.Marshal(tokens)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(bs))
}

func TestAdminToken(t *testing.T) {
	client := rundeck.NewClient(rundeck.DefaultConfig())
	tokens, err := client.Tokens().User("admin")
	if err != nil {
		t.Error(err)
	}

	if len(tokens) < 1 {
		t.Errorf("there should be at least one token")
	}

	bs, err := json.Marshal(tokens)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(bs))
}
