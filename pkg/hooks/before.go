package hooks

// BeforeHook خطاف قبل التنفيذ
type BeforeHook struct {
	manager *HookManager
}

// NewBeforeHook ينشئ خطافاً جديداً قبل التنفيذ
func NewBeforeHook(manager *HookManager) *BeforeHook {
	return &BeforeHook{
		manager: manager,
	}
}

// Register يسجل خطافاً قبل التنفيذ
func (bh *BeforeHook) Register(name string, handler HookHandler, priority int) {
	bh.manager.RegisterBefore(name, handler, priority)
}

// Execute ينفذ الخطافات قبل التنفيذ
func (bh *BeforeHook) Execute(name string, data interface{}) (interface{}, error) {
	return bh.manager.Execute(name, HookBefore, data)
}

// Remove يزيل خطافاً
func (bh *BeforeHook) Remove(name string, handler HookHandler) {
	bh.manager.Remove(name, handler)
}
