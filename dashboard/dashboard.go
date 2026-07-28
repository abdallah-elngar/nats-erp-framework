package dashboard

import (
	"net/http"

	"github.com/nats-framework/nats/pkg/response"
)

// Dashboard يمثل لوحة التحكم
type Dashboard struct {
	widgets []Widget
}

// NewDashboard ينشئ لوحة تحكم جديدة
func NewDashboard() *Dashboard {
	return &Dashboard{
		widgets: make([]Widget, 0),
	}
}

// AddWidget يضيف أداة إلى لوحة التحكم
func (d *Dashboard) AddWidget(widget Widget) {
	d.widgets = append(d.widgets, widget)
}

// Render يعرض لوحة التحكم
func (d *Dashboard) Render(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"widgets": d.widgets,
		"title":   "Dashboard",
	}

	response.Success(w, data)
}

// GetWidgets يعيد جميع الأدوات
func (d *Dashboard) GetWidgets() []Widget {
	return d.widgets
}

// Widget تمثل أداة في لوحة التحكم
type Widget struct {
	ID       string                 `json:"id"`
	Title    string                 `json:"title"`
	Type     string                 `json:"type"`
	Position int                    `json:"position"`
	Data     map[string]interface{} `json:"data"`
}

// NewWidget ينشئ أداة جديدة
func NewWidget(id, title, widgetType string) *Widget {
	return &Widget{
		ID:    id,
		Title: title,
		Type:  widgetType,
		Data:  make(map[string]interface{}),
	}
}
