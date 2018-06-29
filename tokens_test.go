package rundeck_test

import (
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
}

func TestAdminToken(t *testing.T) {
	client := rundeck.NewClient(rundeck.DefaultConfig())
	tokens, err := client.Tokens().User("admin")
	if err != nil {
		t.Error(err)
	}

	if len(tokens) != 1 {
		t.Errorf("there should be exactly one one token")
	}
}

func TestGetToken(t *testing.T) {
	client := rundeck.NewClient(rundeck.DefaultConfig())
	tokens, err := client.Tokens().List()
	if err != nil {
		t.Error(err)
	}

	if len(tokens) < 1 {
		t.Errorf("there should be at least one token")
	}

	token, err := client.Tokens().Get(tokens[0].ID)
	if err != nil {
		t.Error(err)
	}

	if token.ID != tokens[0].ID {
		t.Errorf("retieved token should be the same as the first token in the list")
	}
}

func TestCreateAndDeleteToken(t *testing.T) {
	client := rundeck.NewClient(rundeck.DefaultConfig())

	duration := "5s"
	token, err := client.Tokens().Create("test", []string{"admin"}, &duration)
	if err != nil {
		t.Error(err)
	}

	tokens, err := client.Tokens().List()
	if err != nil {
		t.Error(err)
	}

	found := false
	for _, listToken := range tokens {
		if listToken.ID == token.ID {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("created token wasn't discovered from a List call")
	}

	err = client.Tokens().Delete(token.ID)
	if err != nil {
		t.Error(err)
	}
}

func TestLargeDurationOnCreateToken(t *testing.T) {
	client := rundeck.NewClient(rundeck.DefaultConfig())
	duration := "30d"

	token, err := client.Tokens().Create("admin", []string{"admin"}, &duration)
	if err != nil {
		t.Error(err)
	}

	if token != nil {
		err = client.Tokens().Delete(token.ID)
		if err != nil {
			t.Error(err)
		}
	}
}
