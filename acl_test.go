package rundeck_test

import (
	"testing"

	"github.com/andrewmeissner/go-rundeck"
)

var testPolicy = `description: Admin project level access control. Applies to resources within a specific project.
context:
  project: '.*' # all projects
for:
  resource:
    - equals:
        kind: job
      allow: [create] # allow create jobs
    - equals:
        kind: node
      allow: [read,create,update,refresh] # allow refresh node sources
    - equals:
        kind: event
      allow: [read,create] # allow read/create events
  adhoc:
    - allow: [read,run,runAs,kill,killAs] # allow running/killing adhoc jobs
  job:
    - allow: [create,read,update,delete,run,runAs,kill,killAs] # allow create/read/write/delete/run/kill of all jobs
  node:
    - allow: [read,run] # allow read/run for nodes
by:
  group: admin
`

var updatedPolicy = `description: Admin project level access control. Applies to resources within a specific project.
context:
  project: 'Test Project'
for:
  resource:
    - equals:
        kind: job
      allow: [create] # allow create jobs
    - equals:
        kind: node
      allow: [read,create,update,refresh] # allow refresh node sources
    - equals:
        kind: event
      allow: [read,create] # allow read/create events
  adhoc:
    - allow: [read,run,runAs,kill,killAs] # allow running/killing adhoc jobs
  job:
    - allow: [create,read,update,delete,run,runAs,kill,killAs] # allow create/read/write/delete/run/kill of all jobs
  node:
    - allow: [read,run] # allow read/run for nodes
by:
  group: admin
`

func TestCreateACL(t *testing.T) {
	if err := rundeck.NewClient(rundeck.DefaultConfig()).ACL().Create("test", []byte(testPolicy)); err != nil {
		t.Error(err)
	}
}

func TestUpdateACL(t *testing.T) {
	if err := rundeck.NewClient(rundeck.DefaultConfig()).ACL().Update("test", []byte(testPolicy)); err != nil {
		t.Error(err)
	}
}

func TestListACL(t *testing.T) {
	client := rundeck.NewClient(rundeck.DefaultConfig())
	aclRes, err := client.ACL().List()
	if err != nil {
		t.Error(err)
	}

	if len(aclRes.Resources) == 0 {
		t.Errorf("expected more than 0 acl resources")
	}
}

func TestDeleteACL(t *testing.T) {
	client := rundeck.NewClient(rundeck.DefaultConfig())
	if err := client.ACL().Delete("test.aclpolicy"); err != nil {
		t.Error(err)
	}

	aclRes, err := client.ACL().List()
	if err != nil {
		t.Error(err)
	}

	if len(aclRes.Resources) != 0 {
		t.Errorf("deleted resources should be empty")
	}
}
