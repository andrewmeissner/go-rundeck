package rundeck

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// LogLevel pertains to job log levels
type LogLevel string

const (
	JobLogLevelDebug   LogLevel = "DEBUG"
	JobLogLevelVerbose LogLevel = "VERBOSE"
	JobLogLevelInfo    LogLevel = "INFO"
	JobLogLevelWarn    LogLevel = "WARN"
	JobLogLevelError   LogLevel = "ERROR"
)

// JobFormat specifies the content type
type JobFormat string

const (
	JobFormatXML  JobFormat = "xml"
	JobFormatYAML JobFormat = "yaml"
)

// DuplicateOption instructs the job importer how to handle duplicate jobs
type DuplicateOption string

const (
	DuplicateOptionSkip   DuplicateOption = "skip"
	DuplicateOptionCreate DuplicateOption = "create"
	DuplicateOptionUpdate DuplicateOption = "update"
)

// UUIDOption instructs the job importer how to handle duplicate job uuids
type UUIDOption string

const (
	UUIDOptionPreserve UUIDOption = "preserve"
	UUIDOptionRemove   UUIDOption = "remove"
)

// ToggleKind informs if executions or scheudles are being targeted
type ToggleKind string

const (
	ToggleKindExecution ToggleKind = "execution"
	ToggleKindSchedule  ToggleKind = "schedule"
)

// FileState informs the state of the uploaded file
type FileState string

const (
	FileStateTemp     FileState = "temp"
	FileStateDeleted  FileState = "deleted"
	FileStateExpired  FileState = "expired"
	FileStateRetained FileState = "retained"
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
	LogLevel  LogLevel
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
	LogLevel    LogLevel          `json:"loglevel,omitempty"`
	AsUser      string            `json:"asUser,omitempty"`
	Options     map[string]string `json:"options,omitempty"`
	FailedNodes bool              `json:"failedNodes,omitempty"`
}

// ExportJobsInput are the optional parameters for the export endpoint
type ExportJobsInput struct {
	Format    JobFormat
	IDList    []string
	GroupPath string
	JobFilter string
}

// ImportJobsInput are the parameters for the import endpoint
type ImportJobsInput struct {
	FileFormat      JobFormat
	DuplicateOption DuplicateOption
	UUIDOption      UUIDOption
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

// UploadFileResponse is the response from uploading a file for a job.  It will contain the fileKey.
type UploadFileResponse struct {
	Total   int               `json:"total"`
	Options map[string]string `json:"options"`
}

// UploadedFilesResponse returns the files uploaded for a particular job
type UploadedFilesResponse struct {
	PagingInfo
	File []*FileOption `json:"files"`
}

// FileOption is the metadata about a file that was uploaded for a job option
type FileOption struct {
	ID             string    `json:"id"`
	User           string    `json:"user"`
	FileState      string    `json:"fileState"`
	SHA            string    `json:"sha"`
	JobID          string    `json:"jobId"`
	DateCreated    time.Time `json:"dateCreated"`
	ServerNodeUUID string    `json:"serverNodeUUID"`
	FileName       *string   `json:"fileName"`
	Size           int64     `json:"size"`
	ExpirationDate time.Time `json:"expirationDate"`
	ExecID         *int64    `json:"execId"`
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

	res, err := j.c.checkResponseOK(j.c.get(uri))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

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

	res, err := j.c.checkResponseOK(j.c.post(uri, body))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

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

	res, err := j.c.checkResponseOK(j.c.post(uri, body))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

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
		if input.Format == JobFormatYAML {
			format = JobFormatYAML
		}
		query.Add("format", string(format))

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
		query.Add("format", string(JobFormatXML))
	}

	uri.RawQuery = query.Encode()

	res, err := j.c.checkResponseOK(j.c.get(uri.String()))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

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
	if input.FileFormat == JobFormatYAML {
		fileFormat = JobFormatYAML
	}
	query.Add("fileformat", string(fileFormat))

	if input.DuplicateOption != "" {
		query.Add("dupeOption", string(input.DuplicateOption))
	}

	if input.UUIDOption != "" {
		query.Add("uuidOption", string(input.UUIDOption))
	}

	uri.RawQuery = query.Encode()

	res, err := j.c.checkResponseOK(j.c.post(uri.String(), bytes.NewReader(input.RawContent)))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var resp ImportJobsResponse
	return &resp, json.NewDecoder(res.Body).Decode(&resp)
}

// GetDefinition returns a job definition as a slice of bytces in either xml or yaml
func (j *Jobs) GetDefinition(id string, format *JobFormat) ([]byte, error) {
	rawURL := j.c.RundeckAddr + "/job/" + id

	returnFormat := JobFormatXML
	if format != nil && *format == JobFormatYAML {
		returnFormat = JobFormatYAML
	}

	uri, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	query := uri.Query()
	query.Add("format", string(returnFormat))
	uri.RawQuery = query.Encode()

	res, err := j.c.checkResponseOK(j.c.get(uri.String()))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}

