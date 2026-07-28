package queue

import (
	"sync"
)

// WorkerPool يمثل مجموعة من العمال
type WorkerPool struct {
	queue   *Queue
	workers int
	wg      sync.WaitGroup
	running bool
	mu      sync.Mutex
}

// NewWorkerPool ينشئ مجموعة عمال جديدة
func NewWorkerPool(queue *Queue, workers int) *WorkerPool {
	return &WorkerPool{
		queue:   queue,
		workers: workers,
		running: false,
	}
}

// Start يبدأ تشغيل العمال
func (wp *WorkerPool) Start() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if wp.running {
		return
	}

	wp.running = true
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// Stop يوقف تشغيل العمال
func (wp *WorkerPool) Stop() {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if !wp.running {
		return
	}

	wp.running = false
	wp.wg.Wait()
}

// worker دالة العامل
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()

	for {
		if !wp.running {
			return
		}

		job, err := wp.queue.Pop()
		if err != nil {
			continue
		}

		if job == nil {
			continue
		}

		wp.queue.processJob(job)
	}
}

// GetQueue يعيد الطابور
func (wp *WorkerPool) GetQueue() *Queue {
	return wp.queue
}

// IsRunning يعيد حالة التشغيل
func (wp *WorkerPool) IsRunning() bool {
	wp.mu.Lock()
	defer wp.mu.Unlock()
	return wp.running
}
