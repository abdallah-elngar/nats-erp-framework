package database

import (
	"fmt"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/nats-framework/nats/pkg/config"
)

var (
	instance  *Manager
	once      sync.Once
	defaultDB *gorm.DB
)

// Manager يدير اتصالات قاعدة البيانات
type Manager struct {
	connections map[string]*gorm.DB
	defaultName string
	config      *config.DatabaseConfig
	logger      logger.Interface
	mu          sync.RWMutex
}

// NewManager ينشئ مدير قاعدة بيانات جديد
func NewManager(cfg *config.DatabaseConfig) (*Manager, error) {
	var err error
	once.Do(func() {
		instance, err = newManager(cfg)
		if err == nil && instance != nil {
			// تعيين قاعدة البيانات الافتراضية
			if db, ok := instance.connections[instance.defaultName]; ok {
				defaultDB = db
			}
		}
	})
	return instance, err
}

// newManager ينشئ مدير قاعدة بيانات جديد (داخلية)
func newManager(cfg *config.DatabaseConfig) (*Manager, error) {
	m := &Manager{
		connections: make(map[string]*gorm.DB),
		defaultName: cfg.Default,
		config:      cfg,
		logger:      logger.Default.LogMode(logger.Info),
	}

	// تهيئة الاتصالات
	for name, connCfg := range cfg.Connections {
		db, err := m.connect(connCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to %s: %w", name, err)
		}
		m.connections[name] = db
	}

	// التحقق من وجود الاتصال الافتراضي
	if _, ok := m.connections[cfg.Default]; !ok {
		return nil, fmt.Errorf("default connection %s not found", cfg.Default)
	}

	return m, nil
}

// connect ينشئ اتصال بقاعدة البيانات
func (m *Manager) connect(cfg config.ConnectionConfig) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch cfg.Driver {
	case "postgres", "postgresql":
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database)
		dialector = postgres.Open(dsn)

	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
		dialector = mysql.Open(dsn)

	case "sqlite":
		dialector = sqlite.Open(cfg.Database)

	case "sqlserver":
		dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
			cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
		dialector = sqlserver.Open(dsn)

	default:
		return nil, fmt.Errorf("unsupported driver: %s", cfg.Driver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: m.logger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, err
	}

	// إعداد تجمع الاتصالات
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// DB يعيد الاتصال الافتراضي
func (m *Manager) DB() *gorm.DB {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.connections[m.defaultName]
}

// Connection يعيد اتصالاً باسم محدد
func (m *Manager) Connection(name string) *gorm.DB {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if db, ok := m.connections[name]; ok {
		return db
	}
	return m.connections[m.defaultName]
}

// Close يغلق جميع الاتصالات
func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for name, db := range m.connections {
		sqlDB, err := db.DB()
		if err != nil {
			return fmt.Errorf("failed to get sql.DB for %s: %w", name, err)
		}
		if err := sqlDB.Close(); err != nil {
			return fmt.Errorf("failed to close %s: %w", name, err)
		}
	}
	return nil
}

// ============================================
// دوال عامة (Global Functions)
// ============================================

// DB يعيد قاعدة البيانات الافتراضية (دالة عامة)
// func DB() *gorm.DB {
// 	if defaultDB == nil {
// 		// محاولة تهيئة قاعدة البيانات من الإعدادات
// 		cfg, err := config.Load()
// 		if err != nil {
// 			return nil
// 		}
// 		_, err = NewManager(&cfg.Database)
// 		if err != nil {
// 			return nil
// 		}
// 	}
// 	return defaultDB
// }

// SetDB يضع قاعدة البيانات الافتراضية (للاستخدام في الاختبارات)
func SetDB(db *gorm.DB) {
	defaultDB = db
}

// GetManager يعيد مدير قاعدة البيانات الحالي
func GetManager() *Manager {
	return instance
}
func DB() *gorm.DB {
	if instance == nil {
		return nil
	}
	return instance.DB()
}
