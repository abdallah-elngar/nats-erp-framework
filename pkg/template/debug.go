package template

import (
	"fmt"
	"time"
)

type Debugger struct {
	enabled bool
	start   time.Time
	logs    []string
}

func NewDebugger(enabled bool) *Debugger {
	return &Debugger{
		enabled: enabled,
		start:   time.Now(),
		logs:    make([]string, 0),
	}
}

func (d *Debugger) Log(format string, args ...interface{}) {
	if !d.enabled {
		return
	}
	msg := fmt.Sprintf(format, args...)
	timestamp := time.Now().Format("15:04:05.000")
	d.logs = append(d.logs, fmt.Sprintf("[%s] %s", timestamp, msg))
	fmt.Printf("🐛 %s\n", msg)
}

func (d *Debugger) LogPerformance(name string, duration time.Duration) {
	if !d.enabled {
		return
	}
	d.Log("⏱️ Template '%s' rendered in %v", name, duration)
}
