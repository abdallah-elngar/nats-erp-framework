package events

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Event يمثل حدثاً
type Event interface {
	GetName() string
	GetData() interface{}
	GetTimestamp() time.Time
}

// BaseEvent حدث أساسي
type BaseEvent struct {
	Name      string
	Data      interface{}
	Timestamp time.Time
}

// GetName يعيد اسم الحدث
func (e *BaseEvent) GetName() string {
	return e.Name
}

// GetData يعيد بيانات الحدث
func (e *BaseEvent) GetData() interface{} {
	return e.Data
}

// GetTimestamp يعيد وقت الحدث
func (e *BaseEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

// NewEvent ينشئ حدثاً جديداً
func NewEvent(name string, data interface{}) *BaseEvent {
	return &BaseEvent{
		Name:      name,
		Data:      data,
		Timestamp: time.Now(),
	}
}

// ListenerFunc دالة مستمع للحدث
type ListenerFunc func(event Event) error

// EventManager يدير الأحداث
type EventManager struct {
	listeners map[string][]ListenerFunc
	mu        sync.RWMutex
	async     bool
}

// NewEventManager ينشئ مدير أحداث جديد
func NewEventManager(async bool) *EventManager {
	return &EventManager{
		listeners: make(map[string][]ListenerFunc),
		async:     async,
	}
}

// Listen يسجل مستمعاً للحدث
func (em *EventManager) Listen(name string, listener ListenerFunc) {
	em.mu.Lock()
	defer em.mu.Unlock()

	em.listeners[name] = append(em.listeners[name], listener)
}

// Emit يصدر حدثاً
func (em *EventManager) Emit(event Event) error {
	name := event.GetName()

	em.mu.RLock()
	listeners, ok := em.listeners[name]
	em.mu.RUnlock()

	if !ok {
		return nil
	}

	if em.async {
		// تنفيذ غير متزامن
		go em.executeListeners(listeners, event)
		return nil
	}

	// تنفيذ متزامن
	for _, listener := range listeners {
		if err := listener(event); err != nil {
			return err
		}
	}

	return nil
}

// EmitAsync يصدر حدثاً بشكل غير متزامن
func (em *EventManager) EmitAsync(event Event) {
	go em.Emit(event)
}

// EmitWithContext يصدر حدثاً مع سياق
func (em *EventManager) EmitWithContext(ctx context.Context, event Event) error {
	// سيتم تنفيذها لاحقاً
	return em.Emit(event)
}

// executeListeners ينفذ المستمعين
func (em *EventManager) executeListeners(listeners []ListenerFunc, event Event) {
	for _, listener := range listeners {
		// استعادة من الأخطاء
		func() {
			defer func() {
				if r := recover(); r != nil {
					// تسجيل الخطأ
					fmt.Printf("Event listener panic: %v\n", r)
				}
			}()
			listener(event)
		}()
	}
}

// RemoveListener يزيل مستمعاً
func (em *EventManager) RemoveListener(name string, listener ListenerFunc) {
	em.mu.Lock()
	defer em.mu.Unlock()

	listeners, ok := em.listeners[name]
	if !ok {
		return
	}

	// البحث عن المستمع وإزالته
	for i, l := range listeners {
		if fmt.Sprintf("%p", l) == fmt.Sprintf("%p", listener) {
			em.listeners[name] = append(listeners[:i], listeners[i+1:]...)
			break
		}
	}
}

// ClearListeners يزيل جميع المستمعين
func (em *EventManager) ClearListeners(name string) {
	em.mu.Lock()
	defer em.mu.Unlock()

	delete(em.listeners, name)
}

// GetListeners يعيد جميع المستمعين
func (em *EventManager) GetListeners(name string) []ListenerFunc {
	em.mu.RLock()
	defer em.mu.RUnlock()

	listeners, ok := em.listeners[name]
	if !ok {
		return []ListenerFunc{}
	}

	return listeners
}
