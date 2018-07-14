package rundeck

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	JobLogLevelDebug   = "DEBUG"
	JobLogLevelVerbose = "VERBOSE"
	JobLogLevelInfo    = "INFO"
	JobLogLevelWarn    = "WARN"
	JobLogLevelError   = "ERROR"
)

// Job is information about a Rundeck job
type Job struct {
	ID              string            `json:"id"`
	AverageDuration int64             `json:"averageDuration"`
	Name            string            `json:"name"`
	Group           string            `json:"group"`
	Project         string            `json:"project"`
	Description     string            `json:"description"`
	HREF            string            `json:"href"`
	Permalink       string            `json:"permalink"`
	Options         map[string]string `json:"options"`
	Scheduled       bool              `json:"scheduled"`
	ScheduleEnabled bool              `json:"scheduleEnabled"`
	Enabled         bool              `json:"enabled"`
	ServerNodeUUID  string            `json:"serverNodeUUID"`
	ServerOwner     bool              `json:"serverOwner"`
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
	LogLevel string `json:"loglevel,omiempty"`
	AsUser   string `json:"asUser,omitempty"`
	// TODO: NodeFilters
	// TODO: Filter
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

	// TODO: dbl check this is an execution
	var execution Execution
	return &execution, json.NewDecoder(res.Body).Decode(&execution)
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

	return query.Encode(), nil
}
