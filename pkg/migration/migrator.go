package migration

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Migrator يدير الهجرات
type Migrator struct {
	db         *gorm.DB
	migrations []Migration
	tableName  string
}

// Migration يمثل هجرة
type Migration struct {
	ID        string
	Name      string
	UpSQL     string
	DownSQL   string
	UpFunc    func(*gorm.DB) error
	DownFunc  func(*gorm.DB) error
	CreatedAt time.Time
}

// MigrationRecord سجل الهجرة في قاعدة البيانات
type MigrationRecord struct {
	ID        string `gorm:"primaryKey;size:255"`
	Name      string `gorm:"size:255"`
	Batch     int
	CreatedAt time.Time
}

// NewMigrator ينشئ مدير هجرات جديد
func NewMigrator(db *gorm.DB) *Migrator {
	m := &Migrator{
		db:         db,
		migrations: make([]Migration, 0),
		tableName:  "migration_records", // ✅ تغيير من "migrations" إلى "migration_records"
	}

	// إنشاء جدول الهجرات إذا لم يكن موجوداً
	db.AutoMigrate(&MigrationRecord{})

	return m
}

// Register يسجل هجرة
func (m *Migrator) Register(id, name string, up, down func(*gorm.DB) error) {
	m.migrations = append(m.migrations, Migration{
		ID:       id,
		Name:     name,
		UpFunc:   up,
		DownFunc: down,
	})
}

// RegisterSQL يسجل هجرة SQL
func (m *Migrator) RegisterSQL(id, name, upSQL, downSQL string) {
	m.migrations = append(m.migrations, Migration{
		ID:      id,
		Name:    name,
		UpSQL:   upSQL,
		DownSQL: downSQL,
	})
}

// Run يقوم بتنفيذ الهجرات المعلقة
func (m *Migrator) Run() error {
	executed, err := m.getExecutedMigrations()
	if err != nil {
		return err
	}

	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].ID < m.migrations[j].ID
	})

	for _, migration := range m.migrations {
		// ✅ استخراج اسم الجدول من اسم الهجرة
		tableName := strings.TrimPrefix(migration.Name, "create_")
		tableName = strings.TrimSuffix(tableName, "_table")

		// ✅ التحقق من وجود الجدول الفعلي
		var count int64
		m.db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_name = ?", tableName).Scan(&count)

		// ✅ إذا كان الجدول غير موجود، ننفذ الهجرة (حتى لو كانت مسجلة)
		if count == 0 {
			fmt.Printf("🔄 Table %s not found, executing migration %s\n", tableName, migration.ID)

			if err := m.executeMigration(migration); err != nil {
				return fmt.Errorf("failed to execute migration %s: %w", migration.ID, err)
			}

			// ✅ تحديث السجل (حذف وإعادة إدراج)
			m.db.Where("id = ?", migration.ID).Delete(&MigrationRecord{})
			if err := m.recordMigration(migration); err != nil {
				return err
			}

			fmt.Printf("✅ Migration %s executed successfully\n", migration.ID)
			continue
		}

		// ✅ إذا كانت الهجرة غير مسجلة، ننفذها
		if !executed[migration.ID] {
			fmt.Printf("📝 Running migration: %s (%s)\n", migration.Name, migration.ID)

			if err := m.executeMigration(migration); err != nil {
				return fmt.Errorf("failed to run migration %s: %w", migration.ID, err)
			}

			if err := m.recordMigration(migration); err != nil {
				return err
			}

			fmt.Printf("✅ Migration %s completed\n", migration.ID)
		}
	}

	return nil
}

// Rollback يقوم بالتراجع عن آخر الهجرات
func (m *Migrator) Rollback(step int) error {
	// الحصول على الهجرات المنفذة
	records, err := m.getMigrationRecords()
	if err != nil {
		return err
	}

	if len(records) == 0 {
		fmt.Println("📝 No migrations to rollback")
		return nil
	}

	// ترتيب عكسي
	for i := len(records) - 1; i >= 0 && step > 0; i-- {
		record := records[i]
		step--

		// البحث عن الهجرة
		var migration Migration
		found := false
		for _, m := range m.migrations {
			if m.ID == record.ID {
				migration = m
				found = true
				break
			}
		}

		if !found {
			continue
		}

		fmt.Printf("🔄 Rolling back: %s (%s)\n", migration.Name, migration.ID)

		if migration.DownFunc != nil {
			if err := migration.DownFunc(m.db); err != nil {
				return fmt.Errorf("failed to rollback %s: %w", migration.ID, err)
			}
		} else if migration.DownSQL != "" {
			if err := m.db.Exec(migration.DownSQL).Error; err != nil {
				return fmt.Errorf("failed to rollback %s: %w", migration.ID, err)
			}
		}

		// حذف سجل الهجرة
		if err := m.db.Where("id = ?", migration.ID).Delete(&MigrationRecord{}).Error; err != nil {
			return err
		}

		fmt.Printf("✅ Rollback %s completed\n", migration.ID)
	}

	return nil
}

