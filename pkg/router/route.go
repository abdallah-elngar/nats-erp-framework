package router

import (
	"net/http"
	"strings"
)

// Route يمثل مساراً
type Route struct {
	Method      string
	Pattern     string
	Handler     http.HandlerFunc
	Name        string
	Middlewares []func(http.Handler) http.Handler
}

// NewRoute ينشئ مساراً جديداً
func NewRoute(method, pattern string, handler http.HandlerFunc) *Route {
	return &Route{
		Method:      method,
		Pattern:     pattern,
		Handler:     handler,
		Middlewares: make([]func(http.Handler) http.Handler, 0),
	}
}

// WithName يضيف اسماً للمسار
func (r *Route) WithName(name string) *Route {
	r.Name = name
	return r
}

// WithMiddleware يضيف وسيطاً للمسار
func (r *Route) WithMiddleware(middlewares ...func(http.Handler) http.Handler) *Route {
	r.Middlewares = append(r.Middlewares, middlewares...)
	return r
}

// GetPath يعيد مسار URL
func (r *Route) GetPath(params ...string) string {
	path := r.Pattern
	for i, param := range params {
		path = strings.Replace(path, "{"+string(rune(i))+"}", param, 1)
	}
	return path
}
