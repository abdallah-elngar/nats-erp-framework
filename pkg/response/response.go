package response

import (
	"encoding/json"
	"net/http"
)

// Response يمثل استجابة موحدة
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Code    int         `json:"code,omitempty"`
}

// JSON يرسل استجابة JSON
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Success يرسل استجابة نجاح
func Success(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// SuccessMessage يرسل استجابة نجاح مع رسالة
func SuccessMessage(w http.ResponseWriter, message string, data interface{}) {
	JSON(w, http.StatusOK, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error يرسل استجابة خطأ
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, Response{
		Success: false,
		Error:   message,
		Code:    status,
	})
}

// BadRequest يرسل خطأ 400
func BadRequest(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, message)
}

// Unauthorized يرسل خطأ 401
func Unauthorized(w http.ResponseWriter, message string) {
	Error(w, http.StatusUnauthorized, message)
}

// Forbidden يرسل خطأ 403
func Forbidden(w http.ResponseWriter, message string) {
	Error(w, http.StatusForbidden, message)
}

// NotFound يرسل خطأ 404
func NotFound(w http.ResponseWriter, message string) {
	Error(w, http.StatusNotFound, message)
}

// InternalError يرسل خطأ 500
func InternalError(w http.ResponseWriter, message string) {
	Error(w, http.StatusInternalServerError, message)
}

// ValidationError يرسل خطأ تحقق
func ValidationError(w http.ResponseWriter, errors map[string]string) {
	JSON(w, http.StatusUnprocessableEntity, Response{
		Success: false,
		Error:   "Validation failed",
		Data:    errors,
		Code:    http.StatusUnprocessableEntity,
	})
}

// Paginated يرسل استجابة مرقمة
func Paginated(w http.ResponseWriter, data interface{}, total int64, page, pageSize int) {
	JSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    data,
		"meta": map[string]interface{}{
			"total":       total,
			"page":        page,
			"page_size":   pageSize,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// ============================================
// ✅ دوال إضافية للـ Controllers
// ============================================

// BindJSON يربط بيانات JSON
func BindJSON(r *http.Request, dest interface{}) error {
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(dest)
}

// Created يرسل استجابة 201
func Created(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusCreated, Response{
		Success: true,
		Data:    data,
	})
}

// NoContent يرسل استجابة 204
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// BindForm يربط بيانات النموذج
func BindForm(r *http.Request, dest interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	// ربط بسيط (يمكن توسيعه)
	return nil
}

// BindQuery يربط معلمات الاستعلام
func BindQuery(r *http.Request, dest interface{}) error {
	// ربط بسيط (يمكن توسيعه)
	return nil
}
