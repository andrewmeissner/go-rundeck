package rundeck

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	JobLogLevelDebug      = "DEBUG"
	JobLogLevelVerbose    = "VERBOSE"
	JobLogLevelInfo       = "INFO"
	JobLogLevelWarn       = "WARN"
	JobLogLevelError      = "ERROR"
	JobFormatXML          = "xml"
	JobFormatYAML         = "yaml"
	DuplicateOptionSkip   = "skip"
	DuplicateOptionCreate = "create"
	DuplicateoptionUpdate = "update"
	UUIDOptionPreserve    = "preserve"
	UUIDOptionRemove      = "remove"
)

// Job is information about a Rundeck job
type Job struct {
	ID              string            `json:"id,omitempty"`
	AverageDuration int64             `json:"averageDuration,omitempty"`
	Name            string            `json:"name,omitempty"`
	Group           string            `json:"group,omitempty"`
	Project         string            `json:"project,omitempty"`
	Description     string            `json:"description,omitempty"`
	HREF            string            `json:"href,omitempty"`
	Permalink       string            `json:"permalink,omitempty"`
	Options         map[string]string `json:"options,omitempty"`
	Scheduled       bool              `json:"scheduled,omitempty"`
	ScheduleEnabled bool              `json:"scheduleEnabled,omitempty"`
	Enabled         bool              `json:"enabled,omitempty"`
	ServerNodeUUID  string            `json:"serverNodeUUID,omitempty"`
	ServerOwner     bool              `json:"serverOwner,omitempty"`
	Index           int               `json:"index,omitempty"`
}

// ListJobsInput adds parameters to the endpoint for listing jobs
type ListJobsInput struct {
	IDs                  []string
	GroupPath            string
	JobFilter            string
	JobExactFilter       string
	GroupPathExact       string
	ScheduledFilter      bool
	ServerNodeUUIDFilter string
}

// RunJobInput are the optional paramenters passed when running a job
type RunJobInput struct {
	LogLevel  string
	AsUser    string
	Filters   map[string]string
	RunAtTime time.Time
	Options   map[string]string
}

type runJobInputSerializeable struct {
	LogLevel  string            `json:"loglevel,omitempty"`
	AsUser    string            `json:"asUser,omitempty"`
	Filter    string            `json:"filter,omitempty"`
	RunAtTime time.Time         `json:"runAtTime,omitempty"`
	Options   map[string]string `json:"options,omitempty"`
}

// RetryJobInput are the optional parameters passed when retrying a job based on execution id
type RetryJobInput struct {
	LogLevel    string            `json:"loglevel,omitempty"`
	AsUser      string            `json:"asUser,omitempty"`
	Options     map[string]string `json:"options,omitempty"`
	FailedNodes bool              `json:"failedNodes,omitempty"`
}

// ExportJobsInput are the optional parameters for the export endpoint
type ExportJobsInput struct {
	Format    string
	IDList    []string
	GroupPath string
	JobFilter string
}

// ImportJobsInput are the parameters for the import endpoint
type ImportJobsInput struct {
	FileFormat      string
	DuplicateOption string
	UUIDOption      string
	RawContent      []byte
}

// ImportJobsResponse is the response that comes back from importing jobs
type ImportJobsResponse struct {
	Succeeded []*Job `json:"succeeded,omitempty"`
	Failed    []*Job `json:"failed,omitempty"`
	Skipped   []*Job `json:"skipped,omitempty"`
}

// Jobs is information pertaining to jobs API endpoints
type Jobs struct {
	c *Client
}

// Jobs interacts with endpoints pertaining to jobs
func (c *Client) Jobs() *Jobs {
	return &Jobs{c: c}
}

// List returns a list of jobs
func (j *Jobs) List(project string, input *ListJobsInput) ([]*Job, error) {
	uri, err := j.urlEncodeListInput(j.c.RundeckAddr+"/project/"+project+"/jobs", input)
	if err != nil {
		return nil, err
	}

	res, err := j.c.get(uri)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, makeError(res.Body)
	}

	var jobs []*Job
	return jobs, json.NewDecoder(res.Body).Decode(&jobs)
}

// Run will execute a job
func (j *Jobs) Run(jobID string, input *RunJobInput) (*Execution, error) {
	uri := j.c.RundeckAddr + "/job/" + jobID + "/run"

	var body io.Reader
	if input != nil {
		bs, err := json.Marshal(j.convertToSerializeable(input))
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(bs)
	}

	res, err := j.c.post(uri, body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, makeError(res.Body)
	}

	var execution Execution
	return &execution, json.NewDecoder(res.Body).Decode(&execution)
}

