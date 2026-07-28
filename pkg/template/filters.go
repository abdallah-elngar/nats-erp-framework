// pkg/template/filters.go (مكتمل)
package template

import (
	"fmt"
	"html/template"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

// FilterFunc دالة الفلتر
type FilterFunc func(value interface{}, args ...interface{}) interface{}

// FilterRegistry سجل الفلاتر
type FilterRegistry struct {
	filters map[string]FilterFunc
	mu      sync.RWMutex
}

// NewFilterRegistry ينشئ سجل فلاتر جديد
func NewFilterRegistry() *FilterRegistry {
	return &FilterRegistry{
		filters: make(map[string]FilterFunc),
	}
}

// Register يسجل فلتراً جديداً
func (r *FilterRegistry) Register(name string, fn FilterFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.filters[name] = fn
}

// Get يعيد فلتراً
func (r *FilterRegistry) Get(name string) (FilterFunc, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	fn, ok := r.filters[name]
	return fn, ok
}

// ============================================
// فلاتر النصوص
// ============================================

// FilterUpper يحول النص إلى أحرف كبيرة
func FilterUpper(value interface{}, args ...interface{}) interface{} {
	str := fmt.Sprintf("%v", value)
	return strings.ToUpper(str)
}

// FilterLower يحول النص إلى أحرف صغيرة
func FilterLower(value interface{}, args ...interface{}) interface{} {
	str := fmt.Sprintf("%v", value)
	return strings.ToLower(str)
}

// FilterTitle يحول النص إلى Title Case
func FilterTitle(value interface{}, args ...interface{}) interface{} {
	str := fmt.Sprintf("%v", value)
	return strings.Title(strings.ToLower(str))
}

// FilterCapitalize يحول الحرف الأول إلى كبير
func FilterCapitalize(value interface{}, args ...interface{}) interface{} {
	str := fmt.Sprintf("%v", value)
	if len(str) == 0 {
		return str
	}
	return strings.ToUpper(str[:1]) + str[1:]
}

// FilterSlug يحول النص إلى Slug صالح للـ URL
func FilterSlug(value interface{}, args ...interface{}) interface{} {
	str := fmt.Sprintf("%v", value)
	str = strings.ToLower(str)
	str = strings.ReplaceAll(str, " ", "-")
	reg := regexp.MustCompile(`[^a-z0-9-]`)
	str = reg.ReplaceAllString(str, "")
	reg2 := regexp.MustCompile(`-+`)
	str = reg2.ReplaceAllString(str, "-")
	str = strings.Trim(str, "-")
	return str
}

// FilterTruncate يقص النص إلى طول محدد
func FilterTruncate(value interface{}, args ...interface{}) interface{} {
	str := fmt.Sprintf("%v", value)
	length := 100
	suffix := "..."

	if len(args) > 0 {
		if l, ok := args[0].(int); ok {
			length = l
		}
	}
	if len(args) > 1 {
		if s, ok := args[1].(string); ok {
			suffix = s
		}
	}

	if len(str) <= length {
		return str
	}

	return str[:length] + suffix
}

// FilterWordCount يحسب عدد الكلمات
func FilterWordCount(value interface{}, args ...interface{}) interface{} {
	str := fmt.Sprintf("%v", value)
	words := strings.Fields(str)
	return len(words)
}

// FilterLineBreaks يحول الأسطر الجديدة إلى <br>
func FilterLineBreaks(value interface{}, args ...interface{}) interface{} {
	str := fmt.Sprintf("%v", value)
	return template.HTML(strings.ReplaceAll(str, "\n", "<br>"))
}

// FilterStripTags يزيل وسوم HTML
func FilterStripTags(value interface{}, args ...interface{}) interface{} {
	str := fmt.Sprintf("%v", value)
	reg := regexp.MustCompile(`<[^>]*>`)
	return reg.ReplaceAllString(str, "")
}

// FilterEscape يهرب النص HTML
func FilterEscape(value interface{}, args ...interface{}) interface{} {
	str := fmt.Sprintf("%v", value)
	return template.HTMLEscapeString(str)
}

// FilterSafe يعتبر النص آمناً (يعرض HTML)
func FilterSafe(value interface{}, args ...interface{}) interface{} {
	str := fmt.Sprintf("%v", value)
	return template.HTML(str)
}

// ============================================
// فلاتر الأرقام
// ============================================

// FilterAdd يجمع رقمين
func FilterAdd(value interface{}, args ...interface{}) interface{} {
	val := toFloat64(value)
	if len(args) > 0 {
		val += toFloat64(args[0])
	}
	return val
}

// FilterSub يطرح رقمين
func FilterSub(value interface{}, args ...interface{}) interface{} {
	val := toFloat64(value)
	if len(args) > 0 {
		val -= toFloat64(args[0])
	}
	return val
}

// FilterMul يضرب رقمين
func FilterMul(value interface{}, args ...interface{}) interface{} {
	val := toFloat64(value)
	if len(args) > 0 {
		val *= toFloat64(args[0])
	}
	return val
}

// FilterDiv يقسم رقمين
func FilterDiv(value interface{}, args ...interface{}) interface{} {
	val := toFloat64(value)
	if len(args) > 0 {
		divisor := toFloat64(args[0])
		if divisor != 0 {
			val /= divisor
		}
	}
	return val
}

// FilterFormatNumber ينسق الأرقام مع فواصل الآلاف
func FilterFormatNumber(value interface{}, args ...interface{}) interface{} {
	val := toFloat64(value)
	return FormatNumber(val)
}

// FilterCurrency ينسق كعملة
func FilterCurrency(value interface{}, args ...interface{}) interface{} {
	val := toFloat64(value)
	currency := "$"
	if len(args) > 0 {
		if c, ok := args[0].(string); ok {
			currency = c
		}
	}
	return fmt.Sprintf("%s%s", currency, FormatNumber(val))
}

// FilterPercentage ينسق كنسبة مئوية
func FilterPercentage(value interface{}, args ...interface{}) interface{} {
	val := toFloat64(value)
	return fmt.Sprintf("%.2f%%", val)
}

// ============================================
// فلاتر التواريخ (جديدة)
// ============================================

// FilterDate ينسق التاريخ
func FilterDate(value interface{}, args ...interface{}) interface{} {
	t, err := parseTime(value)
	if err != nil {
		return value
	}
	format := "2006-01-02"
	if len(args) > 0 {
		if f, ok := args[0].(string); ok {
			format = f
		}
	}
	return t.Format(format)
}

// FilterTime ينسق الوقت
func FilterTime(value interface{}, args ...interface{}) interface{} {
	t, err := parseTime(value)
	if err != nil {
		return value
	}
	format := "15:04:05"
	if len(args) > 0 {
		if f, ok := args[0].(string); ok {
			format = f
		}
	}
	return t.Format(format)
}

// FilterDateTime ينسق التاريخ والوقت
func FilterDateTime(value interface{}, args ...interface{}) interface{} {
	t, err := parseTime(value)
	if err != nil {
		return value
	}
	format := "2006-01-02 15:04:05"
	if len(args) > 0 {
		if f, ok := args[0].(string); ok {
			format = f
		}
	}
	return t.Format(format)
}

// FilterTimeAgo يحسب الوقت المنقضي
func FilterTimeAgo(value interface{}, args ...interface{}) interface{} {
	t, err := parseTime(value)
	if err != nil {
		return value
	}
	return TimeAgo(t)
}

// FilterFormatDateTime ينسق التاريخ والوقت بصيغة مخصصة
func FilterFormatDateTime(value interface{}, args ...interface{}) interface{} {
	t, err := parseTime(value)
	if err != nil {
		return value
	}
	format := "2006-01-02 15:04:05"
	if len(args) > 0 {
		if f, ok := args[0].(string); ok {
			format = f
		}
	}
	return t.Format(format)
}

// ============================================
// فلاتر المصفوفات (جديدة)
// ============================================

// FilterJoin يدمج عناصر المصفوفة
func FilterJoin(value interface{}, args ...interface{}) interface{} {
	// محاولة تحويل إلى مصفوفة من strings
	switch v := value.(type) {
	case []string:
		separator := ", "
		if len(args) > 0 {
			if s, ok := args[0].(string); ok {
				separator = s
			}
		}
		return strings.Join(v, separator)
	case []interface{}:
		separator := ", "
		if len(args) > 0 {
			if s, ok := args[0].(string); ok {
				separator = s
			}
		}
		var parts []string
		for _, item := range v {
			parts = append(parts, fmt.Sprintf("%v", item))
		}
		return strings.Join(parts, separator)
	default:
		return value
	}
}

// FilterSplit يقسم النص إلى مصفوفة
func FilterSplit(value interface{}, args ...interface{}) interface{} {
	str := fmt.Sprintf("%v", value)
	separator := ","
	if len(args) > 0 {
		if s, ok := args[0].(string); ok {
			separator = s
		}
	}
	return strings.Split(str, separator)
}

// FilterReverse يعكس المصفوفة
func FilterReverse(value interface{}, args ...interface{}) interface{} {
	switch v := value.(type) {
	case []interface{}:
		reversed := make([]interface{}, len(v))
		for i, j := 0, len(v)-1; i < len(v); i, j = i+1, j-1 {
			reversed[i] = v[j]
		}
		return reversed
	case []string:
		reversed := make([]string, len(v))
		for i, j := 0, len(v)-1; i < len(v); i, j = i+1, j-1 {
			reversed[i] = v[j]
		}
		return reversed
	default:
		return value
	}
}

// FilterFirst يعيد العنصر الأول
func FilterFirst(value interface{}, args ...interface{}) interface{} {
	switch v := value.(type) {
	case []interface{}:
		if len(v) > 0 {
			return v[0]
		}
	case []string:
		if len(v) > 0 {
			return v[0]
		}
	}
	return nil
}

// FilterLast يعيد العنصر الأخير
func FilterLast(value interface{}, args ...interface{}) interface{} {
	switch v := value.(type) {
	case []interface{}:
		if len(v) > 0 {
			return v[len(v)-1]
		}
	case []string:
		if len(v) > 0 {
			return v[len(v)-1]
		}
	}
	return nil
}

// FilterLength يعيد طول المصفوفة أو النص
func FilterLength(value interface{}, args ...interface{}) interface{} {
	switch v := value.(type) {
	case []interface{}:
		return len(v)
	case []string:
		return len(v)
	case string:
		return len(v)
	default:
		return 0
	}
}

// FilterSlice يعيد جزءاً من المصفوفة
func FilterSlice(value interface{}, args ...interface{}) interface{} {
	switch v := value.(type) {
	case []interface{}:
		start := 0
		end := len(v)

		if len(args) > 0 {
			if s, ok := args[0].(int); ok {
				start = s
			}
		}
		if len(args) > 1 {
			if e, ok := args[1].(int); ok {
				end = e
			}
		}

		if start < 0 {
			start = len(v) + start
		}
		if end < 0 {
			end = len(v) + end
		}
		if start < 0 {
			start = 0
		}
		if end > len(v) {
			end = len(v)
		}
		if start > end {
			return []interface{}{}
		}
		return v[start:end]
	default:
		return value
	}
}

// ============================================
// فلاتر منطقية (جديدة)
// ============================================

// FilterYesNo يحول القيمة المنطقية إلى نعم/لا
func FilterYesNo(value interface{}, args ...interface{}) interface{} {
	if toBool(value) {
		return "Yes"
	}
	return "No"
}

// FilterBool يحول القيمة إلى منطقية
func FilterBool(value interface{}, args ...interface{}) interface{} {
	return toBool(value)
}

// FilterDefault يعيد قيمة افتراضية إذا كانت القيمة فارغة
func FilterDefault(value interface{}, args ...interface{}) interface{} {
	if value == nil || value == "" || value == 0 {
		if len(args) > 0 {
			return args[0]
		}
		return ""
	}
	return value
}

// ============================================
// دوال مساعدة
// ============================================

// toFloat64 يحول القيمة إلى float64
func toFloat64(value interface{}) float64 {
	switch v := value.(type) {
	case int:
		return float64(v)
	case int64:
		return float64(v)
	case float64:
		return v
	case float32:
		return float64(v)
	case string:
		f, _ := strconv.ParseFloat(v, 64)
		return f
	default:
		return 0
	}
}

// toBool يحول القيمة إلى bool
func toBool(value interface{}) bool {
	switch v := value.(type) {
	case bool:
		return v
	case int:
		return v != 0
	case int64:
		return v != 0
	case string:
		return v != "" && v != "0" && v != "false"
	default:
		return false
	}
}

// parseTime يحول القيمة إلى time.Time
func parseTime(value interface{}) (time.Time, error) {
	switch v := value.(type) {
	case time.Time:
		return v, nil
	case string:
		layouts := []string{
			time.RFC3339,
			"2006-01-02 15:04:05",
			"2006-01-02",
			"15:04:05",
			"2006-01-02T15:04:05Z",
		}
		for _, layout := range layouts {
			if t, err := time.Parse(layout, v); err == nil {
				return t, nil
			}
		}
		return time.Time{}, fmt.Errorf("invalid time format: %s", v)
	default:
		return time.Time{}, fmt.Errorf("unsupported type: %T", v)
	}
}

// FormatNumber ينسق الأرقام مع فواصل الآلاف
func FormatNumber(n float64) string {
	parts := strings.Split(fmt.Sprintf("%.2f", n), ".")
	intPart := parts[0]
	decPart := parts[1]

	var result strings.Builder
	for i, c := range intPart {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			result.WriteString(",")
		}
		result.WriteRune(c)
	}

	return result.String() + "." + decPart
}

// TimeAgo يحسب الوقت المنقضي
func TimeAgo(t time.Time) string {
	duration := time.Since(t)

	if duration < time.Minute {
		return "Just now"
	} else if duration < time.Hour {
		minutes := int(duration.Minutes())
		return fmt.Sprintf("%d minute%s ago", minutes, pluralize(minutes))
	} else if duration < 24*time.Hour {
		hours := int(duration.Hours())
		return fmt.Sprintf("%d hour%s ago", hours, pluralize(hours))
	} else if duration < 30*24*time.Hour {
		days := int(duration.Hours() / 24)
		return fmt.Sprintf("%d day%s ago", days, pluralize(days))
	} else if duration < 365*24*time.Hour {
		months := int(duration.Hours() / (24 * 30))
		return fmt.Sprintf("%d month%s ago", months, pluralize(months))
	} else {
		years := int(duration.Hours() / (24 * 365))
		return fmt.Sprintf("%d year%s ago", years, pluralize(years))
	}
}

// pluralize يضيف s للجمع
func pluralize(n int) string {
	if n > 1 {
		return "s"
	}
	return ""
}
