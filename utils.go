package rundeck

import (
	"strconv"
)

func stringValue(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}

func stringReference(v string) *string {
	return &v
}

func intSliceToStringSlice(si []int) []string {
	var ss []string
	for _, i := range si {
		ss = append(ss, strconv.Itoa(i))
	}
	return ss
}
