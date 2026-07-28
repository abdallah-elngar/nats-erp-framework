package template

import (
	"fmt"
	"strings"
	"sync"
)

type TagFunc func(ctx *TagContext) (string, error)

type TagContext struct {
	Content string
	Args    []string
	Data    interface{}
	Engine  *Engine
}

type TagRegistry struct {
	tags map[string]TagFunc
	mu   sync.RWMutex
}

func NewTagRegistry() *TagRegistry {
	return &TagRegistry{
		tags: make(map[string]TagFunc),
	}
}

func (r *TagRegistry) Register(name string, fn TagFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tags[name] = fn
}

func (r *TagRegistry) Get(name string) (TagFunc, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	fn, ok := r.tags[name]
	return fn, ok
}

// ============================================
// وسوم أساسية
// ============================================

func TagIf(ctx *TagContext) (string, error) {
	if len(ctx.Args) == 0 {
		return "", fmt.Errorf("if tag requires a condition")
	}

	condition := ctx.Args[0]
	if evaluateCondition(condition, ctx.Data) {
		return ctx.Content, nil
	}
	return "", nil
}

func TagElse(ctx *TagContext) (string, error) {
	return ctx.Content, nil
}

func TagElif(ctx *TagContext) (string, error) {
	if len(ctx.Args) == 0 {
		return "", fmt.Errorf("elif tag requires a condition")
	}

	condition := ctx.Args[0]
	if evaluateCondition(condition, ctx.Data) {
		return ctx.Content, nil
	}
	return "", nil
}

func TagFor(ctx *TagContext) (string, error) {
	// سيتم تنفيذها لاحقاً
	return ctx.Content, nil
}

func TagEmpty(ctx *TagContext) (string, error) {
	return ctx.Content, nil
}

func TagDebug(ctx *TagContext) (string, error) {
	return fmt.Sprintf("<!-- Debug: %+v -->", ctx.Data), nil
}

func TagDump(ctx *TagContext) (string, error) {
	return fmt.Sprintf("<pre>%+v</pre>", ctx.Data), nil
}

func TagAuth(ctx *TagContext) (string, error) {
	if user, ok := ctx.Data.(map[string]interface{})["user"]; ok && user != nil {
		return ctx.Content, nil
	}
	return "", nil
}

func TagPermission(ctx *TagContext) (string, error) {
	if len(ctx.Args) == 0 {
		return "", fmt.Errorf("permission tag requires a permission name")
	}

	permission := ctx.Args[0]
	if hasPermission(ctx.Data, permission) {
		return ctx.Content, nil
	}
	return "", nil
}

func TagCSRF(ctx *TagContext) (string, error) {
	token := generateCSRFToken()
	return fmt.Sprintf(`<input type="hidden" name="csrf_token" value="%s">`, token), nil
}

// دوال مساعدة
func evaluateCondition(condition string, data interface{}) bool {
	parts := strings.Split(condition, " ")

	if len(parts) == 1 {
		val, _ := getValue(parts[0], data)
		return toBool(val)
	}

	if len(parts) == 3 {
		left, _ := getValue(parts[0], data)
		op := parts[1]
		right := parts[2]

		switch op {
		case "==":
			return fmt.Sprintf("%v", left) == right
		case "!=":
			return fmt.Sprintf("%v", left) != right
		}
	}

	return false
}

func getValue(path string, data interface{}) (interface{}, error) {
	parts := strings.Split(path, ".")
	current := data
	for _, part := range parts {
		if m, ok := current.(map[string]interface{}); ok {
			current = m[part]
		} else {
			return nil, fmt.Errorf("failed to get value: %s", path)
		}
	}
	return current, nil
}

func hasPermission(data interface{}, permission string) bool {
	if user, ok := data.(map[string]interface{})["user"]; ok {
		if perms, ok := user.(map[string]interface{})["permissions"]; ok {
			if list, ok := perms.([]string); ok {
				for _, p := range list {
					if p == permission {
						return true
					}
				}
			}
		}
	}
	return false
}
