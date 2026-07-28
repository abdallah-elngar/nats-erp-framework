package pipes

import (
	"context"
	"encoding/json"
)

// JSONPipe أنبوب لتحويل JSON
type JSONPipe struct{}

// Process ينفذ الأنبوب
func (j *JSONPipe) Process(ctx context.Context, data interface{}) (interface{}, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// UnmarshalPipe أنبوب لفك JSON
type UnmarshalPipe struct {
	target interface{}
}

// NewUnmarshalPipe ينشئ أنبوب فك JSON جديد
func NewUnmarshalPipe(target interface{}) *UnmarshalPipe {
	return &UnmarshalPipe{target: target}
}

// Process ينفذ الأنبوب
func (u *UnmarshalPipe) Process(ctx context.Context, data interface{}) (interface{}, error) {
	bytes, ok := data.([]byte)
	if !ok {
		return nil, nil
	}

	if err := json.Unmarshal(bytes, u.target); err != nil {
		return nil, err
	}

	return u.target, nil
}

// ValidatePipe أنبوب للتحقق
type ValidatePipe struct {
	validator func(data interface{}) error
}

// NewValidatePipe ينشئ أنبوب تحقق جديد
func NewValidatePipe(validator func(data interface{}) error) *ValidatePipe {
	return &ValidatePipe{validator: validator}
}

// Process ينفذ الأنبوب
func (v *ValidatePipe) Process(ctx context.Context, data interface{}) (interface{}, error) {
	if err := v.validator(data); err != nil {
		return nil, err
	}
	return data, nil
}

// TransformPipe أنبوب للتحويل
type TransformPipe struct {
	transform func(data interface{}) (interface{}, error)
}

// NewTransformPipe ينشئ أنبوب تحويل جديد
func NewTransformPipe(transform func(data interface{}) (interface{}, error)) *TransformPipe {
	return &TransformPipe{transform: transform}
}

// Process ينفذ الأنبوب
func (t *TransformPipe) Process(ctx context.Context, data interface{}) (interface{}, error) {
	return t.transform(data)
}

// FilterPipe أنبوب للتصفية
type FilterPipe struct {
	filter func(data interface{}) bool
}

// NewFilterPipe ينشئ أنبوب تصفية جديد
func NewFilterPipe(filter func(data interface{}) bool) *FilterPipe {
	return &FilterPipe{filter: filter}
}

// Process ينفذ الأنبوب
func (f *FilterPipe) Process(ctx context.Context, data interface{}) (interface{}, error) {
	if !f.filter(data) {
		return nil, nil
	}
	return data, nil
}
