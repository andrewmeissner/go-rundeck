package rundeck

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	ExecutionStatusRunning         = "running"
	ExecutionStatusSucceeded       = "succeeded"
	ExecutionStatusFailed          = "failed"
	ExecutionStatusAborted         = "aborted"
	ExecutionStatusTimedout        = "timedout"
	ExecutionStatusFailedWithRetry = "failed-with-retry"
	ExecutionStatusScheduled       = "scheduled"
	ExecutionStatusOther           = "other"
)

// Execution is information regarding an execution
type Execution struct {
	ID              int                `json:"id"`
	HREF            string             `json:"href"`
	Permalink       string             `json:"permalink"`
	Status          string             `json:"status"`
	CustomStatus    string             `json:"customStatus"`
	Project         string             `json:"project"`
	User            string             `json:"user"`
	ServerUUID      string             `json:"serverUUID"`
	DateStarted     ExecutionTimestamp `json:"date-started"`
	DateEnded       ExecutionTimestamp `json:"date-ended"`
	Job             Job                `json:"job"`
	Description     string             `json:"description"`
	ArgString       string             `json:"argstring"`
	Storage         LogStorageMetadata `json:"storage"`
	SuccessfulNodes []string           `json:"successfulNodes"`
	FailedNodes     []string           `json:"failedNodes"`
	Errors          []string           `json:"errors"`
}

// ExecutionTimestamp is basic time metadata
type ExecutionTimestamp struct {
	UnixTime int64     `json:"unixtime"`
	Date     time.Time `json:"date"`
}

// ExecutionsResponse contains paging information as well as executions
type ExecutionsResponse struct {
	PagingInfo
	Executions []*Execution `json:"executions"`
}

// DeleteExecutionsResponse contains information about the success and failures of the delete operation
type DeleteExecutionsResponse struct {
	Failures []struct {
		ID      string `json:"id"`
		Message string `json:"message"`
	} `json:"failures"`
	FailedCount   int  `json:"failedCount"`
	SuccessCount  int  `json:"successCount"`
	AllSuccessful bool `json:"allsuccessful"`
	RequestCount  int  `json:"requestCount"`
}

// Executions is information pertaining to executions API endpoints
type Executions struct {
	c *Client
}

// Executions interacts with endpoints pertaining to executions
func (c *Client) Executions() *Executions {
	return &Executions{c: c}
}

// GetExecutionsForAJob returns the executions pertaining to a certain job
func (e *Executions) GetExecutionsForAJob(id string, status *string, paging *PagingInfo) (*ExecutionsResponse, error) {
	rawURL := e.c.RundeckAddr + "/job/" + id + "/executions"

	uri, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	query := uri.Query()

	if status != nil {
		query.Add("status", stringValue(status))
	}

	if paging != nil {
		if paging.Max != 0 {
			query.Add("max", strconv.FormatInt(int64(paging.Max), 10))
		}

		if paging.Offset != 0 {
			query.Add("offset", strconv.FormatInt(int64(paging.Offset), 10))
		}
	}

	uri.RawQuery = query.Encode()

	res, err := e.c.get(uri.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, makeError(res.Body)
	}

	var executions ExecutionsResponse
	return &executions, json.NewDecoder(res.Body).Decode(&executions)
}

// DeleteExecutions deletes all executions for a job
func (e *Executions) DeleteExecutions(id string) (*DeleteExecutionsResponse, error) {
	rawURL := e.c.RundeckAddr + "/job/" + id + "/executions"

	res, err := e.c.delete(rawURL, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, makeError(res.Body)
	}

	var response DeleteExecutionsResponse
	return &response, json.NewDecoder(res.Body).Decode(&response)
}

// ListRunningExecutions returns running exeuctions for the specified project ("*" for all projects)
func (e *Executions) ListRunningExecutions(project string) (*ExecutionsResponse, error) {
	rawURL := e.c.RundeckAddr + "/project/" + project + "/executions/running"

	res, err := e.c.get(rawURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, makeError(res.Body)
	}

	var executions ExecutionsResponse
	return &executions, json.NewDecoder(res.Body).Decode(&executions)
}
