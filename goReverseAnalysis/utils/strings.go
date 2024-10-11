package utils

import "strings"

// StringManipulator 提供字符串操作功能
type StringManipulator struct{}

func (s StringManipulator) ToUpper(input string) string {
	return strings.ToUpper(input)
}

func (s StringManipulator) ToLower(input string) string {
	return strings.ToLower(input)
}
