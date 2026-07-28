package debug

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

// Debugger يدير وضع التصحيح
type Debugger struct {
	Enabled   bool
	Level     string // "basic", "detailed", "full"
	Output    string // "console", "file", "both"
	FilePath  string
	StartTime time.Time
	Requests  []RequestInfo
}

// RequestInfo معلومات الطلب
type RequestInfo struct {
	Command    string
	Args       []string
	Flags      map[string]interface{}
	StartTime  time.Time
	EndTime    time.Time
	Duration   time.Duration
	Error      string
	StackTrace string
	Memory     uint64
	Goroutines int
}

// NewDebugger ينشئ مصححاً جديداً
func NewDebugger(enabled bool) *Debugger {
	return &Debugger{
		Enabled:   enabled,
		Level:     "basic",
		Output:    "console",
		FilePath:  "storage/logs/debug.log",
		StartTime: time.Now(),
		Requests:  make([]RequestInfo, 0),
	}
}

// Enable يفعل وضع التصحيح
func (d *Debugger) Enable() {
	d.Enabled = true
	d.StartTime = time.Now()
}

// Disable يعطل وضع التصحيح
func (d *Debugger) Disable() {
	d.Enabled = false
}

// LogCommand يسجل أمر CLI
func (d *Debugger) LogCommand(command string, args []string, flags map[string]interface{}) {
	if !d.Enabled {
		return
	}

	info := RequestInfo{
		Command:   command,
		Args:      args,
		Flags:     flags,
		StartTime: time.Now(),
	}

	d.Requests = append(d.Requests, info)
}

// LogEnd يسجل نهاية الأمر
func (d *Debugger) LogEnd(command string, err error) {
	if !d.Enabled {
		return
	}

	for i := len(d.Requests) - 1; i >= 0; i-- {
		if d.Requests[i].Command == command && d.Requests[i].EndTime.IsZero() {
			d.Requests[i].EndTime = time.Now()
			d.Requests[i].Duration = d.Requests[i].EndTime.Sub(d.Requests[i].StartTime)
			d.Requests[i].Memory = getMemoryUsage()
			d.Requests[i].Goroutines = runtime.NumGoroutine()

			if err != nil {
				d.Requests[i].Error = err.Error()
				d.Requests[i].StackTrace = string(debug.Stack())
			}

			d.printInfo(d.Requests[i])
			d.saveToFile(d.Requests[i])
			break
		}
	}
}

// printInfo يطبع معلومات التصحيح
func (d *Debugger) printInfo(info RequestInfo) {
	if d.Output != "file" {
		fmt.Println("\n" + strings.Repeat("=", 60))
		fmt.Println("🐛 DEBUG INFORMATION")
		fmt.Println(strings.Repeat("=", 60))

		fmt.Printf("📋 Command: %s\n", info.Command)
		if len(info.Args) > 0 {
			fmt.Printf("📝 Args: %v\n", strings.Join(info.Args, " "))
		}
		if len(info.Flags) > 0 {
			fmt.Printf("🚩 Flags: %v\n", info.Flags)
		}

		fmt.Printf("\n⏱️  Start Time: %s\n", info.StartTime.Format("2006-01-02 15:04:05.000"))
		fmt.Printf("⏱️  End Time: %s\n", info.EndTime.Format("2006-01-02 15:04:05.000"))
		fmt.Printf("⏱️  Duration: %v\n", info.Duration)

		fmt.Printf("\n💾 Memory Usage: %s\n", FormatBytes(info.Memory))
		fmt.Printf("🧵 Goroutines: %d\n", info.Goroutines)

		if info.Error != "" {
			fmt.Printf("\n❌ Error: %s\n", info.Error)
			if d.Level == "detailed" || d.Level == "full" {
				fmt.Printf("\n📚 Stack Trace:\n%s\n", info.StackTrace)
			}
		}

		fmt.Println(strings.Repeat("=", 60))
	}
}

// saveToFile يحفظ معلومات التصحيح في ملف
func (d *Debugger) saveToFile(info RequestInfo) {
	if d.Output == "console" {
		return
	}

	file, err := os.OpenFile(d.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	data, _ := json.MarshalIndent(info, "", "  ")
	file.WriteString(string(data) + "\n")
}

// Summary يعرض ملخص التصحيح
func (d *Debugger) Summary() {
	if !d.Enabled || len(d.Requests) == 0 {
		return
	}

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("📊 DEBUG SUMMARY")
	fmt.Println(strings.Repeat("=", 60))

	fmt.Printf("📋 Total Commands: %d\n", len(d.Requests))
	fmt.Printf("⏱️  Total Time: %v\n", time.Since(d.StartTime))

	var totalMemory uint64
	var errors int
	for _, req := range d.Requests {
		totalMemory += req.Memory
		if req.Error != "" {
			errors++
		}
	}

	fmt.Printf("💾 Average Memory: %s\n", FormatBytes(totalMemory/uint64(len(d.Requests))))
	fmt.Printf("❌ Errors: %d\n", errors)

	if errors > 0 {
		fmt.Println("\n⚠️  Error Details:")
		for _, req := range d.Requests {
			if req.Error != "" {
				fmt.Printf("  - %s: %s\n", req.Command, req.Error)
			}
		}
	}

	fmt.Println(strings.Repeat("=", 60))
}

// ============================================
// دوال مساعدة (Helper Functions)
// ============================================

// getMemoryUsage يحصل على استخدام الذاكرة
func getMemoryUsage() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc
}

// FormatBytes ينسق البايت إلى KB/MB/GB (دالة عامة)
func FormatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// GetRequestInfo يعيد معلومات الطلب
func (d *Debugger) GetRequestInfo(command string) *RequestInfo {
	for i := len(d.Requests) - 1; i >= 0; i-- {
		if d.Requests[i].Command == command {
			return &d.Requests[i]
		}
	}
	return nil
}

// GetAllRequests يعيد جميع الطلبات
func (d *Debugger) GetAllRequests() []RequestInfo {
	return d.Requests
}
