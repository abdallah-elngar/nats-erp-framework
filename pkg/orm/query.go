package orm

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// QueryConditions شروط الاستعلام
type QueryConditions struct {
	Filters  []Filter
	Sorts    []Sort
	Page     int
	PageSize int
	Selects  []string
	Preloads []string
	Joins    []Join
	Groups   []string
	Havings  []Having
}

// Filter شرط تصفية
type Filter struct {
	Field    string
	Operator string // eq, ne, gt, gte, lt, lte, like, in, not_in
	Value    interface{}
}

// Sort ترتيب
type Sort struct {
	Field     string
	Direction string // asc, desc
}

// Join انضمام
type Join struct {
	Type  string // inner, left, right
	Table string
	On    string
}

// Having شرط HAVING
type Having struct {
	Query string
	Args  []interface{}
}

// ApplyConditions يطبق الشروط على الاستعلام
func ApplyConditions(db *gorm.DB, conditions *QueryConditions) *gorm.DB {
	if conditions == nil {
		return db
	}

	// تطبيق التحديد
	if len(conditions.Selects) > 0 {
		db = db.Select(conditions.Selects)
	}

	// تطبيق التصفية
	for _, filter := range conditions.Filters {
		db = applyFilter(db, filter)
	}

	// تطبيق الترتيب
	for _, sort := range conditions.Sorts {
		db = db.Order(fmt.Sprintf("%s %s", sort.Field, sort.Direction))
	}

	// تطبيق الانضمام
	for _, join := range conditions.Joins {
		joinStr := fmt.Sprintf("%s JOIN %s ON %s", strings.ToUpper(join.Type), join.Table, join.On)
		db = db.Joins(joinStr)
	}

	// تطبيق التحميل المسبق
	for _, preload := range conditions.Preloads {
		db = db.Preload(preload)
	}

	// تطبيق التجميع
	for _, group := range conditions.Groups {
		db = db.Group(group)
	}

	// تطبيق HAVING
	for _, having := range conditions.Havings {
		db = db.Having(having.Query, having.Args...)
	}

	// تطبيق الترقيم
	if conditions.Page > 0 && conditions.PageSize > 0 {
		offset := (conditions.Page - 1) * conditions.PageSize
		db = db.Limit(conditions.PageSize).Offset(offset)
	}

	return db
}

// applyFilter يطبق شرط تصفية واحد
func applyFilter(db *gorm.DB, filter Filter) *gorm.DB {
	switch filter.Operator {
	case "eq":
		return db.Where(fmt.Sprintf("%s = ?", filter.Field), filter.Value)
	case "ne":
		return db.Where(fmt.Sprintf("%s != ?", filter.Field), filter.Value)
	case "gt":
		return db.Where(fmt.Sprintf("%s > ?", filter.Field), filter.Value)
	case "gte":
		return db.Where(fmt.Sprintf("%s >= ?", filter.Field), filter.Value)
	case "lt":
		return db.Where(fmt.Sprintf("%s < ?", filter.Field), filter.Value)
	case "lte":
		return db.Where(fmt.Sprintf("%s <= ?", filter.Field), filter.Value)
	case "like":
		return db.Where(fmt.Sprintf("%s LIKE ?", filter.Field), fmt.Sprintf("%%%v%%", filter.Value))
	case "in":
		return db.Where(fmt.Sprintf("%s IN ?", filter.Field), filter.Value)
	case "not_in":
		return db.Where(fmt.Sprintf("%s NOT IN ?", filter.Field), filter.Value)
	default:
		return db.Where(fmt.Sprintf("%s = ?", filter.Field), filter.Value)
	}
}
