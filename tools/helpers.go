package tools

// map2options takes a map[string]string and returns a slice of maps with "label" and "value" keys.
func Map2options(m map[string]string) []map[string]string {
	var options []map[string]string

	for k, v := range m {
		option := map[string]string{"label": v, "value": k}
		options = append(options, option)
	}

	return options
}