// Reset يقوم بإعادة تعيين جميع الهجرات
func (m *Migrator) Reset() error {
	// ✅ التحقق من وجود الجدول قبل الحذف
	var exists bool
	if err := m.db.Raw("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = ?)", m.tableName).Scan(&exists).Error; err != nil {
		// إذا كان هناك خطأ، حاول الحذف مباشرة
		if err := m.db.Exec("DELETE FROM " + m.tableName).Error; err != nil {
			// تجاهل الخطأ إذا كان الجدول غير موجود
			if !strings.Contains(err.Error(), "does not exist") {
				return err
			}
		}
	} else if exists {
		// الجدول موجود، قم بحذف البيانات
		if err := m.db.Exec("DELETE FROM " + m.tableName).Error; err != nil {
			return err
		}
	}

	// ✅ إعادة ضبط التسلسل إذا كان موجوداً (PostgreSQL)
	if err := m.db.Exec("ALTER SEQUENCE " + m.tableName + "_id_seq RESTART WITH 1").Error; err != nil {
		// تجاهل الخطأ إذا كان التسلسل غير موجود
	}

	// تنفيذ جميع الهجرات
	return m.Run()
}

// getExecutedMigrations يحصل على الهجرات المنفذة
func (m *Migrator) getExecutedMigrations() (map[string]bool, error) {
	var records []MigrationRecord
	if err := m.db.Find(&records).Error; err != nil {
		return nil, err
	}

	executed := make(map[string]bool)
	for _, record := range records {
		executed[record.ID] = true
	}

	return executed, nil
}

// getMigrationRecords يحصل على سجلات الهجرات
func (m *Migrator) getMigrationRecords() ([]MigrationRecord, error) {
	var records []MigrationRecord
	if err := m.db.Order("created_at DESC").Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

// executeMigration ينفذ هجرة
func (m *Migrator) executeMigration(migration Migration) error {
	if migration.UpFunc != nil {
		return migration.UpFunc(m.db)
	}

	if migration.UpSQL != "" {
		return m.db.Exec(migration.UpSQL).Error
	}

	return fmt.Errorf("no up function or SQL for migration %s", migration.ID)
}

// recordMigration يسجل هجرة (يحذف السجل القديم أولاً)
func (m *Migrator) recordMigration(migration Migration) error {
	// ✅ حذف السجل القديم إذا كان موجوداً
	m.db.Where("id = ?", migration.ID).Delete(&MigrationRecord{})

	// ✅ إنشاء سجل جديد
	record := MigrationRecord{
		ID:        migration.ID,
		Name:      migration.Name,
		Batch:     1,
		CreatedAt: time.Now(),
	}
	return m.db.Create(&record).Error
}

// LoadFromDirectory يقوم بتحميل الهجرات من مجلد
func (m *Migrator) LoadFromDirectory(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".go") {
			continue
		}

		// قراءة ملف الهجرة
		path := filepath.Join(dir, file.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		// تحليل الهجرة
		// سيتم تنفيذها لاحقاً
		_ = content
	}

	return nil
}

// GenerateMigrationFile يقوم بإنشاء ملف هجرة
func (m *Migrator) GenerateMigrationFile(name string, upSQL, downSQL string) (string, error) {
	// إنشاء معرف
	timestamp := time.Now().Format("20060102150405")
	id := timestamp + "_" + strings.ReplaceAll(strings.ToLower(name), " ", "_")

	// إنشاء المحتوى
	content := fmt.Sprintf(`package migrations

import (
    "gorm.io/gorm"
)

// %s هي هجرة %s
func init() {
    // تسجيل الهجرة
    // Migrator.Register("%s", "%s", Up_%s, Down_%s)
}

// Up_%s يقوم بتطبيق الهجرة
func Up_%s(db *gorm.DB) error {
    sql := `+"`"+`%s`+"`"+`
    return db.Exec(sql).Error
}

// Down_%s يقوم بالتراجع عن الهجرة
func Down_%s(db *gorm.DB) error {
    sql := `+"`"+`%s`+"`"+`
    return db.Exec(sql).Error
}
`,
		id,      // 1: %s - id في التعليق الأول
		name,    // 2: %s - name في التعليق الأول
		id,      // 3: %s - id في Migrator.Register
		name,    // 4: %s - name في Migrator.Register
		id,      // 5: %s - id في Up_ في Migrator.Register
		id,      // 6: %s - id في Down_ في Migrator.Register
		id,      // 7: %s - id في Up_%s في التعليق
		id,      // 8: %s - id في func Up_%s
		upSQL,   // 9: %s - SQL للـ Up
		id,      // 10: %s - id في Down_%s في التعليق
		id,      // 11: %s - id في func Down_%s
		downSQL, // 12: %s - SQL للـ Down
	)

	// حفظ الملف
	filename := id + ".go"
	path := filepath.Join("database/migrations", filename)

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return "", err
	}

	return path, nil
}
