package template

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

// Handler returns the HTTP handler for static files
func (sm *StaticManager) Handler() http.Handler {
	return http.StripPrefix(sm.config.Prefix, http.FileServer(http.Dir(sm.config.Dir)))
}

// Engine represents the template rendering engine
type Engine struct {
	templates map[string]*template.Template
	layout    string
	funcMap   template.FuncMap
	mu        sync.RWMutex
	dir       string
	fs        *embed.FS
	debug     bool

	filters       *FilterRegistry
	tags          *TagRegistry
	auth          *AuthHelper
	debugger      *Debugger
	cache         *Cache
	config        Config
	staticManager *StaticManager
	staticConfig  StaticConfig
}

// Config specifies the configuration for the Engine
type Config struct {
	Dir          string
	Layout       string
	Debug        bool
	FS           embed.FS
	FuncMap      template.FuncMap
	CacheEnabled bool
	CacheTTL     time.Duration
	AutoReload   bool
	AuthEnabled  bool
	AuthConfig   AuthConfig
	StaticConfig StaticConfig
}

// New creates a new template Engine instance
func New(config Config) *Engine {
	e := &Engine{
		templates: make(map[string]*template.Template),
		layout:    config.Layout,
		funcMap:   config.FuncMap,
		debug:     config.Debug,
		dir:       config.Dir,
		fs:        nil,
		config:    config,
	}

	if config.FS != (embed.FS{}) {
		e.fs = &config.FS
	}

	if e.funcMap == nil {
		e.funcMap = make(template.FuncMap)
	}

	e.filters = NewFilterRegistry()
	e.tags = NewTagRegistry()
	e.debugger = NewDebugger(config.Debug)
	e.cache = NewCache(config.CacheTTL)

	if config.AuthEnabled {
		e.auth = NewAuthHelper(config.AuthConfig)
	}

	e.addDefaultFuncs()
	e.registerDefaultFilters()
	e.registerDefaultTags()

	if config.StaticConfig.Dir != "" {
		e.staticConfig = config.StaticConfig
		e.staticManager = NewStaticManager(config.StaticConfig)
	}

	return e
}

// ServeStatic serves static files
func (e *Engine) ServeStatic(w http.ResponseWriter, r *http.Request) {
	if e.staticManager != nil {
		e.staticManager.Handler().ServeHTTP(w, r)
	}
}

// StaticHandler returns the static file handler
func (e *Engine) StaticHandler() http.Handler {
	if e.staticManager != nil {
		return e.staticManager.Handler()
	}
	return http.NotFoundHandler()
}

