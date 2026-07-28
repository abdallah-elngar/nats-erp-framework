package scheduler

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

// Job يمثل مهمة مجدولة
type Job struct {
	ID       string
	Name     string
	Schedule string
	Handler  func() error
	Active   bool
	LastRun  time.Time
	NextRun  time.Time
	Runs     int
	Errors   int
}

// Scheduler يدير الجدولة
type Scheduler struct {
	cron   *cron.Cron
	jobs   map[string]*Job
	mu     sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc
}

// NewScheduler ينشئ جدولة جديدة
func NewScheduler() *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())

	return &Scheduler{
		cron:   cron.New(cron.WithSeconds()),
		jobs:   make(map[string]*Job),
		ctx:    ctx,
		cancel: cancel,
	}
}

// AddJob يضيف مهمة مجدولة
func (s *Scheduler) AddJob(name, schedule string, handler func() error) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// التحقق من صحة الجدولة
	if _, err := cron.ParseStandard(schedule); err != nil {
		return fmt.Errorf("invalid schedule: %w", err)
	}

	job := &Job{
		ID:       fmt.Sprintf("job_%d", time.Now().UnixNano()),
		Name:     name,
		Schedule: schedule,
		Handler:  handler,
		Active:   true,
	}

	s.jobs[name] = job

	// إضافة المهمة إلى cron
	entryID, err := s.cron.AddFunc(schedule, func() {
		s.executeJob(name)
	})

	if err != nil {
		return err
	}

	// حفظ معرف الدخول
	job.ID = fmt.Sprintf("%d", entryID)

	return nil
}

// executeJob ينفذ مهمة
func (s *Scheduler) executeJob(name string) {
	s.mu.RLock()
	job, ok := s.jobs[name]
	s.mu.RUnlock()

	if !ok || !job.Active {
		return
	}

	job.Runs++
	job.LastRun = time.Now()

	// تحديث وقت التشغيل التالي
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	schedule, _ := parser.Parse(job.Schedule)
	job.NextRun = schedule.Next(job.LastRun)

	// تنفيذ المهمة
	if err := job.Handler(); err != nil {
		job.Errors++
		// تسجيل الخطأ
		fmt.Printf("Job %s error: %v\n", job.Name, err)
	}
}

// Start يبدأ الجدولة
func (s *Scheduler) Start() {
	s.cron.Start()
}

// Stop يوقف الجدولة
func (s *Scheduler) Stop() {
	s.cron.Stop()
	s.cancel()
}

// RemoveJob يزيل مهمة
func (s *Scheduler) RemoveJob(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	job, ok := s.jobs[name]
	if !ok {
		return fmt.Errorf("job not found: %s", name)
	}

	// إزالة من cron
	entryID, _ := strconv.ParseInt(job.ID, 10, 64)
	s.cron.Remove(cron.EntryID(entryID))

	delete(s.jobs, name)

	return nil
}

// GetJobs يعيد جميع المهام
func (s *Scheduler) GetJobs() map[string]*Job {
	s.mu.RLock()
	defer s.mu.RUnlock()

	jobs := make(map[string]*Job)
	for k, v := range s.jobs {
		jobs[k] = v
	}

	return jobs
}

// GetJob يعيد مهمة واحدة
func (s *Scheduler) GetJob(name string) (*Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	job, ok := s.jobs[name]
	if !ok {
		return nil, fmt.Errorf("job not found: %s", name)
	}

	return job, nil
}

// PauseJob يوقف مهمة
func (s *Scheduler) PauseJob(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	job, ok := s.jobs[name]
	if !ok {
		return fmt.Errorf("job not found: %s", name)
	}

	job.Active = false
	return nil
}

// ResumeJob يستأنف مهمة
func (s *Scheduler) ResumeJob(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	job, ok := s.jobs[name]
	if !ok {
		return fmt.Errorf("job not found: %s", name)
	}

	job.Active = true
	return nil
}
