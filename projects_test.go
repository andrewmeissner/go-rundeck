package rundeck_test

import (
	"testing"

	"github.com/andrewmeissner/go-rundeck"
)

func TestCreateProject(t *testing.T) {
	cli := rundeck.NewClient(rundeck.DefaultConfig())

	name := "Test"
	info, err := cli.Projects().Create(&rundeck.CreateProjectInput{Name: name})
	if err != nil {
		t.Error("project creation failed", err)
	}

	if info.Name != name {
		t.Errorf("project name is incorrect.  expected: %s\tactual: %s\n", name, info.Name)
	}

	if err := cli.Projects().Delete(info.Name); err != nil {
		t.Errorf("failed to delete project %s\t%v\n", info.Name, err)
	}
}

func TestCreateWithDescription(t *testing.T) {
	cli := rundeck.NewClient(rundeck.DefaultConfig())

	name := "Test"
	myDescription := "my description string"

	info, err := cli.Projects().Create(&rundeck.CreateProjectInput{
		Name:        name,
		Description: myDescription,
	})
	if err != nil {
		t.Errorf("failed to create project %s\t%v\n", name, err)
	}

	if info.Description != myDescription {
		t.Error("description failed to load properly")
	}

	if err := cli.Projects().Delete(info.Name); err != nil {
		t.Errorf("project deletion failed: %v\n", err)
	}
}
