package product

import "github.com/google/uuid"

type CreateProductRequest struct {
	Name  string  `json:"name" binding:"required"`
	Brand string  `json:"brand" binding:"required"`
	Size  int     `json:"size" binding:"required,gt=0"`
	Price float64 `json:"price" binding:"required,gt=0"`
	Stock int     `json:"stock" binding:"required,gte=0"`
}

type UpdateProductRequest struct {
	ID    *string  `json:"id"`
	Name  *string  `json:"name"`
	Brand *string  `json:"brand"`
	Size  *int     `json:"size" binding:"omitempty,gt=0"`
	Price *float64 `json:"price" binding:"omitempty,gt=0"`
	Stock *int     `json:"stock" binding:"omitempty,gte=0"`
}

type ProductResponse struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Brand string  `json:"brand"`
	Price float64 `json:"price"`
}

type LikeProductRequest struct {
	ID uuid.UUID `json:"productId"`
}
