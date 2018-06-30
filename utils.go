package rundeck

func stringValue(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}

func stringReference(v string) *string {
	return &v
}
