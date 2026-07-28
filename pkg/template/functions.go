package template

import (
	"encoding/json"
	"fmt"
	"html/template"
	"reflect"
	"strings"
	"time"
)

// FuncMap يعيد خريطة الدوال المساعدة
func FuncMap() template.FuncMap {
	return template.FuncMap{
		"upper":        strings.ToUpper,
		"lower":        strings.ToLower,
		"title":        strings.Title,
		"join":         strings.Join,
		"add":          func(a, b int) int { return a + b },
		"sub":          func(a, b int) int { return a - b },
		"mul":          func(a, b int) int { return a * b },
		"div":          func(a, b int) int { return a / b },
		"mod":          func(a, b int) int { return a % b },
		"safe":         func(s string) template.HTML { return template.HTML(s) },
		"safeURL":      func(s string) template.URL { return template.URL(s) },
		"json":         jsonMarshal,
		"time":         func() string { return time.Now().Format("2006-01-02 15:04:05") },
		"formatTime":   formatTime,
		"formatDate":   func(t time.Time) string { return t.Format("2006-01-02") },
		"formatNumber": formatNumber,
		"truncate":     truncate,
		"default":      defaultFunc,
		"eq":           func(a, b interface{}) bool { return a == b },
		"ne":           func(a, b interface{}) bool { return a != b },
		"gt":           func(a, b int) bool { return a > b },
		"lt":           func(a, b int) bool { return a < b },
		"not":          func(v bool) bool { return !v },
		"and":          func(a, b bool) bool { return a && b },
		"or":           func(a, b bool) bool { return a || b },
		"len":          func(v interface{}) int { return reflect.ValueOf(v).Len() },
		"range":        func(start, end int) []int { return makeRange(start, end) },
		"dict":         makeDict,
		"get":          getDict,
		"set":          setDict,
	}
}

// jsonMarshal يحول إلى JSON
func jsonMarshal(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}

// formatTime ينسق الوقت
func formatTime(t time.Time, layout string) string {
	if layout == "" {
		layout = "2006-01-02 15:04:05"
	}
	return t.Format(layout)
}

// formatNumber ينسق الأرقام
func formatNumber(n float64) string {
	return fmt.Sprintf("%.2f", n)
}

// truncate يقطع النص
func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

// defaultFunc يعيد القيمة الافتراضية
func defaultFunc(defaultVal interface{}, val interface{}) interface{} {
	if val == nil {
		return defaultVal
	}
	if v, ok := val.(string); ok && v == "" {
		return defaultVal
	}
	return val
}

// makeRange يصنع نطاقاً من الأرقام
func makeRange(start, end int) []int {
	var r []int
	for i := start; i <= end; i++ {
		r = append(r, i)
	}
	return r
}

// makeDict يصنع قاموساً
func makeDict(values ...interface{}) map[string]interface{} {
	if len(values)%2 != 0 {
		return nil
	}
	dict := make(map[string]interface{})
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			continue
		}
		dict[key] = values[i+1]
	}
	return dict
}

// getDict يحصل على قيمة من القاموس
func getDict(dict map[string]interface{}, key string) interface{} {
	if dict == nil {
		return nil
	}
	return dict[key]
}

// setDict يضع قيمة في القاموس
func setDict(dict map[string]interface{}, key string, value interface{}) map[string]interface{} {
	if dict == nil {
		dict = make(map[string]interface{})
	}
	dict[key] = value
	return dict
}
