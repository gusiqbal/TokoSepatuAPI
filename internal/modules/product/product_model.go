package product

import (
	"github.com/google/uuid"
)

type Product struct {
	ID            uuid.UUID `gorm:"type:char(36);primary_key"`
	Name          string    `gorm:"type:varchar(255);not null"`
	Brand         string    `gorm:"type:varchar(100);not null"`
	Size          int       `gorm:"not null"`
	Price         float64   `gorm:"not null"`
	Stock         int       `gorm:"not null"`
	LastUpdatedAt int64     `gorm:"autoUpdateTime:milli"`
	CreatedAt     int64     `gorm:"autoCreateTime:milli"`
}

type ProductVariant struct {
	ID        uuid.UUID `gorm:"type:char(36);primary_key"`
	ProductID uuid.UUID `gorm:"type:char(36);index"`
	Size      int
	Color     string
	Stock     int
	Product   Product `gorm:"foreignKey:ProductID"`
}

type ProductFavorite struct {
	ID        uuid.UUID `gorm:"type:char(36);primary_key"`
	ProductID uuid.UUID `gorm:"type:char(36);index"`
}

func GetProduct() []any {
	return []any{
		&Product{},
		&ProductVariant{},
	}
}
