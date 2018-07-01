package rundeck

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SystemInfoResponse is Rundeck server information and stats
type SystemInfoResponse struct {
	System SystemInfo `json:"system"`
}

// SystemInfo is information about the rundeck system
type SystemInfo struct {
	Timestamp   Timestamp       `json:"timestamp"`
	Rundeck     Rundeck         `json:"rundeck"`
	Executions  ExecutionMode   `json:"executions"`
	OS          OperatingSystem `json:"os"`
	JVM         JVM             `json:"jvm"`
	Stats       Stats           `json:"stats"`
	Metrics     Metrics         `json:"metrics"`
	ThreadDump  ThreadDump      `json:"threadDump"`
	HealthCheck HealthCheck     `json:"healthcheck"`
}

// Timestamp is time information on the Rundeck server
type Timestamp struct {
	Epoch    int64     `json:"epoch"`
	Unit     string    `json:"unit"`
	DateTime time.Time `json:"datetime"`
}

// Rundeck is information about the build
type Rundeck struct {
	Version    string `json:"version"`
	Build      string `json:"build"`
	Node       string `json:"node"`
	APIVersion int    `json:"apiversion"`
	ServerUUID string `json:"serverUUID"`
}

// ExecutionMode is information about the rundeck server's ability to execute jobs
type ExecutionMode struct {
	Active        bool   `json:"active"`
	ExecutionMode string `json:"executionMode"`
}

// OperatingSystem is information regarding the Rundeck host
type OperatingSystem struct {
	Architecture string `json:"arch"`
	Name         string `json:"name"`
	Version      string `json:"version"`
}

// JVM is information about the JVM of Rundeck
type JVM struct {
	Name                   string `json:"name"`
	Vendor                 string `json:"vendor"`
	Version                string `json:"version"`
	ImplementationVeresion string `json:"implementationVersion"`
}

// Stats are basic Rundeck stats
type Stats struct {
	Uptime    UptimeStats    `json:"uptime"`
	CPU       CPUStats       `json:"cpu"`
	Memory    MemoryStats    `json:"memory"`
	Scheduler SchedulerStats `json:"scheduler"`
	Threads   ThreadStats    `json:"threads"`
}

// UptimeStats are basic stats about system uptime
type UptimeStats struct {
	Duration int64     `json:"duration"`
	Unit     string    `json:"unit"`
	Since    Timestamp `json:"since"`
}

// CPUStats are basic stats about the CPU usage
type CPUStats struct {
	LoadAverage LoadAverageStats `json:"loadAverage"`
	Processors  int              `json:"processors"`
}

// MemoryStats are basic stats about memory usage
type MemoryStats struct {
	Unit  string `json:"unit"`
	Max   int64  `json:"max"`
	Free  int64  `json:"free"`
	Total int64  `json:"total"`
}

// SchedulerStats are stats about the Rundeck scheduler
type SchedulerStats struct {
	Running        int `json:"running"`
	ThreadPoolSize int `json:"threadPoolSize"`
}

// ThreadStats are stats about thread usage
type ThreadStats struct {
	Active int `json:"active"`
}

// LoadAverageStats are stats about the CPU load
type LoadAverageStats struct {
	Unit    string  `json:"unit"`
	Average float64 `json:"average"`
}

// Metrics contains a url to a page regarding metrics
type Metrics struct {
	HREF        string `json:"href"`
	ContentType string `json:"contentType"`
}

// ThreadDump contains a url to a page regarding thread dump information
type ThreadDump struct {
	HREF        string `json:"href"`
	ContentType string `json:"contentType"`
}

// HealthCheck contains a urll to a page regarding health information
type HealthCheck struct {
	HREF        string `json:"href"`
	ContentType string `json:"contentType"`
}

// LogStorageStats is the log storage info and stats
type LogStorageStats struct {
	Enabled         bool   `json:"enabled"`
	PluginName      string `json:"pluginName"`
	SucceededCount  int64  `json:"succeededCount"`
	FailedCount     int64  `json:"failedCount"`
	QueuedCount     int64  `json:"queuedCount"`
	TotalCount      int64  `json:"totalCount"`
	ImcompleteCount int64  `json:"incompleteCount"`
	MissingCount    int64  `json:"missingCount"`
}

// LogStorageMetadata is information about an executions associated log storage
type LogStorageMetadata struct {
	LocalFilesPresent   bool      `json:"localFilesPresent"`
	IncompleteFileTypes []string  `json:"incompleteFileTypes"`
	Queued              bool      `json:"queued"`
	Failed              bool      `json:"failed"`
	Date                time.Time `json:"date"`
}

// IncompleteLogStorageResponse is the response from the associated endpoint
type IncompleteLogStorageResponse struct {
	PagingInfo
	Executions []*Execution `json:"executions"`
}

// System is the information regarding system calls
type System struct {
	c *Client
}

// System interacts with the system api
func (c *Client) System() *System {
	return &System{c: c}
}

// Info retrieves Rundeck server information and stats.
func (s *System) Info() (*SystemInfoResponse, error) {
	url := fmt.Sprintf("%s/system/info", s.c.RundeckAddr)

	res, err := s.c.get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, makeError(res.Body)
	}

	var systemInfo SystemInfoResponse
	return &systemInfo, json.NewDecoder(res.Body).Decode(&systemInfo)
}

// LogStorage returns log storage information and stats
func (s *System) LogStorage() (*LogStorageStats, error) {
	url := fmt.Sprintf("%s/system/logstorage", s.c.RundeckAddr)

	res, err := s.c.get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, makeError(res.Body)
	}

	var logStorage LogStorageStats
	return &logStorage, json.NewDecoder(res.Body).Decode(&logStorage)
}

// IncompleteLogStorage lists executions with incomplete logstorage
func (s *System) IncompleteLogStorage() (*IncompleteLogStorageResponse, error) {
	url := fmt.Sprintf("%s/system/logstorage/incomplete", s.c.RundeckAddr)

	res, err := s.c.get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, makeError(res.Body)
	}

	var incompleteLogStorageResponse IncompleteLogStorageResponse
	return &incompleteLogStorageResponse, json.NewDecoder(res.Body).Decode(&incompleteLogStorageResponse)
}
