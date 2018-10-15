package cmd

func humanizeBool(value bool) string {
	if value {
		return "Yes"
	}
	return "No"
}
