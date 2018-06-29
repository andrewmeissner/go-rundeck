package rundeck

func stringValue(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}
