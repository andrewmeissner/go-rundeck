package rundeck

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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
