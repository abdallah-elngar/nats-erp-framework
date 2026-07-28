package binding

import (
	"encoding/json"
	"io"
)

// JSONBinder يربط بيانات JSON
type JSONBinder struct {
	binder *Binder
}

// NewJSONBinder ينشئ رابط JSON جديد
func NewJSONBinder() *JSONBinder {
	return &JSONBinder{
		binder: NewBinder(),
	}
}

// Bind يربط بيانات JSON
func (jb *JSONBinder) Bind(r io.Reader, dest interface{}) error {
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(dest); err != nil {
		return err
	}
	return jb.binder.Validate(dest)
}

// BindBody يربط نص JSON
func (jb *JSONBinder) BindBody(body []byte, dest interface{}) error {
	if err := json.Unmarshal(body, dest); err != nil {
		return err
	}
	return jb.binder.Validate(dest)
}

// ToJSON يحول البيانات إلى JSON
func (jb *JSONBinder) ToJSON(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}
