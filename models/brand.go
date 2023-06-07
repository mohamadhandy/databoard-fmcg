package models

import "time"

type Brand struct {
	ID        string    `gorm:"primaryKey;type:varchar(255)" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy string    `gorm:"type:varchar(100);not null" json:"created_by"`
	UpdatedBy string    `gorm:"type:varchar(100);not null" json:"updated_by"`
	Status    string    `gorm:"type:varchar(20);not null" json:"status"`
}

func (Brand) TableName() string {
	return "Brand"
}

type BrandRequest struct {
	ID        string `gorm:"primaryKey;type:varchar(255)" json:"id"`
	Name      string `gorm:"type:varchar(100);not null" json:"name"`
	CreatedBy string `gorm:"type:varchar(100);not null" json:"created_by"`
	UpdatedBy string `gorm:"type:varchar(100);not null" json:"updated_by"`
	Status    string `gorm:"type:varchar(20);not null" json:"status"`
}

type RequestableBrandInterface interface {
	// create brand
	ForCreation() Brand
}

func (b *BrandRequest) ForCreation() Brand {
	return Brand{
		ID:        b.ID,
		Name:      b.Name,
		CreatedBy: b.CreatedBy,
		UpdatedBy: b.UpdatedBy,
		Status:    b.Status,
	}
}
