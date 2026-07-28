// pkg/template/context.go
package template

import (
	"context"
	"time"
)

// TemplateContext سياق القالب
type TemplateContext struct {
	ctx     context.Context
	cancel  context.CancelFunc // ✅ تخزين دالة الإلغاء
	data    map[string]interface{}
	engine  *Engine
	started time.Time
}

// NewTemplateContext ينشئ سياق قالب جديد
func NewTemplateContext(engine *Engine) *TemplateContext {
	ctx, cancel := context.WithCancel(context.Background())
	return &TemplateContext{
		ctx:     ctx,
		cancel:  cancel,
		data:    make(map[string]interface{}),
		engine:  engine,
		started: time.Now(),
	}
}

// WithContext يضيف سياق
func (tc *TemplateContext) WithContext(ctx context.Context) *TemplateContext {
	// إلغاء السياق القديم
	if tc.cancel != nil {
		tc.cancel()
	}

	// إنشاء سياق جديد مع دالة إلغاء
	newCtx, cancel := context.WithCancel(ctx)
	tc.ctx = newCtx
	tc.cancel = cancel

	return tc
}

// WithData يضيف بيانات
func (tc *TemplateContext) WithData(data map[string]interface{}) *TemplateContext {
	for k, v := range data {
		tc.data[k] = v
	}
	return tc
}

// WithValue يضيف قيمة
func (tc *TemplateContext) WithValue(key string, value interface{}) *TemplateContext {
	tc.data[key] = value
	return tc
}

// GetData يعيد البيانات
func (tc *TemplateContext) GetData() map[string]interface{} {
	// إضافة بيانات افتراضية
	if tc.engine.auth != nil {
		tc.data = tc.engine.auth.PrepareData(tc.data).(map[string]interface{})
	}

	// إضافة بيانات الوقت
	tc.data["now"] = time.Now()
	tc.data["started"] = tc.started

	return tc.data
}

// Context يعيد السياق
func (tc *TemplateContext) Context() context.Context {
	return tc.ctx
}

// WithTimeout يضيف مهلة (مصحح)
func (tc *TemplateContext) WithTimeout(timeout time.Duration) *TemplateContext {
	// إلغاء السياق القديم
	if tc.cancel != nil {
		tc.cancel()
	}

	// إنشاء سياق جديد مع مهلة ودالة إلغاء
	ctx, cancel := context.WithTimeout(tc.ctx, timeout)
	tc.ctx = ctx
	tc.cancel = cancel

	return tc
}

// WithDeadline يضيف موعد نهائي (مصحح)
func (tc *TemplateContext) WithDeadline(deadline time.Time) *TemplateContext {
	// إلغاء السياق القديم
	if tc.cancel != nil {
		tc.cancel()
	}

	// إنشاء سياق جديد مع موعد نهائي ودالة إلغاء
	ctx, cancel := context.WithDeadline(tc.ctx, deadline)
	tc.ctx = ctx
	tc.cancel = cancel

	return tc
}

// Cancel يلغي السياق
func (tc *TemplateContext) Cancel() {
	if tc.cancel != nil {
		tc.cancel()
	}
}

// Done يعيد قناة الانتهاء
func (tc *TemplateContext) Done() <-chan struct{} {
	return tc.ctx.Done()
}

// Err يعيد الخطأ
func (tc *TemplateContext) Err() error {
	return tc.ctx.Err()
}

// Deadline يعيد الموعد النهائي
func (tc *TemplateContext) Deadline() (time.Time, bool) {
	return tc.ctx.Deadline()
}

// IsTimeout يتحقق مما إذا كان السياق قد انتهى بسبب المهلة
func (tc *TemplateContext) IsTimeout() bool {
	return tc.ctx.Err() == context.DeadlineExceeded
}

// IsCancelled يتحقق مما إذا كان السياق قد ألغي
func (tc *TemplateContext) IsCancelled() bool {
	return tc.ctx.Err() == context.Canceled
}

// Close يغلق السياق ويحرر الموارد (للأستخدام مع defer)
func (tc *TemplateContext) Close() error {
	tc.Cancel()
	return nil
}