// StaticMiddleware returns a middleware for serving static assets
func (e *Engine) StaticMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if e.staticManager != nil && e.staticConfig.Prefix != "" && strings.HasPrefix(r.URL.Path, e.staticConfig.Prefix) {
				e.staticManager.Handler().ServeHTTP(w, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// addDefaultFuncs registers built-in helper functions
func (e *Engine) addDefaultFuncs() {
	e.funcMap["upper"] = strings.ToUpper
	e.funcMap["lower"] = strings.ToLower
	e.funcMap["title"] = strings.Title
	e.funcMap["join"] = strings.Join
	e.funcMap["add"] = func(a, b int) int { return a + b }
	e.funcMap["sub"] = func(a, b int) int { return a - b }
	e.funcMap["mul"] = func(a, b int) int { return a * b }
	e.funcMap["div"] = func(a, b int) int { return a / b }
	e.funcMap["safe"] = func(s string) template.HTML { return template.HTML(s) }
	e.funcMap["safeURL"] = func(s string) template.URL { return template.URL(s) }
	e.funcMap["time"] = func() string { return time.Now().Format("2006-01-02 15:04:05") }
	e.funcMap["formatTime"] = func(t time.Time, layout string) string { return t.Format(layout) }
	e.funcMap["static"] = e.staticFunc
	e.funcMap["dir"] = func(lang interface{}) string {
		if l, ok := lang.(string); ok && strings.HasPrefix(l, "ar") {
			return "rtl"
		}
		return "ltr"
	}

	// Template inheritance helpers
	e.funcMap["extends"] = e.extendsFunc
	e.funcMap["block"] = e.blockFunc
	e.funcMap["endblock"] = e.endblockFunc
	e.funcMap["define"] = e.extendsFunc
	e.funcMap["include"] = e.includeFunc
	e.funcMap["yield"] = e.yieldFunc
	e.funcMap["filter"] = e.filterFunc
	e.funcMap["if"] = e.ifFunc
	e.funcMap["else"] = e.elseFunc
	e.funcMap["elif"] = e.elifFunc
	e.funcMap["for"] = e.forFunc

	// Auth helpers
	e.funcMap["user"] = e.authUserFunc
	e.funcMap["is_authenticated"] = e.isAuthenticatedFunc
	e.funcMap["has_permission"] = e.hasPermissionFunc
	e.funcMap["csrf_token"] = e.csrfTokenFunc

	// Debugging helpers
	e.funcMap["debug"] = e.debugFunc
	e.funcMap["dump"] = e.dumpFunc

	// Utility helpers
	e.funcMap["url"] = e.urlFunc
	e.funcMap["now"] = e.nowFunc
	e.funcMap["format_date"] = e.formatDateFunc
	e.funcMap["truncate"] = e.truncateFunc
	e.funcMap["default"] = e.defaultFunc
	e.funcMap["escape"] = e.escapeFunc
}

func (e *Engine) getCacheKey(name string, data interface{}) string {
	return fmt.Sprintf("%s:%p", name, data)
}

func (e *Engine) extendsFunc(name string) string                          { return "" }
func (e *Engine) blockFunc(name string, data interface{}) template.HTML   { return template.HTML("") }
func (e *Engine) endblockFunc() string                                    { return "" }
func (e *Engine) includeFunc(name string, data interface{}) template.HTML { return template.HTML("") }
func (e *Engine) yieldFunc() template.HTML                                { return template.HTML("") }

func (e *Engine) filterFunc(value interface{}, name string, args ...interface{}) interface{} {
	fn, ok := e.filters.Get(name)
	if !ok {
		return value
	}
	return fn(value, args...)
}

func (e *Engine) ifFunc(condition bool, content interface{}) interface{} {
	if condition {
		return content
	}
	return ""
}
func (e *Engine) elseFunc() string { return "" }
func (e *Engine) elifFunc(condition bool, content interface{}) interface{} {
	if condition {
		return content
	}
	return ""
}
func (e *Engine) forFunc(items interface{}, content interface{}) interface{} {
	return content
}

func (e *Engine) authUserFunc(data interface{}) interface{} {
	if e.auth != nil {
		return e.auth.GetUser(data)
	}
	return nil
}

func (e *Engine) isAuthenticatedFunc(data interface{}) bool {
	if e.auth != nil {
		return e.auth.IsAuthenticated(data)
	}
	return false
}

func (e *Engine) hasPermissionFunc(permission string, data interface{}) bool {
	if e.auth != nil {
		return e.auth.HasPermission(data, permission)
	}
	return false
}

func (e *Engine) csrfTokenFunc() template.HTML {
	return template.HTML(`<input type="hidden" name="csrf_token" value="` + generateCSRFToken() + `">`)
}

func (e *Engine) debugFunc(value ...interface{}) template.HTML {
	if !e.debug {
		return template.HTML("")
	}
	var sb strings.Builder
	for _, v := range value {
		sb.WriteString(fmt.Sprintf("<!-- Debug: %+v -->\n", v))
	}
	return template.HTML(sb.String())
}

func (e *Engine) dumpFunc(value interface{}) template.HTML {
	if !e.debug {
		return template.HTML("")
	}
	return template.HTML(fmt.Sprintf("<pre>%+v</pre>", value))
}

func (e *Engine) urlFunc(name string, params ...interface{}) string {
	return "/" + name
}

func (e *Engine) staticFunc(path string) string {
	prefix := e.staticConfig.Prefix
	if prefix == "" {
		prefix = "/static/"
	}
	if !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}
	return prefix + strings.TrimPrefix(path, "/")
}

func (e *Engine) nowFunc() time.Time {
	return time.Now()
}

func (e *Engine) formatDateFunc(t time.Time, layout string) string {
	if layout == "" {
		layout = "2006-01-02"
	}
	return t.Format(layout)
}

func (e *Engine) truncateFunc(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

func (e *Engine) defaultFunc(defaultVal interface{}, val interface{}) interface{} {
	if val == nil {
		return defaultVal
	}
	if v, ok := val.(string); ok && v == "" {
		return defaultVal
	}
	return val
}

func (e *Engine) escapeFunc(s string) string {
	return template.HTMLEscapeString(s)
}

func (e *Engine) registerDefaultFilters() {
	e.filters.Register("upper", FilterUpper)
	e.filters.Register("lower", FilterLower)
	e.filters.Register("title", FilterTitle)
	e.filters.Register("capitalize", FilterCapitalize)
	e.filters.Register("slug", FilterSlug)
	e.filters.Register("truncate", FilterTruncate)
	e.filters.Register("word_count", FilterWordCount)
	e.filters.Register("line_breaks", FilterLineBreaks)
	e.filters.Register("strip_tags", FilterStripTags)
	e.filters.Register("escape", FilterEscape)
	e.filters.Register("safe", FilterSafe)

	e.filters.Register("add", FilterAdd)
	e.filters.Register("sub", FilterSub)
	e.filters.Register("mul", FilterMul)
	e.filters.Register("div", FilterDiv)
	e.filters.Register("format_number", FilterFormatNumber)
	e.filters.Register("currency", FilterCurrency)
	e.filters.Register("percentage", FilterPercentage)

	e.filters.Register("date", FilterDate)
	e.filters.Register("time", FilterTime)
	e.filters.Register("datetime", FilterDateTime)
	e.filters.Register("time_ago", FilterTimeAgo)
	e.filters.Register("format_datetime", FilterFormatDateTime)

	e.filters.Register("join", FilterJoin)
	e.filters.Register("split", FilterSplit)
	e.filters.Register("reverse", FilterReverse)
	e.filters.Register("first", FilterFirst)
	e.filters.Register("last", FilterLast)
	e.filters.Register("length", FilterLength)
	e.filters.Register("slice", FilterSlice)

	e.filters.Register("yesno", FilterYesNo)
	e.filters.Register("bool", FilterBool)
	e.filters.Register("default", FilterDefault)
}

func (e *Engine) registerDefaultTags() {
	e.tags.Register("if", TagIf)
	e.tags.Register("else", TagElse)
	e.tags.Register("elif", TagElif)
	e.tags.Register("for", TagFor)
	e.tags.Register("empty", TagEmpty)
	e.tags.Register("debug", TagDebug)
	e.tags.Register("dump", TagDump)
	e.tags.Register("auth", TagAuth)
	e.tags.Register("permission", TagPermission)
	e.tags.Register("csrf", TagCSRF)
}

func (e *Engine) getTemplate(name string) (*template.Template, error) {
	e.mu.RLock()
	tmpl, ok := e.templates[name]
	e.mu.RUnlock()

	if ok && !e.debug {
		return tmpl, nil
	}

	tmpl, err := e.loadTemplate(name)
	if err != nil {
		return nil, err
	}

	e.mu.Lock()
	e.templates[name] = tmpl
	e.mu.Unlock()

	return tmpl, nil
}

// loadTemplate يحمل القالب مع التخطيط عبر Go html/template القياسي
func (e *Engine) loadTemplate(name string) (*template.Template, error) {
	mainPath := filepath.Join(e.dir, name)
	if !strings.HasSuffix(mainPath, ".html") {
		mainPath += ".html"
	}

	// 1. قراءة محتوى القالب الابن
	var mainContent []byte
	var err error
	if e.fs != nil {
		mainContent, err = e.fs.ReadFile(mainPath)
	} else {
		mainContent, err = os.ReadFile(mainPath)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to load template %s: %w", name, err)
	}

	content := string(mainContent)

	// 2. إنشاء كائن القالب الرئيسي
	tmpl := template.New(name).Funcs(e.funcMap)

	// 3. إذا كان القالب يوراث من تخطيط أصل (Extends)
	if strings.Contains(content, "extends") {
		layoutName := e.extractLayoutName(content)
		if layoutName != "" {
			// مسح وسم {{ extends "..." }} من القالب الابن
			reExtends := regexp.MustCompile(`{{\s*extends\s+"[^"]+"\s*}}`)
			content = reExtends.ReplaceAllString(content, "")

			// قراءة التخطيط الأب
			layoutContent, err := e.loadLayout(layoutName)
			if err == nil && strings.TrimSpace(layoutContent) != "" {
				// نعمل Parse للتخطيط الأب أولاً
				_, err = tmpl.Parse(layoutContent)
				if err != nil {
					return nil, fmt.Errorf("failed to parse layout %s: %w", layoutName, err)
				}
			}
		}
	}

	// 4. نعمل Parse لمحتوى القالب الابن لتجاوز/إكمال الكتل المحددة في الأب
	tmpl, err = tmpl.Parse(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template content %s: %w", name, err)
	}

	return tmpl, nil
}

func (e *Engine) extractLayoutName(content string) string {
	re := regexp.MustCompile(`{{\s*extends\s+"([^"]+)"\s*}}`)
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// pkg/template/engine.go

// RenderFile يعرض قالباً من مسار محدد
func (e *Engine) RenderFile(w io.Writer, path string, data interface{}) error {
	// قراءة الملف مباشرة
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read template file %s: %w", path, err)
	}

	// إنشاء قالب من المحتوى
	tmpl, err := template.New(filepath.Base(path)).Funcs(e.funcMap).Parse(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// تنفيذ القالب
	return tmpl.Execute(w, data)
}

// inheritTemplate يدمج كتل القالب الابن داخل التخطيط الأب
func (e *Engine) inheritTemplate(content, layout string) string {
	// 1. استخراج الكتل من القالب الابن
	blocks := e.extractBlocks(content)

	result := layout
	// 2. استبدال الكتل المكتوبة في الأب بمحتوى الابن
	for name, blockContent := range blocks {
		// دعم {{ block "name" . }}...{{ end }} أو {{ define "name" }}...{{ end }}
		reBlock := regexp.MustCompile(`{{\s*(?:block|define)\s+"` + regexp.QuoteMeta(name) + `"\s*[^}]*}}[\s\S]*?{{\s*(?:end|endblock)\s*}}`)
		if reBlock.MatchString(result) {
			result = reBlock.ReplaceAllString(result, blockContent)
		}
	}

	return result
}

// extractBlocks يستخرج محتوى الكتل من القالب الابن
func (e *Engine) extractBlocks(content string) map[string]string {
	blocks := make(map[string]string)

	// البحث عن الكتل من نوع block أو define
	re := regexp.MustCompile(`{{\s*(?:block|define)\s+"([^"]+)"\s*[^}]*}}([\s\S]*?){{\s*(?:end|endblock)\s*}}`)
	matches := re.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) > 2 {
			name := strings.TrimSpace(match[1])
			blockContent := strings.TrimSpace(match[2])
			blocks[name] = blockContent
		}
	}

	return blocks
}

func (e *Engine) loadLayout(name string) (string, error) {
	// تجربة التحميل المباشر أولاً ثم التجربة داخل layouts/
	paths := []string{
		filepath.Join(e.dir, name+".html"),
		filepath.Join(e.dir, "layouts", name+".html"),
		filepath.Join(e.dir, name),
	}

	for _, path := range paths {
		var content []byte
		var err error

		if e.fs != nil {
			content, err = e.fs.ReadFile(path)
		} else {
			content, err = os.ReadFile(path)
		}

		if err == nil {
			return string(content), nil
		}
	}

	return "", fmt.Errorf("failed to load layout %s", name)
}

func (e *Engine) ClearCache() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.templates = make(map[string]*template.Template)
	if e.cache != nil {
		e.cache.Clear()
	}
}

