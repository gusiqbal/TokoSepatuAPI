package cart

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
)

// mockCartRepo implements ICartRepository for service tests.
type mockCartRepo struct {
	getOrCreateCartFn  func(ctx context.Context, userID uuid.UUID) (*Cart, error)
	getCartWithItemsFn func(ctx context.Context, userID uuid.UUID) (*CartResponse, error)
	addItemFn          func(ctx context.Context, cartID, variantID uuid.UUID, quantity int) error
	updateItemFn       func(ctx context.Context, cartItemID uuid.UUID, quantity int) error
	removeItemFn       func(ctx context.Context, cartItemID uuid.UUID) error
	clearCartFn        func(ctx context.Context, cartID uuid.UUID) error
	getRawCartItemsFn  func(ctx context.Context, userID uuid.UUID) (uuid.UUID, []CartItem, error)
}

func (m *mockCartRepo) GetOrCreateCart(ctx context.Context, userID uuid.UUID) (*Cart, error) {
	if m.getOrCreateCartFn != nil {
		return m.getOrCreateCartFn(ctx, userID)
	}
	return &Cart{ID: uuid.New(), UserID: userID}, nil
}
func (m *mockCartRepo) GetCartWithItems(ctx context.Context, userID uuid.UUID) (*CartResponse, error) {
	if m.getCartWithItemsFn != nil {
		return m.getCartWithItemsFn(ctx, userID)
	}
	return &CartResponse{}, nil
}
func (m *mockCartRepo) AddItem(ctx context.Context, cartID, variantID uuid.UUID, quantity int) error {
	if m.addItemFn != nil {
		return m.addItemFn(ctx, cartID, variantID, quantity)
	}
	return nil
}
func (m *mockCartRepo) UpdateItem(ctx context.Context, cartItemID uuid.UUID, quantity int) error {
	if m.updateItemFn != nil {
		return m.updateItemFn(ctx, cartItemID, quantity)
	}
	return nil
}
func (m *mockCartRepo) RemoveItem(ctx context.Context, cartItemID uuid.UUID) error {
	if m.removeItemFn != nil {
		return m.removeItemFn(ctx, cartItemID)
	}
	return nil
}
func (m *mockCartRepo) ClearCart(ctx context.Context, cartID uuid.UUID) error {
	if m.clearCartFn != nil {
		return m.clearCartFn(ctx, cartID)
	}
	return nil
}
func (m *mockCartRepo) GetRawCartItems(ctx context.Context, userID uuid.UUID) (uuid.UUID, []CartItem, error) {
	if m.getRawCartItemsFn != nil {
		return m.getRawCartItemsFn(ctx, userID)
	}
	return uuid.Nil, nil, nil
}

func newTestCartService(repo ICartRepository) *CartService {
	return NewCartService(repo)
}

// --- GetCart ---

