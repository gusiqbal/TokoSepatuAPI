package product

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

// mockProductService implements IProductService for handler tests.
type mockProductService struct {
	createSepatuFn  func(ctx context.Context, req *CreateProductRequest) error
	getSepatuFn     func(ctx context.Context) ([]Product, error)
	getSepatuByIDFn func(ctx context.Context, id uuid.UUID) (*ProductDetailResponse, error)
	deleteSepatuFn  func(ctx context.Context, id *string) error
	updateSepatuFn  func(ctx context.Context, req *UpdateProductRequest, id uuid.UUID) error
	likeProductFn   func(ctx context.Context, req *LikeProductRequest) error
}

func (m *mockProductService) CreateSepatu(ctx context.Context, req *CreateProductRequest) error {
	if m.createSepatuFn != nil {
		return m.createSepatuFn(ctx, req)
	}
	return nil
}
func (m *mockProductService) GetSepatu(ctx context.Context) ([]Product, error) {
	if m.getSepatuFn != nil {
		return m.getSepatuFn(ctx)
	}
	return []Product{}, nil
}
func (m *mockProductService) GetSepatuByID(ctx context.Context, id uuid.UUID) (*ProductDetailResponse, error) {
	if m.getSepatuByIDFn != nil {
		return m.getSepatuByIDFn(ctx, id)
	}
	return nil, nil
}
func (m *mockProductService) DeleteSepatu(ctx context.Context, id *string) error {
	if m.deleteSepatuFn != nil {
		return m.deleteSepatuFn(ctx, id)
	}
	return nil
}
func (m *mockProductService) UpdateSepatu(ctx context.Context, req *UpdateProductRequest, id uuid.UUID) error {
	if m.updateSepatuFn != nil {
		return m.updateSepatuFn(ctx, req, id)
	}
	return nil
}
func (m *mockProductService) LikeProduct(ctx context.Context, req *LikeProductRequest) error {
	if m.likeProductFn != nil {
		return m.likeProductFn(ctx, req)
	}
	return nil
}

func postProductJSON(r *gin.Engine, path string, body any) *httptest.ResponseRecorder {
	b, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, path, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	return w
}

// --- CreateSepatu handler ---

func TestCreateSepatuHandler_Success(t *testing.T) {
	ctrl := NewProductController(&mockProductService{})
	r := gin.New()
	r.POST("/sepatu", ctrl.CreateSepatu)

	w := postProductJSON(r, "/sepatu", gin.H{
		"name": "Nike Air", "brand": "Nike", "size": 42, "price": 1500000, "stock": 10,
	})
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateSepatuHandler_InvalidJSON(t *testing.T) {
	ctrl := NewProductController(&mockProductService{})
	r := gin.New()
	r.POST("/sepatu", ctrl.CreateSepatu)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/sepatu", bytes.NewBufferString("notjson"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestCreateSepatuHandler_ServiceError(t *testing.T) {
	svc := &mockProductService{
		createSepatuFn: func(_ context.Context, _ *CreateProductRequest) error {
			return errors.New("db error")
		},
	}
	ctrl := NewProductController(svc)
	r := gin.New()
	r.POST("/sepatu", ctrl.CreateSepatu)

	w := postProductJSON(r, "/sepatu", gin.H{
		"name": "Nike", "brand": "Nike", "size": 42, "price": 100, "stock": 1,
	})
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

// --- GetSepatu handler ---

func TestGetSepatuHandler_Success(t *testing.T) {
	svc := &mockProductService{
		getSepatuFn: func(_ context.Context) ([]Product, error) {
			return []Product{{Name: "Nike"}, {Name: "Adidas"}}, nil
		},
	}
	ctrl := NewProductController(svc)
	r := gin.New()
	r.GET("/sepatu", ctrl.GetSepatu)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/sepatu", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestGetSepatuHandler_ServiceError(t *testing.T) {
	svc := &mockProductService{
		getSepatuFn: func(_ context.Context) ([]Product, error) {
			return nil, errors.New("db error")
		},
	}
	ctrl := NewProductController(svc)
	r := gin.New()
	r.GET("/sepatu", ctrl.GetSepatu)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/sepatu", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

// --- GetSepatuByID handler ---

func TestGetSepatuByIDHandler_ValidID(t *testing.T) {
	id := uuid.New()
	svc := &mockProductService{
		getSepatuByIDFn: func(_ context.Context, gotID uuid.UUID) (*ProductDetailResponse, error) {
			return &ProductDetailResponse{ID: gotID.String(), Name: "Nike"}, nil
		},
	}
	ctrl := NewProductController(svc)
	r := gin.New()
	r.GET("/sepatu/:id", ctrl.GetSepatuByID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/sepatu/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetSepatuByIDHandler_InvalidID(t *testing.T) {
	ctrl := NewProductController(&mockProductService{})
	r := gin.New()
	r.GET("/sepatu/:id", ctrl.GetSepatuByID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/sepatu/not-a-uuid", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestGetSepatuByIDHandler_NotFound(t *testing.T) {
	svc := &mockProductService{
		getSepatuByIDFn: func(_ context.Context, _ uuid.UUID) (*ProductDetailResponse, error) {
			return nil, errors.New("not found")
		},
	}
	ctrl := NewProductController(svc)
	r := gin.New()
	r.GET("/sepatu/:id", ctrl.GetSepatuByID)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/sepatu/"+uuid.New().String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

// --- UpdateSepatu handler ---

func TestUpdateSepatuHandler_Success(t *testing.T) {
	id := uuid.New().String()
	ctrl := NewProductController(&mockProductService{})
	r := gin.New()
	r.PUT("/sepatu/:id", ctrl.UpdateSepatu)

	b, _ := json.Marshal(gin.H{"id": id, "name": "Updated"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/sepatu/"+id, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestUpdateSepatuHandler_InvalidJSON(t *testing.T) {
	ctrl := NewProductController(&mockProductService{})
	r := gin.New()
	r.PUT("/sepatu/:id", ctrl.UpdateSepatu)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/sepatu/"+uuid.New().String(), bytes.NewBufferString("bad"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// --- LikeProduct handler ---

func TestLikeProductHandler_Success(t *testing.T) {
	ctrl := NewProductController(&mockProductService{})
	r := gin.New()
	r.POST("/sepatu/like", ctrl.LikeProduct)

	w := postProductJSON(r, "/sepatu/like", gin.H{"productId": uuid.New().String()})
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestLikeProductHandler_ServiceError(t *testing.T) {
	svc := &mockProductService{
		likeProductFn: func(_ context.Context, _ *LikeProductRequest) error {
			return errors.New("db error")
		},
	}
	ctrl := NewProductController(svc)
	r := gin.New()
	r.POST("/sepatu/like", ctrl.LikeProduct)

	w := postProductJSON(r, "/sepatu/like", gin.H{"productId": uuid.New().String()})
	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}
