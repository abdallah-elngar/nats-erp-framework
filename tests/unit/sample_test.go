package unit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSample اختبار عينة
func TestSample(t *testing.T) {
	assert.True(t, true)
	assert.Equal(t, 1, 1)
}

// TestMath اختبار العمليات الحسابية
func TestMath(t *testing.T) {
	t.Run("addition", func(t *testing.T) {
		result := 1 + 1
		assert.Equal(t, 2, result)
	})

	t.Run("subtraction", func(t *testing.T) {
		result := 5 - 3
		assert.Equal(t, 2, result)
	})
}

// TestValidation اختبار التحقق
func TestValidation(t *testing.T) {
	t.Run("email validation", func(t *testing.T) {
		valid := "test@example.com"
		invalid := "not-an-email"

		assert.Contains(t, valid, "@")
		assert.NotContains(t, invalid, "@")
	})
}
