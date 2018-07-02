package rundeck_test

import (
	"testing"

	rundeck "github.com/andrewmeissner/go-rundeck"
)

func TestSystemInfo(t *testing.T) {
	client := rundeck.NewClient(rundeck.DefaultConfig())
	systemInfo, err := client.System().Info()
	if err != nil {
		t.Errorf("system info wasn't supposed to return an error")
	}

	if systemInfo == nil {
		t.Errorf("system info wasn't supposed to be nil")
	}

	if systemInfo.System.Metrics.HREF == "" {
		t.Errorf("unexported struct attributes failed (not really a test)")
	}
}

func TestLogStorage(t *testing.T) {
	client := rundeck.NewClient(rundeck.DefaultConfig())
	logStorage, err := client.System().LogStorage()
	if err != nil {
		t.Error(err)
	}

	if logStorage == nil {
		t.Errorf("logstraoge should have been populated and not nil")
	}
}

func TestIncompleteLogStorage(t *testing.T) {
	client := rundeck.NewClient(rundeck.DefaultConfig())
	incLogStorage, err := client.System().IncompleteLogStorage()
	if err != nil {
		t.Error(err)
	}

	if incLogStorage == nil {
		t.Errorf("incomplete log storage should not have been nil")
	}
}

func TestResumeIncLogStorage(t *testing.T) {
	client := rundeck.NewClient(rundeck.DefaultConfig())
	res, err := client.System().ResumeIncompleteLogStorage()
	if err != nil {
		t.Error(err)
	}

	if !res.Resumed {
		t.Errorf("resumed should have been true for incomplete log storage on a bare-bones rundeck")
	}
}
