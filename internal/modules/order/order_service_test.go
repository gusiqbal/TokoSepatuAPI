package order

import (
	"context"
	"errors"
	"learnapirest/internal/modules/cart"
	"testing"

	"github.com/google/uuid"
)

// mockOrderRepo implements IOrderRepository.
type mockOrderRepo struct {
	createOrderFn       func(ctx context.Context, params CreateOrderParams) (*Order, error)
	getOrdersByUserIDFn func(ctx context.Context, userID uuid.UUID) ([]OrderResponse, error)
	getOrderByIDFn      func(ctx context.Context, orderID, userID uuid.UUID) (*OrderResponse, error)
}

func (m *mockOrderRepo) CreateOrder(ctx context.Context, params CreateOrderParams) (*Order, error) {
	if m.createOrderFn != nil {
		return m.createOrderFn(ctx, params)
	}
	return &Order{ID: uuid.New()}, nil
}
func (m *mockOrderRepo) GetOrdersByUserID(ctx context.Context, userID uuid.UUID) ([]OrderResponse, error) {
	if m.getOrdersByUserIDFn != nil {
		return m.getOrdersByUserIDFn(ctx, userID)
	}
	return []OrderResponse{}, nil
}
func (m *mockOrderRepo) GetOrderByID(ctx context.Context, orderID, userID uuid.UUID) (*OrderResponse, error) {
	if m.getOrderByIDFn != nil {
		return m.getOrderByIDFn(ctx, orderID, userID)
	}
	return nil, nil
}

// mockCartRepoForOrder implements cart.ICartRepository for order service tests.
type mockCartRepoForOrder struct {
	getRawCartItemsFn func(ctx context.Context, userID uuid.UUID) (uuid.UUID, []cart.CartItem, error)
	clearCartFn       func(ctx context.Context, cartID uuid.UUID) error
}

func (m *mockCartRepoForOrder) GetOrCreateCart(ctx context.Context, userID uuid.UUID) (*cart.Cart, error) {
	return &cart.Cart{ID: uuid.New(), UserID: userID}, nil
}
func (m *mockCartRepoForOrder) GetCartWithItems(ctx context.Context, userID uuid.UUID) (*cart.CartResponse, error) {
	return &cart.CartResponse{}, nil
}
func (m *mockCartRepoForOrder) AddItem(ctx context.Context, cartID, variantID uuid.UUID, quantity int) error {
	return nil
}
func (m *mockCartRepoForOrder) UpdateItem(ctx context.Context, cartItemID uuid.UUID, quantity int) error {
	return nil
}
func (m *mockCartRepoForOrder) RemoveItem(ctx context.Context, cartItemID uuid.UUID) error {
	return nil
}
func (m *mockCartRepoForOrder) ClearCart(ctx context.Context, cartID uuid.UUID) error {
	if m.clearCartFn != nil {
		return m.clearCartFn(ctx, cartID)
	}
	return nil
}
func (m *mockCartRepoForOrder) GetRawCartItems(ctx context.Context, userID uuid.UUID) (uuid.UUID, []cart.CartItem, error) {
	if m.getRawCartItemsFn != nil {
		return m.getRawCartItemsFn(ctx, userID)
	}
	return uuid.Nil, nil, nil
}

func newTestOrderService(orderRepo IOrderRepository, cartRepo cart.ICartRepository) *OrderService {
	return NewOrderService(orderRepo, cartRepo)
}

// --- CreateOrderFromCart ---

func TestCreateOrderFromCart_EmptyCart(t *testing.T) {
	cartRepo := &mockCartRepoForOrder{
		getRawCartItemsFn: func(_ context.Context, _ uuid.UUID) (uuid.UUID, []cart.CartItem, error) {
			return uuid.New(), []cart.CartItem{}, nil
		},
	}
	svc := newTestOrderService(&mockOrderRepo{}, cartRepo)
	_, err := svc.CreateOrderFromCart(context.Background(), uuid.New(), &CreateOrderRequest{
		PaymentMethod: "cash", ShippingAddress: "Jl. Test",
	})
	if err == nil {
		t.Error("expected error for empty cart")
	}
	if err.Error() != "cart is empty" {
		t.Errorf("expected 'cart is empty', got %q", err.Error())
	}
}

func TestCreateOrderFromCart_CartError(t *testing.T) {
	cartRepo := &mockCartRepoForOrder{
		getRawCartItemsFn: func(_ context.Context, _ uuid.UUID) (uuid.UUID, []cart.CartItem, error) {
			return uuid.Nil, nil, errors.New("db error")
		},
	}
	svc := newTestOrderService(&mockOrderRepo{}, cartRepo)
	_, err := svc.CreateOrderFromCart(context.Background(), uuid.New(), &CreateOrderRequest{
		PaymentMethod: "cash", ShippingAddress: "Jl. Test",
	})
	if err == nil {
		t.Error("expected error when cart fetch fails")
	}
}

