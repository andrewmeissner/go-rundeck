package rundeck

import (
	"bytes"
	"encoding/json"
	"errors"
)

// AdhocCommandStringInput ...
type AdhocCommandStringInput struct {
	Exec string `json:"exec"`
	AdhocOptions
}

// AdhocScriptInput ...
type AdhocScriptInput struct {
	Script string `json:"script"`
	AdhocOptions
	AdhocScriptOptions
}

// AdhocURLInput ...
type AdhocURLInput struct {
	URL string `json:"url"`
	AdhocOptions
	AdhocScriptOptions
}

type AdhocScriptOptions struct {
	ArgString             string `json:"argString,omitempty"`
	ScriptInterpreter     string `json:"scriptInterpreter,omitempty"`
	InterpreterArgsQuoted bool   `json:"interpreterArgsQuoted,omitempty"`
	FileExtension         string `json:"fileExtension,omitempty"`
}

type AdhocOptions struct {
	Project         string `json:"project"`
	NodeThreadcount int    `json:"nodeThreadcount,omitempty"`
	NodeKeepGoing   bool   `json:"nodeKeepgoing,omitempty"`
	AsUser          string `json:"asUser,omitempty"`
	Filter          string `json:"filter,omitempty"`
}

// AdhocCommandResponse ...
type AdhocCommandResponse struct {
	Message   string    `json:"message"`
	Execution Execution `json:"execution"`
}

// AdhocAPI interacts with the adhoc endpoints
type AdhocAPI struct {
	c *Client
}

// Adhoc restuns and adhoc api client
func (c *Client) Adhoc() *AdhocAPI {
	return &AdhocAPI{c: c}
}

// RunCommandString runs an adhoc command
func (a *AdhocAPI) RunCommandString(input *AdhocCommandStringInput) (*AdhocCommandResponse, error) {
	if input == nil {
		return nil, errors.New("input cannot be nil")
	}

	if input.Project == "" {
		return nil, errors.New("input.Project cannot be empty")
	}

	if input.Exec == "" {
		return nil, errors.New("input.Exec cannot be empty")
	}

	rawURL := a.c.RundeckAddr + "/project/" + input.Project + "/run/command"

	bs, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	res, err := a.c.checkResponseOK(a.c.post(rawURL, bytes.NewReader(bs)))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var output AdhocCommandResponse
	return &output, json.NewDecoder(res.Body).Decode(&output)
}

// RunScript runs a script
func (a *AdhocAPI) RunScript(input *AdhocScriptInput) (*AdhocCommandResponse, error) {
	if input == nil {
		return nil, errors.New("input cannot be nil")
	}

	if input.Project == "" {
		return nil, errors.New("input.Project cannot be empty")
	}

	if input.Script == "" {
		return nil, errors.New("input.Script cannot be empty")
	}

	rawURL := a.c.RundeckAddr + "/project/" + input.Project + "/run/script"

	bs, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	res, err := a.c.checkResponseOK(a.c.post(rawURL, bytes.NewReader(bs)))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var output AdhocCommandResponse
	return &output, json.NewDecoder(res.Body).Decode(&output)
}

// RunURL runs a script downloaded from a url
func (a *AdhocAPI) RunURL(input *AdhocURLInput) (*AdhocCommandResponse, error) {
	if input == nil {
		return nil, errors.New("input cannot be nil")
	}

	if input.Project == "" {
		return nil, errors.New("input.Project cannot be empty")
	}

	if input.URL == "" {
		return nil, errors.New("input.URL cannot be empty")
	}

	rawURL := a.c.RundeckAddr + "/project/" + input.Project + "/run/url"

	bs, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	res, err := a.c.checkResponseOK(a.c.post(rawURL, bytes.NewReader(bs)))
	if err != nil {
		return nil, err
	}

	var output AdhocCommandResponse
	return &output, json.NewDecoder(res.Body).Decode(&output)
}
