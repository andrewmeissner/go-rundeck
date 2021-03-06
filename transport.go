package rundeck

import "net/http"

type rundeckTransport struct {
	apiToken            string
	underlyingTransport http.RoundTripper
}

func (t *rundeckTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Header.Get("Accept") == "" {
		req.Header.Add("Accept", "application/json")
	}
	if req.Header.Get("Content-Type") == "" {
		req.Header.Add("Content-Type", "application/json")
	}
	req.Header.Add("X-Rundeck-Auth-Token", t.apiToken)
	return t.underlyingTransport.RoundTrip(req)
}