// DeleteDefinition deletes a job definition
func (j *Jobs) DeleteDefinition(id string) error {
	rawURL := j.c.RundeckAddr + "/job/" + id

	res, err := j.c.checkResponseNoContent(j.c.delete(rawURL, nil))
	if err != nil {
		return err
	}
	defer res.Body.Close()

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

	res, err := j.c.checkResponseOK(j.c.delete(rawURL, bytes.NewReader(bs)))
	if err != nil {
		return nil, err
	}

	var response BulkModifyResponse
	return &response, json.NewDecoder(res.Body).Decode(&response)
}

// ToggleExecutionsOrSchedules toggles the executions or schedules of the supplied job
func (j *Jobs) ToggleExecutionsOrSchedules(id string, enabled bool, toggleKind ToggleKind) (*SuccessResponse, error) {
	if toggleKind != ToggleKindExecution && toggleKind != ToggleKindSchedule {
		return nil, errors.New(`toggleKind must be "execution" or "schedule"`)
	}

	rawURL := j.c.RundeckAddr + "/job/" + id + "/" + string(toggleKind)

	if enabled {
		rawURL += "/enable"
	} else {
		rawURL += "/disable"
	}

	res, err := j.c.checkResponseOK(j.c.post(rawURL, nil))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var success SuccessResponse
	return &success, json.NewDecoder(res.Body).Decode(&success)
}

// BulkToggleExecutionsOrSchedules toggles the execution or scheudle value of the suppplied job ids
func (j *Jobs) BulkToggleExecutionsOrSchedules(input *BulkModifyInput, enabled bool, toggleKind ToggleKind) (*BulkModifyResponse, error) {
	if input == nil {
		return nil, fmt.Errorf("input cannot be nil")
	}

	if toggleKind != ToggleKindExecution && toggleKind != ToggleKindSchedule {
		return nil, errors.New(`toggleKind must be "execution" or "schedule"`)
	}

	rawURL := j.c.RundeckAddr + "/jobs/" + string(toggleKind)

	if enabled {
		rawURL += "/enable"
	} else {
		rawURL += "/disable"
	}

	bs, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	res, err := j.c.checkResponseOK(j.c.post(rawURL, bytes.NewReader(bs)))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var bulk BulkModifyResponse
	return &bulk, json.NewDecoder(res.Body).Decode(&bulk)
}

// GetMetadata returns basic information about a job
func (j *Jobs) GetMetadata(id string) (*Job, error) {
	rawURL := j.c.RundeckAddr + "/job/" + id + "/info"

	res, err := j.c.checkResponseOK(j.c.get(rawURL))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var job Job
	return &job, json.NewDecoder(res.Body).Decode(&job)
}

// UploadFileForJobOption uploads a file to rundeck for a job option and returns the file key
func (j *Jobs) UploadFileForJobOption(id, optionName string, content []byte, fileName *string) (*UploadFileResponse, error) {
	rawURL := j.c.RundeckAddr + "/job/" + id + "/input/file"

	uri, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	query := uri.Query()
	query.Add("optionName", optionName)
	if fileName != nil {
		query.Add("fileName", stringValue(fileName))
	}
	uri.RawQuery = query.Encode()

	headers := map[string]string{
		"Content-Type": "application/octet-stream",
	}

	res, err := j.c.checkResponseOK(j.c.postWithAdditionalHeaders(uri.String(), headers, bytes.NewReader(content)))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var uploadResponse UploadFileResponse
	return &uploadResponse, json.NewDecoder(res.Body).Decode(&uploadResponse)
}

// ListFilesUploadedForJob returns files that were uploaded for a particular job
func (j *Jobs) ListFilesUploadedForJob(id string, fileState *FileState, max *int) (*UploadedFilesResponse, error) {
	rawURL := j.c.RundeckAddr + "/job/" + id + "/input/files"

	uri, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	query := uri.Query()
	if fileState != nil {
		query.Add("fileState", string(*fileState))
	}
	if max != nil {
		query.Add("max", strconv.Itoa(math.MaxInt8))
	}
	uri.RawQuery = query.Encode()

	res, err := j.c.checkResponseOK(j.c.get(uri.String()))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var files UploadedFilesResponse
	return &files, json.NewDecoder(res.Body).Decode(&files)
}

// FileInfo returns information about an uploaded file
func (j *Jobs) FileInfo(id string) (*FileOption, error) {
	rawURL := j.c.RundeckAddr + "/jobs/file/" + id

	res, err := j.c.checkResponseOK(j.c.get(rawURL))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var info FileOption
	return &info, json.NewDecoder(res.Body).Decode(&info)
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
	serializeable.LogLevel = string(input.LogLevel)
	if input.Options != nil {
		serializeable.Options = input.Options
	}
	serializeable.RunAtTime = input.RunAtTime

	if input.Filters != nil {
		serializeable.Filter = j.c.convertFiltersToSerializeableFormat(input.Filters)
	}

	return &serializeable
}
