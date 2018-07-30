package rundeck_test

import (
	"testing"

	rundeck "github.com/andrewmeissner/go-rundeck"
)

func TestInfo(t *testing.T) {
	project := "TestProject"

	client := rundeck.NewClient(nil)

	client.Projects().Create(&rundeck.CreateProjectInput{
		Name: "project",
	})

	info, err := client.Projects().GetInfo(project)
	if err != nil {
		t.Error(err)
	}

	if info.Name != project {
		t.Errorf("somehow the project's name changed...")
	}

	client.Projects().Delete(project)
}
