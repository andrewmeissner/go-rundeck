package rundeck

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"strings"
)

// Status indicates the general status of an operation
type Status string

const (
	StatusSuccessful Status = "successful"
	StatusFailed     Status = "failed"
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
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Config      map[string]string `json:"config,omitempty"`
}

// ProjectConfigKeyPair is a single key pair from a project configuration
type ProjectConfigKeyPair struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// ArchiveExportInput ...
type ArchiveExportInput struct {
	ExecutionIDs     []int
	ExportAll        bool
	ExportJobs       bool
	ExportExecutions bool
	ExportConfigs    bool
	ExportReadmes    bool
	ExportAcls       bool
}

// ArchiveExportAsyncStatusResponse struct
type ArchiveExportAsyncStatusResponse struct {
	Token      string `json:"token"`
	Ready      bool   `json:"ready"`
	Percentage int    `json:"int"`
}

// ArchiveImportInput are option parameters for importing a project archive
type ArchiveImportInput struct {
	JobUUIDOption    UUIDOption
	ImportExecutions bool
	ImportConfig     bool
	ImportACL        bool
}

// ArchiveImportResponse ...
type ArchiveImportResponse struct {
	ImportStatus    Status   `json:"import_status"`
	Errors          []string `json:"errors"`
	ExecutionErrors []string `json:"execution_errors"`
	ACLErrors       []string `json:"acl_errors"`
}

// NodeEntryAnything represents everything else that can be added to the node entry map
type NodeEntryAnything map[string]string

