package routes

import (
    "github.com/go-chi/chi/v5"
)

// RegisterRoutes يسجل مسارات التطبيق
func RegisterRoutes(r *chi.Mux) {
    // مسارات {{.AppTitle}}
    r.Route("/{{.AppName}}", func(r chi.Router) {
        // {{range .Models}}
        // مسارات {{.Name}}
        r.Get("/{{.Name | lower}}", {{.Name}}Controller.Index)
        r.Get("/{{.Name | lower}}/{id}", {{.Name}}Controller.Show)
        r.Post("/{{.Name | lower}}", {{.Name}}Controller.Create)
        r.Put("/{{.Name | lower}}/{id}", {{.Name}}Controller.Update)
        r.Delete("/{{.Name | lower}}/{id}", {{.Name}}Controller.Delete)
        // {{end}}
    })
}