package strs

import (
	"net/url"
	"strings"
)

// Strs  Strs的包装
type Strs struct{}

// Substr 子字符串
func (s Strs) Substr(str string, start, end int) string {
	if start >= 0 && start <= end && end <= len(str) {
		return str[start:end]
	}
	return str
}

// Split 分割字符串
func (s Strs) Split(str string, sep string) []string {
	return strings.Split(str, sep)
}

// QueryEscape 分割字符串
func (s Strs) QueryEscape(str string) string {
	return url.QueryEscape(str)
}
