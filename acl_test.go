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
	client := rundeck.NewClient(rundeck.DefaultConfig())

	if err := client.ACL().Create("test", []byte(testPolicy)); err != nil {
		t.Error(err)
	}

	client.SetAPIToken("badToken")

	if err := client.ACL().Create("test", []byte("")); err == nil {
		t.Errorf("badToken was submitted to create acl - this should return an unauthorized rundeck body")
	}
}

func TestGetACL(t *testing.T) {
	client := rundeck.NewClient(rundeck.DefaultConfig())
	bs, err := client.ACL().Get("test")
	if err != nil {
		t.Error(err)
	}

	testBS := []byte(testPolicy)

	if len(bs) != len(testBS) {
		t.Errorf("policy from GET is different length than policy from UPDATE")
	}

	for i := 0; i < len(bs); i++ {
		if bs[i] != testBS[i] {
			t.Errorf("byte comparison of policies yielded different results")
		}
	}

	client.SetAPIToken("badToken")
	_, badTokenErr := client.ACL().Get("test")
	if badTokenErr == nil {
		t.Errorf("badToken was supplied to get acl - this should be a rundeck error")
	}
}

func TestUpdateACL(t *testing.T) {
	client := rundeck.NewClient(rundeck.DefaultConfig())
	if err := client.ACL().Update("test", []byte(testPolicy)); err != nil {
		t.Error(err)
	}

	client.SetAPIToken("badToken")
	if err := client.ACL().Update("test", []byte("")); err == nil {
		t.Errorf("badToken was supplied to update acl - this should fail with rundeck unauthorized")
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

	client.SetAPIToken("badToken")
	aclBadRes, err := client.ACL().List()
	if err == nil {
		t.Errorf("badToken was submitted to acl list - this should return a rundeck unauthorized error")
	}

	if aclBadRes != nil {
		t.Errorf("badToken was submitted to acl list - this should return a rundeck unauthorized error")
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

	client.SetAPIToken("badToken")
	if err := client.ACL().Delete("test"); err == nil {
		t.Errorf("badToken was supplied to delete acl - this should be a rundeck auth failure")
	}
}
