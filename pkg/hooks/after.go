package hooks

// AfterHook خطاف بعد التنفيذ
type AfterHook struct {
	manager *HookManager
}

// NewAfterHook ينشئ خطافاً جديداً بعد التنفيذ
func NewAfterHook(manager *HookManager) *AfterHook {
	return &AfterHook{
		manager: manager,
	}
}

// Register يسجل خطافاً بعد التنفيذ
func (ah *AfterHook) Register(name string, handler HookHandler, priority int) {
	ah.manager.RegisterAfter(name, handler, priority)
}

// Execute ينفذ الخطافات بعد التنفيذ
func (ah *AfterHook) Execute(name string, data interface{}) (interface{}, error) {
	return ah.manager.Execute(name, HookAfter, data)
}

// Remove يزيل خطافاً
func (ah *AfterHook) Remove(name string, handler HookHandler) {
	ah.manager.Remove(name, handler)
}
