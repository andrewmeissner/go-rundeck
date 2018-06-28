package rundeck

import (
	"encoding/json"
	"io"
)

// Error is what Rundeck returns given a bad API call
type Error struct {
	ErrorPresent bool   `json:"error"`
	APIVersion   int    `json:"apiversion"`
	ErrorCode    string `json:"errorCode"`
	Message      string `json:"message"`
}

// Error marshals the struct into a json string
func (e Error) Error() string {
	bs, err := json.Marshal(e)
	if err != nil {
		return err.Error()
	}
	return string(bs)
}

func makeError(body io.Reader) Error {
	var err Error
	goErr := json.NewDecoder(body).Decode(&err)
	if goErr != nil {
		return Error{
			ErrorPresent: true,
			APIVersion:   -1,
			ErrorCode:    "decode.failure",
			Message:      goErr.Error(),
		}
	}
	return err
}
