package routes

import (
    "github.com/go-chi/chi/v5"

    "github.com/nats-framework/nats/apps/sale/controllers"
)

// RegisterRoutes يسجل مسارات التطبيق
func RegisterRoutes(r *chi.Mux) {
    // مسارات Sale
    r.Get("/sales", controllers.NewSaleController().Index)
    r.Get("/sales/{id}", controllers.NewSaleController().Show)
    r.Post("/sales", controllers.NewSaleController().Create)
    r.Put("/sales/{id}", controllers.NewSaleController().Update)
    r.Delete("/sales/{id}", controllers.NewSaleController().Delete)
}
