package hooks

import (
	"fmt"
	"sync"
)

// HookType نوع الخطاف
type HookType string

const (
	HookBefore HookType = "before"
	HookAfter  HookType = "after"
	HookAround HookType = "around"
)

// Hook يمثل خطافاً
type Hook struct {
	Name     string
	Type     HookType
	Priority int
	Handler  HookHandler
}

// HookHandler دالة معالجة الخطاف
type HookHandler func(data interface{}) (interface{}, error)

// HookManager يدير الخطافات
type HookManager struct {
	hooks map[string][]Hook
	mu    sync.RWMutex
}

// NewHookManager ينشئ مدير خطافات جديد
func NewHookManager() *HookManager {
	return &HookManager{
		hooks: make(map[string][]Hook),
	}
}

// Register يسجل خطافاً
func (hm *HookManager) Register(name string, hookType HookType, handler HookHandler, priority int) {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	hook := Hook{
		Name:     name,
		Type:     hookType,
		Priority: priority,
		Handler:  handler,
	}

	hm.hooks[name] = append(hm.hooks[name], hook)

	// ترتيب حسب الأولوية
	hooks := hm.hooks[name]
	for i := 0; i < len(hooks); i++ {
		for j := i + 1; j < len(hooks); j++ {
			if hooks[i].Priority > hooks[j].Priority {
				hooks[i], hooks[j] = hooks[j], hooks[i]
			}
		}
	}
}

// Execute ينفذ الخطافات
func (hm *HookManager) Execute(name string, hookType HookType, data interface{}) (interface{}, error) {
	hm.mu.RLock()
	hooks, ok := hm.hooks[name]
	hm.mu.RUnlock()

	if !ok {
		return data, nil
	}

	var err error
	result := data

	for _, hook := range hooks {
		if hook.Type != hookType && hookType != HookAround {
			continue
		}

		if hookType == HookAround {
			// الخطافات التي حول التنفيذ
			result, err = hook.Handler(result)
		} else {
			// الخطافات قبل أو بعد
			result, err = hook.Handler(result)
		}

		if err != nil {
			return result, err
		}
	}

	return result, nil
}

// RegisterBefore يسجل خطافاً قبل التنفيذ
func (hm *HookManager) RegisterBefore(name string, handler HookHandler, priority int) {
	hm.Register(name, HookBefore, handler, priority)
}

// RegisterAfter يسجل خطافاً بعد التنفيذ
func (hm *HookManager) RegisterAfter(name string, handler HookHandler, priority int) {
	hm.Register(name, HookAfter, handler, priority)
}

// RegisterAround يسجل خطافاً حول التنفيذ
func (hm *HookManager) RegisterAround(name string, handler HookHandler, priority int) {
	hm.Register(name, HookAround, handler, priority)
}

// Remove يزيل خطافاً
func (hm *HookManager) Remove(name string, handler HookHandler) {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	hooks, ok := hm.hooks[name]
	if !ok {
		return
	}

	for i, hook := range hooks {
		if fmt.Sprintf("%p", hook.Handler) == fmt.Sprintf("%p", handler) {
			hm.hooks[name] = append(hooks[:i], hooks[i+1:]...)
			break
		}
	}
}

// Clear يزيل جميع الخطافات
func (hm *HookManager) Clear(name string) {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	delete(hm.hooks, name)
}

// GetHooks يعيد جميع الخطافات
func (hm *HookManager) GetHooks(name string) []Hook {
	hm.mu.RLock()
	defer hm.mu.RUnlock()

	hooks, ok := hm.hooks[name]
	if !ok {
		return []Hook{}
	}

	return hooks
}
