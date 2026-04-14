package order

import (
	"context"
	"errors"
	"learnapirest/internal/modules/cart"

	"github.com/google/uuid"
)

type OrderService struct {
	repo     *OrderRepository
	cartRepo *cart.CartRepository
}

func NewOrderService(repo *OrderRepository, cartRepo *cart.CartRepository) *OrderService {
	return &OrderService{repo: repo, cartRepo: cartRepo}
}

func (s *OrderService) CreateOrderFromCart(ctx context.Context, userID uuid.UUID, req *CreateOrderRequest) (*OrderResponse, error) {
	cartID, items, err := s.cartRepo.GetRawCartItems(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, errors.New("cart is empty")
	}

	orderItems := make([]OrderItemParams, 0, len(items))
	for _, item := range items {
		orderItems = append(orderItems, OrderItemParams{
			ProductVariantID: item.ProductVariantID,
			Quantity:         item.Quantity,
			PriceAtPurchase:  item.ProductVariant.Product.Price,
		})
	}

	order, err := s.repo.CreateOrder(ctx, CreateOrderParams{
		UserID:          userID,
		PaymentMethod:   req.PaymentMethod,
		ShippingAddress: req.ShippingAddress,
		Items:           orderItems,
	})
	if err != nil {
		return nil, err
	}

	// Clear the cart after successful order creation
	_ = s.cartRepo.ClearCart(ctx, cartID)

	return s.repo.GetOrderByID(ctx, order.ID, userID)
}

func (s *OrderService) GetOrderHistory(ctx context.Context, userID uuid.UUID) ([]OrderResponse, error) {
	return s.repo.GetOrdersByUserID(ctx, userID)
}

func (s *OrderService) GetOrderDetail(ctx context.Context, userID uuid.UUID, orderID uuid.UUID) (*OrderResponse, error) {
	return s.repo.GetOrderByID(ctx, orderID, userID)
}
