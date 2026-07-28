package queue

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// JobBuilder مساعد لبناء المهام
type JobBuilder struct {
	job *Job
}

// NewJobBuilder ينشئ مساعد بناء مهام جديد
func NewJobBuilder(name string) *JobBuilder {
	return &JobBuilder{
		job: &Job{
			ID:         uuid.New().String(),
			Name:       name,
			Payload:    make(map[string]interface{}),
			Status:     "pending",
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
			MaxRetries: 3,
		},
	}
}

// WithPayload يضيف بيانات المهمة
func (jb *JobBuilder) WithPayload(payload map[string]interface{}) *JobBuilder {
	jb.job.Payload = payload
	return jb
}

// WithRetries يضيف عدد المحاولات
func (jb *JobBuilder) WithRetries(maxRetries int) *JobBuilder {
	jb.job.MaxRetries = maxRetries
	return jb
}

// WithDelay يضيف تأخيراً
func (jb *JobBuilder) WithDelay(delay time.Duration) *JobBuilder {
	// سيتم تنفيذها باستخدام Redis Delayed Queue
	return jb
}

// Build يبني المهمة
func (jb *JobBuilder) Build() *Job {
	return jb.job
}

// ToJSON يحول المهمة إلى JSON
func (jb *JobBuilder) ToJSON() ([]byte, error) {
	return json.Marshal(jb.job)
}
