package utilstime

import "strings"

func clean(s string) string {
	if idx := strings.Index(s, ":"); idx != -1 {
		i := strings.Trim(s[idx:], ":")
		i = strings.Join(strings.Fields(strings.TrimSpace(i)), " ")
		return i
	}
	return s
}