func TestGetCart_Success(t *testing.T) {
	expected := &CartResponse{ID: uuid.New().String(), TotalPrice: 100}
	repo := &mockCartRepo{
		getCartWithItemsFn: func(_ context.Context, _ uuid.UUID) (*CartResponse, error) {
			return expected, nil
		},
	}
	svc := newTestCartService(repo)
	cart, err := svc.GetCart(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cart.TotalPrice != 100 {
		t.Errorf("expected total 100, got %v", cart.TotalPrice)
	}
}

func TestGetCart_RepoError(t *testing.T) {
	repo := &mockCartRepo{
		getCartWithItemsFn: func(_ context.Context, _ uuid.UUID) (*CartResponse, error) {
			return nil, errors.New("db error")
		},
	}
	svc := newTestCartService(repo)
	_, err := svc.GetCart(context.Background(), uuid.New())
	if err == nil {
		t.Error("expected error from repo")
	}
}

// --- AddItem ---

func TestAddItem_Success(t *testing.T) {
	cartID := uuid.New()
	variantID := uuid.New()
	addCalled := false
	repo := &mockCartRepo{
		getOrCreateCartFn: func(_ context.Context, _ uuid.UUID) (*Cart, error) {
			return &Cart{ID: cartID}, nil
		},
		addItemFn: func(_ context.Context, cID, vID uuid.UUID, qty int) error {
			if cID != cartID {
				t.Errorf("expected cartID %v, got %v", cartID, cID)
			}
			if vID != variantID {
				t.Errorf("expected variantID %v, got %v", variantID, vID)
			}
			if qty != 2 {
				t.Errorf("expected qty 2, got %d", qty)
			}
			addCalled = true
			return nil
		},
	}
	svc := newTestCartService(repo)
	err := svc.AddItem(context.Background(), uuid.New(), &AddToCartRequest{
		ProductVariantID: variantID, Quantity: 2,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !addCalled {
		t.Error("expected repo.AddItem to be called")
	}
}

func TestAddItem_CartCreationError(t *testing.T) {
	repo := &mockCartRepo{
		getOrCreateCartFn: func(_ context.Context, _ uuid.UUID) (*Cart, error) {
			return nil, errors.New("cannot create cart")
		},
	}
	svc := newTestCartService(repo)
	err := svc.AddItem(context.Background(), uuid.New(), &AddToCartRequest{
		ProductVariantID: uuid.New(), Quantity: 1,
	})
	if err == nil {
		t.Error("expected error when cart creation fails")
	}
}

// --- UpdateItem ---

func TestUpdateItem_Success(t *testing.T) {
	itemID := uuid.New()
	called := false
	repo := &mockCartRepo{
		updateItemFn: func(_ context.Context, id uuid.UUID, qty int) error {
			if id != itemID {
				t.Errorf("expected itemID %v, got %v", itemID, id)
			}
			called = true
			return nil
		},
	}
	svc := newTestCartService(repo)
	err := svc.UpdateItem(context.Background(), itemID, &UpdateCartItemRequest{Quantity: 3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected repo.UpdateItem to be called")
	}
}

func TestUpdateItem_RepoError(t *testing.T) {
	repo := &mockCartRepo{
		updateItemFn: func(_ context.Context, _ uuid.UUID, _ int) error {
			return errors.New("insufficient stock")
		},
	}
	svc := newTestCartService(repo)
	err := svc.UpdateItem(context.Background(), uuid.New(), &UpdateCartItemRequest{Quantity: 100})
	if err == nil {
		t.Error("expected error from repo")
	}
}

// --- RemoveItem ---

func TestRemoveItem_Success(t *testing.T) {
	itemID := uuid.New()
	called := false
	repo := &mockCartRepo{
		removeItemFn: func(_ context.Context, id uuid.UUID) error {
			if id != itemID {
				t.Errorf("expected itemID %v, got %v", itemID, id)
			}
			called = true
			return nil
		},
	}
	svc := newTestCartService(repo)
	err := svc.RemoveItem(context.Background(), itemID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected repo.RemoveItem to be called")
	}
}

// --- ClearCart ---

func TestClearCart_Success(t *testing.T) {
	cartID := uuid.New()
	clearCalled := false
	repo := &mockCartRepo{
		getOrCreateCartFn: func(_ context.Context, _ uuid.UUID) (*Cart, error) {
			return &Cart{ID: cartID}, nil
		},
		clearCartFn: func(_ context.Context, id uuid.UUID) error {
			if id != cartID {
				t.Errorf("expected cartID %v, got %v", cartID, id)
			}
			clearCalled = true
			return nil
		},
	}
	svc := newTestCartService(repo)
	err := svc.ClearCart(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !clearCalled {
		t.Error("expected repo.ClearCart to be called")
	}
}

func TestClearCart_CartError(t *testing.T) {
	repo := &mockCartRepo{
		getOrCreateCartFn: func(_ context.Context, _ uuid.UUID) (*Cart, error) {
			return nil, errors.New("cannot get cart")
		},
	}
	svc := newTestCartService(repo)
	err := svc.ClearCart(context.Background(), uuid.New())
	if err == nil {
		t.Error("expected error when cart lookup fails")
	}
}
