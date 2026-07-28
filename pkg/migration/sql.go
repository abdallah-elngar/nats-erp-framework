package migration

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time" // ✅ إضافة import time
)

// SQLMigrator يدير هجرات SQL
type SQLMigrator struct {
	migrationsDir string
}

// NewSQLMigrator ينشئ مدير هجرات SQL جديد
func NewSQLMigrator(migrationsDir string) *SQLMigrator {
	return &SQLMigrator{
		migrationsDir: migrationsDir,
	}
}

// LoadMigrationFiles يقوم بتحميل ملفات الهجرات
func (sm *SQLMigrator) LoadMigrationFiles() ([]string, error) {
	var files []string

	entries, err := os.ReadDir(sm.migrationsDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			files = append(files, filepath.Join(sm.migrationsDir, entry.Name()))
		}
	}

	// ترتيب الملفات
	sort.Strings(files)

	return files, nil
}

// ReadMigrationFile يقرأ ملف هجرة
func (sm *SQLMigrator) ReadMigrationFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// ParseMigrationFile يحلل ملف هجرة
func (sm *SQLMigrator) ParseMigrationFile(path string) (name string, upSQL string, downSQL string, err error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", "", "", err
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "-- name:") {
			name = strings.TrimSpace(strings.TrimPrefix(line, "-- name:"))
		} else if strings.HasPrefix(line, "-- up:") {
			upSQL = strings.TrimSpace(strings.TrimPrefix(line, "-- up:"))
		} else if strings.HasPrefix(line, "-- down:") {
			downSQL = strings.TrimSpace(strings.TrimPrefix(line, "-- down:"))
		} else if !strings.HasPrefix(line, "--") && name != "" {
			if upSQL == "" {
				upSQL = line
			}
		}
	}

	return name, upSQL, downSQL, nil
}

// CreateMigrationFile يقوم بإنشاء ملف هجرة SQL
func (sm *SQLMigrator) CreateMigrationFile(name, upSQL, downSQL string) (string, error) {
	timestamp := time.Now().Format("20060102150405")
	id := timestamp + "_" + strings.ReplaceAll(strings.ToLower(name), " ", "_")

	content := fmt.Sprintf(`-- name: %s
-- id: %s
-- created: %s

-- up:
%s

-- down:
%s
`, name, id, time.Now().Format("2006-01-02 15:04:05"), upSQL, downSQL)

	filename := id + ".sql"
	path := filepath.Join(sm.migrationsDir, filename)

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return "", err
	}

	return path, nil
}

// GenerateTableMigration يقوم بإنشاء هجرة جدول
func (sm *SQLMigrator) GenerateTableMigration(tableName string, columns []ColumnDefinition) (string, error) {
	var upSQL strings.Builder
	var downSQL strings.Builder

	upSQL.WriteString(fmt.Sprintf("CREATE TABLE %s (\n", tableName))
	for i, col := range columns {
		upSQL.WriteString(fmt.Sprintf("    %s %s", col.Name, col.Type))
		if col.NotNull {
			upSQL.WriteString(" NOT NULL")
		}
		if col.Unique {
			upSQL.WriteString(" UNIQUE")
		}
		if col.PrimaryKey {
			upSQL.WriteString(" PRIMARY KEY")
		}
		if col.Default != "" {
			upSQL.WriteString(fmt.Sprintf(" DEFAULT %s", col.Default))
		}
		if i < len(columns)-1 {
			upSQL.WriteString(",")
		}
		upSQL.WriteString("\n")
	}
	upSQL.WriteString(");")

	downSQL.WriteString(fmt.Sprintf("DROP TABLE IF EXISTS %s;", tableName))

	name := fmt.Sprintf("create_%s_table", tableName)

	return sm.CreateMigrationFile(name, upSQL.String(), downSQL.String())
}

// ColumnDefinition تعريف عمود
type ColumnDefinition struct {
	Name       string
	Type       string
	NotNull    bool
	Unique     bool
	PrimaryKey bool
	Default    string
	Reference  string
}
