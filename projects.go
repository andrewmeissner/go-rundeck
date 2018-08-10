package rundeck

import (
	"bytes"
	"encoding/json"
	"errors"
)

// Project is metadata about a project in rundeck
type Project struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

// ProjectInfo contains the configuration for the project in addition to the metadata
type ProjectInfo struct {
	Project
	Config map[string]string `json:"config,omitempty"`
}

// CreateProjectInput is the payload when posting for project creation
type CreateProjectInput struct {
	Name   string            `json:"name"`
	Config map[string]string `json:"config,omitempty"`
}

// Projects is information pertaining to projects API endpoints
type Projects struct {
	c *Client
}

// Projects interacts with endpoints pertaining to projects
func (c *Client) Projects() *Projects {
	return &Projects{c: c}
}

// List returns a list of the projects
func (p *Projects) List() ([]*Project, error) {
	rawURL := p.c.RundeckAddr + "/projects"

	res, err := p.c.checkResponseOK(p.c.get(rawURL))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var projects []*Project
	return projects, json.NewDecoder(res.Body).Decode(&projects)
}

// Create will make a new project
func (p *Projects) Create(data *CreateProjectInput) (*ProjectInfo, error) {
	if data == nil {
		return nil, errors.New("data cannot be nil")
	}

	rawURL := p.c.RundeckAddr + "/projects"

	bs, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	res, err := p.c.checkResponseCreated(p.c.post(rawURL, bytes.NewReader(bs)))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var info ProjectInfo
	return &info, json.NewDecoder(res.Body).Decode(&info)
}

// GetInfo returns project info
func (p *Projects) GetInfo(project string) (*ProjectInfo, error) {
	rawURL := p.c.RundeckAddr + "/project/" + project

	res, err := p.c.checkResponseOK(p.c.get(rawURL))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var info ProjectInfo
	return &info, json.NewDecoder(res.Body).Decode(&info)
}

// Delete removes an existing project
func (p *Projects) Delete(project string) error {
	rawURL := p.c.RundeckAddr + "/project/" + project

	res, err := p.c.checkResponseNoContent(p.c.delete(rawURL, nil))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
