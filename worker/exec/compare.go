package executor

import "strings"

func Normalize(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
