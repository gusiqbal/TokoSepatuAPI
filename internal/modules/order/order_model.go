package order

import (
	"learnapirest/internal/modules/account"
	"learnapirest/internal/modules/product"

	"github.com/google/uuid"
)

type Order struct {
	ID              uuid.UUID    `gorm:"type:char(36);primary_key"`
	UserID          uuid.UUID    `gorm:"type:char(36); index"`
	Status          string       `gorm:"type:varchar(50);default:'pending'"`
	TotalPrice      float64
	PaymentMethod   string       `gorm:"type:varchar(100)"`
	ShippingAddress string       `gorm:"type:varchar(500)"`
	User            account.User `gorm:"foreignKey:UserID"`
	LastUpdatedAt   int64        `gorm:"autoUpdateTime:milli"`
	CreatedAt       int64        `gorm:"autoCreateTime:milli"`
}

type OrderItem struct {
	ID               uuid.UUID `gorm:"type:char(36);primary_key"`
	OrderID          uuid.UUID `gorm:"type:char(36);index"`
	ProductVariantID uuid.UUID `gorm:"type:char(36);index"`
	Quantity         int
	PriceAtPurchase  float64
	Order            Order                  `gorm:"foreignKey:OrderID"`
	ProductVariant   product.ProductVariant `gorm:"foreignKey:ProductVariantID"`
}

func GetOrder() []any {
	return []any{
		&Order{},
		&OrderItem{},
	}
}
