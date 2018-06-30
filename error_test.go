package rundeck_test

import (
	"testing"

	rundeck "github.com/andrewmeissner/go-rundeck"
)

func TestErrorResponse(t *testing.T) {
	client := rundeck.NewClient(rundeck.DefaultConfig())
	_, err := client.Tokens().Get("bad-id")
	if err == nil {
		t.Errorf("error should have been nil")
	}

	if err.Error() == "" {
		t.Errorf("error should not be an empty string")
	}
}
