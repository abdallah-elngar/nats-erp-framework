package session

import (
	"errors"
	"sync"

	"github.com/gorilla/sessions"
)

// Manager يدير الجلسات المتعددة
type Manager struct {
	stores       map[string]sessions.Store
	defaultStore string
	config       Config
	mu           sync.RWMutex
}

// NewManager ينشئ مدير جلسات جديد
func NewManager(config Config) *Manager {
	return &Manager{
		stores: make(map[string]sessions.Store),
		config: config,
	}
}

// RegisterStore يسجل متجر جلسات
func (m *Manager) RegisterStore(name string, store sessions.Store) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.stores[name] = store
	if m.defaultStore == "" {
		m.defaultStore = name
	}
}

// SetDefaultStore يحدد المتجر الافتراضي
func (m *Manager) SetDefaultStore(name string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if _, ok := m.stores[name]; !ok {
		return errors.New("store not found")
	}

	m.defaultStore = name
	return nil
}

// GetStore يحصل على متجر جلسات
func (m *Manager) GetStore(name string) (sessions.Store, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if name == "" {
		name = m.defaultStore
	}

	store, ok := m.stores[name]
	if !ok {
		return nil, errors.New("store not found")
	}

	return store, nil
}

// GetSessionManager يحصل على مدير جلسات لمتجر
func (m *Manager) GetSessionManager(storeName string) (*SessionManager, error) {
	store, err := m.GetStore(storeName)
	if err != nil {
		return nil, err
	}

	return NewSessionManager(store, m.config), nil
}

// CookieStore يخلق متجر كوكيز
func (m *Manager) CookieStore(keyPairs ...[]byte) sessions.Store {
	return sessions.NewCookieStore(keyPairs...)
}

// FilesystemStore يخلق متجر ملفات
func (m *Manager) FilesystemStore(path string, keyPairs ...[]byte) sessions.Store {
	return sessions.NewFilesystemStore(path, keyPairs...)
}
