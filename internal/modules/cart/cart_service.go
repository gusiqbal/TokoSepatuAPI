package cart

import (
	"context"

	"github.com/google/uuid"
)

type ICartService interface {
	GetCart(ctx context.Context, userID uuid.UUID) (*CartResponse, error)
	AddItem(ctx context.Context, userID uuid.UUID, req *AddToCartRequest) error
	UpdateItem(ctx context.Context, cartItemID uuid.UUID, req *UpdateCartItemRequest) error
	RemoveItem(ctx context.Context, cartItemID uuid.UUID) error
	ClearCart(ctx context.Context, userID uuid.UUID) error
}

type CartService struct {
	repo ICartRepository
}

func NewCartService(repo ICartRepository) *CartService {
	return &CartService{repo: repo}
}

func (s *CartService) GetCart(ctx context.Context, userID uuid.UUID) (*CartResponse, error) {
	return s.repo.GetCartWithItems(ctx, userID)
}

func (s *CartService) AddItem(ctx context.Context, userID uuid.UUID, req *AddToCartRequest) error {
	cart, err := s.repo.GetOrCreateCart(ctx, userID)
	if err != nil {
		return err
	}
	return s.repo.AddItem(ctx, cart.ID, req.ProductVariantID, req.Quantity)
}

func (s *CartService) UpdateItem(ctx context.Context, cartItemID uuid.UUID, req *UpdateCartItemRequest) error {
	return s.repo.UpdateItem(ctx, cartItemID, req.Quantity)
}

func (s *CartService) RemoveItem(ctx context.Context, cartItemID uuid.UUID) error {
	return s.repo.RemoveItem(ctx, cartItemID)
}

func (s *CartService) ClearCart(ctx context.Context, userID uuid.UUID) error {
	cart, err := s.repo.GetOrCreateCart(ctx, userID)
	if err != nil {
		return err
	}
	return s.repo.ClearCart(ctx, cart.ID)
}
