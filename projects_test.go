package rundeck_test

import (
	"testing"

	rundeck "github.com/andrewmeissner/go-rundeck"
)

func TestInfo(t *testing.T) {
	project := "GroundZeroTest"

	client := rundeck.NewClient(nil)

	info, err := client.Projects().GetInfo(project)
	if err != nil {
		t.Error(err)
	}

	if info.Name != project {
		t.Errorf("somehow the project's name changed...")
	}
}
