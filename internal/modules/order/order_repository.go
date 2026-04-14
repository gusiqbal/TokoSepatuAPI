package order

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

type CreateOrderParams struct {
	UserID          uuid.UUID
	PaymentMethod   string
	ShippingAddress string
	Items           []OrderItemParams
}

type OrderItemParams struct {
	ProductVariantID uuid.UUID
	Quantity         int
	PriceAtPurchase  float64
}

func (r *OrderRepository) CreateOrder(ctx context.Context, params CreateOrderParams) (*Order, error) {
	if len(params.Items) == 0 {
		return nil, errors.New("cart is empty")
	}

	var totalPrice float64
	for _, item := range params.Items {
		totalPrice += item.PriceAtPurchase * float64(item.Quantity)
	}

	order := Order{
		ID:              uuid.New(),
		UserID:          params.UserID,
		Status:          "pending",
		TotalPrice:      totalPrice,
		PaymentMethod:   params.PaymentMethod,
		ShippingAddress: params.ShippingAddress,
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&order).Error; err != nil {
			return err
		}

		for _, item := range params.Items {
			orderItem := OrderItem{
				ID:               uuid.New(),
				OrderID:          order.ID,
				ProductVariantID: item.ProductVariantID,
				Quantity:         item.Quantity,
				PriceAtPurchase:  item.PriceAtPurchase,
			}
			if err := tx.Create(&orderItem).Error; err != nil {
				return err
			}

			// Decrement stock on the product variant
			if err := tx.Model(&struct{ ID uuid.UUID }{}).
				Table("product_variants").
				Where("id = ? AND stock >= ?", item.ProductVariantID, item.Quantity).
				Update("stock", gorm.Expr("stock - ?", item.Quantity)).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *OrderRepository) GetOrdersByUserID(ctx context.Context, userID uuid.UUID) ([]OrderResponse, error) {
	var orders []Order
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, err
	}

	result := make([]OrderResponse, 0, len(orders))
	for _, o := range orders {
		var items []OrderItem
		r.db.WithContext(ctx).
			Preload("ProductVariant").
			Preload("ProductVariant.Product").
			Where("order_id = ?", o.ID).
			Find(&items)

		result = append(result, buildOrderResponse(o, items))
	}

	return result, nil
}

func (r *OrderRepository) GetOrderByID(ctx context.Context, orderID uuid.UUID, userID uuid.UUID) (*OrderResponse, error) {
	var order Order
	if err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", orderID, userID).
		First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}

	var items []OrderItem
	r.db.WithContext(ctx).
		Preload("ProductVariant").
		Preload("ProductVariant.Product").
		Where("order_id = ?", order.ID).
		Find(&items)

	res := buildOrderResponse(order, items)
	return &res, nil
}

func buildOrderResponse(o Order, items []OrderItem) OrderResponse {
	itemResponses := make([]OrderItemResponse, 0, len(items))
	for _, item := range items {
		v := item.ProductVariant
		itemResponses = append(itemResponses, OrderItemResponse{
			ID:               item.ID.String(),
			ProductVariantID: v.ID.String(),
			ProductName:      v.Product.Name,
			Brand:            v.Product.Brand,
			Color:            v.Color,
			Size:             v.Size,
			Quantity:         item.Quantity,
			PriceAtPurchase:  item.PriceAtPurchase,
			Subtotal:         item.PriceAtPurchase * float64(item.Quantity),
		})
	}

	return OrderResponse{
		ID:              o.ID.String(),
		Status:          o.Status,
		TotalPrice:      o.TotalPrice,
		PaymentMethod:   o.PaymentMethod,
		ShippingAddress: o.ShippingAddress,
		Items:           itemResponses,
		CreatedAt:       o.CreatedAt,
	}
}
