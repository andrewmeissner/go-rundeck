package rundeck

import (
	"bytes"
	"encoding/json"
	"net/url"
	"strconv"
	"time"
)

// ExecutionStatus ensure a constant is used in parameters
type ExecutionStatus string

const (
	ExecutionStatusRunning         ExecutionStatus = "running"
	ExecutionStatusSucceeded       ExecutionStatus = "succeeded"
	ExecutionStatusFailed          ExecutionStatus = "failed"
	ExecutionStatusAborted         ExecutionStatus = "aborted"
	ExecutionStatusTimedout        ExecutionStatus = "timedout"
	ExecutionStatusFailedWithRetry ExecutionStatus = "failed-with-retry"
	ExecutionStatusScheduled       ExecutionStatus = "scheduled"
	ExecutionStatusOther           ExecutionStatus = "other"
)

// ExecutionType ensure a constant is used in parameters
type ExecutionType string

const (
	ExecutionTypeScheduled     ExecutionType = "scheduled"
	ExecutionTypeUser          ExecutionType = "user"
	ExecutionTypeUserScheduled ExecutionType = "user-scheduled"
)

// Boolean often times have more than 2 values
type Boolean int

const (
	BooleanDefault Boolean = iota
	BooleanFalse
	BooleanTrue
)

// ExecutionState is the state of the execution
type ExecutionState string

