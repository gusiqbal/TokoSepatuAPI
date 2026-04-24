package order

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

// mockOrderService implements IOrderService for handler tests.
type mockOrderService struct {
	createOrderFromCartFn func(ctx context.Context, userID uuid.UUID, req *CreateOrderRequest) (*OrderResponse, error)
	getOrderHistoryFn     func(ctx context.Context, userID uuid.UUID) ([]OrderResponse, error)
	getOrderDetailFn      func(ctx context.Context, userID, orderID uuid.UUID) (*OrderResponse, error)
}

func (m *mockOrderService) CreateOrderFromCart(ctx context.Context, userID uuid.UUID, req *CreateOrderRequest) (*OrderResponse, error) {
	if m.createOrderFromCartFn != nil {
		return m.createOrderFromCartFn(ctx, userID, req)
	}
	return &OrderResponse{}, nil
}
func (m *mockOrderService) GetOrderHistory(ctx context.Context, userID uuid.UUID) ([]OrderResponse, error) {
	if m.getOrderHistoryFn != nil {
		return m.getOrderHistoryFn(ctx, userID)
	}
	return []OrderResponse{}, nil
}
func (m *mockOrderService) GetOrderDetail(ctx context.Context, userID, orderID uuid.UUID) (*OrderResponse, error) {
	if m.getOrderDetailFn != nil {
		return m.getOrderDetailFn(ctx, userID, orderID)
	}
	return &OrderResponse{}, nil
}

func setUserID(userID uuid.UUID) gin.HandlerFunc {
	return func(c *gin.Context) { c.Set("userId", userID) }
}

// --- CreateOrder handler ---

func TestCreateOrderHandler_Unauthorized(t *testing.T) {
	ctrl := NewOrderController(&mockOrderService{})
	r := gin.New()
	r.POST("/orders", ctrl.CreateOrder)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/orders", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestCreateOrderHandler_InvalidJSON(t *testing.T) {
	ctrl := NewOrderController(&mockOrderService{})
	r := gin.New()
	r.POST("/orders", setUserID(uuid.New()), ctrl.CreateOrder)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/orders", bytes.NewBufferString("bad"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreateOrderHandler_Success(t *testing.T) {
	orderID := uuid.New()
	svc := &mockOrderService{
		createOrderFromCartFn: func(_ context.Context, _ uuid.UUID, _ *CreateOrderRequest) (*OrderResponse, error) {
			return &OrderResponse{ID: orderID.String(), Status: "pending"}, nil
		},
	}
	ctrl := NewOrderController(svc)
	r := gin.New()
	r.POST("/orders", setUserID(uuid.New()), ctrl.CreateOrder)

	b, _ := json.Marshal(gin.H{"paymentMethod": "transfer", "shippingAddress": "Jl. Test 1"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateOrderHandler_ServiceError(t *testing.T) {
	svc := &mockOrderService{
		createOrderFromCartFn: func(_ context.Context, _ uuid.UUID, _ *CreateOrderRequest) (*OrderResponse, error) {
			return nil, errors.New("cart is empty")
		},
	}
	ctrl := NewOrderController(svc)
	r := gin.New()
	r.POST("/orders", setUserID(uuid.New()), ctrl.CreateOrder)

	b, _ := json.Marshal(gin.H{"paymentMethod": "cash", "shippingAddress": "Jl. Test"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// --- GetOrderHistory handler ---

func TestGetOrderHistoryHandler_Unauthorized(t *testing.T) {
	ctrl := NewOrderController(&mockOrderService{})
	r := gin.New()
	r.GET("/orders", ctrl.GetOrderHistory)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/orders", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestGetOrderHistoryHandler_Success(t *testing.T) {
	svc := &mockOrderService{
		getOrderHistoryFn: func(_ context.Context, _ uuid.UUID) ([]OrderResponse, error) {
			return []OrderResponse{{ID: "order-1"}, {ID: "order-2"}}, nil
		},
	}
	ctrl := NewOrderController(svc)
	r := gin.New()
	r.GET("/orders", setUserID(uuid.New()), ctrl.GetOrderHistory)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/orders", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetOrderHistoryHandler_ServiceError(t *testing.T) {
	svc := &mockOrderService{
		getOrderHistoryFn: func(_ context.Context, _ uuid.UUID) ([]OrderResponse, error) {
			return nil, errors.New("db error")
		},
	}
	ctrl := NewOrderController(svc)
	r := gin.New()
	r.GET("/orders", setUserID(uuid.New()), ctrl.GetOrderHistory)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/orders", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

// --- GetOrderDetail handler ---

func TestGetOrderDetailHandler_Unauthorized(t *testing.T) {
	ctrl := NewOrderController(&mockOrderService{})
	r := gin.New()
	r.GET("/orders/:id", ctrl.GetOrderDetail)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/orders/"+uuid.New().String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestGetOrderDetailHandler_InvalidID(t *testing.T) {
	ctrl := NewOrderController(&mockOrderService{})
	r := gin.New()
	r.GET("/orders/:id", setUserID(uuid.New()), ctrl.GetOrderDetail)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/orders/not-a-uuid", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestGetOrderDetailHandler_NotFound(t *testing.T) {
	svc := &mockOrderService{
		getOrderDetailFn: func(_ context.Context, _, _ uuid.UUID) (*OrderResponse, error) {
			return nil, errors.New("order not found")
		},
	}
	ctrl := NewOrderController(svc)
	r := gin.New()
	r.GET("/orders/:id", setUserID(uuid.New()), ctrl.GetOrderDetail)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/orders/"+uuid.New().String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestGetOrderDetailHandler_Success(t *testing.T) {
	orderID := uuid.New()
	svc := &mockOrderService{
		getOrderDetailFn: func(_ context.Context, _, oID uuid.UUID) (*OrderResponse, error) {
			return &OrderResponse{ID: oID.String(), Status: "pending"}, nil
		},
	}
	ctrl := NewOrderController(svc)
	r := gin.New()
	r.GET("/orders/:id", setUserID(uuid.New()), ctrl.GetOrderDetail)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/orders/"+orderID.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}
