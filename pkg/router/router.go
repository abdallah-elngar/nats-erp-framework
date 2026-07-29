package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/nats-framework/nats/pkg/template"
)

// Router يمثل الموجه
type Router struct {
	chi.Router
	engine *template.Engine
	config RouterConfig
}

// RouterConfig إعدادات الموجه
type RouterConfig struct {
	Debug        bool
	StaticPrefix string
	StaticDir    string
}

// NewRouter ينشئ موجه جديد
func NewRouter(engine *template.Engine, config RouterConfig) *Router {
	r := chi.NewRouter()

	// إضافة ميدلوير عامة
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60))

	router := &Router{
		Router: r,
		engine: engine,
		config: config,
	}

	// إضافة المسارات الأساسية
	router.registerBaseRoutes()

	return router
}

// registerBaseRoutes يسجل المسارات الأساسية
func (r *Router) registerBaseRoutes() {
	// خدمة الملفات الثابتة
	if r.engine != nil {
		r.Router.Handle(r.config.StaticPrefix+"*", r.engine.StaticHandler())
	}

	// صفحة الصحة
	r.Router.Get("/health", r.healthHandler)

	// الصفحة الرئيسية (إعادة توجيه)
	r.Router.Get("/", func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, "/admin/developer", http.StatusFound)
	})
}

// healthHandler يعالج طلب الصحة
func (r *Router) healthHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","message":"NATS Framework is running"}`))
}

// ✅ ChiRouter يعيد الـ chi.Router الأصلي (للتوافق مع الكود القديم)
func (r *Router) ChiRouter() *chi.Mux {
	if mux, ok := r.Router.(*chi.Mux); ok {
		return mux
	}
	return nil
}

// Route ينشئ مساراً فرعياً
func (r *Router) Route(pattern string, fn func(r *Router)) *Router {
	subRouter := &Router{
		Router: chi.NewRouter(),
		engine: r.engine,
		config: r.config,
	}

	fn(subRouter)
	r.Router.Mount(pattern, subRouter.Router)
	return subRouter
}

// Group ينشئ مجموعة مسارات
func (r *Router) Group(fn func(r *Router)) *Router {
	group := &Router{
		Router: chi.NewRouter(),
		engine: r.engine,
		config: r.config,
	}
	fn(group)
	r.Router.Mount("/", group.Router)
	return group
}

// With يضيف ميدلوير للمجموعة
func (r *Router) With(middlewares ...func(http.Handler) http.Handler) *Router {
	return &Router{
		Router: r.Router.With(middlewares...),
		engine: r.engine,
		config: r.config,
	}
}

// Get يضيف مسار GET
func (r *Router) Get(pattern string, handler http.HandlerFunc) {
	r.Router.Get(pattern, handler)
}

// Post يضيف مسار POST
func (r *Router) Post(pattern string, handler http.HandlerFunc) {
	r.Router.Post(pattern, handler)
}

// Put يضيف مسار PUT
func (r *Router) Put(pattern string, handler http.HandlerFunc) {
	r.Router.Put(pattern, handler)
}

// Delete يضيف مسار DELETE
func (r *Router) Delete(pattern string, handler http.HandlerFunc) {
	r.Router.Delete(pattern, handler)
}

// Patch يضيف مسار PATCH
func (r *Router) Patch(pattern string, handler http.HandlerFunc) {
	r.Router.Patch(pattern, handler)
}

// Options يضيف مسار OPTIONS
func (r *Router) Options(pattern string, handler http.HandlerFunc) {
	r.Router.Options(pattern, handler)
}

// Head يضيف مسار HEAD
func (r *Router) Head(pattern string, handler http.HandlerFunc) {
	r.Router.Head(pattern, handler)
}

// Use يضيف ميدلوير
func (r *Router) Use(middlewares ...func(http.Handler) http.Handler) {
	r.Router.Use(middlewares...)
}

// Mount يثبت مسارات فرعية
func (r *Router) Mount(pattern string, handler http.Handler) {
	r.Router.Mount(pattern, handler)
}

// Render يعرض قالباً
func (r *Router) Render(w http.ResponseWriter, name string, data interface{}) error {
	if r.engine == nil {
		return nil
	}
	return r.engine.RenderWriter(w, name, data)
}

// StaticPath يعيد مسار الملف الثابت
func (r *Router) StaticPath(path string) string {
	return r.config.StaticPrefix + path
}

// GetEngine يعيد محرك القوالب
func (r *Router) GetEngine() *template.Engine {
	return r.engine
}
