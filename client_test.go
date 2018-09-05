package rundeck_test

import (
	"net/url"
	"strings"
	"testing"

	"github.com/andrewmeissner/go-rundeck"
)

func TestSetAPITokenForClient(t *testing.T) {
	cli := rundeck.NewClient(nil)

	token := "dev-token"
	cli.SetAPIToken(token)

	if cli.Config.RundeckAuthToken != token {
		t.Errorf("setting token failed.  expected: %s\tactual: %s\n", token, cli.Config.RundeckAuthToken)
	}
}

func TestExtraSlashesOnRundeckAddr(t *testing.T) {
	cli := rundeck.NewClient(&rundeck.Config{
		APIVersion:       rundeck.APIVersion24,
		RundeckAuthToken: "dev-token",
		ServerURL:        "http://localhost:4440/",
	})

	uri, err := url.Parse(cli.RundeckAddr)
	if err != nil {
		t.Error("failed to parse rundeck addr", err)
	}

	if strings.Contains(uri.Path, "//") {
		t.Error("failed to sanitize rundeck addr")
	}
}
