package rundeck_test

import (
	"fmt"
	"testing"

	"github.com/andrewmeissner/go-rundeck"
)

func TestListProjects(t *testing.T) {
	cli := rundeck.NewClient(nil)

	numProjects := 3
	projects := make([]*rundeck.ProjectInfo, numProjects)

	for i := 0; i < numProjects; i++ {
		name := fmt.Sprintf("Test-%d", i+1)
		project, err := cli.Projects().Create(&rundeck.CreateProjectInput{
			Name: name,
		})
		if err != nil {
			t.Errorf("failed to create project %s\t%v\n", name, err)
		}
		projects[i] = project
	}

	listedProjects, err := cli.Projects().List()
	if err != nil {
		t.Error(err)
	}

	if len(listedProjects) != numProjects {
		t.Errorf("number of projects from list call was wrong.  execpted: %d\tactual: %d\n", numProjects, len(listedProjects))
	}

	for i, project := range listedProjects {
		name := fmt.Sprintf("Test-%d", i+1)
		if project.Name != name {
			t.Errorf("names didn't match.  expected: %s\tactual%s\n", name, project.Name)
		}
		if err := cli.Projects().Delete(project.Name); err != nil {
			t.Errorf("failed to delete project %s\t%v\n", project.Name, err)
		}
	}
}

func TestCreateProject(t *testing.T) {
	cli := rundeck.NewClient(nil)

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
	cli := rundeck.NewClient(nil)

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

	if err := cli.Projects().Delete(name); err != nil {
		t.Errorf("project deletion failed: %v\n", err)
	}

	info2, err := cli.Projects().Create(&rundeck.CreateProjectInput{
		Name: name,
		Config: map[string]string{
			"project.description": myDescription,
		},
	})
	if err != nil {
		t.Errorf("failed to create secondary project for desc test\t%v\n", err)
	}

	if info2.Description != myDescription {
		t.Error("failed on secondary desc compromise")
	}

	cli.Projects().Delete(name)
}

func TestCreateProjectWithNilData(t *testing.T) {
	cli := rundeck.NewClient(nil)

	project, err := cli.Projects().Create(nil)
	if err == nil {
		t.Errorf("project cannot be created with nil input, but nil was provided: %v\n", err)
	}

	if project != nil {
		t.Errorf("project creation with nil input should yield a nil project.  project == nil? %v\n", project == nil)
	}
}

func TestGetProjectInfo(t *testing.T) {
	cli := rundeck.NewClient(nil)

	name := "Test"

	projectInfo, err := cli.Projects().Create(&rundeck.CreateProjectInput{Name: name})
	if err != nil {
		t.Error("failed to create project", err)
	}

	info, err := cli.Projects().GetInfo(projectInfo.Name)
	if err != nil {
		t.Error("failed to get project info", err)
	}

	info1 := *projectInfo
	info2 := *info

	if info1.Project != info2.Project {
		t.Error("projects do not match")
	}

	if err := cli.Projects().Delete(name); err != nil {
		t.Error("project failed to delete", err)
	}
}

func TestProjectDeletion(t *testing.T) {
	cli := rundeck.NewClient(nil)

	name := "TestForDeletion"

	project, err := cli.Projects().Create(&rundeck.CreateProjectInput{Name: name})
	if err != nil {
		t.Error("failed to create project", name, err)
	}

	if err := cli.Projects().Delete(project.Name); err != nil {
		t.Error("failed to delete project", project.Name, err)
	}
}

func TestGetConfiguration(t *testing.T) {
	cli := rundeck.NewClient(nil)

	name := "Test"

	project, err := cli.Projects().Create(&rundeck.CreateProjectInput{
		Name: name,
	})
	if err != nil {
		t.Error("failed to create project for config retieval", err)
	}

	projectConfig, err := cli.Projects().Configuration(project.Name)
	if err != nil {
		t.Error("failed to retieve project configuration", err)
	}

	if projectConfig == nil {
		t.Error("project config is nil, when it probably shouldn't be")
	}

	if err := cli.Projects().Delete(project.Name); err != nil {
		t.Error("project failed to delete", err)
	}
}

