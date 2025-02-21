package strings

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
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

func Uint64ToString(value uint64) string {
	return fmt.Sprintf("%d", value)
}

func StringToInt64(str string) int64 {
	val, _ := strconv.ParseInt(str, 10, 64)
	return val
}

func StringToUint64(str string) uint64 {
	val, _ := strconv.ParseUint(str, 10, 64)
	return val
}

func IsEmpty(s *string) bool {
	return s == nil || *s == ""
}

func IsNotEmpty(s *string) bool {
	return s != nil && *s != ""
}

func Int64SliceToStringSlice(values []int64) (valueStrs []string) {
	valueStrs = make([]string, 0, len(values))
	for _, value := range values {
		valueStrs = append(valueStrs, Int64ToString(value))
	}
	return
}

func StringSliceToInt64Slice(valueStrs []string) (values []int64) {
	values = make([]int64, 0, len(valueStrs))
	for _, valueStr := range valueStrs {
		values = append(values, StringToInt64(valueStr))
	}
	return
}

// Int64ToBase32 将 int64 类型的整数转换为 32 进制字符串
func Int64ToBase32(num int64) string {
	// 定义 32 进制字符集
	base32Chars := "0123456789abcdefghijklmnopqrstuv"
	// 处理 num 为 0 的特殊情况
	if num == 0 {
		return "0"
	}
	// 用于存储转换结果的字符串构建器
	var result strings.Builder
	// 循环取模，直到 num 变为 0
	for num > 0 {
		// 计算当前位对应的 32 进制字符的索引
		remainder := num % 32
		// 将对应的字符添加到结果构建器中
		result.WriteByte(base32Chars[remainder])
		// 更新 num 的值，进行下一轮取模
		num /= 32
	}
	// 反转结果字符串
	reversed := result.String()
	return reverseString(reversed)
}

// reverseString 反转字符串
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
