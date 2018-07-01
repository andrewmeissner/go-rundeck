package rundeck

import "time"

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
