package dto

import "time"

// CreateSaleRequest طلب إنشاء Sale
type CreateSaleRequest struct {
    SaName string `validate:"required" json:"sa_name"`
    Price float64 `json:"price"`
    Quantity int `json:"quantity"`
}

// UpdateSaleRequest طلب تحديث Sale
type UpdateSaleRequest struct {
    SaName string `json:"sa_name"`
    Price float64 `json:"price"`
    Quantity int `json:"quantity"`
}

// SaleResponse استجابة Sale
type SaleResponse struct {
    ID        uint      `json:"id"`
    SaName string `json:"sa_name"`
    Price float64 `json:"price"`
    Quantity int `json:"quantity"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
