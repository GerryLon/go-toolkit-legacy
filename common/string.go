package common

import (
	"math/rand"
	"strings"
	"time"
)

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

func RandomString(length int) string {
	if length < 1 {
		return ""
	}
	source := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_="
	bytes := []byte(source)
	result := make([]byte, 0)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := len(bytes)
	for i := length - 1; i >= 0; i-- {
		result = append(result, bytes[r.Intn(n)])
	}
	return string(result)
}
