package cart

import "github.com/google/uuid"

type AddToCartRequest struct {
	ProductVariantID uuid.UUID `json:"productVariantId" binding:"required"`
	Quantity         int       `json:"quantity" binding:"required,gt=0"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" binding:"required,gt=0"`
}

type CartItemResponse struct {
	ID               string  `json:"id"`
	ProductVariantID string  `json:"productVariantId"`
	ProductName      string  `json:"productName"`
	Brand            string  `json:"brand"`
	Color            string  `json:"color"`
	Size             int     `json:"size"`
	Price            float64 `json:"price"`
	Quantity         int     `json:"quantity"`
	Subtotal         float64 `json:"subtotal"`
}

type CartResponse struct {
	ID         string             `json:"id"`
	Items      []CartItemResponse `json:"items"`
	TotalPrice float64            `json:"totalPrice"`
}
