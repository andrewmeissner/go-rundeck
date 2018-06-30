package rundeck_test

import (
	"fmt"
	"testing"
	"time"

	rundeck "github.com/andrewmeissner/go-rundeck"
)

func TestUsersList(t *testing.T) {
	client := rundeck.NewClient(rundeck.DefaultConfig())
	users, err := client.Users().List()
	if err != nil {
		t.Error(err)
	}

	if len(users) == 0 {
		t.Errorf("users list was empty, but there should have been at least 1 entry")
	}
}

func TestUsersGet(t *testing.T) {
	client := rundeck.NewClient(rundeck.DefaultConfig())
	currentUser, err := client.Users().Get(nil)
	if err != nil {
		t.Error(err)
	}

	if currentUser.Login != "admin" {
		t.Errorf("expected the admin user as the supplied token SHOULD have been the admin token (from README)")
	}

	testUserLogin := "test"
	testUser, err := client.Users().Get(&testUserLogin)
	if err != nil {
		t.Error(err)
	}

	if testUser.Login != testUserLogin {
		t.Errorf("test user was supposed to be fetched")
	}
}

func TestModifyUser(t *testing.T) {
	client := rundeck.NewClient(rundeck.DefaultConfig())
	testUserLogin := "test"
	testUser, err := client.Users().Get(&testUserLogin)
	if err != nil {
		t.Error(err)
	}

	newFirstName := fmt.Sprintf("%d", time.Now().Unix())

	input := rundeck.ModifyUserInput{
		Email:     testUser.Email,
		FirstName: newFirstName,
		LastName:  testUser.LastName,
	}

	newProfile, err := client.Users().Modify(&testUserLogin, &input)
	if err != nil {
		t.Error(err)
	}

	if testUser.FirstName == newProfile.FirstName {
		t.Errorf("first names were the same when they should have been different")
	}
}
