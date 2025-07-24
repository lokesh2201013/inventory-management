package models

import (
	"time"
	"github.com/google/uuid"
)

type User struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"user_id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Password  string    `gorm:"not null" json:"password"`
	Email     string    `gorm:"unique;not null" json:"email"  validate:"required,email"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}


type Product struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	UserID      uuid.UUID `gorm:"type:uuid;not null" json:"user_id"` 
	Name        string    `gorm:"not null" json:"name"`
	Type        string    `json:"type"`
	SKU         string    `gorm:"unique;not null" json:"sku"`
	ImageURL    string    `json:"image_url"`
	Description string    `json:"description"`
	Quantity    int       `gorm:"not null" json:"quantity"`
	Price       float64   `gorm:"not null" json:"price"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
