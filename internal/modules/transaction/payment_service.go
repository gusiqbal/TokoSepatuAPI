package transaction

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"learnapirest/internal/modules/order"
	"net/http"

	stripe "github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/paymentintent"
	"github.com/stripe/stripe-go/v82/webhook"

	"github.com/google/uuid"
)

type IPaymentService interface {
	CreateCheckout(ctx context.Context, userID uuid.UUID, orderID uuid.UUID) (*CheckoutResponse, error)
	HandleWebhook(r *http.Request) error
	GetPaymentStatus(ctx context.Context, userID uuid.UUID, orderID uuid.UUID) (*PaymentStatusResponse, error)
}

type PaymentService struct {
	paymentRepo   IPaymentRepository
	orderRepo     order.IOrderRepository
	stripeKey     string
	webhookSecret string
}

func NewPaymentService(
	paymentRepo IPaymentRepository,
	orderRepo order.IOrderRepository,
	stripeKey string,
	webhookSecret string,
) *PaymentService {
	return &PaymentService{
		paymentRepo:   paymentRepo,
		orderRepo:     orderRepo,
		stripeKey:     stripeKey,
		webhookSecret: webhookSecret,
	}
}

func (s *PaymentService) CreateCheckout(ctx context.Context, userID uuid.UUID, orderID uuid.UUID) (*CheckoutResponse, error) {
	// Verify order belongs to user and get amount
	orderResp, err := s.orderRepo.GetOrderByID(ctx, orderID, userID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	// Prevent duplicate payment
	existing, _ := s.paymentRepo.GetPaymentByOrderID(ctx, orderID)
	if existing != nil && existing.Status == "succeeded" {
		return nil, errors.New("order already paid")
	}

	amountInCents := int64(orderResp.TotalPrice * 100)
	currency := "usd"

	stripe.Key = s.stripeKey
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amountInCents),
		Currency: stripe.String(currency),
		Metadata: map[string]string{
			"order_id": orderID.String(),
			"user_id":  userID.String(),
		},
	}

	pi, err := paymentintent.New(params)
	if err != nil {
		return nil, errors.New("failed to create payment intent: " + err.Error())
	}

	payment := &Payment{
		ID:                    uuid.New(),
		OrderID:               orderID,
		StripePaymentIntentID: pi.ID,
		Amount:                amountInCents,
		Currency:              currency,
		Status:                "pending",
	}

	if existing != nil {
		// Update existing payment record if it was previously failed
		_ = s.paymentRepo.UpdatePaymentStatus(ctx, existing.StripePaymentIntentID, "cancelled")
	}

	if err := s.paymentRepo.CreatePayment(ctx, payment); err != nil {
		return nil, err
	}

	return &CheckoutResponse{
		ClientSecret:    pi.ClientSecret,
		PaymentIntentID: pi.ID,
		Amount:          amountInCents,
		Currency:        currency,
	}, nil
}

func (s *PaymentService) HandleWebhook(r *http.Request) error {
	const maxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(nil, r.Body, maxBodyBytes)

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		return errors.New("failed to read request body")
	}

	sig := r.Header.Get("Stripe-Signature")
	event, err := webhook.ConstructEvent(payload, sig, s.webhookSecret)
	if err != nil {
		return errors.New("webhook signature verification failed: " + err.Error())
	}

	switch event.Type {
	case "payment_intent.succeeded":
		var pi stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
			return err
		}
		_ = s.paymentRepo.UpdatePaymentStatus(r.Context(), pi.ID, "succeeded")
		if orderIDStr, ok := pi.Metadata["order_id"]; ok {
			if orderID, err := uuid.Parse(orderIDStr); err == nil {
				_ = s.orderRepo.UpdateOrderStatus(r.Context(), orderID, "paid")
			}
		}

	case "payment_intent.payment_failed":
		var pi stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
			return err
		}
		_ = s.paymentRepo.UpdatePaymentStatus(r.Context(), pi.ID, "failed")
		if orderIDStr, ok := pi.Metadata["order_id"]; ok {
			if orderID, err := uuid.Parse(orderIDStr); err == nil {
				_ = s.orderRepo.UpdateOrderStatus(r.Context(), orderID, "payment_failed")
			}
		}
	}

	return nil
}

func (s *PaymentService) GetPaymentStatus(ctx context.Context, userID uuid.UUID, orderID uuid.UUID) (*PaymentStatusResponse, error) {
	// Verify order belongs to user
	if _, err := s.orderRepo.GetOrderByID(ctx, orderID, userID); err != nil {
		return nil, errors.New("order not found")
	}

	payment, err := s.paymentRepo.GetPaymentByOrderID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	return &PaymentStatusResponse{
		OrderID:               orderID.String(),
		StripePaymentIntentID: payment.StripePaymentIntentID,
		Amount:                payment.Amount,
		Currency:              payment.Currency,
		Status:                payment.Status,
	}, nil
}
