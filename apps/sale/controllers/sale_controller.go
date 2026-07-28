package controllers

import (
    "net/http"
    "strconv"

    "github.com/go-chi/chi/v5"

    "github.com/nats-framework/nats/apps/sale/dto"
    "github.com/nats-framework/nats/apps/sale/services"
    "github.com/nats-framework/nats/pkg/response"
)

// SaleController متحكم Sale
type SaleController struct {
    service *services.SaleService
}

// NewSaleController ينشئ متحكم Sale جديد
func NewSaleController() *SaleController {
    return &SaleController{
        service: services.NewSaleService(),
    }
}

// Index يعيد قائمة sales
func (c *SaleController) Index(w http.ResponseWriter, r *http.Request) {
    items, err := c.service.GetAll()
    if err != nil {
        response.InternalError(w, err.Error())
        return
    }
    response.Success(w, items)
}

// Show يعيد Sale واحداً
func (c *SaleController) Show(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        response.BadRequest(w, "Invalid ID")
        return
    }

    item, err := c.service.GetByID(uint(id))
    if err != nil {
        response.NotFound(w, "Sale not found")
        return
    }

    response.Success(w, item)
}

// Create ينشئ Sale جديداً
func (c *SaleController) Create(w http.ResponseWriter, r *http.Request) {
    var req dto.CreateSaleRequest
    if err := response.BindJSON(r, &req); err != nil {
        response.BadRequest(w, err.Error())
        return
    }

    item, err := c.service.Create(req)
    if err != nil {
        response.BadRequest(w, err.Error())
        return
    }

    response.Created(w, item)
}

// Update يحدث Sale
func (c *SaleController) Update(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        response.BadRequest(w, "Invalid ID")
        return
    }

    var req dto.UpdateSaleRequest
    if err := response.BindJSON(r, &req); err != nil {
        response.BadRequest(w, err.Error())
        return
    }

    item, err := c.service.Update(uint(id), req)
    if err != nil {
        response.BadRequest(w, err.Error())
        return
    }

    response.Success(w, item)
}

// Delete يحذف Sale
func (c *SaleController) Delete(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        response.BadRequest(w, "Invalid ID")
        return
    }

    if err := c.service.Delete(uint(id)); err != nil {
        response.InternalError(w, err.Error())
        return
    }

    response.NoContent(w)
}
