package binding

import (
	"net/url"
	"reflect"
)

// FormBinder يربط بيانات النموذج
type FormBinder struct {
	binder *Binder
}

// NewFormBinder ينشئ رابط نموذج جديد
func NewFormBinder() *FormBinder {
	return &FormBinder{
		binder: NewBinder(),
	}
}

// Bind يربط بيانات النموذج
func (fb *FormBinder) Bind(form url.Values, dest interface{}) error {
	if err := fb.binder.bindFormData(form, dest); err != nil {
		return err
	}
	return fb.binder.Validate(dest)
}

// BindValues يربط خريطة القيم
func (fb *FormBinder) BindValues(values map[string]string, dest interface{}) error {
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
			continue
		}

		// الحصول على القيمة
		value, ok := values[tag]
		if !ok {
			continue
		}

		// تعيين القيمة
		if err := setFieldValue(field, value); err != nil {
			return err
		}
	}

	return fb.binder.Validate(dest)
}
