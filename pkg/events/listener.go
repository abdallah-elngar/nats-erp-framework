package events

import (
	"fmt"
	"reflect"
	"sync"
)

// Listener مستمع أحداث
type Listener struct {
	manager *EventManager
	name    string
}

// NewListener ينشئ مستمع أحداث جديد
func NewListener(manager *EventManager, name string) *Listener {
	return &Listener{
		manager: manager,
		name:    name,
	}
}

// Listen يسجل مستمعاً
func (l *Listener) Listen(eventName string, fn interface{}) error {
	listener, err := l.toListenerFunc(fn)
	if err != nil {
		return err
	}

	l.manager.Listen(eventName, listener)
	return nil
}

// ListenAsync يسجل مستمعاً غير متزامن
func (l *Listener) ListenAsync(eventName string, fn interface{}) error {
	listener, err := l.toListenerFunc(fn)
	if err != nil {
		return err
	}

	asyncListener := func(event Event) error {
		go listener(event)
		return nil
	}

	l.manager.Listen(eventName, asyncListener)
	return nil
}

// ListenOnce يسجل مستمعاً لمرة واحدة
func (l *Listener) ListenOnce(eventName string, fn interface{}) error {
	listener, err := l.toListenerFunc(fn)
	if err != nil {
		return err
	}

	var once sync.Once

	// ✅ الحل: تعريف دالة خارجية تستدعي نفسها
	// نستخدم متغيراً من النوع ListenerFunc ثم نعيّنه بعد تعريفه
	var onceListener ListenerFunc

	onceListener = func(event Event) error {
		var execErr error
		once.Do(func() {
			execErr = listener(event)
			// إزالة المستمع بعد التنفيذ
			if execErr == nil {
				l.manager.RemoveListener(eventName, onceListener)
			}
		})
		return execErr
	}

	l.manager.Listen(eventName, onceListener)
	return nil
}

// toListenerFunc يحول دالة إلى ListenerFunc
func (l *Listener) toListenerFunc(fn interface{}) (ListenerFunc, error) {
	v := reflect.ValueOf(fn)
	if v.Kind() != reflect.Func {
		return nil, fmt.Errorf("expected function, got %s", v.Kind())
	}

	return func(event Event) error {
		results := v.Call([]reflect.Value{reflect.ValueOf(event)})

		if len(results) == 0 {
			return nil
		}

		if len(results) == 1 {
			if err, ok := results[0].Interface().(error); ok {
				return err
			}
		}

		return nil
	}, nil
}

// RemoveListener يزيل مستمعاً
func (l *Listener) RemoveListener(eventName string, listener ListenerFunc) {
	l.manager.RemoveListener(eventName, listener)
}
