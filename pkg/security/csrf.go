package security

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
	"time"
)

// CSRFProtection يحمي من CSRF
type CSRFProtection struct {
	tokens map[string]time.Time
	mu     sync.RWMutex
	maxAge time.Duration
}

// NewCSRFProtection ينشئ حماية CSRF جديدة
func NewCSRFProtection(maxAge time.Duration) *CSRFProtection {
	return &CSRFProtection{
		tokens: make(map[string]time.Time),
		maxAge: maxAge,
	}
}

// GenerateToken يولد توكن CSRF
func (c *CSRFProtection) GenerateToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	token := base64.URLEncoding.EncodeToString(bytes)

	c.mu.Lock()
	c.tokens[token] = time.Now().Add(c.maxAge)
	c.mu.Unlock()

	return token
}

// ValidateToken يتحقق من صحة توكن CSRF
func (c *CSRFProtection) ValidateToken(token string) error {
	c.mu.RLock()
	expiry, ok := c.tokens[token]
	c.mu.RUnlock()

	if !ok {
		return errors.New("invalid CSRF token")
	}

	if time.Now().After(expiry) {
		c.mu.Lock()
		delete(c.tokens, token)
		c.mu.Unlock()
		return errors.New("CSRF token expired")
	}

	return nil
}

// RemoveToken يزيل توكن CSRF
func (c *CSRFProtection) RemoveToken(token string) {
	c.mu.Lock()
	delete(c.tokens, token)
	c.mu.Unlock()
}

// Cleanup يزيل التوكنات المنتهية
func (c *CSRFProtection) Cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for token, expiry := range c.tokens {
		if now.After(expiry) {
			delete(c.tokens, token)
		}
	}
}
