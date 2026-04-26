package transaction

import "github.com/google/uuid"

type Payment struct {
	ID                    uuid.UUID `gorm:"type:char(36);primary_key"`
	OrderID               uuid.UUID `gorm:"type:char(36);uniqueIndex"`
	StripePaymentIntentID string    `gorm:"type:varchar(255)"`
	Amount                int64
	Currency              string `gorm:"type:varchar(10);default:'usd'"`
	Status                string `gorm:"type:varchar(50);default:'pending'"`
	CreatedAt             int64  `gorm:"autoCreateTime:milli"`
	UpdatedAt             int64  `gorm:"autoUpdateTime:milli"`
}

func GetPayment() []any {
	return []any{&Payment{}}
}
