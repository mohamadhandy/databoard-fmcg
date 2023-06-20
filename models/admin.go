package models

import "time"

type Admin struct {
	ID          string    `gorm:"primaryKey;type:varchar(255)" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Password    string    `gorm:"type:varchar(255);not null" json:"password"`
	Role        string    `gorm:"type:varchar(255);not null" json:"role"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `gorm:"type:varchar(255);not null" json:"email"`
	Status      string    `gorm:"type:varchar(20);not null" json:"status"`
	PhoneNumber string    `gorm:"type:varchar(20);not null" json:"phone_number"`
}

type AdminResponse struct {
	// Singular
	ID          string `json:"id"`
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Phonenumber string `json:"phone_number" binding:"required"`
	Status      string `json:"status" binding:"required"`
	Role        string `json:"role" binding:"required"`
}

type AdminRequest struct {
	// Singular
	ID          string `json:"id"`
	Name        string `json:"name" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Phonenumber string `json:"phone_number" binding:"required"`
	Status      string `json:"status" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Role        string `json:"role" binding:"required"`
	// coming soon NewPassword
	NewPassword string `json:"new_password"`
	// Plural
	Page  uint `json:"page"`
	Limit uint `json:"limit"`
}

func (Admin) TableName() string {
	return "Admin"
}

type RequestableAdministratorInterface interface {
	// create admin
	ForCreation() Admin
	ForList() (uint, uint)
}

func (req *AdminRequest) ForList() (uint, uint) {
	return req.Page, req.Limit
}

func (req *AdminRequest) ForCreation() Admin {
	return Admin{
		ID:          req.ID,
		Name:        req.Name,
		Email:       req.Email,
		PhoneNumber: req.Phonenumber,
		Status:      req.Status,
		Role:        req.Role,
		Password:    req.Password,
	}
}