// Retry retries a job based on an execution id
func (j *Jobs) Retry(jobID string, execID int64, input *RetryJobInput) (*Execution, error) {
	uri := j.c.RundeckAddr + "/job/" + jobID + "/retry/" + strconv.FormatInt(execID, 10)

	var body io.Reader
	if input != nil {
		bs, err := json.Marshal(input)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(bs)
	}

	res, err := j.c.post(uri, body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, makeError(res.Body)
	}

	var execution Execution
	return &execution, json.NewDecoder(res.Body).Decode(&execution)
}

// Export exports a projects jobs defintions
func (j *Jobs) Export(project string, input *ExportJobsInput) ([]byte, error) {
	rawURL := j.c.RundeckAddr + "/project/" + project + "/jobs/export"

	uri, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	query := uri.Query()

	if input != nil {
		format := JobFormatXML
		if strings.ToLower(input.Format) == JobFormatYAML {
			format = JobFormatYAML
		}
		query.Add("format", format)

		if input.GroupPath != "" {
			query.Add("groupPath", input.GroupPath)
		}

		if input.IDList != nil && len(input.IDList) > 0 {
			query.Add("idlist", strings.Join(input.IDList, ","))
		}

		if input.JobFilter != "" {
			query.Add("jobFilter", input.JobFilter)
		}
	} else {
		query.Add("format", JobFormatXML)
	}

	uri.RawQuery = query.Encode()

	res, err := j.c.get(uri.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, makeError(res.Body)
	}

	return ioutil.ReadAll(res.Body)
}

// Import imports job definitions
func (j *Jobs) Import(project string, input *ImportJobsInput) (*ImportJobsResponse, error) {
	if input == nil {
		return nil, fmt.Errorf("input cannot be nil as ImportJobsInput.RawContent is required to import anything")
	}
	rawURL := j.c.RundeckAddr + "/project/" + project + "/jobs/import"

	uri, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	query := uri.Query()

	fileFormat := JobFormatXML
	if strings.ToLower(input.FileFormat) == JobFormatYAML {
		fileFormat = JobFormatYAML
	}
	query.Add("fileformat", fileFormat)

	if input.DuplicateOption != "" {
		query.Add("dupeOption", input.DuplicateOption)
	}

	if input.UUIDOption != "" {
		query.Add("uuidOption", input.UUIDOption)
	}

	uri.RawQuery = query.Encode()

	res, err := j.c.post(uri.String(), bytes.NewReader(input.RawContent))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, makeError(res.Body)
	}

	var resp ImportJobsResponse
	return &resp, json.NewDecoder(res.Body).Decode(&resp)
}

func (j *Jobs) urlEncodeListInput(rawURL string, input *ListJobsInput) (string, error) {
	uri, err := url.Parse(rawURL)
	if err != nil {
		return rawURL, err
	}

	query := uri.Query()

	if input != nil {
		if input.GroupPath != "" {
			query.Add("groupPath", input.GroupPath)
		}

		if input.GroupPathExact != "" {
			query.Add("groupPathExact", input.GroupPathExact)
		}

		if input.IDs != nil && len(input.IDs) > 0 {
			query.Add("idlist", strings.Join(input.IDs, ","))
		}

		if input.JobExactFilter != "" {
			query.Add("jobExactFilter", input.JobExactFilter)
		}

		if input.JobFilter != "" {
			query.Add("jobFilter", input.JobFilter)
		}

		if input.ScheduledFilter {
			query.Add("scheduledFilter", strconv.FormatBool(input.ScheduledFilter))
		}

		if input.ServerNodeUUIDFilter != "" {
			query.Add("serverNodeUUIDFilter", input.ServerNodeUUIDFilter)
		}
	}

	uri.RawQuery = query.Encode()

	return uri.String(), nil
}

func (j *Jobs) convertToSerializeable(input *RunJobInput) *runJobInputSerializeable {
	var serializeable runJobInputSerializeable
	serializeable.AsUser = input.AsUser
	serializeable.LogLevel = input.LogLevel
	if input.Options != nil {
		serializeable.Options = input.Options
	}
	serializeable.RunAtTime = input.RunAtTime

	if input.Filters != nil {
		serializeable.Filter = j.c.convertFiltersToSerializeableFormat(input.Filters)
	}

	return &serializeable
}