func TestConfigureExistingProject(t *testing.T) {
	cli := rundeck.NewClient(nil)

	name := "Test"

	project, err := cli.Projects().Create(&rundeck.CreateProjectInput{Name: name})
	if err != nil {
		t.Error("failed to create project", err)
	}

	description := "updated description"
	config := project.Config
	config["project.description"] = description

	newConfig, err := cli.Projects().Configure(project.Name, config)
	if err != nil {
		t.Error("failed to reconfigure project", err)
	}

	if newConfig == nil {
		t.Error("new config should not be nil")
	}

	val, exists := newConfig["project.description"]
	if !exists {
		t.Error("project.description should exist in config")
	}

	if val != description {
		t.Errorf("project.description doesn't match.  expected: %s\tactual: %s\n", description, val)
	}

	if err := cli.Projects().Delete(project.Name); err != nil {
		t.Error("failed to delete project", err)
	}
}

func TestGetConfigurationKey(t *testing.T) {
	cli := rundeck.NewClient(nil)

	name := "TestGetConfigKey"
	description := "test get config key description"

	project, err := cli.Projects().Create(&rundeck.CreateProjectInput{
		Name:        name,
		Description: description,
	})
	if err != nil {
		t.Error("failed to create project", name, err)
	}

	key := "project.description"
	val, err := cli.Projects().GetConfigKey(project.Name, key)
	if err != nil {
		t.Error("failed to get config key", err)
	}

	if val.Key != key {
		t.Errorf("wrong key pair was returned.  expected: %s\tactual: %s\n", key, val.Key)
	}

	if val.Value != description {
		t.Errorf("values doesn't match.  expected: %s\nactual: %s\n", description, val.Value)
	}

	_, err = cli.Projects().GetConfigKey(project.Name, "keyThatVeryClearlyDoesNotExist")
	if err == nil {
		t.Error("obvious bad config key should fail", err)
	}

	if err := cli.Projects().Delete(project.Name); err != nil {
		t.Error("failed to delete project", err)
	}
}

func TestSetConfigurationKey(t *testing.T) {
	cli := rundeck.NewClient(nil)

	name := "TestSetConfigKey"

	project, err := cli.Projects().Create(&rundeck.CreateProjectInput{Name: name})
	if err != nil {
		t.Error("project creation failed for", name, err)
	}

	description := "setting the description"
	newConfigKeyPair := rundeck.ProjectConfigKeyPair{
		Key:   "project.description",
		Value: description,
	}

	if _, err := cli.Projects().SetConfigKey(project.Name, nil); err == nil {
		t.Error(err)
	}

	setConfigKeyPair, err := cli.Projects().SetConfigKey(project.Name, &newConfigKeyPair)
	if err != nil {
		t.Error("failed to set config key pair", err)
	}

	if newConfigKeyPair != *setConfigKeyPair {
		t.Errorf("set and returned key pairs to not match")
	}

	if err := cli.Projects().Delete(project.Name); err != nil {
		t.Error("project failed to delete", err)
	}
}

func TestDeleteConfigKeyPair(t *testing.T) {
	cli := rundeck.NewClient(nil)

	name := "TestForDeletingConfigKeyPair"
	configKey := "config.to.delete"

	project, err := cli.Projects().Create(&rundeck.CreateProjectInput{
		Name: name,
		Config: map[string]string{
			configKey: "delete me",
		},
	})
	if err != nil {
		t.Error("failed to create project", err)
	}

	if err := cli.Projects().DeleteConfigKey(project.Name, configKey); err != nil {
		t.Error("failed to delete config key", configKey, err)
	}

	configuration, err := cli.Projects().Configuration(project.Name)
	if err != nil {
		t.Error("failed to get configuation from project", err)
	}

	if _, exists := configuration[configKey]; exists {
		t.Error(configKey, " should not exist after deletion")
	}

	if err := cli.Projects().Delete(project.Name); err != nil {
		t.Error("failed to delete project", err)
	}
}

func TestListProjectResources(t *testing.T) {
	cli := rundeck.NewClient(nil)

	name := "Test"

	projectInfo, err := cli.Projects().Create(&rundeck.CreateProjectInput{Name: name})
	if err != nil {
		t.Error("failed to create project", err)
	}

	resources, err := cli.Projects().ListResources(projectInfo.Name, map[string]string{"name": "ubuntu"})
	if err != nil {
		t.Error("resources were not nil", err)
	}

	if resources == nil || len(resources) < 1 {
		t.Error("rundeck will always provide 1 resource - the OS", err)
	}

	if err := cli.Projects().Delete(projectInfo.Name); err != nil {
		t.Error("failed to delete project", err)
	}
}

