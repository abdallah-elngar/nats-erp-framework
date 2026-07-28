package signals

// Sender مرسل الإشارات
type Sender struct {
	manager *SignalManager
	name    string
}

// NewSender ينشئ مرسلاً جديداً
func NewSender(manager *SignalManager, name string) *Sender {
	return &Sender{
		manager: manager,
		name:    name,
	}
}

// Send يرسل إشارة
func (s *Sender) Send(signalName string, data interface{}) error {
	signal := NewSignal(signalName, data, s.name)
	return s.manager.Emit(signal)
}

// SendAsync يرسل إشارة بشكل غير متزامن
func (s *Sender) SendAsync(signalName string, data interface{}) {
	signal := NewSignal(signalName, data, s.name)
	s.manager.EmitAsync(signal)
}

// SendTo يرسل إشارة لمعالج محدد
func (s *Sender) SendTo(signalName string, data interface{}, handlerName string) error {
	signal := NewSignal(signalName, data, s.name)
	// إضافة معلومات المستهدف
	signal.Data = map[string]interface{}{
		"data":   data,
		"target": handlerName,
	}
	return s.manager.Emit(signal)
}
