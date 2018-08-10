package rundeck

import (
	"encoding/json"
)

// LogStore contains information about logstorage in the system API
type LogStore struct {
	c *Client
}

// LogStore interacts with the according API
func (c *Client) LogStore() *LogStore {
	return &LogStore{c: c}
}

// LogStorage returns log storage information and stats
func (l *LogStore) LogStorage() (*LogStorageStats, error) {
	url := l.c.RundeckAddr + "/system/logstorage"

	res, err := l.c.checkResponseOK(l.c.get(url))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var logStorage LogStorageStats
	return &logStorage, json.NewDecoder(res.Body).Decode(&logStorage)
}

// IncompleteLogStorage lists executions with incomplete logstorage
func (l *LogStore) IncompleteLogStorage() (*IncompleteLogStorageResponse, error) {
	url := l.c.RundeckAddr + "/system/logstorage/incomplete"

	res, err := l.c.checkResponseOK(l.c.get(url))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var incompleteLogStorageResponse IncompleteLogStorageResponse
	return &incompleteLogStorageResponse, json.NewDecoder(res.Body).Decode(&incompleteLogStorageResponse)
}

// ResumeIncompleteLogStorage resumes processing incomplete log storage uploads
func (l *LogStore) ResumeIncompleteLogStorage() (*ResumedIncompleteLogStorageResponse, error) {
	url := l.c.RundeckAddr + "/system/logstorage/incomplete/resume"

	res, err := l.c.checkResponseOK(l.c.post(url, nil))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var resumed ResumedIncompleteLogStorageResponse
	return &resumed, json.NewDecoder(res.Body).Decode(&resumed)
}
