package engine

import "fmt"

// LifecycleStage يمثل مرحلة في دورة حياة التطبيق
type LifecycleStage string

const (
	StageBooting  LifecycleStage = "booting"
	StageBooted   LifecycleStage = "booted"
	StageRunning  LifecycleStage = "running"
	StageShutdown LifecycleStage = "shutdown"
)

// Lifecycle يدير دورة حياة التطبيق
type Lifecycle struct {
	stage  LifecycleStage
	hooks  map[LifecycleStage][]func() error
	engine *Engine
}

// NewLifecycle ينشئ مدير دورة حياة جديد
func NewLifecycle(engine *Engine) *Lifecycle {
	return &Lifecycle{
		stage:  StageBooting,
		hooks:  make(map[LifecycleStage][]func() error),
		engine: engine,
	}
}

// RegisterHook يسجل خطافاً في مرحلة معينة
func (l *Lifecycle) RegisterHook(stage LifecycleStage, hook func() error) {
	l.hooks[stage] = append(l.hooks[stage], hook)
}

// Boot يقوم بتشغيل التطبيق
func (l *Lifecycle) Boot() error {
	// تنفيذ خطافات التمهيد
	for _, hook := range l.hooks[StageBooting] {
		if err := hook(); err != nil {
			return fmt.Errorf("boot hook failed: %w", err)
		}
	}

	l.stage = StageBooted
	return nil
}

// Run يقوم بتشغيل التطبيق
func (l *Lifecycle) Run() error {
	l.stage = StageRunning
	return nil
}

// Shutdown يقوم بإيقاف التطبيق
func (l *Lifecycle) Shutdown() error {
	l.stage = StageShutdown
	return nil
}

// GetStage يعيد المرحلة الحالية
func (l *Lifecycle) GetStage() LifecycleStage {
	return l.stage
}
