package signals

import (
	"fmt"
	"sync"
	"time"
)

// Signal يمثل إشارة
type Signal struct {
	Name      string
	Data      interface{}
	Sender    string
	Timestamp time.Time
}

// NewSignal ينشئ إشارة جديدة
func NewSignal(name string, data interface{}, sender string) *Signal {
	return &Signal{
		Name:      name,
		Data:      data,
		Sender:    sender,
		Timestamp: time.Now(),
	}
}

// SignalHandler دالة معالجة الإشارة
type SignalHandler func(signal *Signal) error

// SignalManager يدير الإشارات
type SignalManager struct {
	handlers map[string][]SignalHandler
	mu       sync.RWMutex
	async    bool
}

// NewSignalManager ينشئ مدير إشارات جديد
func NewSignalManager(async bool) *SignalManager {
	return &SignalManager{
		handlers: make(map[string][]SignalHandler),
		async:    async,
	}
}

// Connect يربط معالجاً بإشارة
func (sm *SignalManager) Connect(name string, handler SignalHandler) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.handlers[name] = append(sm.handlers[name], handler)
}

// Disconnect يفصل معالجاً عن إشارة
func (sm *SignalManager) Disconnect(name string, handler SignalHandler) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	handlers, ok := sm.handlers[name]
	if !ok {
		return
	}

	for i, h := range handlers {
		if fmt.Sprintf("%p", h) == fmt.Sprintf("%p", handler) {
			sm.handlers[name] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}
}

// Emit يصدر إشارة
func (sm *SignalManager) Emit(signal *Signal) error {
	sm.mu.RLock()
	handlers, ok := sm.handlers[signal.Name]
	sm.mu.RUnlock()

	if !ok {
		return nil
	}

	if sm.async {
		// تنفيذ غير متزامن
		go sm.executeHandlers(handlers, signal)
		return nil
	}

	// تنفيذ متزامن
	for _, handler := range handlers {
		if err := handler(signal); err != nil {
			return err
		}
	}

	return nil
}

// EmitAsync يصدر إشارة بشكل غير متزامن
func (sm *SignalManager) EmitAsync(signal *Signal) {
	go sm.Emit(signal)
}

// executeHandlers ينفذ المعالجات
func (sm *SignalManager) executeHandlers(handlers []SignalHandler, signal *Signal) {
	for _, handler := range handlers {
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("Signal handler panic: %v\n", r)
				}
			}()
			handler(signal)
		}()
	}
}

// GetHandlers يعيد جميع المعالجات
func (sm *SignalManager) GetHandlers(name string) []SignalHandler {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	handlers, ok := sm.handlers[name]
	if !ok {
		return []SignalHandler{}
	}

	return handlers
}

// Clear يزيل جميع المعالجات
func (sm *SignalManager) Clear(name string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	delete(sm.handlers, name)
}