func (e *Engine) Reload() error {
	e.ClearCache()
	return nil
}

func generateCSRFToken() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func (e *Engine) RenderWriter(w http.ResponseWriter, name string, data interface{}) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := e.Render(w, name, data)
	if err != nil {
		// 🛠️ بدلاً من الصفحة البيضاء، اطبع الخطأ مباشرة لمعرفته
		http.Error(w, fmt.Sprintf("Template Error: %v", err), http.StatusInternalServerError)
	}
	return err
}

func (e *Engine) Render(w io.Writer, name string, data interface{}) error {
	start := time.Now()
	defer func() {
		if e.debug {
			e.debugger.LogPerformance(name, time.Since(start))
		}
	}()

	// 1️⃣ التحقق من التخزين المؤقت
	if e.config.CacheEnabled {
		cacheKey := e.getCacheKey(name, data)
		if cached, found := e.cache.Get(cacheKey); found {
			if bytesData, ok := cached.([]byte); ok && len(bytesData) > 0 {
				if e.debug {
					e.debugger.Log("✅ Cache hit for: %s", name)
				}
				_, err := w.Write(bytesData)
				return err
			}
		}
	}

	// 2️⃣ إعادة التحميل التلقائي في وضع التطوير
	if e.config.AutoReload && e.debug {
		e.Reload()
	}

	// 3️⃣ جلب القالب
	tmpl, err := e.getTemplate(name)
	if err != nil {
		errStr := fmt.Sprintf("Template Load Error [%s]: %v", name, err)
		if e.debug {
			w.Write([]byte("<pre style='color:red;background:#fee;padding:15px;'>" + errStr + "</pre>"))
		}
		return fmt.Errorf("%s", errStr)
	}

	// 4️⃣ تجهيز بيانات المصادقة
	if e.auth != nil {
		data = e.auth.PrepareData(data)
	}

	// 5️⃣ تنفيذ القالب بشكل آمن
	var buf strings.Builder
	execErr := tmpl.Execute(&buf, data)

	// إذا فشل التنفيذ المباشر وكان هناك layout محدد، نحاول الـ Fallback
	if execErr != nil && e.layout != "" {
		buf.Reset() // 🧹 مسح أي نص جزئي كُتب قبل حدوث الخطأ
		execErr = tmpl.ExecuteTemplate(&buf, e.layout, data)
	}

	// إذا فشل التنفيذ في الحالتين
	if execErr != nil {
		errStr := fmt.Sprintf("Template Render Error [%s]: %v", name, execErr)
		if e.debug {
			// 🚀 بدلاً من الصفحة البيضاء، يتم طباعة الخطأ فوراً على المتصفح!
			w.Write([]byte("<pre style='color:red;background:#fee;padding:15px;border:1px solid red;'>" + errStr + "</pre>"))
		}
		return fmt.Errorf("%s", errStr)
	}

	result := buf.String()

	// التأكد من أن النتيجة ليست فارغة قبل وضعها في الكاش
	if strings.TrimSpace(result) == "" {
		errStr := fmt.Sprintf("Template [%s] rendered an empty string!", name)
		if e.debug {
			w.Write([]byte("<pre style='color:orange;background:#fff3cd;padding:15px;'>" + errStr + "</pre>"))
		}
		return fmt.Errorf("%s", errStr)
	}

	// 6️⃣ حفظ النتيجة في الكاش
	if e.config.CacheEnabled {
		e.cache.Set(e.getCacheKey(name, data), []byte(result))
	}

	_, err = w.Write([]byte(result))
	return err
}
