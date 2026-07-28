package binding

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Binder يقوم بربط البيانات
type Binder struct {
	validator *validator.Validate
}

// NewBinder ينشئ رابطاً جديداً
func NewBinder() *Binder {
	return &Binder{
		validator: validator.New(),
	}
}

// BindJSON يربط بيانات JSON
func (b *Binder) BindJSON(r *http.Request, dest interface{}) error {
	// قراءة الجسم
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(dest); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	// التحقق من البيانات
	return b.Validate(dest)
}

// BindForm يربط بيانات النموذج
func (b *Binder) BindForm(r *http.Request, dest interface{}) error {
	// تحليل النموذج
	if err := r.ParseForm(); err != nil {
		return fmt.Errorf("failed to parse form: %w", err)
	}

	// ربط الحقول
	if err := b.bindFormData(r.Form, dest); err != nil {
		return err
	}

	// التحقق من البيانات
	return b.Validate(dest)
}

// BindQuery يربط معلمات الاستعلام
func (b *Binder) BindQuery(r *http.Request, dest interface{}) error {
	// ربط معلمات الاستعلام
	if err := b.bindQueryData(r.URL.Query(), dest); err != nil {
		return err
	}

	// التحقق من البيانات
	return b.Validate(dest)
}

// BindMultipart يربط بيانات متعددة الأجزاء
func (b *Binder) BindMultipart(r *http.Request, dest interface{}) error {
	// تحليل النموذج المتعدد الأجزاء
	if err := r.ParseMultipartForm(32 << 20); err != nil { // 32MB
		return fmt.Errorf("failed to parse multipart form: %w", err)
	}

	// ربط البيانات
	if err := b.bindFormData(r.Form, dest); err != nil {
		return err
	}

	// ربط الملفات
	if err := b.bindFiles(r.MultipartForm.File, dest); err != nil {
		return err
	}

	// التحقق من البيانات
	return b.Validate(dest)
}

// Validate يقوم بالتحقق من البيانات
func (b *Binder) Validate(data interface{}) error {
	return b.validator.Struct(data)
}

// bindFormData يربط بيانات النموذج
func (b *Binder) bindFormData(form map[string][]string, dest interface{}) error {
	v := reflect.ValueOf(dest)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// الحصول على اسم الحقل
		tag := fieldType.Tag.Get("form")
		if tag == "" {
			tag = fieldType.Tag.Get("json")
		}
		if tag == "" {
			tag = strings.ToLower(fieldType.Name)
		}

		// تخطي الحقول المستبعدة
		if tag == "-" || tag == "" {
			continue
		}

		// الحصول على القيمة
		values, ok := form[tag]
		if !ok || len(values) == 0 {
			continue
		}

		// تعيين القيمة
		if err := setFieldValue(field, values[0]); err != nil {
			return err
		}
	}

	return nil
}

// bindQueryData يربط معلمات الاستعلام
func (b *Binder) bindQueryData(query map[string][]string, dest interface{}) error {
	v := reflect.ValueOf(dest)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// الحصول على اسم الحقل
		tag := fieldType.Tag.Get("query")
		if tag == "" {
			tag = strings.ToLower(fieldType.Name)
		}

		// تخطي الحقول المستبعدة
		if tag == "-" || tag == "" {
			continue
		}

		// الحصول على القيمة
		values, ok := query[tag]
		if !ok || len(values) == 0 {
			continue
		}

		// تعيين القيمة
		if err := setFieldValue(field, values[0]); err != nil {
			return err
		}
	}

	return nil
}

// bindFiles يربط الملفات
func (b *Binder) bindFiles(files map[string][]*multipart.FileHeader, dest interface{}) error {
	v := reflect.ValueOf(dest)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// الحصول على اسم الحقل
		tag := fieldType.Tag.Get("form")
		if tag == "" {
			tag = fieldType.Tag.Get("file")
		}
		if tag == "" {
			continue
		}

		// تخطي الحقول المستبعدة
		if tag == "-" || tag == "" {
			continue
		}

		// الحصول على الملفات
		fileHeaders, ok := files[tag]
		if !ok || len(fileHeaders) == 0 {
			continue
		}

		// تعيين الملفات
		if err := setFileValue(field, fileHeaders); err != nil {
			return err
		}
	}

	return nil
}

// setFieldValue يضع قيمة في حقل
func setFieldValue(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)

	case reflect.Bool:
		v, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(v)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(v)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetUint(v)

	case reflect.Float32, reflect.Float64:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		field.SetFloat(v)

	case reflect.Struct:
		// معالجة الأنواع الخاصة مثل time.Time
		// سيتم تنفيذها لاحقاً

	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}

	return nil
}

// setFileValue يضع قيمة ملف في حقل
func setFileValue(field reflect.Value, fileHeaders []*multipart.FileHeader) error {
	// سيتم تنفيذها لاحقاً
	return nil
}
