package rundeck

// PagingInfo contains information relating to a paginated response
type PagingInfo struct {
	Total  int `json:"total"`
	Max    int `json:"max"`
	Offset int `json:"offset"`
}
