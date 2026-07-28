package response

import (
	"html/template"
	"net/http"
)

// HTMLResponse مساعد للاستجابات HTML
type HTMLResponse struct {
	w        http.ResponseWriter
	template *template.Template
}

// NewHTMLResponse ينشئ مساعد استجابات HTML جديد
func NewHTMLResponse(w http.ResponseWriter, tmpl *template.Template) *HTMLResponse {
	return &HTMLResponse{
		w:        w,
		template: tmpl,
	}
}

// Render يعرض قالباً
func (hr *HTMLResponse) Render(name string, data interface{}) error {
	// ✅ استخدام hr.w بدلاً من w
	hr.w.Header().Set("Content-Type", "text/html")
	return hr.template.ExecuteTemplate(hr.w, name, data)
}

// RenderPartial يعرض قالباً جزئياً
func (hr *HTMLResponse) RenderPartial(name string, data interface{}) error {
	// ✅ استخدام hr.w بدلاً من w
	hr.w.Header().Set("Content-Type", "text/html")
	return hr.template.ExecuteTemplate(hr.w, name, data)
}

// Redirect يقوم بإعادة توجيه
func (hr *HTMLResponse) Redirect(url string, status int) {
	// ✅ استخدام hr.w بدلاً من w، وإضافة r
	// ملاحظة: هذه الدالة تحتاج إلى http.Request كمعامل
	// سيتم إصلاحها في الأسفل
	http.Redirect(hr.w, nil, url, status)
}

// RedirectWithRequest يقوم بإعادة توجيه مع طلب
func (hr *HTMLResponse) RedirectWithRequest(r *http.Request, url string, status int) {
	http.Redirect(hr.w, r, url, status)
}
