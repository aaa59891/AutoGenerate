package utils

import "strings"

func ToLowerFirst(s string) string{
	return strings.ToLower(string(s[0])) + string(s[1:])
}
