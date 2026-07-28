package utils

import (
	"time"
)

// TimeFormat تنسيقات الوقت
const (
	TimeFormatRFC3339  = time.RFC3339
	TimeFormatISO      = "2006-01-02T15:04:05"
	TimeFormatDate     = "2006-01-02"
	TimeFormatTime     = "15:04:05"
	TimeFormatDateTime = "2006-01-02 15:04:05"
)

// FormatTime ينسق الوقت
func FormatTime(t time.Time, format string) string {
	return t.Format(format)
}

// ParseTime يحلل الوقت
func ParseTime(value, format string) (time.Time, error) {
	return time.Parse(format, value)
}

// Now يعيد الوقت الحالي
func Now() time.Time {
	return time.Now()
}

// UTCNow يعيد الوقت الحالي بتوقيت UTC
func UTCNow() time.Time {
	return time.Now().UTC()
}

// StartOfDay يعيد بداية اليوم
func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// EndOfDay يعيد نهاية اليوم
func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}

// DaysBetween يعيد عدد الأيام بين تاريخين
func DaysBetween(start, end time.Time) int {
	start = StartOfDay(start)
	end = StartOfDay(end)
	return int(end.Sub(start).Hours() / 24)
}

// AddDays يضيف أياماً
func AddDays(t time.Time, days int) time.Time {
	return t.Add(time.Duration(days) * 24 * time.Hour)
}

// IsToday يتحقق من أن التاريخ هو اليوم
func IsToday(t time.Time) bool {
	now := Now()
	return t.Year() == now.Year() && t.Month() == now.Month() && t.Day() == now.Day()
}
