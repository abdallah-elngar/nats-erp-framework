package logger

import (
	"fmt"
	"log/slog"
	"strings"
	"time"
)

// Formatter ينسق السجلات
type Formatter struct {
	format string
}

// NewFormatter ينشئ منسقاً جديداً
func NewFormatter(format string) *Formatter {
	return &Formatter{format: format}
}

// Format ينسق السجل
func (f *Formatter) Format(level slog.Level, msg string, args ...interface{}) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	levelStr := strings.ToUpper(level.String())

	switch f.format {
	case "json":
		return f.formatJSON(timestamp, levelStr, msg, args)
	case "text":
		return f.formatText(timestamp, levelStr, msg, args)
	default:
		return fmt.Sprintf("[%s] %s: %s %v\n", timestamp, levelStr, msg, args)
	}
}

// formatJSON ينسق بتنسيق JSON
func (f *Formatter) formatJSON(timestamp, level, msg string, args []interface{}) string {
	var fields string
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			fields += fmt.Sprintf(",\"%v\":\"%v\"", args[i], args[i+1])
		}
	}

	return fmt.Sprintf(`{"time":"%s","level":"%s","msg":"%s"%s}`+"\n", timestamp, level, msg, fields)
}

// formatText ينسق بتنسيق نصي
func (f *Formatter) formatText(timestamp, level, msg string, args []interface{}) string {
	var fields string
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			fields += fmt.Sprintf(" %v=%v", args[i], args[i+1])
		}
	}

	return fmt.Sprintf("[%s] %s: %s%s\n", timestamp, level, msg, fields)
}
