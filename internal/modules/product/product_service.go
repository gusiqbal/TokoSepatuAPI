package product

import (
	"context"

	"github.com/google/uuid"
)

type IProductService interface {
	CreateSepatu(ctx context.Context, sepatus *CreateProductRequest) error
	GetSepatu(ctx context.Context) ([]Product, error)
	GetSepatuByID(ctx context.Context, id uuid.UUID) (*ProductDetailResponse, error)
	DeleteSepatu(ctx context.Context, id *string) error
	UpdateSepatu(ctx context.Context, sepatu *UpdateProductRequest, id uuid.UUID) error
}

type ProductService struct {
	repo *ProductRepoSitory
}

func NewProductService(repo *ProductRepoSitory) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func (s *ProductService) CreateSepatu(ctx context.Context, sepatus *CreateProductRequest) error {
	return s.repo.CreateProduct(ctx, sepatus)
}

func (s *ProductService) GetSepatu(ctx context.Context) ([]Product, error) {
	return s.repo.GetProduct(ctx)
}

func (s *ProductService) GetSepatuByID(ctx context.Context, id uuid.UUID) (*ProductDetailResponse, error) {
	return s.repo.GetProductByID(ctx, id)
}

func (s *ProductService) DeleteSepatu(ctx context.Context, id *string) error {
	return s.repo.DeleteProduct(ctx, id)
}

func (s *ProductService) UpdateSepatu(ctx context.Context, sepatu *UpdateProductRequest, id uuid.UUID) error {
	return s.repo.UpdateProduct(ctx, sepatu, id)
}

func (s *ProductService) LikeProduct(ctx context.Context, req *LikeProductRequest) error {
	return s.repo.LikeProduct(ctx, req)
}
