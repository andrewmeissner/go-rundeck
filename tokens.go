package rundeck

import (
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
