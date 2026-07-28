package repository

import (
    "gorm.io/gorm"

    "github.com/nats-framework/nats/apps/sale/models"
    "github.com/nats-framework/nats/pkg/database"
)

// SaleRepository مستودع Sale
type SaleRepository struct {
    db *gorm.DB
}

// NewSaleRepository ينشئ مستودع Sale جديد
func NewSaleRepository() *SaleRepository {
    return &SaleRepository{
        db: database.DB(),
    }
}

// Create ينشئ Sale جديداً
func (r *SaleRepository) Create(item *models.Sale) error {
    return r.db.Create(item).Error
}

// FindByID يبحث عن Sale بالمعرف
func (r *SaleRepository) FindByID(id uint) (*models.Sale, error) {
    var item models.Sale
    err := r.db.First(&item, id).Error
    if err != nil {
        return nil, err
    }
    return &item, nil
}

// FindAll يعيد جميع sales
func (r *SaleRepository) FindAll() ([]models.Sale, error) {
    var items []models.Sale
    err := r.db.Find(&items).Error
    return items, err
}

// Update يحدث Sale
func (r *SaleRepository) Update(item *models.Sale) error {
    return r.db.Save(item).Error
}

// Delete يحذف Sale
func (r *SaleRepository) Delete(id uint) error {
    return r.db.Delete(&models.Sale{}, id).Error
}

// Exists يتحقق من وجود Sale
func (r *SaleRepository) Exists(query string, args ...interface{}) (bool, error) {
    var count int64
    err := r.db.Model(&models.Sale{}).Where(query, args...).Count(&count).Error
    return count > 0, err
}
