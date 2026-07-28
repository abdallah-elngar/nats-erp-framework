package localization

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// Loader يحمل ملفات الترجمة
type Loader struct {
	path string
}

// NewLoader ينشئ محملاً جديداً
func NewLoader(path string) *Loader {
	return &Loader{path: path}
}

// Load يحمل ملف ترجمة
func (l *Loader) Load(locale string) (map[string]string, error) {
	filePath := filepath.Join(l.path, locale+".json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var translations map[string]string
	if err := json.Unmarshal(data, &translations); err != nil {
		return nil, err
	}

	return translations, nil
}

// LoadAll يحمل جميع ملفات الترجمة
func (l *Loader) LoadAll() (map[string]map[string]string, error) {
	result := make(map[string]map[string]string)

	files, err := os.ReadDir(l.path)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		locale := strings.TrimSuffix(file.Name(), ".json")
		translations, err := l.Load(locale)
		if err != nil {
			return nil, err
		}

		result[locale] = translations
	}

	return result, nil
}

// Save يحفظ ملف ترجمة
func (l *Loader) Save(locale string, translations map[string]string) error {
	filePath := filepath.Join(l.path, locale+".json")
	data, err := json.MarshalIndent(translations, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, 0644)
}
