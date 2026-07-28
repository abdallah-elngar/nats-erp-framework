// {{.Model.Name}}Controller - CRUD Operations

// Index يعرض قائمة {{.Model.Name | lower}}s
func (c *{{.Model.Name}}Controller) Index(w http.ResponseWriter, r *http.Request) {
    items, err := c.service.GetAll()
    if err != nil {
        response.Error(w, http.StatusInternalServerError, err.Error())
        return
    }
    response.Success(w, items)
}

// Show يعرض {{.Model.Name}} واحداً
func (c *{{.Model.Name}}Controller) Show(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        response.Error(w, http.StatusBadRequest, "Invalid ID")
        return
    }

    item, err := c.service.GetByID(uint(id))
    if err != nil {
        response.Error(w, http.StatusNotFound, "{{.Model.Name}} not found")
        return
    }

    response.Success(w, item)
}

// Create ينشئ {{.Model.Name}} جديداً
func (c *{{.Model.Name}}Controller) Create(w http.ResponseWriter, r *http.Request) {
    var req Create{{.Model.Name}}Request
    if err := response.BindJSON(r, &req); err != nil {
        response.Error(w, http.StatusBadRequest, err.Error())
        return
    }

    item, err := c.service.Create(req)
    if err != nil {
        response.Error(w, http.StatusInternalServerError, err.Error())
        return
    }

    response.Created(w, item)
}

// Update يحدث {{.Model.Name}}
func (c *{{.Model.Name}}Controller) Update(w http.ResponseWriter, r *http.Request) {
    id, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        response.Error(w, http.StatusBadRequest, "Invalid ID")
        return
    }

    var req Update{{.Model.Name}}Request
    if err := response.BindJSON(r, &req); err != nil {
        response.Error(w, http.StatusBadRequest, err.Error())
        return
    }

    item, err := c.service.Update(uint(id), req)
    if err != nil {
        response.Error(w, http.StatusInternalServerError, err.Error())
        return
    }

    response.Success(w, item)
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

    response.NoContent(w)
}