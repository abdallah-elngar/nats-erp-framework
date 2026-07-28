package services

import (
    "errors"

    "github.com/nats-framework/nats/apps/sale/dto"
    "github.com/nats-framework/nats/apps/sale/models"
    "github.com/nats-framework/nats/apps/sale/repository"
)

// SaleService خدمة Sale
type SaleService struct {
    repo *repository.SaleRepository
}

// NewSaleService ينشئ خدمة Sale جديدة
func NewSaleService() *SaleService {
    return &SaleService{
        repo: repository.NewSaleRepository(),
    }
}

// GetAll يعيد جميع sales
func (s *SaleService) GetAll() ([]models.Sale, error) {
    return s.repo.FindAll()
}

// GetByID يعيد Sale بالمعرف
func (s *SaleService) GetByID(id uint) (*models.Sale, error) {
    return s.repo.FindByID(id)
}

// Create ينشئ Sale جديداً
func (s *SaleService) Create(req dto.CreateSaleRequest) (*models.Sale, error) {
    item := &models.Sale{
        // TODO: Map fields from req to model
    }

    if err := s.repo.Create(item); err != nil {
        return nil, err
    }

    return item, nil
}

// Update يحدث Sale
func (s *SaleService) Update(id uint, req dto.UpdateSaleRequest) (*models.Sale, error) {
    item, err := s.repo.FindByID(id)
    if err != nil {
        return nil, errors.New("record not found")
    }

    // TODO: Update fields from req to model

    if err := s.repo.Update(item); err != nil {
        return nil, err
    }

    return item, nil
}

// Delete يحذف Sale
func (s *SaleService) Delete(id uint) error {
    return s.repo.Delete(id)
}
