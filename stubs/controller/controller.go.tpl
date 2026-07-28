package controllers

import (
    "net/http"
    "strconv"

    "github.com/nats-framework/nats/pkg/response"
    "github.com/nats-framework/nats/apps/{{.AppName}}/models"
    "github.com/nats-framework/nats/apps/{{.AppName}}/services"
)

// {{.Model.Name}}Controller متحكم {{.Model.Name}}
type {{.Model.Name}}Controller struct {
    service *services.{{.Model.Name}}Service
}

// New{{.Model.Name}}Controller ينشئ متحكماً جديداً
func New{{.Model.Name}}Controller() *{{.Model.Name}}Controller {
    return &{{.Model.Name}}Controller{
        service: services.New{{.Model.Name}}Service(),
    }
}

// Index يعرض قائمة {{.Model.Name | lower}}s
func (c *{{.Model.Name}}Controller) Index(w http.ResponseWriter, r *http.Request) {
    {{.Model.Name | lower}}s, err := c.service.GetAll()
    if err != nil {
        response.Error(w, http.StatusInternalServerError, err.Error())
        return
    }

    response.JSON(w, http.StatusOK, {{.Model.Name | lower}}s)
}

// Show يعرض {{.Model.Name}} واحداً
func (c *{{.Model.Name}}Controller) Show(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        response.Error(w, http.StatusBadRequest, "Invalid ID")
        return
    }

    {{.Model.Name | lower}}, err := c.service.GetByID(uint(id))
    if err != nil {
        response.Error(w, http.StatusNotFound, "{{.Model.Name}} not found")
        return
    }

    response.JSON(w, http.StatusOK, {{.Model.Name | lower}})
}

// Create ينشئ {{.Model.Name}} جديداً
func (c *{{.Model.Name}}Controller) Create(w http.ResponseWriter, r *http.Request) {
    var {{.Model.Name | lower}} models.{{.Model.Name}}
    if err := response.BindJSON(r, &{{.Model.Name | lower}}); err != nil {
        response.Error(w, http.StatusBadRequest, err.Error())
        return
    }

    if err := c.service.Create(&{{.Model.Name | lower}}); err != nil {
        response.Error(w, http.StatusInternalServerError, err.Error())
        return
    }

    response.JSON(w, http.StatusCreated, {{.Model.Name | lower}})
}

// Update يحدث {{.Model.Name}}
func (c *{{.Model.Name}}Controller) Update(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        response.Error(w, http.StatusBadRequest, "Invalid ID")
        return
    }

    var {{.Model.Name | lower}} models.{{.Model.Name}}
    if err := response.BindJSON(r, &{{.Model.Name | lower}}); err != nil {
        response.Error(w, http.StatusBadRequest, err.Error())
        return
    }

    {{.Model.Name | lower}}.ID = uint(id)
    if err := c.service.Update(&{{.Model.Name | lower}}); err != nil {
        response.Error(w, http.StatusInternalServerError, err.Error())
        return
    }

    response.JSON(w, http.StatusOK, {{.Model.Name | lower}})
}

// Delete يحذف {{.Model.Name}}
func (c *{{.Model.Name}}Controller) Delete(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        response.Error(w, http.StatusBadRequest, "Invalid ID")
        return
    }

    if err := c.service.Delete(uint(id)); err != nil {
        response.Error(w, http.StatusInternalServerError, err.Error())
        return
    }

    response.JSON(w, http.StatusNoContent, nil)
}