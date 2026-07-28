package template

import (
	"embed"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

// StaticManager مدير الملفات الثابتة
type StaticManager struct {
	dir         string
	fs          *embed.FS
	prefix      string
	cache       *Cache
	debug       bool
	mu          sync.RWMutex
	fileModTime map[string]time.Time
	etags       map[string]string
	config      StaticConfig
}

// StaticConfig إعدادات الملفات الثابتة
type StaticConfig struct {
	Dir          string        // مجلد الملفات الثابتة
	FS           *embed.FS     // نظام الملفات المضمن (اختياري)
	Prefix       string        // بادئة URL (مثل /static/)
	CacheEnabled bool          // تفعيل التخزين المؤقت
	CacheTTL     time.Duration // مدة التخزين المؤقت
	Debug        bool          // وضع التصحيح
	Compress     bool          // ضغط الملفات
	Minify       bool          // تصغير الملفات
}

// NewStaticManager ينشئ مدير ملفات ثابتة جديد
func NewStaticManager(config StaticConfig) *StaticManager {
	m := &StaticManager{
		dir:         config.Dir,
		fs:          config.FS,
		prefix:      config.Prefix,
		debug:       config.Debug,
		config:      config,
		fileModTime: make(map[string]time.Time),
		etags:       make(map[string]string),
	}

	if config.Prefix == "" {
		m.prefix = "/static/"
	}

	m.cache = NewCache(config.CacheTTL)

	return m
}

// ServeFile يخدم ملفاً ثابتاً
func (m *StaticManager) ServeFile(w http.ResponseWriter, r *http.Request, path string) {
	// التحقق من التخزين المؤقت
	if m.cache != nil && m.config.CacheEnabled {
		cacheKey := "static:" + path
		if cached, found := m.cache.Get(cacheKey); found {
			if m.debug {
				fmt.Printf("🐛 Cache hit for: %s\n", path)
			}
			w.Header().Set("Cache-Control", "public, max-age=31536000")
			w.Write(cached.([]byte))
			return
		}
	}

	// قراءة الملف
	content, modTime, err := m.readFile(path)
	if err != nil {
		if m.debug {
			fmt.Printf("🐛 Failed to read file %s: %v\n", path, err)
		}
		http.NotFound(w, r)
		return
	}

	// تعيين نوع المحتوى
	contentType := m.getContentType(path)
	w.Header().Set("Content-Type", contentType)

	// تعيين رؤوس التخزين المؤقت
	if m.cache != nil && m.config.CacheEnabled {
		w.Header().Set("Cache-Control", "public, max-age=31536000")
		w.Header().Set("Last-Modified", modTime.Format(http.TimeFormat))

		// ETag
		etag := m.generateETag(path, modTime)
		w.Header().Set("ETag", etag)

		// التحقق من If-None-Match
		if match := r.Header.Get("If-None-Match"); match == etag {
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}

	// معالجة الملف (ضغط، تصغير)
	content = m.processFile(path, content)

	// تخزين في الكاش
	if m.cache != nil && m.config.CacheEnabled {
		m.cache.Set("static:"+path, content)
	}

	// إرسال الملف
	w.Write(content)
}

// readFile يقرأ ملفاً
func (m *StaticManager) readFile(path string) ([]byte, time.Time, error) {
	var content []byte
	var err error
	var modTime time.Time

	if m.fs != nil {
		// قراءة من النظام المضمن
		content, err = m.fs.ReadFile(path)
		if err != nil {
			return nil, time.Time{}, err
		}

		// الحصول على معلومات الملف
		file, err := m.fs.Open(path)
		if err != nil {
			return content, time.Now(), nil
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			return content, time.Now(), nil
		}
		modTime = stat.ModTime()

		return content, modTime, nil
	}

	// قراءة من نظام الملفات المحلي
	fullPath := filepath.Join(m.dir, path)
	content, err = os.ReadFile(fullPath)
	if err != nil {
		return nil, time.Time{}, err
	}

	info, err := os.Stat(fullPath)
	if err != nil {
		return content, time.Now(), nil
	}

	return content, info.ModTime(), nil
}

// getContentType يعيد نوع المحتوى
func (m *StaticManager) getContentType(path string) string {
	ext := filepath.Ext(path)
	switch ext {
	case ".html":
		return "text/html; charset=utf-8"
	case ".css":
		return "text/css; charset=utf-8"
	case ".js":
		return "application/javascript; charset=utf-8"
	case ".json":
		return "application/json; charset=utf-8"
	case ".xml":
		return "application/xml; charset=utf-8"
	case ".pdf":
		return "application/pdf"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".svg":
		return "image/svg+xml"
	case ".ico":
		return "image/x-icon"
	case ".woff":
		return "font/woff"
	case ".woff2":
		return "font/woff2"
	case ".ttf":
		return "font/ttf"
	case ".eot":
		return "application/vnd.ms-fontobject"
	default:
		return mime.TypeByExtension(ext)
	}
}

// generateETag يولد ETag للملف
func (m *StaticManager) generateETag(path string, modTime time.Time) string {
	key := path + ":" + modTime.String()
	if etag, ok := m.etags[key]; ok {
		return etag
	}

	etag := fmt.Sprintf(`"%x"`, modTime.UnixNano())
	m.etags[key] = etag
	return etag
}

// processFile يعالج الملف (ضغط، تصغير)
func (m *StaticManager) processFile(path string, content []byte) []byte {
	if !m.config.Compress && !m.config.Minify {
		return content
	}

	ext := filepath.Ext(path)

	// تصغير الملفات
	if m.config.Minify {
		switch ext {
		case ".css":
			content = MinifyCSS(content)
		case ".js":
			content = MinifyJS(content)
		case ".html":
			content = MinifyHTML(content)
		}
	}

	// ضغط الملفات (يمكن إضافة Gzip لاحقاً)
	if m.config.Compress {
		// يمكن إضافة ضغط Gzip هنا
	}

	return content
}

// ============================================
// دوال التصغير (Minification)
// ============================================

// MinifyCSS يصغر ملف CSS
func MinifyCSS(content []byte) []byte {
	str := string(content)

	// إزالة التعليقات
	str = regexp.MustCompile(`/\*.*?\*/`).ReplaceAllString(str, "")

	// إزالة المسافات الزائدة
	str = regexp.MustCompile(`\s+`).ReplaceAllString(str, " ")

	// إزالة المسافات بين الأقواس
	str = regexp.MustCompile(`\s*{\s*`).ReplaceAllString(str, "{")
	str = regexp.MustCompile(`\s*}\s*`).ReplaceAllString(str, "}")
	str = regexp.MustCompile(`\s*;\s*`).ReplaceAllString(str, ";")
	str = regexp.MustCompile(`\s*:\s*`).ReplaceAllString(str, ":")
	str = regexp.MustCompile(`\s*,\s*`).ReplaceAllString(str, ",")

	// إزالة المسافات قبل وبعد
	str = strings.TrimSpace(str)

	return []byte(str)
}

// MinifyJS يصغر ملف JavaScript
func MinifyJS(content []byte) []byte {
	str := string(content)

	// إزالة التعليقات سطر واحد
	str = regexp.MustCompile(`//.*?$`).ReplaceAllString(str, "")

	// إزالة التعليقات متعددة الأسطر
	str = regexp.MustCompile(`/\*.*?\*/`).ReplaceAllString(str, "")

	// إزالة المسافات الزائدة
	str = regexp.MustCompile(`\s+`).ReplaceAllString(str, " ")

	// إزالة المسافات بين الرموز
	str = regexp.MustCompile(`\s*([{}();:,])\s*`).ReplaceAllString(str, "$1")

	// إزالة المسافات قبل وبعد العوامل
	str = regexp.MustCompile(`\s*([=!<>]=?)\s*`).ReplaceAllString(str, "$1")

	// إزالة المسافات الزائدة
	str = strings.TrimSpace(str)

	return []byte(str)
}

// MinifyHTML يصغر ملف HTML
func MinifyHTML(content []byte) []byte {
	str := string(content)

	// إزالة التعليقات
	str = regexp.MustCompile(`<!--.*?-->`).ReplaceAllString(str, "")

	// إزالة المسافات الزائدة بين الوسوم
	str = regexp.MustCompile(`>\s+<`).ReplaceAllString(str, "><")

	// إزالة المسافات الزائدة
	str = regexp.MustCompile(`\s+`).ReplaceAllString(str, " ")

	// إزالة المسافات قبل وبعد
	str = strings.TrimSpace(str)

	return []byte(str)
}

// ============================================
// BundleManager - تجميع الملفات
// ============================================

// BundleConfig إعدادات التجميع
type BundleConfig struct {
	CSSFiles  []string
	JSFiles   []string
	OutputDir string
	Minify    bool
	Compress  bool
}

// BundleManager مدير التجميع
type BundleManager struct {
	config  BundleConfig
	manager *StaticManager
}

// NewBundleManager ينشئ مدير تجميع جديد
func NewBundleManager(config BundleConfig, manager *StaticManager) *BundleManager {
	return &BundleManager{
		config:  config,
		manager: manager,
	}
}

// BundleCSS يجمع ملفات CSS
func (bm *BundleManager) BundleCSS() ([]byte, error) {
	var result strings.Builder

	for _, file := range bm.config.CSSFiles {
		content, _, err := bm.manager.readFile(file)
		if err != nil {
			continue
		}
		result.Write(content)
		result.WriteString("\n")
	}

	content := []byte(result.String())

	if bm.config.Minify {
		content = MinifyCSS(content)
	}

	return content, nil
}

// BundleJS يجمع ملفات JavaScript
func (bm *BundleManager) BundleJS() ([]byte, error) {
	var result strings.Builder

	for _, file := range bm.config.JSFiles {
		content, _, err := bm.manager.readFile(file)
		if err != nil {
			continue
		}
		result.Write(content)
		result.WriteString("\n")
	}

	content := []byte(result.String())

	if bm.config.Minify {
		content = MinifyJS(content)
	}

	return content, nil
}

// pkg/template/static.go - أضف هذه الدوال

// GetStaticPath يعيد مسار الملف الثابت
func (m *StaticManager) GetStaticPath(path string) string {
	return m.prefix + path
}

// StaticFunc دالة مساعدة للقوالب
func (m *StaticManager) StaticFunc(path string) string {
	return m.GetStaticPath(path)
}
