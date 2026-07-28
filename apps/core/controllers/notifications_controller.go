package controllers

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

// NotificationsController متحكم الإشعارات
type NotificationsController struct{}

// NewNotificationsController ينشئ متحكم إشعارات جديد
func NewNotificationsController() *NotificationsController {
    return &NotificationsController{}
}

// Stream يبث الإشعارات عبر Server-Sent Events (SSE)
func (c *NotificationsController) Stream(w http.ResponseWriter, r *http.Request) {
    // ✅ إعدادات SSE
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")
    w.Header().Set("Access-Control-Allow-Origin", "*")

    flusher, ok := w.(http.Flusher)
    if !ok {
        http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
        return
    }

    // إرسال إشعار ترحيب
    c.sendEvent(w, "connected", map[string]string{
        "message": "Connected to notifications stream",
        "time":    time.Now().Format("15:04:05"),
    })
    flusher.Flush()

    // محاكاة إرسال إشعارات دورية
    ticker := time.NewTicker(30 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-r.Context().Done():
            return
        case <-ticker.C:
            // إرسال إشعار اختبار
            c.sendEvent(w, "notification", map[string]string{
                "message": "System is running smoothly",
                "time":    time.Now().Format("15:04:05"),
                "type":    "info",
            })
            flusher.Flush()
        }
    }
}

// sendEvent يرسل حدث SSE
func (c *NotificationsController) sendEvent(w http.ResponseWriter, event string, data interface{}) {
    jsonData, _ := json.Marshal(data)
    fmt.Fprintf(w, "event: %s\n", event)
    fmt.Fprintf(w, "data: %s\n\n", jsonData)
}

// GetNotifications يعيد الإشعارات السابقة (API)
func (c *NotificationsController) GetNotifications(w http.ResponseWriter, r *http.Request) {
    notifications := []map[string]interface{}{
        {
            "id":        1,
            "message":   "Welcome to NATS ERP!",
            "type":      "success",
            "timestamp": time.Now().Add(-5 * time.Minute).Format(time.RFC3339),
            "read":      false,
        },
        {
            "id":        2,
            "message":   "System initialized successfully",
            "type":      "info",
            "timestamp": time.Now().Add(-10 * time.Minute).Format(time.RFC3339),
            "read":      true,
        },
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "data":    notifications,
    })
}

// MarkAsRead يحدد إشعار كمقروء
func (c *NotificationsController) MarkAsRead(w http.ResponseWriter, r *http.Request) {
    // TODO: تنفيذ تحديد الإشعار كمقروء
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "success": true,
        "message": "Notification marked as read",
    })
}