package orm

import (
	"time"

	"gorm.io/gorm"
)

// ORM هو غلاف لـ GORM مع وظائف إضافية
type ORM struct {
	db *gorm.DB
}

// New ينشئ ORM جديد
func New(db *gorm.DB) *ORM {
	return &ORM{db: db}
}

// DB يعيد اتصال GORM الأصلي
func (o *ORM) DB() *gorm.DB {
	return o.db
}

// Model يمثل نموذج قاعدة البيانات
type Model interface {
	TableName() string
	GetID() uint
}

// BaseModel نموذج أساسي مع حقول مشتركة
type BaseModel struct {
	ID        uint           `gorm:"primaryKey"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// GetID يعيد معرف النموذج
func (m *BaseModel) GetID() uint {
	return m.ID
}

// QueryBuilder مساعد لبناء الاستعلامات
type QueryBuilder struct {
	db            *gorm.DB
	model         interface{}
	selectFields  []string
	whereClauses  []string
	orderClauses  []string
	groupClauses  []string
	havingClauses []string
	joins         []string
	preloads      []string
	limit         int
	offset        int
}

// NewQueryBuilder ينشئ مساعد استعلام جديد
func NewQueryBuilder(db *gorm.DB, model interface{}) *QueryBuilder {
	return &QueryBuilder{
		db:    db,
		model: model,
	}
}

// Select يضيف حقولاً للاستعلام
func (q *QueryBuilder) Select(fields ...string) *QueryBuilder {
	q.selectFields = append(q.selectFields, fields...)
	return q
}

// Where يضيف شرط WHERE
func (q *QueryBuilder) Where(query string, args ...interface{}) *QueryBuilder {
	q.whereClauses = append(q.whereClauses, query)
	q.db = q.db.Where(query, args...)
	return q
}

// Order يضيف ترتيب
func (q *QueryBuilder) Order(order string) *QueryBuilder {
	q.orderClauses = append(q.orderClauses, order)
	q.db = q.db.Order(order)
	return q
}

// Limit يحدد عدد النتائج
func (q *QueryBuilder) Limit(limit int) *QueryBuilder {
	q.limit = limit
	q.db = q.db.Limit(limit)
	return q
}

// Offset يحدد نقطة البداية
func (q *QueryBuilder) Offset(offset int) *QueryBuilder {
	q.offset = offset
	q.db = q.db.Offset(offset)
	return q
}

// Preload يضيف تحميل مسبق للعلاقات
func (q *QueryBuilder) Preload(relation string) *QueryBuilder {
	q.preloads = append(q.preloads, relation)
	q.db = q.db.Preload(relation)
	return q
}

// Join يضيف JOIN
func (q *QueryBuilder) Join(query string, args ...interface{}) *QueryBuilder {
	q.joins = append(q.joins, query)
	q.db = q.db.Joins(query, args...)
	return q
}

// Group يضيف GROUP BY
func (q *QueryBuilder) Group(group string) *QueryBuilder {
	q.groupClauses = append(q.groupClauses, group)
	q.db = q.db.Group(group)
	return q
}

// Having يضيف HAVING
func (q *QueryBuilder) Having(query string, args ...interface{}) *QueryBuilder {
	q.havingClauses = append(q.havingClauses, query)
	q.db = q.db.Having(query, args...)
	return q
}

// Find ينفذ الاستعلام ويعيد النتائج
func (q *QueryBuilder) Find(dest interface{}) error {
	return q.db.Find(dest).Error
}

// First يعيد النتيجة الأولى
func (q *QueryBuilder) First(dest interface{}) error {
	return q.db.First(dest).Error
}

// Count يحسب عدد النتائج
func (q *QueryBuilder) Count() (int64, error) {
	var count int64
	err := q.db.Count(&count).Error
	return count, err
}

// Paginate يحسب الترقيم
func (q *QueryBuilder) Paginate(page, pageSize int) (*PaginationResult, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	// حساب العدد الإجمالي
	var total int64
	countDB := q.db.Session(&gorm.Session{})
	if err := countDB.Model(q.model).Count(&total).Error; err != nil {
		return nil, err
	}

	// حساب الإزاحة
	offset := (page - 1) * pageSize
	q.db = q.db.Limit(pageSize).Offset(offset)

	return &PaginationResult{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
		Offset:     offset,
	}, nil
}

// PaginationResult نتيجة الترقيم
type PaginationResult struct {
	Page       int
	PageSize   int
	Total      int64
	TotalPages int
	Offset     int
}

// Repository مستودع عام للبيانات
type Repository[T Model] struct {
	db  *gorm.DB
	orm *ORM
}

// NewRepository ينشئ مستودعاً جديداً
func NewRepository[T Model](db *gorm.DB) *Repository[T] {
	return &Repository[T]{
		db:  db,
		orm: New(db),
	}
}

// Create ينشئ سجلاً جديداً
func (r *Repository[T]) Create(entity *T) error {
	return r.db.Create(entity).Error
}

// Update يحدث سجلاً
func (r *Repository[T]) Update(entity *T) error {
	return r.db.Save(entity).Error
}

// Delete يحذف سجلاً (ناعم)
func (r *Repository[T]) Delete(id uint) error {
	var entity T
	return r.db.Delete(&entity, id).Error
}

// DeletePermanently يحذف سجلاً بشكل دائم
func (r *Repository[T]) DeletePermanently(id uint) error {
	var entity T
	return r.db.Unscoped().Delete(&entity, id).Error
}

// FindByID يبحث عن سجل بالمعرف
func (r *Repository[T]) FindByID(id uint) (*T, error) {
	var entity T
	err := r.db.First(&entity, id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// FindAll يعيد جميع السجلات
func (r *Repository[T]) FindAll() ([]T, error) {
	var entities []T
	err := r.db.Find(&entities).Error
	return entities, err
}

// FindWhere يبحث عن سجلات بشرط
func (r *Repository[T]) FindWhere(query string, args ...interface{}) ([]T, error) {
	var entities []T
	err := r.db.Where(query, args...).Find(&entities).Error
	return entities, err
}

// Exists يتحقق من وجود سجل
func (r *Repository[T]) Exists(query string, args ...interface{}) (bool, error) {
	var count int64
	err := r.db.Model(new(T)).Where(query, args...).Count(&count).Error
	return count > 0, err
}

// Transaction ينفذ عملية في معاملة
func (r *Repository[T]) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}

// BulkInsert يقوم بإدخال متعدد
func (r *Repository[T]) BulkInsert(entities []T, batchSize int) error {
	return r.db.CreateInBatches(entities, batchSize).Error
}
