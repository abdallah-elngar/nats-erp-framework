package context

import (
	"context"

	"github.com/nats-framework/nats/pkg/auth"
)

// ============================================
// دوال عامة للسياق (تجنب التعارض)
// ============================================

// GetUserFromRequestContext يحصل على المستخدم من سياق الطلب
func GetUserFromRequestContext(ctx context.Context) (*auth.User, bool) {
	return auth.GetUserFromContext(ctx)
}

// GetUserIDFromRequestContext يحصل على معرف المستخدم من سياق الطلب
func GetUserIDFromRequestContext(ctx context.Context) (uint, bool) {
	return auth.GetUserIDFromContext(ctx)
}

// GetUsernameFromRequestContext يحصل على اسم المستخدم من سياق الطلب
func GetUsernameFromRequestContext(ctx context.Context) string {
	if user, ok := auth.GetUserFromContext(ctx); ok && user != nil {
		return user.Username
	}
	return "system"
}

// GetUserEmailFromRequestContext يحصل على بريد المستخدم من سياق الطلب
func GetUserEmailFromRequestContext(ctx context.Context) string {
	if user, ok := auth.GetUserFromContext(ctx); ok && user != nil {
		return user.Email
	}
	return ""
}

// IsAuthenticatedRequest يتحقق من وجود مستخدم في سياق الطلب
func IsAuthenticatedRequest(ctx context.Context) bool {
	return auth.IsAuthenticated(ctx)
}