const (
	ExecutionStateWaiting              ExecutionState = "WAITING"
	ExecutionStateRunning              ExecutionState = "RUNNING"
	ExecutionStateRunningHandler       ExecutionState = "RUNNING_HANDLER"
	ExecutionStateSucceeded            ExecutionState = "SUCCEEDED"
	ExecutionStateFailed               ExecutionState = "FAILED"
	ExecutionStateAborted              ExecutionState = "ABORTED"
	ExecutionStateNodePartialSucceeded ExecutionState = "NODE_PARTIAL_SUCCEEDED"
	ExecutionStateNodeMixed            ExecutionState = "NODE_MIXED"
	ExecutionStateNotStarted           ExecutionState = "NOT_STARTED"
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

// ExecutionQueryInput are parameters to narrow down the result set of a query operation
type ExecutionQueryInput struct {
	PagingInfo
	Status                ExecutionStatus
	AbortedBy             string
	User                  string
	RecentFilter          string
	OlderFilter           string
	Begin                 *time.Time // unix ms
	End                   *time.Time // unix ms
	AdHoc                 Boolean
	JobIDList             []string
	ExcludeJobIDList      []string
	JobList               []string
	ExcludeJobList        []string
	GroupPath             string
	GroupPathExact        string
	ExcludeGroupPath      string
	ExcludeGroupPathExact string
	JobName               string
	ExcludeJobName        string
	JobNameExact          string
	ExcludeJobNameExact   string
	ExecutionType         ExecutionType
}

// ExecutionStateInfo ...
type ExecutionStateInfo struct {
	StartTime      time.Time      `json:"startTime"`
	EndTime        time.Time      `json:"endTime"`
	UpdateTime     time.Time      `json:"updateTime"`
	ExecutionState ExecutionState `json:"executionState"`
}

// ExecutionWorkflow ...
type ExecutionWorkflow struct {
	StepCount   int                  `json:"stepCount"`
	TargetNodes []string             `json:"targetNodes"`
	Steps       []ExecutionStepState `json:"steps"`
}

// ExecutionStateIndicator ...
type ExecutionStateIndicator struct {
	ExecutionState        ExecutionState `json:"executionState"`
	StepContextIdentifier string         `json:"stepctx"`
}

// ExecutionStepState ...
type ExecutionStepState struct {
	ExecutionStateInfo
	ID                    string                        `json:"id"`
	StepContextIdentifier string                        `json:"stepctx"`
	NodeStep              bool                          `json:"nodeStep"`
	NodeStates            map[string]ExecutionStateInfo `json:"nodeStates"`
	HasSubworkflow        bool                          `json:"hasSubworkflow"`
	Workflow              ExecutionWorkflow             `json:"workflow"`
}

// ExecutionStateResponse ...
type ExecutionStateResponse struct {
	ExecutionStateInfo
	ExecutionWorkflow
	AllNodes    []string                           `json:"allNodes"`
	Nodes       map[string]ExecutionStateIndicator `json:"nodes"`
	ServerNode  string                             `json:"serverNode"`
	ExecutionID int                                `json:"executionId"`
	Completed   bool                               `json:"completed"`
}

// ExecutionsOutputInput ...
type ExecutionsOutputInput struct {
	Node        string
	StepContext string
	Offset      int
	LastLines   int
	LastMod     *time.Time
	Compacted   bool
}

// ExecutionsOutputResponse ...
type ExecutionsOutputResponse struct {
	ID             int            `json:"id"`
	Message        string         `json:"message"`
	Error          string         `json:"error"`
	Unmodified     bool           `json:"unmodified"`
	Empty          bool           `json:"empty"`
	Offset         int            `json:"offset"`
	Completed      bool           `json:"completed"`
	ExecCompleted  bool           `json:"execCompleted"`
	HasFailedNodes bool           `json:"hasFailedNodes"`
	ExecutionState ExecutionState `json:"execState"`
	LastModified   int64          `json:"lastModified"`
	ExecDuration   int64          `json:"execDuration"`
	PercentLoaded  float32        `json:"percentLoaded"`
	TotalSize      int            `json:"totalSize"`
	// Filter?  NodeName StepContext?
	Compacted     bool        `json:"compacted"`
	CompactedAttr string      `json:"compactedAttr"`
	Entries       []*LogEntry `json:"entries"`
}

// LogEntry ...
type LogEntry struct {
	Time         string     `json:"time"`
	AbsoluteTime *time.Time `json:"absolute_time"`
	Level        LogLevel   `json:"level"`
	Log          string     `json:"log"`
	User         string     `json:"user"`
	Command      string     `json:"command"`
	Node         string     `json:"node"`
	StepContext  string     `json:"stepctx"`
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

	res, err := e.c.checkResponseOK(e.c.get(uri.String()))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var executions ExecutionsResponse
	return &executions, json.NewDecoder(res.Body).Decode(&executions)
}

// DeleteExecutions deletes all executions for a job
func (e *Executions) DeleteExecutions(id string) (*DeleteExecutionsResponse, error) {
	rawURL := e.c.RundeckAddr + "/job/" + id + "/executions"

	res, err := e.c.checkResponseOK(e.c.delete(rawURL, nil))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var response DeleteExecutionsResponse
	return &response, json.NewDecoder(res.Body).Decode(&response)
}

// ListRunningExecutions returns running exeuctions for the specified project ("*" for all projects)
func (e *Executions) ListRunningExecutions(project string) (*ExecutionsResponse, error) {
	rawURL := e.c.RundeckAddr + "/project/" + project + "/executions/running"

	res, err := e.c.checkResponseOK(e.c.get(rawURL))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var executions ExecutionsResponse
	return &executions, json.NewDecoder(res.Body).Decode(&executions)
}

// Info returns information about the specific execution
func (e *Executions) Info(id int) (*Execution, error) {
	rawURL := e.c.RundeckAddr + "/exeuction/" + strconv.FormatInt(int64(id), 10)

	res, err := e.c.checkResponseOK(e.c.get(rawURL))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var execution Execution
	return &execution, json.NewDecoder(res.Body).Decode(&execution)
}

// ListInputFiles lists input ifle sused for an execution
func (e *Executions) ListInputFiles(id int) (*UploadedFilesResponse, error) {
	rawURL := e.c.RundeckAddr + "/execution/" + strconv.FormatInt(int64(id), 10) + "/input/files"

	res, err := e.c.checkResponseOK(e.c.get(rawURL))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var files UploadedFilesResponse
	return &files, json.NewDecoder(res.Body).Decode(&files)
}

// Delete deletes an execution by id
func (e *Executions) Delete(id int) error {
	rawURL := e.c.RundeckAddr + "/execution/" + strconv.FormatInt(int64(id), 10)

	res, err := e.c.checkResponseNoContent(e.c.delete(rawURL, nil))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

// BulkDelete deletes a set of executions by their ids
func (e *Executions) BulkDelete(ids []int) (*DeleteExecutionsResponse, error) {
	rawURL := e.c.RundeckAddr + "/executions/delete"

	bs, err := json.Marshal(ids)
	if err != nil {
		return nil, err
	}

	res, err := e.c.checkResponseOK(e.c.delete(rawURL, bytes.NewReader(bs)))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var bulkResponse DeleteExecutionsResponse
	return &bulkResponse, json.NewDecoder(res.Body).Decode(&bulkResponse)
}

// Query queries for executions based on job or execution details
func (e *Executions) Query(project string, input *ExecutionQueryInput) (*ExecutionsResponse, error) {
	rawURL := e.c.RundeckAddr + "/project/" + project + "/executions"

	uri, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	query := uri.Query()

	if input != nil {
		if input.AbortedBy != "" {
			query.Add("abortedByFilter", input.AbortedBy)
		}

		if input.AdHoc == BooleanFalse {
			query.Add("adhoc", "false")
		} else if input.AdHoc == BooleanTrue {
			query.Add("adhoc", "true")
		}

		if input.Begin != nil {
			query.Add("begin", strconv.FormatInt(input.Begin.Unix(), 10))
		}

		if input.End != nil {
			query.Add("end", strconv.FormatInt(input.End.Unix(), 10))
		}

		if input.ExcludeGroupPath != "" {
			query.Add("excludeGroupPath", input.ExcludeGroupPath)
		}

		if input.ExcludeGroupPathExact != "" {
			query.Add("excludeGroupPathExact", input.ExcludeGroupPathExact)
		}

		if input.ExcludeJobList != nil && len(input.ExcludeJobList) > 0 {
			for i := range input.ExcludeJobList {
				query.Add("excludeJobListFilter", input.ExcludeJobList[i])
			}
		}

		if input.ExcludeJobIDList != nil && len(input.ExcludeJobIDList) > 0 {
			for i := range input.ExcludeJobIDList {
				query.Add("excludeJobIdListFilter", input.ExcludeJobIDList[i])
			}
		}

		if input.ExcludeJobList != nil && len(input.ExcludeJobList) > 0 {
			for i := range input.ExcludeJobList {
				query.Add("excludeJobListFilter", input.ExcludeJobList[i])
			}
		}

		if input.ExcludeJobName != "" {
			query.Add("excludeJobFilter", input.ExcludeJobName)
		}

		if input.ExcludeJobNameExact != "" {
			query.Add("excludeJobExactFilter", input.ExcludeJobNameExact)
		}

		if input.ExecutionType != "" {
			query.Add("executionTypeFilter", string(input.ExecutionType))
		}

		if input.GroupPath != "" {
			query.Add("groupPath", input.GroupPath)
		}

		if input.GroupPathExact != "" {
			query.Add("groupPathExact", input.GroupPathExact)
		}

		if input.JobIDList != nil && len(input.JobIDList) > 0 {
			for i := range input.JobIDList {
				query.Add("jobIdListFilter", input.JobIDList[i])
			}
		}

		if input.JobList != nil && len(input.JobList) > 0 {
			for i := range input.JobList {
				query.Add("jobListFilter", input.JobList[i])
			}
		}

		if input.JobName != "" {
			query.Add("jobFilter", input.JobName)
		}

		if input.JobNameExact != "" {
			query.Add("jobExactFilter", input.JobNameExact)
		}

		if input.Max != 0 {
			query.Add("max", strconv.FormatInt(int64(input.Max), 10))
		}

		if input.Offset != 0 {
			query.Add("offset", strconv.FormatInt(int64(input.Offset), 10))
		}

		if input.OlderFilter != "" {
			query.Add("olderFilter", input.OlderFilter)
		}

		if input.RecentFilter != "" {
			query.Add("recentFilter", input.RecentFilter)
		}

		if input.Status != "" {
			query.Add("statusFilter", string(input.Status))
		}

		if input.User != "" {
			query.Add("userFilter", input.User)
		}
	}

	uri.RawQuery = query.Encode()

	res, err := e.c.checkResponseOK(e.c.get(uri.String()))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var executions ExecutionsResponse
	return &executions, json.NewDecoder(res.Body).Decode(&executions)
}

// State gets detailed about the node and step state of an execution by ID. The execution can be currently running or completed.
func (e *Executions) State(id string) (*ExecutionStateResponse, error) {
	rawURL := e.c.RundeckAddr + "/execution/" + id + "/state"

	res, err := e.c.checkResponseOK(e.c.get(rawURL))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var esr ExecutionStateResponse
	return &esr, json.NewDecoder(res.Body).Decode(&esr)
}

// Output gets the output for an execution by ID.
// The execution can be currently running or may have already completed.
// Output can be filtered down to a specific node or workflow step.
func (e *Executions) Output(id string, input *ExecutionsOutputInput) (*ExecutionsOutputResponse, error) {
	rawURL := e.c.RundeckAddr + "/execution/" + id + "/output"

	uri, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	query := uri.Query()

	if input != nil {
		if input.Node != "" {
			uri.Path += "/node/" + input.Node
		}

		if input.StepContext != "" {
			uri.Path += "/step/" + input.StepContext
		}

		if input.Compacted {
			query.Add("compacted", "true")
		}

		if input.LastLines != 0 {
			query.Add("lastlines", strconv.FormatInt(int64(input.LastLines), 10))
		}

		if input.LastMod != nil {
			query.Add("lastmod", strconv.FormatInt(int64(input.LastMod.Unix()), 10))
		}

		if input.Offset != 0 {
			query.Add("offset", strconv.FormatInt(int64(input.Offset), 10))
		}
	}

	uri.RawQuery = query.Encode()

	res, err := e.c.checkResponseOK(e.c.get(uri.String()))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var output ExecutionsOutputResponse
	return &output, json.NewDecoder(res.Body).Decode(&output)
}
