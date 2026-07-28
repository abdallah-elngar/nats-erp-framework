package pipes

import (
	"context"
	"sync"
)

// Pipe يمثل أنبوباً لمعالجة البيانات
type Pipe interface {
	Process(ctx context.Context, data interface{}) (interface{}, error)
}

// PipeFunc دالة معالجة الأنبوب
type PipeFunc func(ctx context.Context, data interface{}) (interface{}, error)

// Process ينفذ الأنبوب
func (f PipeFunc) Process(ctx context.Context, data interface{}) (interface{}, error) {
	return f(ctx, data)
}

// Pipeline يمثل سلسلة من الأنابيب
type Pipeline struct {
	pipes []Pipe
	mu    sync.RWMutex
}

// NewPipeline ينشئ سلسلة أنابيب جديدة
func NewPipeline(pipes ...Pipe) *Pipeline {
	return &Pipeline{
		pipes: pipes,
	}
}

// Add يضيف أنبوباً
func (p *Pipeline) Add(pipe Pipe) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.pipes = append(p.pipes, pipe)
}

// Process ينفذ السلسلة
func (p *Pipeline) Process(ctx context.Context, data interface{}) (interface{}, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var err error
	result := data

	for _, pipe := range p.pipes {
		result, err = pipe.Process(ctx, result)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// PipelineBuilder مساعد لبناء السلسلة
type PipelineBuilder struct {
	pipes []Pipe
}

// NewPipelineBuilder ينشئ مساعد بناء جديد
func NewPipelineBuilder() *PipelineBuilder {
	return &PipelineBuilder{
		pipes: make([]Pipe, 0),
	}
}

// Pipe يضيف أنبوباً
func (pb *PipelineBuilder) Pipe(pipe Pipe) *PipelineBuilder {
	pb.pipes = append(pb.pipes, pipe)
	return pb
}

// PipeFunc يضيف دالة أنبوب
func (pb *PipelineBuilder) PipeFunc(fn func(ctx context.Context, data interface{}) (interface{}, error)) *PipelineBuilder {
	pb.pipes = append(pb.pipes, PipeFunc(fn))
	return pb
}

// Build يبني السلسلة
func (pb *PipelineBuilder) Build() *Pipeline {
	return NewPipeline(pb.pipes...)
}
