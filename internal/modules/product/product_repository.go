package product

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IProductRepository interface {
	CreateProduct(ctx context.Context, sepatu *CreateProductRequest) error
	GetProduct(ctx context.Context) ([]Product, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (*ProductDetailResponse, error)
	DeleteProduct(ctx context.Context, id *string) error
	UpdateProduct(ctx context.Context, sepatuUpdate *UpdateProductRequest, id uuid.UUID) error
}

type ProductRepoSitory struct {
	repo IProductRepository
	db   *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepoSitory {
	return &ProductRepoSitory{
		db: db,
	}
}

func (r *ProductRepoSitory) CreateProduct(ctx context.Context, request *CreateProductRequest) error {
	now := time.Now().Unix()

	product := Product{
		ID:            uuid.New(),
		Name:          request.Name,
		Brand:         request.Brand,
		Price:         request.Price,
		Stock:         request.Stock,
		Size:          request.Size,
		CreatedAt:     now,
		LastUpdatedAt: now,
	}

	return r.db.WithContext(ctx).Create(&product).Error
}

func (r *ProductRepoSitory) GetProduct(ctx context.Context) ([]Product, error) {
	var sepatus []Product
	err := r.db.WithContext(ctx).Find(&sepatus).Error
	return sepatus, err
}

func (r *ProductRepoSitory) DeleteProduct(ctx context.Context, id *string) error {
	return r.db.WithContext(ctx).Delete(&Product{}, "id = ?", id).Error
}

func (r *ProductRepoSitory) UpdateProduct(ctx context.Context, sepatuUpdate *UpdateProductRequest, id uuid.UUID) error {
	var sepatus Product
	if err := r.db.WithContext(ctx).First(&sepatus, "id = ?", id).Error; err != nil {
		return err
	}

	if sepatuUpdate.Name != nil {
		sepatus.Name = *sepatuUpdate.Name
	}
	if sepatuUpdate.Brand != nil {
		sepatus.Brand = *sepatuUpdate.Brand
	}
	if sepatuUpdate.Price != nil {
		sepatus.Price = *sepatuUpdate.Price
	}
	if sepatuUpdate.Size != nil {
		sepatus.Size = *sepatuUpdate.Size
	}
	if sepatuUpdate.Stock != nil {
		sepatus.Stock = *sepatuUpdate.Stock
	}

	sepatus.LastUpdatedAt = time.Now().Unix()

	return r.db.Save(&sepatus).Error
}

func (r *ProductRepoSitory) GetProductByID(ctx context.Context, id uuid.UUID) (*ProductDetailResponse, error) {
	var product Product
	if err := r.db.WithContext(ctx).First(&product, "id = ?", id).Error; err != nil {
		return nil, err
	}

	var variants []ProductVariant
	r.db.WithContext(ctx).Where("product_id = ?", id).Find(&variants)

	variantResponses := make([]ProductVariantResponse, len(variants))
	for i, v := range variants {
		variantResponses[i] = ProductVariantResponse{
			ID:    v.ID.String(),
			Size:  v.Size,
			Color: v.Color,
			Stock: v.Stock,
		}
	}

	return &ProductDetailResponse{
		ID:        product.ID.String(),
		Name:      product.Name,
		Brand:     product.Brand,
		Size:      product.Size,
		Price:     product.Price,
		Stock:     product.Stock,
		CreatedAt: product.CreatedAt,
		Variants:  variantResponses,
	}, nil
}

func (r *ProductRepoSitory) LikeProduct(ctx context.Context, req *LikeProductRequest) error {
	var likedProduct ProductFavorite

	likedProduct.ID = uuid.New()
	likedProduct.ProductID = req.ID

	return r.db.WithContext(ctx).Save(&likedProduct).Error
}
