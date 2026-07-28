package validation

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

// Validator يقوم بالتحقق من البيانات
type Validator struct {
	validate *validator.Validate
	messages map[string]string
}

// NewValidator ينشئ محققاً جديداً
func NewValidator() *Validator {
	v := validator.New()

	// تسجيل محققات مخصصة
	v.RegisterValidation("phone", validatePhone)
	v.RegisterValidation("password", validatePassword)
	v.RegisterValidation("date", validateDate)
	v.RegisterValidation("datetime", validateDateTime)
	v.RegisterValidation("json", validateJSON)

	// تسجيل رسائل مخصصة
	messages := map[string]string{
		"required": "هذا الحقل مطلوب",
		"email":    "البريد الإلكتروني غير صحيح",
		"min":      "القيمة صغيرة جداً",
		"max":      "القيمة كبيرة جداً",
		"len":      "الطول غير صحيح",
		"phone":    "رقم الهاتف غير صحيح",
		"password": "كلمة المرور ضعيفة",
		"date":     "التاريخ غير صحيح",
		"datetime": "التاريخ والوقت غير صحيحين",
		"json":     "JSON غير صحيح",
	}

	return &Validator{
		validate: v,
		messages: messages,
	}
}

// Validate يقوم بالتحقق من البيانات
func (v *Validator) Validate(data interface{}) error {
	return v.validate.Struct(data)
}

// ValidateVar يقوم بالتحقق من متغير
func (v *Validator) ValidateVar(value interface{}, tag string) error {
	return v.validate.Var(value, tag)
}

// GetValidationErrors يحصل على أخطاء التحقق
func (v *Validator) GetValidationErrors(err error) map[string]string {
	if err == nil {
		return nil
	}

	errors := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, ve := range validationErrors {
			field := ve.Field()
			tag := ve.Tag()
			param := ve.Param()

			// الحصول على رسالة الخطأ
			message := v.getErrorMessage(field, tag, param)
			errors[field] = message
		}
	}

	return errors
}

// getErrorMessage يحصل على رسالة الخطأ
func (v *Validator) getErrorMessage(field, tag, param string) string {
	// البحث عن رسالة مخصصة
	key := field + "." + tag
	if msg, ok := v.messages[key]; ok {
		return msg
	}

	// البحث عن رسالة عامة
	if msg, ok := v.messages[tag]; ok {
		return strings.ReplaceAll(msg, "{param}", param)
	}

	return fmt.Sprintf("%s غير صالح", field)
}

// AddCustomMessage يضيف رسالة مخصصة
func (v *Validator) AddCustomMessage(field, tag, message string) {
	key := field + "." + tag
	v.messages[key] = message
}

// AddValidation يضيف محققاً مخصصاً
func (v *Validator) AddValidation(tag string, fn validator.Func) {
	v.validate.RegisterValidation(tag, fn)
}

// ============================================
// محققات مخصصة
// ============================================

// validatePhone يتحقق من رقم الهاتف
func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	if phone == "" {
		return true
	}

	// رقم هاتف مصري: 01XXXXXXXXX
	// رقم هاتف دولي: +201XXXXXXXXX
	re := regexp.MustCompile(`^(\+201|01)[0-9]{9}$`)
	return re.MatchString(phone)
}

// validatePassword يتحقق من كلمة المرور
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if password == "" {
		return true
	}

	// على الأقل 8 أحرف، حرف كبير، حرف صغير، رقم، رمز خاص
	if len(password) < 8 {
		return false
	}

	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasNumber = true
		default:
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// validateDate يتحقق من التاريخ
func validateDate(fl validator.FieldLevel) bool {
	date := fl.Field().String()
	if date == "" {
		return true
	}

	// محاولة تحليل التاريخ بتنسيقات مختلفة
	formats := []string{
		"2006-01-02",
		"02-01-2006",
		"02/01/2006",
		"2006/01/02",
	}

	for _, format := range formats {
		if _, err := time.Parse(format, date); err == nil {
			return true
		}
	}

	return false
}

// validateDateTime يتحقق من التاريخ والوقت
func validateDateTime(fl validator.FieldLevel) bool {
	datetime := fl.Field().String()
	if datetime == "" {
		return true
	}

	// محاولة تحليل التاريخ والوقت بتنسيقات مختلفة
	formats := []string{
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z",
		"2006-01-02T15:04:05+07:00",
	}

	for _, format := range formats {
		if _, err := time.Parse(format, datetime); err == nil {
			return true
		}
	}

	return false
}

// validateJSON يتحقق من JSON
func validateJSON(fl validator.FieldLevel) bool {
	jsonStr := fl.Field().String()
	if jsonStr == "" {
		return true
	}

	var js interface{}
	return json.Unmarshal([]byte(jsonStr), &js) == nil
}
