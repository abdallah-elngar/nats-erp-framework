package validation

import (
	"strconv"
	"strings"
)

// Rule يمثل قاعدة تحقق
type Rule struct {
	Tag     string
	Message string
	Param   string
}

// ValidationRules يمثل مجموعة من قواعد التحقق
type ValidationRules map[string][]Rule

// NewValidationRules ينشئ مجموعة قواعد جديدة
func NewValidationRules() ValidationRules {
	return make(ValidationRules)
}

// Add يضيف قاعدة تحقق
func (r ValidationRules) Add(field, tag, message, param string) {
	r[field] = append(r[field], Rule{
		Tag:     tag,
		Message: message,
		Param:   param,
	})
}

// Required يضيف قاعدة مطلوب
func (r ValidationRules) Required(field string) {
	r.Add(field, "required", "هذا الحقل مطلوب", "")
}

// Email يضيف قاعدة بريد إلكتروني
func (r ValidationRules) Email(field string) {
	r.Add(field, "email", "البريد الإلكتروني غير صحيح", "")
}

// Min يضيف قاعدة حد أدنى
func (r ValidationRules) Min(field string, min int) {
	r.Add(field, "min", "القيمة صغيرة جداً (الحد الأدنى: {param})", strconv.Itoa(min))
}

// Max يضيف قاعدة حد أقصى
func (r ValidationRules) Max(field string, max int) {
	r.Add(field, "max", "القيمة كبيرة جداً (الحد الأقصى: {param})", strconv.Itoa(max))
}

// MinLength يضيف قاعدة طول أدنى
func (r ValidationRules) MinLength(field string, min int) {
	r.Add(field, "min", "النص قصير جداً (الحد الأدنى: {param})", strconv.Itoa(min))
}

// MaxLength يضيف قاعدة طول أقصى
func (r ValidationRules) MaxLength(field string, max int) {
	r.Add(field, "max", "النص طويل جداً (الحد الأقصى: {param})", strconv.Itoa(max))
}

// Phone يضيف قاعدة رقم هاتف
func (r ValidationRules) Phone(field string) {
	r.Add(field, "phone", "رقم الهاتف غير صحيح", "")
}

// Password يضيف قاعدة كلمة مرور
func (r ValidationRules) Password(field string) {
	r.Add(field, "password", "كلمة المرور ضعيفة (8 أحرف على الأقل، حرف كبير، حرف صغير، رقم، رمز خاص)", "")
}

// Date يضيف قاعدة تاريخ
func (r ValidationRules) Date(field string) {
	r.Add(field, "date", "التاريخ غير صحيح", "")
}

// DateTime يضيف قاعدة تاريخ ووقت
func (r ValidationRules) DateTime(field string) {
	r.Add(field, "datetime", "التاريخ والوقت غير صحيحين", "")
}

// JSON يضيف قاعدة JSON
func (r ValidationRules) JSON(field string) {
	r.Add(field, "json", "JSON غير صحيح", "")
}

// Unique يضيف قاعدة فريدة
func (r ValidationRules) Unique(field string) {
	r.Add(field, "unique", "يجب أن يكون فريداً", "")
}

// Exists يضيف قاعدة موجود
func (r ValidationRules) Exists(field string) {
	r.Add(field, "exists", "غير موجود", "")
}

// ValidateStruct يقوم بالتحقق من هيكل باستخدام القواعد
func (r ValidationRules) ValidateStruct(data interface{}) error {
	// سيتم تنفيذها باستخدام validator
	return nil
}

// GetValidationTags يحصل على علامات التحقق
func (r ValidationRules) GetValidationTags() map[string]string {
	tags := make(map[string]string)

	for field, rules := range r {
		var tagParts []string
		for _, rule := range rules {
			if rule.Param != "" {
				tagParts = append(tagParts, rule.Tag+"="+rule.Param)
			} else {
				tagParts = append(tagParts, rule.Tag)
			}
		}
		tags[field] = strings.Join(tagParts, ",")
	}

	return tags
}
