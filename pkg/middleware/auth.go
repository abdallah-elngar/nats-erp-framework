package middleware

import (
	"context"

	"github.com/nats-framework/nats/pkg/auth"
)

// GetUserFromContext يحصل على المستخدم من السياق
func GetUserFromContext(ctx context.Context) (*auth.Claims, bool) {
	user, ok := ctx.Value("user").(*auth.Claims)
	return user, ok
}

// GetUserIDFromContext يحصل على معرف المستخدم من السياق
func GetUserIDFromContext(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value("user_id").(uint)
	return userID, ok
}

// GetTokenFromContext يحصل على التوكن من السياق
func GetTokenFromContext(ctx context.Context) (string, bool) {
	token, ok := ctx.Value("token").(string)
	return token, ok
}
