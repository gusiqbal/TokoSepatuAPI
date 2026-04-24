package cart

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// mockCartService implements ICartService for handler tests.
type mockCartService struct {
	getCartFn    func(ctx context.Context, userID uuid.UUID) (*CartResponse, error)
	addItemFn    func(ctx context.Context, userID uuid.UUID, req *AddToCartRequest) error
	updateItemFn func(ctx context.Context, cartItemID uuid.UUID, req *UpdateCartItemRequest) error
	removeItemFn func(ctx context.Context, cartItemID uuid.UUID) error
	clearCartFn  func(ctx context.Context, userID uuid.UUID) error
}

func (m *mockCartService) GetCart(ctx context.Context, userID uuid.UUID) (*CartResponse, error) {
	if m.getCartFn != nil {
		return m.getCartFn(ctx, userID)
	}
	return &CartResponse{}, nil
}
func (m *mockCartService) AddItem(ctx context.Context, userID uuid.UUID, req *AddToCartRequest) error {
	if m.addItemFn != nil {
		return m.addItemFn(ctx, userID, req)
	}
	return nil
}
func (m *mockCartService) UpdateItem(ctx context.Context, cartItemID uuid.UUID, req *UpdateCartItemRequest) error {
	if m.updateItemFn != nil {
		return m.updateItemFn(ctx, cartItemID, req)
	}
	return nil
}
func (m *mockCartService) RemoveItem(ctx context.Context, cartItemID uuid.UUID) error {
	if m.removeItemFn != nil {
		return m.removeItemFn(ctx, cartItemID)
	}
	return nil
}
func (m *mockCartService) ClearCart(ctx context.Context, userID uuid.UUID) error {
	if m.clearCartFn != nil {
		return m.clearCartFn(ctx, userID)
	}
	return nil
}

func withUserID(userID uuid.UUID) gin.HandlerFunc {
	return func(c *gin.Context) { c.Set("userId", userID) }
}

// --- GetCart handler ---

func TestGetCartHandler_Unauthorized(t *testing.T) {
	ctrl := NewCartController(&mockCartService{})
	r := gin.New()
	r.GET("/cart", ctrl.GetCart)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/cart", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestGetCartHandler_Success(t *testing.T) {
	userID := uuid.New()
	svc := &mockCartService{
		getCartFn: func(_ context.Context, _ uuid.UUID) (*CartResponse, error) {
			return &CartResponse{ID: "cart-1", TotalPrice: 50000}, nil
		},
	}
	ctrl := NewCartController(svc)
	r := gin.New()
	r.GET("/cart", withUserID(userID), ctrl.GetCart)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/cart", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetCartHandler_ServiceError(t *testing.T) {
	svc := &mockCartService{
		getCartFn: func(_ context.Context, _ uuid.UUID) (*CartResponse, error) {
			return nil, errors.New("db error")
		},
	}
	ctrl := NewCartController(svc)
	r := gin.New()
	r.GET("/cart", withUserID(uuid.New()), ctrl.GetCart)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/cart", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

// --- AddItem handler ---

func TestAddItemHandler_Unauthorized(t *testing.T) {
	ctrl := NewCartController(&mockCartService{})
	r := gin.New()
	r.POST("/cart/items", ctrl.AddItem)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/cart/items", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAddItemHandler_InvalidJSON(t *testing.T) {
	ctrl := NewCartController(&mockCartService{})
	r := gin.New()
	r.POST("/cart/items", withUserID(uuid.New()), ctrl.AddItem)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/cart/items", bytes.NewBufferString("bad"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestAddItemHandler_Success(t *testing.T) {
	ctrl := NewCartController(&mockCartService{})
	r := gin.New()
	r.POST("/cart/items", withUserID(uuid.New()), ctrl.AddItem)

	b, _ := json.Marshal(gin.H{"productVariantId": uuid.New().String(), "quantity": 2})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/cart/items", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

// --- UpdateItem handler ---

func TestUpdateItemHandler_InvalidID(t *testing.T) {
	ctrl := NewCartController(&mockCartService{})
	r := gin.New()
	r.PUT("/cart/items/:id", ctrl.UpdateItem)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/cart/items/not-a-uuid", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestUpdateItemHandler_InvalidJSON(t *testing.T) {
	ctrl := NewCartController(&mockCartService{})
	r := gin.New()
	r.PUT("/cart/items/:id", ctrl.UpdateItem)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/cart/items/"+uuid.New().String(), bytes.NewBufferString("bad"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestUpdateItemHandler_Success(t *testing.T) {
	ctrl := NewCartController(&mockCartService{})
	r := gin.New()
	r.PUT("/cart/items/:id", ctrl.UpdateItem)

	b, _ := json.Marshal(gin.H{"quantity": 3})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/cart/items/"+uuid.New().String(), bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

// --- RemoveItem handler ---

func TestRemoveItemHandler_InvalidID(t *testing.T) {
	ctrl := NewCartController(&mockCartService{})
	r := gin.New()
	r.DELETE("/cart/items/:id", ctrl.RemoveItem)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/cart/items/not-a-uuid", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestRemoveItemHandler_Success(t *testing.T) {
	ctrl := NewCartController(&mockCartService{})
	r := gin.New()
	r.DELETE("/cart/items/:id", ctrl.RemoveItem)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/cart/items/"+uuid.New().String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

// --- ClearCart handler ---

func TestClearCartHandler_Unauthorized(t *testing.T) {
	ctrl := NewCartController(&mockCartService{})
	r := gin.New()
	r.DELETE("/cart", ctrl.ClearCart)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/cart", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestClearCartHandler_Success(t *testing.T) {
	ctrl := NewCartController(&mockCartService{})
	r := gin.New()
	r.DELETE("/cart", withUserID(uuid.New()), ctrl.ClearCart)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/cart", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestClearCartHandler_ServiceError(t *testing.T) {
	svc := &mockCartService{
		clearCartFn: func(_ context.Context, _ uuid.UUID) error {
			return errors.New("db error")
		},
	}
	ctrl := NewCartController(svc)
	r := gin.New()
	r.DELETE("/cart", withUserID(uuid.New()), ctrl.ClearCart)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/cart", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}
