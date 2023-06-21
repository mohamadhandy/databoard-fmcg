package models

import "time"

type Product struct {
	ID         string    `gorm:"primaryKey;type:varchar(255)" json:"id"`
	Name       string    `gorm:"type:varchar(255);not null" json:"name"`
	BrandId    string    `gorm:"type:varchar(255);not null" json:"brand_id"`
	CategoryId string    `gorm:"type:varchar(255);not null" json:"category_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	CreatedBy  string    `gorm:"type:varchar(100);not null" json:"created_by"`
	UpdatedBy  string    `gorm:"type:varchar(100);not null" json:"updated_by"`
	SKU        string    `gorm:"type:varchar(9);not null" json:"sku"`
	Status     string    `gorm:"type:varchar(20);not null" json:"status"`
}

type ProductRequest struct {
	ID         string `gorm:"primaryKey;type:varchar(255)" json:"id"`
	Name       string `gorm:"type:varchar(255);not null" json:"name"`
	BrandId    string `gorm:"type:varchar(255);not null" json:"brand_id"`
	CategoryId string `gorm:"type:varchar(255);not null" json:"category_id"`
}

func (Product) TableName() string {
	return "Product"
}
