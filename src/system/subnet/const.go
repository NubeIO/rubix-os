package checksubnet

import "regexp"

var (
	// IPV4Pattern - ip v4 regular expression
	IPV4Pattern, _ = regexp.Compile("\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}")
)
