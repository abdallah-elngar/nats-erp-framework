package session

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
)

// Session يمثل جلسة مستخدم
type Session struct {
	ID        string
	UserID    uint
	Data      map[string]interface{}
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

// SessionManager يدير الجلسات
type SessionManager struct {
	store      sessions.Store
	cookieName string
	maxAge     int
}

// Config يمثل إعدادات الجلسات
type Config struct {
	CookieName string
	MaxAge     int
	Secure     bool
	HttpOnly   bool
	Path       string
	Domain     string
}

// NewSessionManager ينشئ مدير جلسات جديد
func NewSessionManager(store sessions.Store, config Config) *SessionManager {
	return &SessionManager{
		store:      store,
		cookieName: config.CookieName,
		maxAge:     config.MaxAge,
	}
}

// Start يبدأ جلسة جديدة
func (sm *SessionManager) Start(w http.ResponseWriter, r *http.Request) (*Session, error) {
	session, err := sm.store.Get(r, sm.cookieName)
	if err != nil {
		return nil, err
	}

	// إنشاء جلسة جديدة
	s := &Session{
		Data:      make(map[string]interface{}),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Duration(sm.maxAge) * time.Second),
	}

	// حفظ معرف الجلسة
	if session.ID != "" {
		s.ID = session.ID
	}

	// حفظ الجلسة
	return s, sm.Save(w, r, s)
}

// Get يحصل على الجلسة
func (sm *SessionManager) Get(r *http.Request) (*Session, error) {
	session, err := sm.store.Get(r, sm.cookieName)
	if err != nil {
		return nil, err
	}

	if session.IsNew {
		return nil, errors.New("session not found")
	}

	s := &Session{
		ID:   session.ID,
		Data: make(map[string]interface{}),
	}

	// قراءة البيانات
	if data, ok := session.Values["data"]; ok {
		if bytes, ok := data.([]byte); ok {
			json.Unmarshal(bytes, &s.Data)
		}
	}

	if userID, ok := session.Values["user_id"]; ok {
		s.UserID = userID.(uint)
	}

	if expiresAt, ok := session.Values["expires_at"]; ok {
		s.ExpiresAt = expiresAt.(time.Time)
	}

	return s, nil
}

// Save يحفظ الجلسة
func (sm *SessionManager) Save(w http.ResponseWriter, r *http.Request, s *Session) error {
	session, err := sm.store.Get(r, sm.cookieName)
	if err != nil {
		return err
	}

	// تحديث البيانات
	s.UpdatedAt = time.Now()

	// حفظ البيانات
	data, err := json.Marshal(s.Data)
	if err != nil {
		return err
	}

	session.Values["data"] = data
	session.Values["user_id"] = s.UserID
	session.Values["expires_at"] = s.ExpiresAt

	session.Options = &sessions.Options{
		MaxAge:   sm.maxAge,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
	}

	return session.Save(r, w)
}

// Destroy يلغي الجلسة
func (sm *SessionManager) Destroy(w http.ResponseWriter, r *http.Request) error {
	session, err := sm.store.Get(r, sm.cookieName)
	if err != nil {
		return err
	}

	session.Options.MaxAge = -1
	return session.Save(r, w)
}

// Set يضع قيمة في الجلسة
func (sm *SessionManager) Set(w http.ResponseWriter, r *http.Request, s *Session, key string, value interface{}) error {
	s.Data[key] = value
	return sm.Save(w, r, s)
}

// GetValue يحصل على قيمة من الجلسة
func (sm *SessionManager) GetValue(s *Session, key string) interface{} {
	return s.Data[key]
}

// DeleteValue يحذف قيمة من الجلسة
func (sm *SessionManager) DeleteValue(w http.ResponseWriter, r *http.Request, s *Session, key string) error {
	delete(s.Data, key)
	return sm.Save(w, r, s)
}
