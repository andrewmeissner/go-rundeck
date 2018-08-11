package rundeck

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"strings"
)

const aclPolicySuffix = ".aclpolicy"

// ListACLsResponse is the response body from the system/acl endpoint
type ListACLsResponse struct {
	Path      string         `json:"path"`
	Type      string         `json:"type"`
	HREF      string         `json:"href"`
	Resources []*ACLResource `json:"resources"`
}

// ACLResource is an element in the resources portion of the ListACLsReponse payload.
type ACLResource struct {
	Path string `json:"path"`
	Type string `json:"type"`
	Name string `json:"name"`
	HREF string `json:"href"`
}

// ACL in information about API submitted ACLs, not hard-disk ACLs
type ACL struct {
	c *Client
}

// ACL interacts with the ACL API
func (c *Client) ACL() *ACL {
	return &ACL{c: c}
}

// List returns an overview of the API submitted ACLs
func (a *ACL) List() (*ListACLsResponse, error) {
	url := a.c.RundeckAddr + "/system/acl/"

	res, err := a.c.checkResponseOK(a.c.get(url))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var listACLs ListACLsResponse
	return &listACLs, json.NewDecoder(res.Body).Decode(&listACLs)
}

// Get retrieves the YAML text of the ACL Policy file.  The contents of the file as a []byte will be returned.
func (a *ACL) Get(name string) ([]byte, error) {
	url := a.c.RundeckAddr + "/system/acl/" + a.sanitizeACLName(name)

	res, err := a.c.checkResponseOK(a.c.getWithAdditionalHeaders(url, map[string]string{"Accept": "text/plain"}))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}

// Create is used to create an ACL policy
func (a *ACL) Create(name string, policy []byte) error {
	url := a.c.RundeckAddr + "/system/acl/" + a.sanitizeACLName(name)

	res, err := a.c.checkResponseCreated(a.c.postWithAdditionalHeaders(url, map[string]string{"Content-Type": "text/plain"}, bytes.NewReader(policy)))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

// Update updates an existing acl policy
func (a *ACL) Update(name string, policy []byte) error {
	url := a.c.RundeckAddr + "/system/acl/" + a.sanitizeACLName(name)

	res, err := a.c.checkResponseOK(a.c.putWithAdditionalHeaders(url, map[string]string{"Content-Type": "text/plain"}, bytes.NewReader(policy)))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

// Delete removes an ACL polciy file
func (a *ACL) Delete(name string) error {
	url := a.c.RundeckAddr + "/system/acl/" + a.sanitizeACLName(name)

	res, err := a.c.checkResponseNoContent(a.c.delete(url, nil))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func (a *ACL) sanitizeACLName(name string) string {
	if !strings.HasSuffix(name, aclPolicySuffix) {
		name += aclPolicySuffix
	}
	return name
}
