package models

import "time"

type Category struct {
	ID        string    `gorm:"primaryKey;type:varchar(255)" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy string    `gorm:"type:varchar(100);not null" json:"created_by"`
	UpdatedBy string    `gorm:"type:varchar(100);not null" json:"updated_by"`
}

type CategoryRequest struct {
	ID        string `gorm:"primaryKey;type:varchar(255)" json:"id"`
	Name      string `gorm:"type:varchar(100);not null" json:"name"`
	CreatedBy string `gorm:"type:varchar(100);not null" json:"created_by"`
	UpdatedBy string `gorm:"type:varchar(100);not null" json:"updated_by"`
}

func (Category) TableName() string {
	return "Category"
}
