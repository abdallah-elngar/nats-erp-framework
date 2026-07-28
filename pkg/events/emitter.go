package events

import (
	"time"
)

// Emitter باعث أحداث
type Emitter struct {
	manager *EventManager
	name    string
}

// NewEmitter ينشئ باعث أحداث جديد
func NewEmitter(manager *EventManager, name string) *Emitter {
	return &Emitter{
		manager: manager,
		name:    name,
	}
}

// Emit يصدر حدثاً
func (e *Emitter) Emit(eventName string, data interface{}) error {
	event := NewEvent(eventName, data)
	return e.manager.Emit(event)
}

// EmitAsync يصدر حدثاً بشكل غير متزامن
func (e *Emitter) EmitAsync(eventName string, data interface{}) {
	event := NewEvent(eventName, data)
	e.manager.EmitAsync(event)
}

// EmitDelayed يصدر حدثاً بعد تأخير
func (e *Emitter) EmitDelayed(eventName string, data interface{}, delay time.Duration) {
	go func() {
		time.Sleep(delay)
		e.Emit(eventName, data)
	}()
}
