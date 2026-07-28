// pkg/template/static_bundle.go
package template

import (
	"bytes"
	"compress/gzip"
)

// compress يضغط المحتوى
func (bm *BundleManager) compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)

	_, err := gz.Write(data)
	if err != nil {
		return nil, err
	}

	err = gz.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// WriteToFile يكتب المحتوى إلى ملف
func (bm *BundleManager) WriteToFile(filename string, data []byte) error {
	// تنفيذ الكتابة إلى ملف
	return nil
}