func TestCreateOrderFromCart_Success(t *testing.T) {
	userID := uuid.New()
	cartID := uuid.New()
	orderID := uuid.New()

	cartRepo := &mockCartRepoForOrder{
		getRawCartItemsFn: func(_ context.Context, _ uuid.UUID) (uuid.UUID, []cart.CartItem, error) {
			items := []cart.CartItem{
				{
					ID:               uuid.New(),
					CartID:           cartID,
					ProductVariantID: uuid.New(),
					Quantity:         2,
				},
			}
			return cartID, items, nil
		},
	}
	orderRepo := &mockOrderRepo{
		createOrderFn: func(_ context.Context, _ CreateOrderParams) (*Order, error) {
			return &Order{ID: orderID, UserID: userID, Status: "pending"}, nil
		},
		getOrderByIDFn: func(_ context.Context, oID, _ uuid.UUID) (*OrderResponse, error) {
			return &OrderResponse{ID: oID.String(), Status: "pending"}, nil
		},
	}
	svc := newTestOrderService(orderRepo, cartRepo)
	resp, err := svc.CreateOrderFromCart(context.Background(), userID, &CreateOrderRequest{
		PaymentMethod: "transfer", ShippingAddress: "Jl. Merdeka 1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ID != orderID.String() {
		t.Errorf("expected order ID %v, got %v", orderID.String(), resp.ID)
	}
}

func TestCreateOrderFromCart_CreateOrderError(t *testing.T) {
	cartRepo := &mockCartRepoForOrder{
		getRawCartItemsFn: func(_ context.Context, _ uuid.UUID) (uuid.UUID, []cart.CartItem, error) {
			return uuid.New(), []cart.CartItem{{ID: uuid.New(), Quantity: 1}}, nil
		},
	}
	orderRepo := &mockOrderRepo{
		createOrderFn: func(_ context.Context, _ CreateOrderParams) (*Order, error) {
			return nil, errors.New("transaction failed")
		},
	}
	svc := newTestOrderService(orderRepo, cartRepo)
	_, err := svc.CreateOrderFromCart(context.Background(), uuid.New(), &CreateOrderRequest{
		PaymentMethod: "cash", ShippingAddress: "Jl. Test",
	})
	if err == nil {
		t.Error("expected error when order creation fails")
	}
}

// --- GetOrderHistory ---

func TestGetOrderHistory_Success(t *testing.T) {
	expected := []OrderResponse{{ID: "order-1"}, {ID: "order-2"}}
	orderRepo := &mockOrderRepo{
		getOrdersByUserIDFn: func(_ context.Context, _ uuid.UUID) ([]OrderResponse, error) {
			return expected, nil
		},
	}
	svc := newTestOrderService(orderRepo, &mockCartRepoForOrder{})
	orders, err := svc.GetOrderHistory(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(orders) != 2 {
		t.Errorf("expected 2 orders, got %d", len(orders))
	}
}

func TestGetOrderHistory_RepoError(t *testing.T) {
	orderRepo := &mockOrderRepo{
		getOrdersByUserIDFn: func(_ context.Context, _ uuid.UUID) ([]OrderResponse, error) {
			return nil, errors.New("db error")
		},
	}
	svc := newTestOrderService(orderRepo, &mockCartRepoForOrder{})
	_, err := svc.GetOrderHistory(context.Background(), uuid.New())
	if err == nil {
		t.Error("expected error from repo")
	}
}

// --- GetOrderDetail ---

func TestGetOrderDetail_Success(t *testing.T) {
	orderID := uuid.New()
	userID := uuid.New()
	orderRepo := &mockOrderRepo{
		getOrderByIDFn: func(_ context.Context, oID, _ uuid.UUID) (*OrderResponse, error) {
			return &OrderResponse{ID: oID.String(), Status: "pending"}, nil
		},
	}
	svc := newTestOrderService(orderRepo, &mockCartRepoForOrder{})
	resp, err := svc.GetOrderDetail(context.Background(), userID, orderID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ID != orderID.String() {
		t.Errorf("expected order ID %v, got %v", orderID.String(), resp.ID)
	}
}

func TestGetOrderDetail_NotFound(t *testing.T) {
	orderRepo := &mockOrderRepo{
		getOrderByIDFn: func(_ context.Context, _, _ uuid.UUID) (*OrderResponse, error) {
			return nil, errors.New("order not found")
		},
	}
	svc := newTestOrderService(orderRepo, &mockCartRepoForOrder{})
	_, err := svc.GetOrderDetail(context.Background(), uuid.New(), uuid.New())
	if err == nil {
		t.Error("expected error when order not found")
	}
}
