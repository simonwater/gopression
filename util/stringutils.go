package util

import (
	"strings"
)

func Clean(str string) string {
	return strings.TrimSpace(str)
}

func Trim(str string) string {
	return strings.TrimSpace(str)
}

func IsNotEmpty(str string) bool {
	return str != ""
}

func IsEmpty(str string) bool {
	return str == ""
}

func IsBlank(str string) bool {
	return strings.TrimSpace(str) == ""
}

func IsNotBlank(str string) bool {
	return !IsBlank(str)
}

func Equals(str1, str2 string) bool {
	return str1 == str2
}

func EqualsIgnoreCase(str1, str2 string) bool {
	return strings.EqualFold(str1, str2)
}

func GetUTF8Bytes(str string) []byte {
	return []byte(str)
}

func GetUTF8String(bytesArr []byte) string {
	return string(bytesArr)
}
