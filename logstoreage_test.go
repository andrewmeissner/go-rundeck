package rundeck_test

import (
	"testing"

	rundeck "github.com/andrewmeissner/go-rundeck"
)

func TestLogStorage(t *testing.T) {
	client := rundeck.NewClient(rundeck.DefaultConfig())
	logStorage, err := client.LogStore().LogStorage()
	if err != nil {
		t.Error(err)
	}

	if logStorage == nil {
		t.Errorf("logstraoge should have been populated and not nil")
	}
}

func TestIncompleteLogStorage(t *testing.T) {
	client := rundeck.NewClient(rundeck.DefaultConfig())
	incLogStorage, err := client.LogStore().IncompleteLogStorage()
	if err != nil {
		t.Error(err)
	}

	if incLogStorage == nil {
		t.Errorf("incomplete log storage should not have been nil")
	}
}

func TestResumeIncLogStorage(t *testing.T) {
	client := rundeck.NewClient(rundeck.DefaultConfig())
	res, err := client.LogStore().ResumeIncompleteLogStorage()
	if err != nil {
		t.Error(err)
	}

	if !res.Resumed {
		t.Errorf("resumed should have been true for incomplete log storage on a bare-bones rundeck")
	}
}
