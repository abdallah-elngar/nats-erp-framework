package template

import (
	"bytes"
	"testing"
	"time"
)

func TestEngine(t *testing.T) {
	config := Config{
		Dir:          "testdata/templates",
		Layout:       "main",
		Debug:        true,
		CacheEnabled: false,
		AuthEnabled:  true,
		AuthConfig: AuthConfig{
			LoginURL:  "/login",
			LogoutURL: "/logout",
			UserKey:   "user",
		},
	}
	
	engine := New(config)
	
	// اختبار تحميل القالب
	tmpl, err := engine.getTemplate("pages/index")
	if err != nil {
		t.Fatalf("Failed to load template: %v", err)
	}
	
	if tmpl == nil {
		t.Fatal("Template is nil")
	}
}

func TestRender(t *testing.T) {
	config := Config{
		Dir:          "testdata/templates",
		Layout:       "main",
		Debug:        true,
		CacheEnabled: false,
	}
	
	engine := New(config)
	
	data := map[string]interface{}{
		"Title": "Test Page",
		"User": map[string]interface{}{
			"Username": "admin",
			"Email":    "admin@example.com",
		},
		"Items": []string{"Item 1", "Item 2", "Item 3"},
	}
	
	var buf bytes.Buffer
	err := engine.Render(&buf, "pages/index", data)
	if err != nil {
		t.Fatalf("Failed to render: %v", err)
	}
	
	output := buf.String()
	t.Logf("Output: %s", output)
	
	// التحقق من وجود المحتوى
	if !bytes.Contains([]byte(output), []byte("Test Page")) {
		t.Error("Output does not contain expected title")
	}
}

func TestFilters(t *testing.T) {
	engine := New(Config{Debug: true})
	
	tests := []struct {
		name     string
		filter   string
		input    interface{}
		args     []interface{}
		expected interface{}
	}{
		{"Upper", "upper", "hello", nil, "HELLO"},
		{"Lower", "lower", "HELLO", nil, "hello"},
		{"Title", "title", "hello world", nil, "Hello World"},
		{"Add", "add", 10, []interface{}{5}, 15},
		{"Sub", "sub", 10, []interface{}{3}, 7},
		{"Mul", "mul", 10, []interface{}{2}, 20},
		{"Div", "div", 10, []interface{}{2}, 5},
		{"Truncate", "truncate", "Hello World", []interface{}{5}, "Hello..."},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn, ok := engine.filters.Get(tt.filter)
			if !ok {
				t.Skipf("Filter %s not found", tt.filter)
			}
			
			result := fn(tt.input, tt.args...)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestCache(t *testing.T) {
	cache := NewCache(1 * time.Minute)
	
	// اختبار التخزين
	cache.Set("key1", "value1")
	
	// اختبار الاسترجاع
	val, ok := cache.Get("key1")
	if !ok {
		t.Error("Cache item not found")
	}
	if val != "value1" {
		t.Errorf("Expected value1, got %v", val)
	}
	
	// اختبار انتهاء المدة
	cache.ttl = 0
	_, ok = cache.Get("key1")
	if ok {
		t.Error("Cache item should have expired")
	}
	
	// اختبار المسح
	cache.Clear()
	_, ok = cache.Get("key1")
	if ok {
		t.Error("Cache item should have been cleared")
	}
}

func TestAuthHelper(t *testing.T) {
	config := AuthConfig{
		Enabled:    true,
		LoginURL:   "/login",
		LogoutURL:  "/logout",
		UserKey:    "user",
	}
	
	auth := NewAuthHelper(config)
	
	// اختبار البيانات
	data := map[string]interface{}{
		"user": map[string]interface{}{
			"id":       1,
			"username": "admin",
		},
		"permissions": []string{"users.create", "users.edit"},
	}
	
	// اختبار PrepareData
	prepared := auth.PrepareData(data)
	if m, ok := prepared.(map[string]interface{}); ok {
		if _, exists := m["auth"]; !exists {
			t.Error("Auth data not added to context")
		}
	}
	
	// اختبار IsAuthenticated
	if !auth.IsAuthenticated(data) {
		t.Error("User should be authenticated")
	}
	
	// اختبار GetUser
	user := auth.GetUser(data)
	if user == nil {
		t.Error("User should not be nil")
	}
	
	// اختبار HasPermission
	if !auth.HasPermission(data, "users.create") {
		t.Error("User should have permission users.create")
	}
	if auth.HasPermission(data, "users.delete") {
		t.Error("User should not have permission users.delete")
	}
}

func TestInheritance(t *testing.T) {
	config := Config{
		Dir:    "testdata/templates",
		Layout: "main",
		Debug:  true,
	}
	
	engine := New(config)
	im := NewInheritanceManager(engine)
	
	// اختبار استخراج اسم التخطيط
	content := `{{ extends "main" }}
{{ block "title" . }}Test Page{{ endblock }}
{{ block "content" . }}<p>Hello World</p>{{ endblock }}`
	
	layoutName := im.extractLayoutName(content)
	if layoutName != "main" {
		t.Errorf("Expected main, got %s", layoutName)
	}
	
	// اختبار استخراج الكتل
	blocks := im.extractBlocks(content)
	if len(blocks) != 2 {
		t.Errorf("Expected 2 blocks, got %d", len(blocks))
	}
	
	if blocks["title"] != "Test Page" {
		t.Errorf("Expected Test Page, got %s", blocks["title"])
	}
	
	if blocks["content"] != "<p>Hello World</p>" {
		t.Errorf("Expected <p>Hello World</p>, got %s", blocks["content"])
	}
}