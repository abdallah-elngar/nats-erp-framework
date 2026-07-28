// pkg/template/static_fs.go
package template

import (
	"io/fs"
	"os"
	"path/filepath"
)

// FileSystem نظام الملفات المتقدم
type FileSystem struct {
	dir  string
	fs   fs.FS
	root string
}

// NewFileSystem ينشئ نظام ملفات جديد
func NewFileSystem(dir string, embedFS fs.FS) *FileSystem {
	return &FileSystem{
		dir:  dir,
		fs:   embedFS,
		root: "/",
	}
}

// Open يفتح ملفاً
func (fsys *FileSystem) Open(name string) (fs.File, error) {
	if fsys.fs != nil {
		return fsys.fs.Open(name)
	}
	return os.Open(filepath.Join(fsys.dir, name))
}

// ReadFile يقرأ ملفاً
func (fsys *FileSystem) ReadFile(name string) ([]byte, error) {
	if fsys.fs != nil {
		return fs.ReadFile(fsys.fs, name)
	}
	return os.ReadFile(filepath.Join(fsys.dir, name))
}

// ReadDir يقرأ مجلداً
func (fsys *FileSystem) ReadDir(name string) ([]fs.DirEntry, error) {
	if fsys.fs != nil {
		return fs.ReadDir(fsys.fs, name)
	}
	return os.ReadDir(filepath.Join(fsys.dir, name))
}

// Stat يعيد معلومات الملف
func (fsys *FileSystem) Stat(name string) (fs.FileInfo, error) {
	if fsys.fs != nil {
		f, err := fsys.fs.Open(name)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		return f.Stat()
	}
	return os.Stat(filepath.Join(fsys.dir, name))
}

// Glob يعيد الملفات المطابقة للنمط
func (fsys *FileSystem) Glob(pattern string) ([]string, error) {
	if fsys.fs != nil {
		return fs.Glob(fsys.fs, pattern)
	}
	return filepath.Glob(filepath.Join(fsys.dir, pattern))
}

// Walk يسير في المجلد
func (fsys *FileSystem) Walk(root string, fn fs.WalkDirFunc) error {
	if fsys.fs != nil {
		return fs.WalkDir(fsys.fs, root, fn)
	}
	return filepath.WalkDir(filepath.Join(fsys.dir, root), func(path string, d os.DirEntry, err error) error {
		relPath, _ := filepath.Rel(fsys.dir, path)
		return fn(relPath, d, err)
	})
}

// Exists يتحقق من وجود ملف
func (fsys *FileSystem) Exists(name string) bool {
	_, err := fsys.Stat(name)
	return err == nil
}

// IsDir يتحقق مما إذا كان مجلداً
func (fsys *FileSystem) IsDir(name string) bool {
	info, err := fsys.Stat(name)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// ListFiles يعيد قائمة الملفات
func (fsys *FileSystem) ListFiles(path string) ([]string, error) {
	var files []string

	err := fsys.Walk(path, func(p string, d fs.DirEntry, err error) error {
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

// ListDirs يعيد قائمة المجلدات
func (fsys *FileSystem) ListDirs(path string) ([]string, error) {
	var dirs []string

	err := fsys.Walk(path, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() && p != path {
			dirs = append(dirs, p)
		}
		return nil
	})

	return dirs, err
}
