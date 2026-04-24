package product

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
)

// mockProductRepo implements IProductRepository.
type mockProductRepo struct {
	createProductFn  func(ctx context.Context, req *CreateProductRequest) error
	getProductFn     func(ctx context.Context) ([]Product, error)
	getProductByIDFn func(ctx context.Context, id uuid.UUID) (*ProductDetailResponse, error)
	deleteProductFn  func(ctx context.Context, id *string) error
	updateProductFn  func(ctx context.Context, req *UpdateProductRequest, id uuid.UUID) error
	likeProductFn    func(ctx context.Context, req *LikeProductRequest) error
}

func (m *mockProductRepo) CreateProduct(ctx context.Context, req *CreateProductRequest) error {
	if m.createProductFn != nil {
		return m.createProductFn(ctx, req)
	}
	return nil
}
func (m *mockProductRepo) GetProduct(ctx context.Context) ([]Product, error) {
	if m.getProductFn != nil {
		return m.getProductFn(ctx)
	}
	return []Product{}, nil
}
func (m *mockProductRepo) GetProductByID(ctx context.Context, id uuid.UUID) (*ProductDetailResponse, error) {
	if m.getProductByIDFn != nil {
		return m.getProductByIDFn(ctx, id)
	}
	return nil, nil
}
func (m *mockProductRepo) DeleteProduct(ctx context.Context, id *string) error {
	if m.deleteProductFn != nil {
		return m.deleteProductFn(ctx, id)
	}
	return nil
}
func (m *mockProductRepo) UpdateProduct(ctx context.Context, req *UpdateProductRequest, id uuid.UUID) error {
	if m.updateProductFn != nil {
		return m.updateProductFn(ctx, req, id)
	}
	return nil
}
func (m *mockProductRepo) LikeProduct(ctx context.Context, req *LikeProductRequest) error {
	if m.likeProductFn != nil {
		return m.likeProductFn(ctx, req)
	}
	return nil
}

func newTestProductService(repo IProductRepository) *ProductService {
	return NewProductService(repo)
}

// --- CreateSepatu ---

func TestCreateSepatu_Success(t *testing.T) {
	svc := newTestProductService(&mockProductRepo{})
	err := svc.CreateSepatu(context.Background(), &CreateProductRequest{
		Name: "Nike Air", Brand: "Nike", Size: 42, Price: 1500000, Stock: 10,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateSepatu_RepoError(t *testing.T) {
	repo := &mockProductRepo{
		createProductFn: func(_ context.Context, _ *CreateProductRequest) error {
			return errors.New("db error")
		},
	}
	svc := newTestProductService(repo)
	err := svc.CreateSepatu(context.Background(), &CreateProductRequest{
		Name: "Nike", Brand: "Nike", Size: 42, Price: 100, Stock: 1,
	})
	if err == nil {
		t.Error("expected error to propagate from repo")
	}
}

// --- GetSepatu ---

func TestGetSepatu_Success(t *testing.T) {
	expected := []Product{{Name: "Nike Air"}, {Name: "Adidas"}}
	repo := &mockProductRepo{
		getProductFn: func(_ context.Context) ([]Product, error) {
			return expected, nil
		},
	}
	svc := newTestProductService(repo)
	products, err := svc.GetSepatu(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(products) != 2 {
		t.Errorf("expected 2 products, got %d", len(products))
	}
}

func TestGetSepatu_RepoError(t *testing.T) {
	repo := &mockProductRepo{
		getProductFn: func(_ context.Context) ([]Product, error) {
			return nil, errors.New("db error")
		},
	}
	svc := newTestProductService(repo)
	_, err := svc.GetSepatu(context.Background())
	if err == nil {
		t.Error("expected error from repo")
	}
}

// --- GetSepatuByID ---

func TestGetSepatuByID_Success(t *testing.T) {
	id := uuid.New()
	detail := &ProductDetailResponse{ID: id.String(), Name: "Nike Air"}
	repo := &mockProductRepo{
		getProductByIDFn: func(_ context.Context, gotID uuid.UUID) (*ProductDetailResponse, error) {
			if gotID != id {
				t.Errorf("expected id %v, got %v", id, gotID)
			}
			return detail, nil
		},
	}
	svc := newTestProductService(repo)
	got, err := svc.GetSepatuByID(context.Background(), id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "Nike Air" {
		t.Errorf("expected name Nike Air, got %q", got.Name)
	}
}

func TestGetSepatuByID_NotFound(t *testing.T) {
	repo := &mockProductRepo{
		getProductByIDFn: func(_ context.Context, _ uuid.UUID) (*ProductDetailResponse, error) {
			return nil, errors.New("not found")
		},
	}
	svc := newTestProductService(repo)
	_, err := svc.GetSepatuByID(context.Background(), uuid.New())
	if err == nil {
		t.Error("expected error for not found product")
	}
}

// --- DeleteSepatu ---

func TestDeleteSepatu_Success(t *testing.T) {
	called := false
	repo := &mockProductRepo{
		deleteProductFn: func(_ context.Context, _ *string) error {
			called = true
			return nil
		},
	}
	svc := newTestProductService(repo)
	id := "some-id"
	err := svc.DeleteSepatu(context.Background(), &id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected repo.DeleteProduct to be called")
	}
}

// --- UpdateSepatu ---

func TestUpdateSepatu_Success(t *testing.T) {
	called := false
	repo := &mockProductRepo{
		updateProductFn: func(_ context.Context, _ *UpdateProductRequest, _ uuid.UUID) error {
			called = true
			return nil
		},
	}
	svc := newTestProductService(repo)
	err := svc.UpdateSepatu(context.Background(), &UpdateProductRequest{}, uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected repo.UpdateProduct to be called")
	}
}

// --- LikeProduct ---

func TestLikeProduct_Success(t *testing.T) {
	called := false
	repo := &mockProductRepo{
		likeProductFn: func(_ context.Context, _ *LikeProductRequest) error {
			called = true
			return nil
		},
	}
	svc := newTestProductService(repo)
	err := svc.LikeProduct(context.Background(), &LikeProductRequest{ID: uuid.New()})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected repo.LikeProduct to be called")
	}
}

func TestLikeProduct_RepoError(t *testing.T) {
	repo := &mockProductRepo{
		likeProductFn: func(_ context.Context, _ *LikeProductRequest) error {
			return errors.New("db error")
		},
	}
	svc := newTestProductService(repo)
	err := svc.LikeProduct(context.Background(), &LikeProductRequest{})
	if err == nil {
		t.Error("expected error to propagate")
	}
}
