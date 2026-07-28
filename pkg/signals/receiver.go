package signals

import (
	"sync"
)

// Receiver مستقبل الإشارات
type Receiver struct {
	manager *SignalManager
	name    string
}

// NewReceiver ينشئ مستقبلاً جديداً
func NewReceiver(manager *SignalManager, name string) *Receiver {
	return &Receiver{
		manager: manager,
		name:    name,
	}
}

// Receive يسجل مستقبلاً لإشارة
func (r *Receiver) Receive(signalName string, handler SignalHandler) {
	wrappedHandler := func(signal *Signal) error {
		if data, ok := signal.Data.(map[string]interface{}); ok {
			data["receiver"] = r.name
		}
		return handler(signal)
	}

	r.manager.Connect(signalName, wrappedHandler)
}

// ReceiveOnce يسجل مستقبلاً لمرة واحدة
func (r *Receiver) ReceiveOnce(signalName string, handler SignalHandler) {
	var once sync.Once

	// ✅ تعريف onceHandler كمتغير ثم تعيينه
	var onceHandler SignalHandler

	onceHandler = func(signal *Signal) error {
		var err error
		once.Do(func() {
			if data, ok := signal.Data.(map[string]interface{}); ok {
				data["receiver"] = r.name
			}
			err = handler(signal)
			r.manager.Disconnect(signalName, onceHandler)
		})
		return err
	}

	r.manager.Connect(signalName, onceHandler)
}

// StopListening يتوقف عن الاستماع لإشارة
func (r *Receiver) StopListening(signalName string, handler SignalHandler) {
	r.manager.Disconnect(signalName, handler)
}

// StopAll يتوقف عن الاستماع لجميع الإشارات
func (r *Receiver) StopAll(signalName string) {
	r.manager.Clear(signalName)
}

// GetHandlers يعيد المعالجات المسجلة
func (r *Receiver) GetHandlers(signalName string) []SignalHandler {
	return r.manager.GetHandlers(signalName)
}
