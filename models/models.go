package models

import (
	"time"
	"github.com/google/uuid"
)

type User struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"user_id" example:"bfc5b2b1-bc0e-4f2b-8c18-7c7a47fdc9c4"`
	Username  string    `gorm:"unique;not null" json:"username" example:"john_doe"`
	Password  string    `gorm:"not null" json:"password" example:"strongPassword123"`
	Email     string    `gorm:"unique;not null" json:"email" validate:"required,email" example:"john@example.com"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at" example:"2025-07-25T14:00:00Z"`
}


type LoginRequest struct {
    Username string `json:"username" example:"john_doe"`
    Password string `json:"password" example:"strongPassword123"`
}

type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id" example:"2c8a21e3-c882-4b40-9f27-35413e5e64e7"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"user_id" example:"bfc5b2b1-bc0e-4f2b-8c18-7c7a47fdc9c4"`
	Name        string    `gorm:"not null" json:"name" example:"Red T-Shirt"`
	Type        string    `json:"type" example:"Clothing"`
	SKU         string    `gorm:"not null" json:"sku" example:"RTS-XL-001"`
	ImageURL    string    `json:"image_url" example:"https://example.com/images/redshirt.png"`
	Description string    `json:"description" example:"A bright red cotton t-shirt"`
	Quantity    int       `gorm:"not null" json:"quantity" example:"42"`
	Price       float64   `gorm:"not null" json:"price" example:"19.99"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at" example:"2025-07-25T14:00:00Z"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at" example:"2025-07-25T14:30:00Z"`
}


type QuantityUpdateRequest struct {
	Quantity int `json:"quantity" example:"5"`
}