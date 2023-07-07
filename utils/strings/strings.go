package strings

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func SliceContains(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// ReverseRunes reverse a string runes and return the result
func ReverseRunes(str string) (output string) {
	if str == "" {
		return
	}

	runes := []rune(str)

	buf := make([]rune, 0, len(runes))
	for i := len(runes) - 1; i >= 0; i-- {
		buf = append(buf, runes[i])
	}
	output = string(buf)
	return
}

// ReverseString reverse a string content and return the result
func ReverseString(str string) (output string) {
	if str == "" {
		return
	}
	buf := make([]byte, 0, len(str))
	for i := len(str) - 1; i >= 0; i-- {
		buf = append(buf, str[i])
	}
	output = string(buf)
	return
}

func TruncateRuneToSize(str string, maxAllowedSize int) string {
	chars := []rune(str)
	if len(chars) <= maxAllowedSize {
		return str
	}
	// otherwise
	return string(chars[:maxAllowedSize])
}

func ToJsonString(src any) string {
	return string(ToJsonBytes(src))
}

func ToJsonBytes(src any) []byte {
	data, _ := json.Marshal(src)
	return data
}

func ToJsonStringPretty(src any) string {
	dst, _ := json.MarshalIndent(src, "", "\t")
	return string(dst)
}

func Int64ToString(value int64) string {
	return fmt.Sprintf("%d", value)
}

func StringToInt64(str string) int64 {
	val, _ := strconv.ParseInt(str, 10, 64)
	return val
}

func IsEmpty(s *string) bool {
	return s == nil || *s == ""
}

func IsNotEmpty(s *string) bool {
	return s != nil && *s != ""
}
