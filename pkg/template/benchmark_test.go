package template

import (
	"bytes"
	"testing"
	"time"
)

func BenchmarkEngine_Render(b *testing.B) {
	config := Config{
		Dir:          "testdata/templates",
		Layout:       "main",
		Debug:        false,
		CacheEnabled: true,
	}

	engine := New(config)

	data := map[string]interface{}{
		"Title": "Benchmark Page",
		"User": map[string]interface{}{
			"Username": "admin",
			"Email":    "admin@example.com",
		},
		"Items": []string{"Item 1", "Item 2", "Item 3"},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		engine.Render(&buf, "pages/index", data)
	}
}

func BenchmarkFilters(b *testing.B) {
	engine := New(Config{Debug: true})

	filters := []string{"upper", "lower", "title", "truncate", "add", "mul"}
	inputs := []interface{}{"hello world", 10, 20, "This is a long text"}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, name := range filters {
			if fn, ok := engine.filters.Get(name); ok {
				for _, input := range inputs {
					fn(input)
				}
			}
		}
	}
}

func BenchmarkCache(b *testing.B) {
	cache := NewCache(1 * time.Minute)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := "key" + string(rune(i))
		cache.Set(key, "value"+string(rune(i)))
		cache.Get(key)
	}
}
