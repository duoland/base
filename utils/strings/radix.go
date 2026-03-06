package strings

import (
	"fmt"
	"strings"
)

const Base32Chars = "0123456789abcdefghijklmnopqrstuv"
const Base62Chars = "0123456789abcdefghijklmnopqrstuvABCDEFGHIJKLMNOPQRSTUV"

// Int64ToBase32 将 int64 类型的整数转换为 32 进制字符串
func Int64ToBase32(num int64) string {
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
		result.WriteByte(Base32Chars[remainder])
		// 更新 num 的值，进行下一轮取模
		num /= 32
	}
	// 反转结果字符串
	reversed := result.String()
	return reverseString(reversed)
}

// Base32ToInt64 将 32 进制字符串转换为 int64 类型的整数
func Base32ToInt64(base32Str string) (int64, error) {
	var result int64
	for _, char := range base32Str {
		// 查找字符在 32 进制字符集中的位置
		index := strings.IndexRune(Base32Chars, char)
		if index == -1 {
			return 0, fmt.Errorf("invalid base32 character: %c", char)
		}
		// 更新结果值，通过乘以 32 加上当前字符对应的数值
		result = result*32 + int64(index)
	}
	return result, nil
}

// Int64ToBase62 将 int64 类型的整数转换为 62 进制字符串
func Int64ToBase62(num int64) string {
	// 处理 num 为 0 的特殊情况
	if num == 0 {
		return "0"
	}
	// 用于存储转换结果的字符串构建器
	var result strings.Builder
	// 循环取模，直到 num 变为 0
	for num > 0 {
		// 计算当前位对应的 62 进制字符的索引
		remainder := num % 62
		// 将对应的字符添加到结果构建器中
		result.WriteByte(Base62Chars[remainder])
		// 更新 num 的值，进行下一轮取模
		num /= 62
	}
	// 反转结果字符串
	reversed := result.String()
	return reverseString(reversed)
}

// Base62ToInt64 将 62 进制字符串转换为 int64 类型的整数
func Base62ToInt64(base62Str string) (int64, error) {
	var result int64
	for _, char := range base62Str {
		// 查找字符在 62 进制字符集中的位置
		index := strings.IndexRune(Base62Chars, char)
		if index == -1 {
			return 0, fmt.Errorf("invalid base62 character: %c", char)
		}
		// 更新结果值，通过乘以 62 加上当前字符对应的数值
		result = result*62 + int64(index)
	}
	return result, nil
}
