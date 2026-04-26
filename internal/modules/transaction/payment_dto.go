package transaction

type CreateCheckoutRequest struct {
	OrderID string `json:"orderId" binding:"required,uuid"`
}

type CheckoutResponse struct {
	ClientSecret          string `json:"clientSecret"`
	PaymentIntentID       string `json:"paymentIntentId"`
	Amount                int64  `json:"amount"`
	Currency              string `json:"currency"`
}

type PaymentStatusResponse struct {
	OrderID               string `json:"orderId"`
	StripePaymentIntentID string `json:"stripePaymentIntentId"`
	Amount                int64  `json:"amount"`
	Currency              string `json:"currency"`
	Status                string `json:"status"`
}
