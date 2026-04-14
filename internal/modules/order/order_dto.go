package order

type CreateOrderRequest struct {
	PaymentMethod   string `json:"paymentMethod" binding:"required"`
	ShippingAddress string `json:"shippingAddress" binding:"required"`
}

type OrderItemResponse struct {
	ID               string  `json:"id"`
	ProductVariantID string  `json:"productVariantId"`
	ProductName      string  `json:"productName"`
	Brand            string  `json:"brand"`
	Color            string  `json:"color"`
	Size             int     `json:"size"`
	Quantity         int     `json:"quantity"`
	PriceAtPurchase  float64 `json:"priceAtPurchase"`
	Subtotal         float64 `json:"subtotal"`
}

type OrderResponse struct {
	ID              string              `json:"id"`
	Status          string              `json:"status"`
	TotalPrice      float64             `json:"totalPrice"`
	PaymentMethod   string              `json:"paymentMethod"`
	ShippingAddress string              `json:"shippingAddress"`
	Items           []OrderItemResponse `json:"items"`
	CreatedAt       int64               `json:"createdAt"`
}
