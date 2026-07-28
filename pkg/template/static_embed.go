// pkg/template/static_embed.go
package template

import (
	"embed"
	"io/fs"
)

// StaticEmbed يضيف دعم تضمين الملفات الثابتة
type StaticEmbed struct {
	fs      embed.FS
	prefix  string
	manager *StaticManager
}

// NewStaticEmbed ينشئ كائن تضمين جديد
func NewStaticEmbed(fs embed.FS, prefix string) *StaticEmbed {
	return &StaticEmbed{
		fs:     fs,
		prefix: prefix,
	}
}

// WithManager يضيف مدير الملفات الثابتة
func (se *StaticEmbed) WithManager(manager *StaticManager) *StaticEmbed {
	se.manager = manager
	return se
}

// FS يعيد نظام الملفات
func (se *StaticEmbed) FS() fs.FS {
	return se.fs
}

// ReadFile يقرأ ملفاً
func (se *StaticEmbed) ReadFile(path string) ([]byte, error) {
	return se.fs.ReadFile(path)
}

// ReadDir يقرأ مجلداً
func (se *StaticEmbed) ReadDir(path string) ([]fs.DirEntry, error) {
	return se.fs.ReadDir(path)
}

// Walk يسير في المجلد
func (se *StaticEmbed) Walk(root string, fn fs.WalkDirFunc) error {
	return fs.WalkDir(se.fs, root, fn)
}

// Exists يتحقق من وجود ملف
func (se *StaticEmbed) Exists(path string) bool {
	_, err := se.fs.ReadFile(path)
	return err == nil
}

// ListFiles يعيد قائمة الملفات
func (se *StaticEmbed) ListFiles(path string) ([]string, error) {
	var files []string

	err := se.Walk(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			files = append(files, p)
		}
		return nil
	})

	return files, err
}
