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

func TestSetExecutionMode(t *testing.T) {
	client := rundeck.NewClient(rundeck.DefaultConfig())
	_, err := client.System().SetExecutionMode("invalidMode")
	if err == nil {
		t.Errorf("set execution mode was passed an invalid mode - this should error")
	}

	res, err := client.System().SetExecutionMode(rundeck.ExecutionModePassive)
	if err != nil {
		t.Error(err)
	}
	if res == nil {
		t.Errorf("response should not be nil")
	}
	if res.Active {
		t.Errorf("active should be false on a passive set")
	}
	if res.ExecutionMode != string(rundeck.ExecutionModePassive) {
		t.Errorf("setting execution to passive failed")
	}

	res, err = client.System().SetExecutionMode(rundeck.ExecutionModeActive)
	if err != nil {
		t.Error(err)
	}
	if !res.Active {
		t.Errorf("active should be true on an active set")
	}
	if res.ExecutionMode != string(rundeck.ExecutionModeActive) {
		t.Errorf("setting execution to active failed")
	}
}
