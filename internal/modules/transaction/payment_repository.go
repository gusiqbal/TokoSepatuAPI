package transaction

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IPaymentRepository interface {
	CreatePayment(ctx context.Context, payment *Payment) error
	GetPaymentByOrderID(ctx context.Context, orderID uuid.UUID) (*Payment, error)
	GetPaymentByStripeID(ctx context.Context, stripePaymentIntentID string) (*Payment, error)
	UpdatePaymentStatus(ctx context.Context, stripePaymentIntentID string, status string) error
}

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) CreatePayment(ctx context.Context, payment *Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

func (r *PaymentRepository) GetPaymentByOrderID(ctx context.Context, orderID uuid.UUID) (*Payment, error) {
	var payment Payment
	if err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&payment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment not found")
		}
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepository) GetPaymentByStripeID(ctx context.Context, stripePaymentIntentID string) (*Payment, error) {
	var payment Payment
	if err := r.db.WithContext(ctx).Where("stripe_payment_intent_id = ?", stripePaymentIntentID).First(&payment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payment not found")
		}
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepository) UpdatePaymentStatus(ctx context.Context, stripePaymentIntentID string, status string) error {
	return r.db.WithContext(ctx).
		Model(&Payment{}).
		Where("stripe_payment_intent_id = ?", stripePaymentIntentID).
		Update("status", status).Error
}
