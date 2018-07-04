package rundeck

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// TakeoverScheduleInput is the payload to takeover jobs in cluster mode
type TakeoverScheduleInput struct {
	Server  *TakeoverServer   `json:"server,omitempty"`
	Project *string           `json:"project,omitempty"`
	Job     *TakeoverJobInput `json:"job,omitempty"`
}

// TakeoverServer ...
type TakeoverServer struct {
	UUID string `json:"uuid,omitempty"`
	All  bool   `json:"all,omitempty"`
}

// TakeoverJobInput ...
type TakeoverJobInput struct {
	ID string `json:"id,omitempty"`
}

// TakeoverScheduleResponse is the response body from the endpoint
type TakeoverScheduleResponse struct {
	TakeoverSchedule TakeoverSchedule `json:"takeoverSchedule"`
	Self             struct {
		Server TakeoverServer `json:"server,omitempty"`
	} `json:"self,omitempty"`
	Message    string `json:"message"`
	APIVersion int    `json:"apiversion"`
	Success    bool   `json:"success"`
}

// TakeoverSchedule is the result from issuing a takeover call
type TakeoverSchedule struct {
	Jobs    TakeoverJobs   `json:"jobs,omitempty"`
	Server  TakeoverServer `json:"server,omitempty"`
	Project string         `json:"project,omitempty"`
}

// TakeoverJobs ...
type TakeoverJobs struct {
	Failed     []TakeoverJob `json:"failed"`
	Successful []TakeoverJob `json:"successful"`
	Total      int           `json:"total"`
}

// TakeoverJob ...
type TakeoverJob struct {
	HREF          string `json:"href"`
	Permalink     string `json:"permalink"`
	ID            string `json:"id"`
	PreviousOwner string `json:"previous-owner"`
}

// ClusterScheduler contains information relating to cluster mode in Rundeck
type ClusterScheduler struct {
	c *Client
}

// ClusterScheduler interacts with the cluster endpoints on Rundeck
func (c *Client) ClusterScheduler() *ClusterScheduler {
	return &ClusterScheduler{c: c}
}

// TakeoverSchedule tells the Rundeck server in cluster mode to claim scheduled jobs from another cluster server
func (cs *ClusterScheduler) TakeoverSchedule(input *TakeoverScheduleInput) (*TakeoverScheduleResponse, error) {
	if input == nil {
		return nil, fmt.Errorf("input cannot be nil")
	}
	url := cs.c.RundeckAddr + "/scheduler/takeover"

	bs, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	res, err := cs.c.put(url, bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, makeError(res.Body)
	}

	var takeover TakeoverScheduleResponse
	return &takeover, json.NewDecoder(res.Body).Decode(&takeover)
}

// ListScheduledJobs lists scheduled jobs with the schedule owned by the server with the specified uuid.
// If uuid is nil, then the client server will be used.
func (cs *ClusterScheduler) ListScheduledJobs(uuid *string) ([]*Job, error) {
	url := cs.c.RundeckAddr + "/scheduler"

	if uuid != nil {
		url += "/server/" + stringValue(uuid)
	}

	url += "/jobs"

	res, err := cs.c.get(url)
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
