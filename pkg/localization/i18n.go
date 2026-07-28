package localization

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// I18n يدير الترجمة
type I18n struct {
	translations  map[string]map[string]string
	defaultLocale string
	mu            sync.RWMutex
}

// NewI18n ينشئ مدير ترجمة جديد
func NewI18n(defaultLocale string) *I18n {
	return &I18n{
		translations:  make(map[string]map[string]string),
		defaultLocale: defaultLocale,
	}
}

// LoadLocale يحمل ملف ترجمة
func (i *I18n) LoadLocale(locale, path string) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var translations map[string]string
	if err := json.Unmarshal(data, &translations); err != nil {
		return err
	}

	i.translations[locale] = translations
	return nil
}

// LoadDirectory يحمل جميع ملفات الترجمة من مجلد
func (i *I18n) LoadDirectory(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		locale := strings.TrimSuffix(file.Name(), ".json")
		path := filepath.Join(dir, file.Name())
		if err := i.LoadLocale(locale, path); err != nil {
			return err
		}
	}

	return nil
}

// Translate يترجم نصاً
func (i *I18n) Translate(locale, key string, params ...interface{}) string {
	i.mu.RLock()
	defer i.mu.RUnlock()

	// البحث عن الترجمة
	translations, ok := i.translations[locale]
	if !ok {
		translations, ok = i.translations[i.defaultLocale]
		if !ok {
			return key
		}
	}

	text, ok := translations[key]
	if !ok {
		return key
	}

	// تطبيق المعاملات
	if len(params) > 0 {
		return fmt.Sprintf(text, params...)
	}

	return text
}

// GetLocale يعيد الترجمة للغة محددة
func (i *I18n) GetLocale(locale string) map[string]string {
	i.mu.RLock()
	defer i.mu.RUnlock()

	translations, ok := i.translations[locale]
	if !ok {
		return i.translations[i.defaultLocale]
	}

	return translations
}

// AddTranslation يضيف ترجمة
func (i *I18n) AddTranslation(locale, key, value string) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if _, ok := i.translations[locale]; !ok {
		i.translations[locale] = make(map[string]string)
	}

	i.translations[locale][key] = value
}
