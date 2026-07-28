package utils

import (
	"io"
	"os"
	"path/filepath"
)

// FileExists يتحقق من وجود ملف
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// CreateFile يخلق ملفاً
func CreateFile(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	return file.Close()
}

// ReadFile يقرأ ملفاً
func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// WriteFile يكتب ملفاً
func WriteFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// CopyFile ينسخ ملفاً
func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// DeleteFile يحذف ملفاً
func DeleteFile(path string) error {
	return os.Remove(path)
}

// FileSize يعيد حجم الملف
func FileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}