func TestSynchronousArchiveExport(t *testing.T) {
	cli := rundeck.NewClient(nil)

	name := "Test"
	description := "test project for synchronous archive"

	project, err := cli.Projects().Create(&rundeck.CreateProjectInput{
		Name:        name,
		Description: description,
	})
	if err != nil {
		t.Error("project failed to create", name, err)
	}

	adhocInput := rundeck.AdhocCommandStringInput{
		Exec: "pwd",
		AdhocOptions: rundeck.AdhocOptions{
			Project: project.Name,
		},
	}
	adhocResp, err := cli.Adhoc().RunCommandString(&adhocInput)
	if err != nil {
		t.Error("failed to run adhoc command", err)
	}

	exportInput := rundeck.ArchiveExportInput{
		ExecutionIDs:     []int{adhocResp.Execution.ID},
		ExportAcls:       true,
		ExportAll:        true,
		ExportConfigs:    true,
		ExportExecutions: true,
		ExportJobs:       true,
		ExportReadmes:    true,
	}
	_, err = cli.Projects().ArchiveExport(project.Name, &exportInput)
	if err != nil {
		t.Error("failed to export all", err)
	}

	if err := cli.Projects().Delete(project.Name); err != nil {
		t.Error("project failed to delete", err)
	}
}

func TestProjectArchiveImport(t *testing.T) {
	cli := rundeck.NewClient(nil)

	name := "Test"
	description := "test project for archive import"

	project, err := cli.Projects().Create(&rundeck.CreateProjectInput{
		Name:        name,
		Description: description,
	})
	if err != nil {
		t.Error("project creation failed", err)
	}

	content, err := cli.Projects().ArchiveExport(project.Name, &rundeck.ArchiveExportInput{ExportAll: true})
	if err != nil {
		t.Error("project export failed", err)
	}

	newProject, err := cli.Projects().Create(&rundeck.CreateProjectInput{
		Name:        "Test-New",
		Description: "new description",
	})

	response, err := cli.Projects().ArchiveImport(newProject.Name, content, &rundeck.ArchiveImportInput{
		ImportACL:        true,
		ImportConfig:     true,
		ImportExecutions: true,
		JobUUIDOption:    rundeck.UUIDOptionRemove,
	})
	if err != nil {
		t.Error("failed to import archive to new project", err)
	}

	if response.ImportStatus != rundeck.StatusSuccessful {
		t.Error("import failed", response.ImportStatus)
		t.Error("acl errors:", response.ACLErrors)
		t.Error("errors:", response.Errors)
		t.Error("execution errors:", response.ExecutionErrors)
	}

	for _, projectName := range []string{project.Name, newProject.Name} {
		if err := cli.Projects().Delete(projectName); err != nil {
			t.Error("failed to delete project", projectName, err)
		}
	}
}

func TestAsyncArchiveExport(t *testing.T) {
	cli := rundeck.NewClient(nil)

	projectInput := rundeck.CreateProjectInput{
		Name:        "Test",
		Description: "test for asynchronous archive export",
	}
	project, err := cli.Projects().Create(&projectInput)
	if err != nil {
		t.Error("project creation failed", err)
	}

	archiveExportInput := rundeck.ArchiveExportInput{ExportAll: true}
	asyncStatusResponse, err := cli.Projects().ArchiveExportAsync(project.Name, &archiveExportInput)
	if err != nil {
		t.Error("failed to initiate async archive export", err)
	}
	token := asyncStatusResponse.Token

	// ensures the call is made at least one time
	status, err := cli.Projects().ArchiveExportAsyncStatus(project.Name, token)
	if err != nil {
		t.Error("failed to get status from seconary call", err)
	}

	ready := status.Ready
	for !ready {
		s, err := cli.Projects().ArchiveExportAsyncStatus(project.Name, token)
		if err != nil {
			t.Error("failed to get status from call", err)
		}
		ready = s.Ready
	}

	_, err = cli.Projects().ArchiveExportAsyncDownload(project.Name, token)
	if err != nil {
		t.Error("failed to download async archive", err)
	}

	if err := cli.Projects().Delete(project.Name); err != nil {
		t.Error("failed to delete project", err)
	}
}
