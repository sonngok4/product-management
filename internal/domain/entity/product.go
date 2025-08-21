package entity

import (
	"time"

	"gorm.io/gorm"
)

// Product represents a product entity in the domain layer
type Product struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Name        string         `json:"name" gorm:"size:255;not null" validate:"required,min=3,max=255"`
	Description string         `json:"description" gorm:"type:text"`
	Price       float64        `json:"price" gorm:"type:decimal(10,2);not null" validate:"required,min=0"`
	Stock       int            `json:"stock" gorm:"default:0" validate:"min=0"`
	Category    string         `json:"category" gorm:"size:100"`
	ImageURL    string         `json:"image_url" gorm:"size:500"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName returns the table name for Product entity
func (Product) TableName() string {
	return "products"
}

// BeforeCreate is a GORM hook that runs before creating a product
func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
	}
	p.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate is a GORM hook that runs before updating a product
func (p *Product) BeforeUpdate(tx *gorm.DB) error {
	p.UpdatedAt = time.Now()
	return nil
}

// Validate performs basic validation on the product entity
func (p *Product) Validate() error {
	if p.Name == "" {
		return ErrProductNameRequired
	}
	if len(p.Name) < 3 {
		return ErrProductNameTooShort
	}
	if len(p.Name) > 255 {
		return ErrProductNameTooLong
	}
	if p.Price < 0 {
		return ErrProductPriceInvalid
	}
	if p.Stock < 0 {
		return ErrProductStockInvalid
	}
	return nil
}
