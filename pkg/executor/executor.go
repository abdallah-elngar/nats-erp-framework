package executor

import (
	"bytes"
	"os/exec"
	"strings"
)

// CommandResult نتيجة تنفيذ الأمر
type CommandResult struct {
	Output  string
	Error   string
	Success bool
}

// ExecuteCommand ينفذ أمراً
func ExecuteCommand(command string) CommandResult {
	var result CommandResult

	// تقسيم الأمر إلى أجزاء
	parts := strings.Fields(command)
	if len(parts) == 0 {
		result.Error = "Empty command"
		return result
	}

	cmd := exec.Command(parts[0], parts[1:]...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	result.Output = stdout.String()

	if err != nil {
		result.Error = stderr.String()
		if result.Error == "" {
			result.Error = err.Error()
		}
		result.Success = false
	} else {
		result.Success = true
		if result.Output == "" {
			result.Output = "✅ Command executed successfully"
		}
	}

	return result
}

// ExecuteNatsCommand ينفذ أمر Nats
func ExecuteNatsCommand(args ...string) CommandResult {
	// البحث عن مسار nats
	natsPath := "./nats"
	// يمكن إضافة منطق للبحث عن المسار الصحيح

	cmd := exec.Command(natsPath, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	result := CommandResult{
		Output: stdout.String(),
	}

	if err != nil {
		result.Error = stderr.String()
		if result.Error == "" {
			result.Error = err.Error()
		}
		result.Success = false
	} else {
		result.Success = true
		if result.Output == "" {
			result.Output = "✅ Command executed successfully"
		}
	}

	return result
}
