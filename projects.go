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

// ArchiveExportInput ...
type ArchiveExportInput struct {
	ExecutionIDs     []string
	ExoprtAll        bool
	ExportJobs       bool
	ExportExecutions bool
	ExportConfigs    bool
	ExportReadmes    bool
	ExportAcls       bool
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

// Configuration retrieves the project configuration data
func (p *Projects) Configuration(project string) (map[string]string, error) {
	rawURL := p.c.RundeckAddr + "/project/" + project + "/config"

	res, err := p.c.checkResponseOK(p.c.get(rawURL))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var config map[string]string
	return config, json.NewDecoder(res.Body).Decode(&config)
}

// Configure modifies the project configuration data
func (p *Projects) Configure(project string, config map[string]string) (map[string]string, error) {
	rawURL := p.c.RundeckAddr + "/project/" + project + "/config"

	bs, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	res, err := p.c.checkResponseOK(p.c.put(rawURL, bytes.NewReader(bs)))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var conf map[string]string
	return conf, json.NewDecoder(res.Body).Decode(&conf)
}

// GetConfigKey retieves the value
func (p *Projects) GetConfigKey(project, key string) (map[string]string, error) {
	rawURL := p.c.RundeckAddr + "/project/" + project + "/config/" + key

	res, err := p.c.checkResponseOK(p.c.get(rawURL))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]string
	return result, json.NewDecoder(res.Body).Decode(&result)
}

// SetConfigKey modifies the value
func (p *Projects) SetConfigKey(project, key, value string) (map[string]string, error) {
	rawURL := p.c.RundeckAddr + "/project/" + project + "/config/" + key

	bs, err := json.Marshal(map[string]string{key: value})
	if err != nil {
		return nil, err
	}

	res, err := p.c.checkResponseOK(p.c.put(rawURL, bytes.NewReader(bs)))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]string
	return result, json.NewDecoder(res.Body).Decode(&result)
}

// DeleteConfigKey removes the key
func (p *Projects) DeleteConfigKey(project, key string) error {
	rawURL := p.c.RundeckAddr + "/project/" + project + "/config/" + key

	_, err := p.c.checkResponseNoContent(p.c.delete(rawURL, nil))
	if err != nil {
		return err
	}
	return nil
}

// ArchiveExport exports a zip archive of the project synchronously
func (p *Projects) ArchiveExport(project string, input *ArchiveExportInput) error {
	return errors.New("not implemented yet")
}
