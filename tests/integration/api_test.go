package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nats-framework/nats/pkg/engine"
)

// TestAPI اختبار واجهات برمجة التطبيقات
func TestAPI(t *testing.T) {
	// إنشاء محرك التطبيق
	app := engine.New()

	// تحميل الإعدادات
	err := app.LoadConfig()
	require.NoError(t, err)

	// تهيئة قاعدة البيانات
	err = app.InitDatabase()
	require.NoError(t, err)

	// تحميل التطبيقات
	err = app.LoadApps()
	require.NoError(t, err)

	// تشغيل الخادم
	err = app.Run()
	require.NoError(t, err)

	// إنشاء طلب اختبار
	req, err := http.NewRequest("GET", "/admin/stats", nil)
	require.NoError(t, err)

	// تنفيذ الطلب
	w := httptest.NewRecorder()
	app.Server.Handler.ServeHTTP(w, req)

	// التحقق من الاستجابة
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.True(t, response["success"].(bool))
}

// TestAuthAPI اختبار واجهات مصادقة
func TestAuthAPI(t *testing.T) {
	// سيتم تنفيذها لاحقاً
}
