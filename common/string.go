package common

import "strings"

func Trim(s string) string {
	return strings.TrimSpace(s)
}

func IsEmptyString(s string) bool {
	return len(Trim(s)) == 0
}

func IsDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func IsAlpha(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z'
}

func IsAlphaDigit(c byte) bool {
	return IsDigit(c) || IsAlpha(c)
}
