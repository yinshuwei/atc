package strs

import (
	"html/template"
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

// SplitI 分割字符串,并返回[]interface{}
func (s Strs) SplitI(str string, sep string) []interface{} {
	var result []interface{}
	strs := strings.Split(str, sep)
	for _, str := range strs {
		result = append(result, str)
	}
	return result
}

// At 取Item
func (s Strs) At(arr []string, index int) string {
	if len(arr) > index && index >= 0 {
		return arr[index]
	}
	return ""
}

// Trim 去两头
func (s Strs) Trim(str string, cutset string) string {
	return strings.Trim(str, cutset)
}

// TrimSpace 去两头空格
func (s Strs) TrimSpace(str string) string {
	return strings.TrimSpace(str)
}

// QueryEscape URL编码
func (s Strs) QueryEscape(str string) string {
	return url.QueryEscape(str)
}

// String to string
func (s Strs) String(v interface{}) string {
	return v.(string)
}

// ReplaceAll Replace All
func (s Strs) ReplaceAll(v interface{}, old, new string) string {
	r := v.(string)
	return strings.Replace(r, old, new, -1)
}

// HTML to html
func (s Strs) HTML(v interface{}) template.HTML {
	r := v.(string)
	return template.HTML(r)
}

// CSS to css
func (s Strs) CSS(v interface{}) template.CSS {
	r := v.(string)
	return template.CSS(r)
}

// Pre string to html
func (s Strs) Pre(v interface{}) template.HTML {
	r := v.(string)
	r = strings.Replace(r, "\n", "<br/>", -1)
	r = strings.Replace(r, " ", "&nbsp;&nbsp;&nbsp;", -1)
	return template.HTML(r)
}
