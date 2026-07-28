package auth

import (
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
)

// SessionManager يدير الجلسات
type SessionManager struct {
	store      sessions.Store
	cookieName string
	maxAge     int
}

// SessionConfig يمثل إعدادات الجلسات
type SessionConfig struct {
	CookieName string
	MaxAge     int
	Secure     bool
	HttpOnly   bool
	Path       string
}

// NewSessionManager ينشئ مدير جلسات جديد
func NewSessionManager(store sessions.Store, config SessionConfig) *SessionManager {
	return &SessionManager{
		store:      store,
		cookieName: config.CookieName,
		maxAge:     config.MaxAge,
	}
}

// Get يحصل على الجلسة
func (sm *SessionManager) Get(r *http.Request) (*sessions.Session, error) {
	session, err := sm.store.Get(r, sm.cookieName)
	if err != nil {
		return nil, err
	}
	return session, nil
}

// Set يضع قيمة في الجلسة
func (sm *SessionManager) Set(w http.ResponseWriter, r *http.Request, key string, value interface{}) error {
	session, err := sm.store.Get(r, sm.cookieName)
	if err != nil {
		return err
	}

	session.Values[key] = value
	session.Options = &sessions.Options{
		MaxAge: sm.maxAge,
		Path:   "/",
	}

	return session.Save(r, w)
}

// GetValue يحصل على قيمة من الجلسة
func (sm *SessionManager) GetValue(r *http.Request, key string) interface{} {
	session, err := sm.store.Get(r, sm.cookieName)
	if err != nil {
		return nil
	}
	return session.Values[key]
}

// Delete يحذف قيمة من الجلسة
func (sm *SessionManager) Delete(w http.ResponseWriter, r *http.Request, key string) error {
	session, err := sm.store.Get(r, sm.cookieName)
	if err != nil {
		return err
	}

	delete(session.Values, key)
	return session.Save(r, w)
}

// Destroy يلغي الجلسة
func (sm *SessionManager) Destroy(w http.ResponseWriter, r *http.Request) error {
	session, err := sm.store.Get(r, sm.cookieName)
	if err != nil {
		return err
	}

	// إلغاء الجلسة
	session.Options.MaxAge = -1
	return session.Save(r, w)
}

// IsAuthenticated يتحقق من المصادقة
func (sm *SessionManager) IsAuthenticated(r *http.Request) bool {
	userID := sm.GetValue(r, "user_id")
	return userID != nil
}

// GetUserID يحصل على معرف المستخدم
func (sm *SessionManager) GetUserID(r *http.Request) (uint, error) {
	userID := sm.GetValue(r, "user_id")
	if userID == nil {
		return 0, errors.New("user not authenticated")
	}

	if id, ok := userID.(uint); ok {
		return id, nil
	}

	return 0, errors.New("invalid user id")
}
