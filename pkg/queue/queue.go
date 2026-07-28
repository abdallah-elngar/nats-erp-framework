package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// Job يمثل مهمة في الطابور
type Job struct {
	ID         string                 `json:"id"`
	Name       string                 `json:"name"`
	Payload    map[string]interface{} `json:"payload"`
	Status     string                 `json:"status"`
	CreatedAt  time.Time              `json:"created_at"`
	UpdatedAt  time.Time              `json:"updated_at"`
	Retries    int                    `json:"retries"`
	MaxRetries int                    `json:"max_retries"`
	Error      string                 `json:"error,omitempty"`
}

// JobHandler دالة معالجة المهمة
type JobHandler func(job *Job) error

// Queue يمثل طابوراً
type Queue struct {
	name     string
	handlers map[string]JobHandler
	client   *redis.Client
	mu       sync.RWMutex
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewQueue ينشئ طابوراً جديداً
func NewQueue(name string, client *redis.Client) *Queue {
	ctx, cancel := context.WithCancel(context.Background())

	return &Queue{
		name:     name,
		handlers: make(map[string]JobHandler),
		client:   client,
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Register يسجل معالجاً لمهمة
func (q *Queue) Register(jobName string, handler JobHandler) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.handlers[jobName] = handler
}

// Push يضيف مهمة إلى الطابور
func (q *Queue) Push(job *Job) error {
	data, err := json.Marshal(job)
	if err != nil {
		return err
	}

	return q.client.LPush(q.ctx, q.name, data).Err()
}

// Pop يسترجع مهمة من الطابور
func (q *Queue) Pop() (*Job, error) {
	result, err := q.client.RPop(q.ctx, q.name).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var job Job
	if err := json.Unmarshal([]byte(result), &job); err != nil {
		return nil, err
	}

	return &job, nil
}

// Process يعالج المهام في الطابور
func (q *Queue) Process() {
	for {
		select {
		case <-q.ctx.Done():
			return
		default:
			job, err := q.Pop()
			if err != nil {
				time.Sleep(time.Second)
				continue
			}

			if job == nil {
				time.Sleep(time.Millisecond * 100)
				continue
			}

			q.processJob(job)
		}
	}
}

// processJob يعالج مهمة واحدة
func (q *Queue) processJob(job *Job) {
	q.mu.RLock()
	handler, ok := q.handlers[job.Name]
	q.mu.RUnlock()

	if !ok {
		job.Status = "failed"
		job.Error = fmt.Sprintf("no handler for job: %s", job.Name)
		return
	}

	job.Status = "processing"
	job.UpdatedAt = time.Now()

	// تنفيذ المهمة
	if err := handler(job); err != nil {
		job.Status = "failed"
		job.Error = err.Error()
		job.Retries++

		// إعادة المحاولة
		if job.Retries < job.MaxRetries {
			job.Status = "pending"
			q.Push(job)
		}
	} else {
		job.Status = "completed"
	}

	job.UpdatedAt = time.Now()
}

// Stop يوقف معالجة الطابور
func (q *Queue) Stop() {
	q.cancel()
}

// GetStats يعيد إحصائيات الطابور
func (q *Queue) GetStats() (int64, error) {
	return q.client.LLen(q.ctx, q.name).Result()
}

// Clear يفرغ الطابور
func (q *Queue) Clear() error {
	return q.client.Del(q.ctx, q.name).Err()
}

// ProcessJobs يعالج المهام بشكل غير متزامن
func (q *Queue) ProcessJobs(workers int) {
	for i := 0; i < workers; i++ {
		go q.Process()
	}
}
