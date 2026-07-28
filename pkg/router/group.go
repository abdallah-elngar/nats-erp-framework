package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Group يمثل مجموعة مسارات
type Group struct {
	chi.Router
	prefix string
	parent *Router
}

// NewGroup ينشئ مجموعة جديدة
func NewGroup(prefix string, parent *Router) *Group {
	return &Group{
		Router: chi.NewRouter(),
		prefix: prefix,
		parent: parent,
	}
}

// Get يضيف مسار GET
func (g *Group) Get(pattern string, handler http.HandlerFunc) {
	g.Router.Get(pattern, handler)
}

// Post يضيف مسار POST
func (g *Group) Post(pattern string, handler http.HandlerFunc) {
	g.Router.Post(pattern, handler)
}

// Put يضيف مسار PUT
func (g *Group) Put(pattern string, handler http.HandlerFunc) {
	g.Router.Put(pattern, handler)
}

// Delete يضيف مسار DELETE
func (g *Group) Delete(pattern string, handler http.HandlerFunc) {
	g.Router.Delete(pattern, handler)
}

// Patch يضيف مسار PATCH
func (g *Group) Patch(pattern string, handler http.HandlerFunc) {
	g.Router.Patch(pattern, handler)
}

// Use يضيف وسيطاً
func (g *Group) Use(middlewares ...func(http.Handler) http.Handler) {
	g.Router.Use(middlewares...)
}

// Resource يضيف مسارات CRUD
func (g *Group) Resource(pattern string, handler interface{}) {
	// سيتم تنفيذها لاحقاً
}
