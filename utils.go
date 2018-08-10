package rundeck

import "net/http"

func stringValue(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}

func stringReference(v string) *string {
	return &v
}

func checkResponseOK(res *http.Response, err error) (*http.Response, error) {
	return checkResponse(res, http.StatusOK, err)
}

func checkResponse(res *http.Response, statusCode int, err error) (*http.Response, error) {
	if err != nil {
		return nil, err
	}
	if res.StatusCode != statusCode {
		return nil, makeError(res.Body)
	}
	return res, nil
}
