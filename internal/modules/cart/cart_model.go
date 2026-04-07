package cart

import (
	"learnapirest/internal/modules/product"
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	ID        uuid.UUID `gorm:"type:char(36);primary_key"`
	UserID    uuid.UUID `gorm:"type:char(36);uniqueIndex"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Items     []CartItem `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE"`
}

type CartItem struct {
	ID               uuid.UUID              `gorm:"type:char(36);primary_key"`
	CartID           uuid.UUID              `gorm:"type:char(36);index"`
	ProductVariantID uuid.UUID              `gorm:"type:char(36);index"`
	Quantity         int                    `gorm:"default:1"`
	ProductVariant   product.ProductVariant `gorm:"foreignKey:ProductVariantID"`
}

func GetCart() []any {
	return []any{
		&Cart{},
		&CartItem{},
	}
}