// NodeEntry contains some specific entries in the node entry map
type NodeEntry struct {
	Nodename    string `json:"nodename"`
	Hostname    string `json:"hostname"`
	Username    string `json:"username"`
	Description string `json:"description"`
	Tags        string `json:"tags"`
	OSFamily    string `json:"osFamily"`
	OSArch      string `json:"osArch"`
	OSName      string `json:"osName"`
	OSVersion   string `json:"osVersion"`
	EditURL     string `json:"editUrl"`
	RemoteURL   string `json:"remoteUrl"`
	NodeEntryAnything
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
func (p *Projects) GetConfigKey(project, key string) (*ProjectConfigKeyPair, error) {
	rawURL := p.c.RundeckAddr + "/project/" + project + "/config/" + key

	res, err := p.c.checkResponseOK(p.c.get(rawURL))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result ProjectConfigKeyPair
	return &result, json.NewDecoder(res.Body).Decode(&result)
}

// SetConfigKey modifies the value
func (p *Projects) SetConfigKey(project string, keyPair *ProjectConfigKeyPair) (*ProjectConfigKeyPair, error) {
	if keyPair == nil {
		return nil, errors.New("keyPair cannot be nil when setting a config key")
	}

	bs, err := json.Marshal(keyPair)
	if err != nil {
		return nil, err
	}

	rawURL := p.c.RundeckAddr + "/project/" + project + "/config/" + keyPair.Key
	res, err := p.c.checkResponseOK(p.c.put(rawURL, bytes.NewReader(bs)))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result ProjectConfigKeyPair
	return &result, json.NewDecoder(res.Body).Decode(&result)
}

// DeleteConfigKey removes the key
func (p *Projects) DeleteConfigKey(project, key string) error {
	rawURL := p.c.RundeckAddr + "/project/" + project + "/config/" + key
	_, err := p.c.checkResponseNoContent(p.c.delete(rawURL, nil))
	return err
}

// ArchiveExport exports a zip archive of the project synchronously
func (p *Projects) ArchiveExport(project string, input *ArchiveExportInput) ([]byte, error) {
	rawURL := p.c.RundeckAddr + "/project/" + project + "/export"

	uri, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	uri.RawQuery = p.encodeArchiveExportInput(uri.Query(), input)

	res, err := p.c.checkResponseOK(p.c.get(uri.String()))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}

// ArchiveExportAsync exports a zip archive of the project asynchronously
func (p *Projects) ArchiveExportAsync(project string, input *ArchiveExportInput) (*ArchiveExportAsyncStatusResponse, error) {
	rawURL := p.c.RundeckAddr + "/project/" + project + "/export/async"

	uri, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	uri.RawQuery = p.encodeArchiveExportInput(uri.Query(), input)

	res, err := p.c.checkResponseOK(p.c.get(uri.String()))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var status ArchiveExportAsyncStatusResponse
	return &status, json.NewDecoder(res.Body).Decode(&status)
}

// ArchiveExportAsyncStatus gets the status of the async archive
func (p *Projects) ArchiveExportAsyncStatus(project, token string) (*ArchiveExportAsyncStatusResponse, error) {
	rawURL := p.c.RundeckAddr + "/project/" + project + "/export/status/" + token

	res, err := p.c.checkResponseOK(p.c.get(rawURL))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var status ArchiveExportAsyncStatusResponse
	return &status, json.NewDecoder(res.Body).Decode(&status)
}

// ArchiveExportAsyncDownload downloads the finished artifact
func (p *Projects) ArchiveExportAsyncDownload(project, token string) ([]byte, error) {
	status, err := p.ArchiveExportAsyncStatus(project, token)
	if err != nil {
		return nil, err
	}

	if !status.Ready {
		return nil, fmt.Errorf("archive is only %d%% complete", status.Percentage)
	}

	rawURL := p.c.RundeckAddr + "/project/" + project + "/export/download/" + token
	res, err := p.c.checkResponseOK(p.c.get(rawURL))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}

// ArchiveImport imports a zip archive into the project
func (p *Projects) ArchiveImport(project string, content []byte, input *ArchiveImportInput) (*ArchiveImportResponse, error) {
	rawURL := p.c.RundeckAddr + "/project/" + project + "/import"

	uri, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	query := uri.Query()

	if input != nil {
		if input.JobUUIDOption == UUIDOptionRemove || input.JobUUIDOption == UUIDOptionPreserve {
			query.Add("jobUuidOption", string(input.JobUUIDOption))
		}

		if input.ImportACL {
			query.Add("importACL", "true")
		}

		if input.ImportConfig {
			query.Add("importConfig", "true")
		}

		if input.ImportExecutions {
			query.Add("importExecutions", "true")
		}
	}

	uri.RawQuery = query.Encode()

	headers := map[string]string{
		"Content-Type": "application/zip",
	}

	res, err := p.c.checkResponseOK(p.c.putWithAdditionalHeaders(uri.String(), headers, bytes.NewReader(content)))
	if err != nil {
		return nil, err
	}

	var response ArchiveImportResponse
	return &response, json.NewDecoder(res.Body).Decode(&response)
}

// ListResources lists resources for a given project
func (p *Projects) ListResources(project string, nodeFilters map[string]string) (map[string]*NodeEntry, error) {
	rawURL := p.c.RundeckAddr + "/project/" + project + "/resources"

	uri, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	query := uri.Query()

	if nodeFilters != nil && len(nodeFilters) > 0 {
		query.Add("filter", p.c.convertFiltersToSerializeableFormat(nodeFilters))
	}

	uri.RawQuery = query.Encode()

	res, err := p.c.checkResponseOK(p.c.get(uri.String()))
	if err != nil {
		return nil, err
	}

	var entries map[string]*NodeEntry
	return entries, json.NewDecoder(res.Body).Decode(&entries)
}

func (p *Projects) encodeArchiveExportInput(query url.Values, input *ArchiveExportInput) string {
	if input != nil {
		if input.ExecutionIDs != nil && len(input.ExecutionIDs) > 0 {
			query.Add("executionIds", strings.Join(intSliceToStringSlice(input.ExecutionIDs), ","))
		}

		if input.ExportAcls {
			query.Add("exportAcls", "true")
		}

		if input.ExportAll {
			query.Add("exportAll", "true")
		}

		if input.ExportConfigs {
			query.Add("exportConfigs", "true")
		}

		if input.ExportExecutions {
			query.Add("exportExecutions", "true")
		}

		if input.ExportJobs {
			query.Add("exportJobs", "true")
		}

		if input.ExportReadmes {
			query.Add("exportReadmes", "true")
		}
	}
	return query.Encode()
}
