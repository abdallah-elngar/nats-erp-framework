package template

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Loader محمل القوالب
type Loader struct {
	dir   string
	fs    *embed.FS // ✅ تغيير إلى مؤشر
	debug bool
}

// NewLoader ينشئ محملاً جديداً
func NewLoader(dir string, fs embed.FS, debug bool) *Loader {
	l := &Loader{
		dir:   dir,
		fs:    nil,
		debug: debug,
	}

	// ✅ التحقق من وجود FS
	if fs != (embed.FS{}) {
		l.fs = &fs
	}

	return l
}

// LoadAll يحمل جميع القوالب
func (l *Loader) LoadAll() (map[string]string, error) {
	templates := make(map[string]string)

	var err error
	if l.fs != nil {
		err = l.loadFromFS(*l.fs, ".", templates)
	} else {
		err = l.loadFromDir(l.dir, templates)
	}

	if err != nil {
		return nil, err
	}

	return templates, nil
}

// loadFromDir يحمل القوالب من المجلد
func (l *Loader) loadFromDir(dir string, templates map[string]string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			name := info.Name()
			if name == "layouts" || name == "partials" {
				return nil
			}
			return nil
		}

		if strings.HasSuffix(info.Name(), ".html") {
			relPath, err := filepath.Rel(l.dir, path)
			if err != nil {
				return err
			}

			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			templates[relPath] = string(content)
		}

		return nil
	})
}

// loadFromFS يحمل القوالب من نظام الملفات المضمن
func (l *Loader) loadFromFS(fs embed.FS, dir string, templates map[string]string) error {
	entries, err := fs.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())

		if entry.IsDir() {
			if entry.Name() == "layouts" || entry.Name() == "partials" {
				continue
			}
			if err := l.loadFromFS(fs, path, templates); err != nil {
				return err
			}
		} else if strings.HasSuffix(entry.Name(), ".html") {
			content, err := fs.ReadFile(path)
			if err != nil {
				return err
			}
			templates[path] = string(content)
		}
	}

	return nil
}

// LoadLayout يحمل التخطيط
func (l *Loader) LoadLayout(name string) (string, error) {
	path := filepath.Join(l.dir, "layouts", name)
	path = strings.TrimSuffix(path, ".html") + ".html"

	var content []byte
	var err error

	if l.fs != nil {
		content, err = l.fs.ReadFile(path)
	} else {
		content, err = os.ReadFile(path)
	}

	if err != nil {
		return "", fmt.Errorf("failed to load layout %s: %w", name, err)
	}

	return string(content), nil
}

// LoadPartial يحمل قالباً جزئياً
func (l *Loader) LoadPartial(name string) (string, error) {
	path := filepath.Join(l.dir, "partials", name)
	path = strings.TrimSuffix(path, ".html") + ".html"

	var content []byte
	var err error

	if l.fs != nil {
		content, err = l.fs.ReadFile(path)
	} else {
		content, err = os.ReadFile(path)
	}

	if err != nil {
		return "", fmt.Errorf("failed to load partial %s: %w", name, err)
	}

	return string(content), nil
}
