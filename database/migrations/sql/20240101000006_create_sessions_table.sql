-- name: create_sessions_table
-- id: 20240101000000
-- created: 2024-01-01 00:00:00
-- description: Create sessions table for user session management

-- ============================================
-- UP: إنشاء جدول الجلسات
-- ============================================
-- up:
CREATE TABLE IF NOT EXISTS sessions (
    -- المعرف الفريد للجلسة (UUID or custom)
    id VARCHAR(64) PRIMARY KEY,
    
    -- معرف المستخدم (ربط مع جدول users)
    user_id BIGINT NOT NULL,
    
    -- بيانات الجلسة المشفرة (JSON)
    data TEXT,
    
    -- معلومات الجهاز
    ip VARCHAR(45),
    user_agent VARCHAR(255),
    
    -- تاريخ انتهاء الجلسة (للتنظيف التلقائي)
    expires_at TIMESTAMP NOT NULL,
    
    -- حقول التدقيق
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    
    -- المفتاح الخارجي
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- ============================================
-- الفهارس (Indexes)
-- ============================================

-- للبحث السريع عن جلسات المستخدم
CREATE INDEX idx_sessions_user_id ON sessions(user_id);

-- لتنظيف الجلسات المنتهية تلقائياً
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);

-- للحذف الناعم (Soft Delete)
CREATE INDEX idx_sessions_deleted_at ON sessions(deleted_at);

-- للبحث عن الجلسات النشطة
CREATE INDEX idx_sessions_active ON sessions(user_id, expires_at) WHERE deleted_at IS NULL;

-- ============================================
-- دوال مساعدة (Functions)
-- ============================================

-- تحديث updated_at تلقائياً
CREATE OR REPLACE FUNCTION update_sessions_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- إنشاء Trigger للتحديث التلقائي
CREATE TRIGGER trigger_sessions_updated_at
    BEFORE UPDATE ON sessions
    FOR EACH ROW
    EXECUTE FUNCTION update_sessions_updated_at();

-- ============================================
-- تنظيف الجلسات المنتهية (Cleanup Function)
-- ============================================

-- حذف الجلسات المنتهية تلقائياً
CREATE OR REPLACE FUNCTION cleanup_expired_sessions()
RETURNS void AS $$
BEGIN
    DELETE FROM sessions 
    WHERE expires_at < CURRENT_TIMESTAMP 
    AND deleted_at IS NULL;
END;
$$ LANGUAGE plpgsql;

-- ============================================
-- DOWN: حذف جدول الجلسات
-- ============================================
-- down:
DROP TABLE IF EXISTS sessions CASCADE;