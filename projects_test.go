package rundeck_test

import (
	"testing"

	rundeck "github.com/andrewmeissner/go-rundeck"
)

const (
	projectName = "TestProject"
)

func TestInfo(t *testing.T) {
	client := rundeck.NewClient(nil)

	client.Projects().Create(&rundeck.CreateProjectInput{
		Name: "project",
	})

	info, err := client.Projects().GetInfo(projectName)
	if err != nil {
		t.Error(err)
	}

	if info.Name != projectName {
		t.Errorf("somehow the project's name changed...")
	}

	client.Projects().Delete(projectName)
}
