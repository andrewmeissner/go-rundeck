package rundeck

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// UserProfile is the information relating to a user on Rundeck
type UserProfile struct {
	Login     string `json:"login"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

// ModifyUserInput are the fields necessary to modify a user
type ModifyUserInput struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

// Users the is information regarding user profiles
type Users struct {
	c *Client
}

// Users interacts with the User profile API
func (c *Client) Users() *Users {
	return &Users{c: c}
}

// List returns a list of all the users
func (u *Users) List() ([]*UserProfile, error) {
	url := u.c.RundeckAddr + "/user/list"

	res, err := u.c.checkResponseOK(u.c.get(url))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var users []*UserProfile
	return users, json.NewDecoder(res.Body).Decode(&users)
}

// Get fetches a user profile.
//
// If the login parameter is nil, the profile associated with
// the supplied auth token will be returned.
func (u *Users) Get(login *string) (*UserProfile, error) {
	url := u.c.RundeckAddr + "/user/info"

	if login != nil {
		url += "/" + stringValue(login)
	}

	res, err := u.c.checkResponseOK(u.c.get(url))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var user UserProfile
	return &user, json.NewDecoder(res.Body).Decode(&user)
}

// Modify will modify the user.
//
// If the user parameter is nil, then the user associated with
// the auth token will be modified.
func (u *Users) Modify(login *string, input *ModifyUserInput) (*UserProfile, error) {
	if input == nil {
		return nil, fmt.Errorf("the parameter ModifyUserInput cannot be nil")
	}

	url := u.c.RundeckAddr + "/user/info"

	if login != nil {
		url += "/" + stringValue(login)
	}

	var body io.Reader
	if input != nil {
		bs, err := json.Marshal(input)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(bs)
	}

	res, err := u.c.checkResponseOK(u.c.post(url, body))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var userProfile UserProfile
	return &userProfile, json.NewDecoder(res.Body).Decode(&userProfile)
}
