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
	ToggleKindExecution   = "execution"
	ToggleKindSchedule    = "schedule"
)

// Job is information about a Rundeck job
type Job struct {
	ID                     string            `json:"id,omitempty"`
	AverageDuration        int64             `json:"averageDuration,omitempty"`
	Name                   string            `json:"name,omitempty"`
	Group                  string            `json:"group,omitempty"`
	Project                string            `json:"project,omitempty"`
	Description            string            `json:"description,omitempty"`
	HREF                   string            `json:"href,omitempty"`
	Permalink              string            `json:"permalink,omitempty"`
	Options                map[string]string `json:"options,omitempty"`
	Scheduled              bool              `json:"scheduled,omitempty"`
	ScheduleEnabled        bool              `json:"scheduleEnabled,omitempty"`
	Enabled                bool              `json:"enabled,omitempty"`
	ServerNodeUUID         string            `json:"serverNodeUUID,omitempty"`
	ServerOwner            bool              `json:"serverOwner,omitempty"`
	Index                  int               `json:"index,omitempty"`
	NextScheduledExecution time.Time         `json:"nextScheduledExecution,omitempty"`
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

// BulkModifyInput contains a string slice of ids to modify in bulk
type BulkModifyInput struct {
	IDs []string `json:"ids"`
}

// BulkModifyResponse is the response body from the bulk modify endpoint
type BulkModifyResponse struct {
	RequestCount  int                 `json:"requestCount"`
	AllSuccessful bool                `json:"allsuccessful"`
	Enabled       bool                `json:"enabled,omitempty"`
	Succeeded     []*BulkModifyObject `json:"successful,omitempty"`
	Failed        []*BulkModifyObject `json:"failed,omitempty"`
}

// BulkModifyObject is an object in the BulkModifyResponse slices
type BulkModifyObject struct {
	ID        string `json:"id"`
	ErrorCode string `json:"errorCode"`
	Message   string `json:"message"`
}

// SuccessResponse is a response containing whether or not the api call was a success
type SuccessResponse struct {
	Success bool `json:"success"`
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

// GetDefinition returns a job definition as a slice of bytces in either xml or yaml
func (j *Jobs) GetDefinition(id string, format *string) ([]byte, error) {
	rawURL := j.c.RundeckAddr + "/job/" + id

	returnFormat := JobFormatXML
	if strings.ToLower(stringValue(format)) == JobFormatYAML {
		returnFormat = JobFormatYAML
	}

	uri, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	query := uri.Query()
	query.Add("format", returnFormat)
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

// DeleteDefinition deletes a job definition
func (j *Jobs) DeleteDefinition(id string) error {
	rawURL := j.c.RundeckAddr + "/job/" + id

	res, err := j.c.delete(rawURL, nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		return makeError(res.Body)
	}

	return nil
}

// BulkDelete deletes jobs in bulk by ID
func (j *Jobs) BulkDelete(input *BulkModifyInput) (*BulkModifyResponse, error) {
	if input == nil {
		return nil, fmt.Errorf("bulk delete input cannot be nil")
	}

	rawURL := j.c.RundeckAddr + "/jobs/delete"

	bs, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	res, err := j.c.delete(rawURL, bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, makeError(res.Body)
	}

	var response BulkModifyResponse
	return &response, json.NewDecoder(res.Body).Decode(&response)
}

// ToggleExecutionsOrSchedules toggles the executions or schedules of the supplied job
func (j *Jobs) ToggleExecutionsOrSchedules(id string, enabled bool, toggleKind string) (*SuccessResponse, error) {
	if toggleKind != ToggleKindExecution && toggleKind != ToggleKindSchedule {
		return nil, fmt.Errorf(`toggleKind must be "execution" or "schedule"`)
	}

	rawURL := j.c.RundeckAddr + "/job/" + id + "/" + toggleKind

	if enabled {
		rawURL += "/enable"
	} else {
		rawURL += "/disable"
	}

	res, err := j.c.post(rawURL, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, makeError(res.Body)
	}

	var success SuccessResponse
	return &success, json.NewDecoder(res.Body).Decode(&success)
}

// BulkToggleExecutionsOrSchedules toggles the execution or scheudle value of the suppplied job ids
func (j *Jobs) BulkToggleExecutionsOrSchedules(input *BulkModifyInput, enabled bool, toggleKind string) (*BulkModifyResponse, error) {
	if input == nil {
		return nil, fmt.Errorf("input cannot be nil")
	}

	if toggleKind != ToggleKindExecution && toggleKind != ToggleKindSchedule {
		return nil, fmt.Errorf(`toggleKind must be "execution" or "schedule"`)
	}

	rawURL := j.c.RundeckAddr + "/jobs/" + toggleKind

	if enabled {
		rawURL += "/enable"
	} else {
		rawURL += "/disable"
	}

	res, err := j.c.post(rawURL, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, makeError(res.Body)
	}

	var bulk BulkModifyResponse
	return &bulk, json.NewDecoder(res.Body).Decode(&bulk)
}

// GetMetadata returns basic information about a job
func (j *Jobs) GetMetadata(id string) (*Job, error) {
	rawURL := j.c.RundeckAddr + "/job/" + id + "/info"

	res, err := j.c.get(rawURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, makeError(res.Body)
	}

	var job Job
	return &job, json.NewDecoder(res.Body).Decode(&job)
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
