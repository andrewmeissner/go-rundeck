package rundeck_test

import (
	"testing"

	"github.com/andrewmeissner/go-rundeck"
)

var (
	scriptURL = "https://gist.githubusercontent.com/andrewmeissner/2a970d96a51a374f0e02f713eb445b7b/raw/58f2e1b8dce85e005724e22a24140bbab226e398/test.sh"
	script    = `#!/bin/bash
pwd`
)

func TestRunCommandString(t *testing.T) {
	cli := rundeck.NewClient(nil)

	project, err := cli.Projects().Create(&rundeck.CreateProjectInput{Name: "Test"})
	if err != nil {
		t.Error("failed to create project", err)
	}

	_, err = cli.Adhoc().RunCommandString(nil)
	if err == nil {
		t.Error("err should say input cannot be nil")
	}

	input := rundeck.AdhocCommandStringInput{}

	_, err = cli.Adhoc().RunCommandString(&input)
	if err == nil {
		t.Error("err should say project cannot be empty")
	}

	input.Project = project.Name
	_, err = cli.Adhoc().RunCommandString(&input)
	if err == nil {
		t.Error("err should say that Exec cannot be nil")
	}

	input.Exec = "pwd"
	_, err = cli.Adhoc().RunCommandString(&input)
	if err != nil {
		t.Error("adhoc string command failed", err)
	}

	if err := cli.Projects().Delete(project.Name); err != nil {
		t.Error("failed to delete project", err)
		t.FailNow()
	}
}
