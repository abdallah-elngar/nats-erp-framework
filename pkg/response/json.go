package response

import (
	"encoding/json"
	"net/http"
)

// JSONResponse مساعد للاستجابات JSON
type JSONResponse struct {
	w http.ResponseWriter
}

// NewJSONResponse ينشئ مساعد استجابات JSON جديد
func NewJSONResponse(w http.ResponseWriter) *JSONResponse {
	return &JSONResponse{w: w}
}

// Send يرسل استجابة JSON
func (jr *JSONResponse) Send(status int, data interface{}) {
	JSON(jr.w, status, data)
}

// OK يرسل استجابة 200
func (jr *JSONResponse) OK(data interface{}) {
	Success(jr.w, data)
}

// Created يرسل استجابة 201
func (jr *JSONResponse) Created(data interface{}) {
	JSON(jr.w, http.StatusCreated, Response{
		Success: true,
		Data:    data,
	})
}

// NoContent يرسل استجابة 204
func (jr *JSONResponse) NoContent() {
	jr.w.WriteHeader(http.StatusNoContent)
}

// BadRequest يرسل خطأ 400
func (jr *JSONResponse) BadRequest(message string) {
	BadRequest(jr.w, message)
}

// NotFound يرسل خطأ 404
func (jr *JSONResponse) NotFound(message string) {
	NotFound(jr.w, message)
}

// InternalError يرسل خطأ 500
func (jr *JSONResponse) InternalError(message string) {
	InternalError(jr.w, message)
}

// Stream يرسل تدفق JSON
func (jr *JSONResponse) Stream(data interface{}) error {
	jr.w.Header().Set("Content-Type", "application/json")
	jr.w.Header().Set("Transfer-Encoding", "chunked")

	encoder := json.NewEncoder(jr.w)
	return encoder.Encode(data)
}
