package dashboard

import (
	"time"
)

// WidgetManager يدير أدوات لوحة التحكم
type WidgetManager struct {
	dashboard *Dashboard
}

// NewWidgetManager ينشئ مدير أدوات جديد
func NewWidgetManager(dashboard *Dashboard) *WidgetManager {
	return &WidgetManager{
		dashboard: dashboard,
	}
}

// RegisterStatsWidget يسجل أداة الإحصائيات
func (wm *WidgetManager) RegisterStatsWidget() {
	widget := NewWidget("stats", "System Stats", "stats")
	widget.Data["users"] = wm.getUserCount()
	widget.Data["apps"] = wm.getAppCount()
	widget.Data["records"] = wm.getRecordCount()
	widget.Data["storage"] = wm.getStorageUsage()
	widget.Data["uptime"] = wm.getUptime()

	wm.dashboard.AddWidget(*widget)
}

// RegisterActivityWidget يسجل أداة النشاط
func (wm *WidgetManager) RegisterActivityWidget() {
	widget := NewWidget("activity", "Recent Activity", "activity")
	widget.Data["activities"] = wm.getRecentActivity()

	wm.dashboard.AddWidget(*widget)
}

// RegisterChartsWidget يسجل أداة الرسوم البيانية
func (wm *WidgetManager) RegisterChartsWidget() {
	widget := NewWidget("charts", "Analytics", "charts")
	widget.Data["chart_data"] = wm.getChartData()

	wm.dashboard.AddWidget(*widget)
}

// getUserCount يعيد عدد المستخدمين
func (wm *WidgetManager) getUserCount() int {
	// استعلام من قاعدة البيانات
	return 0
}

// getAppCount يعيد عدد التطبيقات
func (wm *WidgetManager) getAppCount() int {
	// استعلام من قاعدة البيانات
	return 0
}

// getRecordCount يعيد عدد السجلات
func (wm *WidgetManager) getRecordCount() int {
	// استعلام من قاعدة البيانات
	return 0
}

// getStorageUsage يعيد استخدام التخزين
func (wm *WidgetManager) getStorageUsage() string {
	return "0 MB"
}

// getUptime يعيد وقت التشغيل
func (wm *WidgetManager) getUptime() string {
	return "0h 0m"
}

// getRecentActivity يعيد النشاطات الأخيرة
func (wm *WidgetManager) getRecentActivity() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"action": "System started",
			"time":   time.Now().Format("2006-01-02 15:04:05"),
		},
	}
}

// getChartData يعيد بيانات الرسوم البيانية
func (wm *WidgetManager) getChartData() map[string]interface{} {
	return map[string]interface{}{
		"labels": []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"},
		"values": []int{10, 20, 15, 30, 25, 40},
	}
}
