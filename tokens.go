package rundeck

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Token is the information regarding a user token
type Token struct {
	User       string    `json:"user"`
	ID         string    `json:"id"`
	Creator    string    `json:"creator"`
	Expiration time.Time `json:"expiration"`
	Roles      []string  `json:"roles"`
	Expired    bool      `json:"expired"`
}

// Tokens is used to perform token specific API operations
type Tokens struct {
	c *Client
}

// Tokens is used to return the client for token specific API operations
func (c *Client) Tokens() *Tokens {
	return &Tokens{c: c}
}

// List returns all tokens
func (t *Tokens) List() ([]*Token, error) {
	url := fmt.Sprintf("%s/tokens", t.c.RundeckAddr)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := t.c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, makeError(res.Body)
	}

	var tokenList []*Token
	err = json.NewDecoder(res.Body).Decode(&tokenList)
	if err != nil {
		return nil, err
	}

	return tokenList, nil
}

// User returns the tokens associated with the supplied user
func (t *Tokens) User(user string) ([]*Token, error) {
	url := fmt.Sprintf("%s/tokens/%s", t.c.RundeckAddr, user)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := t.c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, makeError(res.Body)
	}

	var tokenList []*Token
	err = json.NewDecoder(res.Body).Decode(&tokenList)
	if err != nil {
		return nil, err
	}

	return tokenList, nil
}

// Get returns the token by the supplied id
func (t *Tokens) Get(id string) (*Token, error) {
	url := fmt.Sprintf("%s/token/%s", t.c.RundeckAddr, id)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := t.c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, makeError(res.Body)
	}

	var token Token
	err = json.NewDecoder(res.Body).Decode(&token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

// Create creates a token based on the supplied config.
//
// If duration is nil, Rundeck will use the configured default.
// NOTE: the duration needs to be something that rundeck can understand.
// Unfortunately, this isn't a go parseable duration.  "120d" is understood by Rundeck
// while "2880h0m0s" is not (what time.Duration.String() returns for the equivalence).
func (t *Tokens) Create(user string, roles []string, duration *string) (*Token, error) {
	url := fmt.Sprintf("%s/tokens", t.c.RundeckAddr)

	payload := map[string]interface{}{
		"user":  user,
		"roles": roles,
	}

	if duration != nil {
		payload["duration"] = stringValue(duration)
	}

	bs, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}

	res, err := t.c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		return nil, makeError(res.Body)
	}

	var token Token
	err = json.NewDecoder(res.Body).Decode(&token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

// Delete deletes a token
func (t *Tokens) Delete(id string) error {
	url := fmt.Sprintf("%s/token/%s", t.c.RundeckAddr, id)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	res, err := t.c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		return makeError(res.Body)
	}

	return nil
}
