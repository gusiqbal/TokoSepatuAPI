package cart

import (
	"context"
	"errors"
	"learnapirest/internal/modules/product"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ICartRepository interface {
	GetOrCreateCart(ctx context.Context, userID uuid.UUID) (*Cart, error)
	GetCartWithItems(ctx context.Context, userID uuid.UUID) (*CartResponse, error)
	AddItem(ctx context.Context, cartID uuid.UUID, variantID uuid.UUID, quantity int) error
	UpdateItem(ctx context.Context, cartItemID uuid.UUID, quantity int) error
	RemoveItem(ctx context.Context, cartItemID uuid.UUID) error
	ClearCart(ctx context.Context, cartID uuid.UUID) error
	GetRawCartItems(ctx context.Context, userID uuid.UUID) (uuid.UUID, []CartItem, error)
}

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{db: db}
}

func (r *CartRepository) GetOrCreateCart(ctx context.Context, userID uuid.UUID) (*Cart, error) {
	var cart Cart
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&cart).Error
	if err == nil {
		return &cart, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	cart = Cart{
		ID:     uuid.New(),
		UserID: userID,
	}
	if err := r.db.WithContext(ctx).Create(&cart).Error; err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *CartRepository) GetCartWithItems(ctx context.Context, userID uuid.UUID) (*CartResponse, error) {
	cart, err := r.GetOrCreateCart(ctx, userID)
	if err != nil {
		return nil, err
	}

	var items []CartItem
	r.db.WithContext(ctx).
		Preload("ProductVariant").
		Preload("ProductVariant.Product").
		Where("cart_id = ?", cart.ID).
		Find(&items)

	itemResponses := make([]CartItemResponse, 0, len(items))
	var total float64

	for _, item := range items {
		v := item.ProductVariant
		subtotal := v.Product.Price * float64(item.Quantity)
		total += subtotal

		itemResponses = append(itemResponses, CartItemResponse{
			ID:               item.ID.String(),
			ProductVariantID: v.ID.String(),
			ProductName:      v.Product.Name,
			Brand:            v.Product.Brand,
			Color:            v.Color,
			Size:             v.Size,
			Price:            v.Product.Price,
			Quantity:         item.Quantity,
			Subtotal:         subtotal,
		})
	}

	return &CartResponse{
		ID:         cart.ID.String(),
		Items:      itemResponses,
		TotalPrice: total,
	}, nil
}

func (r *CartRepository) AddItem(ctx context.Context, cartID uuid.UUID, variantID uuid.UUID, quantity int) error {
	// Check if the variant exists
	var variant product.ProductVariant
	if err := r.db.WithContext(ctx).First(&variant, "id = ?", variantID).Error; err != nil {
		return errors.New("product variant not found")
	}
	if variant.Stock < quantity {
		return errors.New("insufficient stock")
	}

	// Check if item already in cart — increment quantity if so
	var existing CartItem
	err := r.db.WithContext(ctx).
		Where("cart_id = ? AND product_variant_id = ?", cartID, variantID).
		First(&existing).Error

	if err == nil {
		newQty := existing.Quantity + quantity
		if variant.Stock < newQty {
			return errors.New("insufficient stock")
		}
		return r.db.WithContext(ctx).Model(&existing).Update("quantity", newQty).Error
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	item := CartItem{
		ID:               uuid.New(),
		CartID:           cartID,
		ProductVariantID: variantID,
		Quantity:         quantity,
	}
	return r.db.WithContext(ctx).Create(&item).Error
}

func (r *CartRepository) UpdateItem(ctx context.Context, cartItemID uuid.UUID, quantity int) error {
	var item CartItem
	if err := r.db.WithContext(ctx).Preload("ProductVariant").First(&item, "id = ?", cartItemID).Error; err != nil {
		return errors.New("cart item not found")
	}
	if item.ProductVariant.Stock < quantity {
		return errors.New("insufficient stock")
	}
	return r.db.WithContext(ctx).Model(&item).Update("quantity", quantity).Error
}

func (r *CartRepository) RemoveItem(ctx context.Context, cartItemID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&CartItem{}, "id = ?", cartItemID).Error
}

func (r *CartRepository) ClearCart(ctx context.Context, cartID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&CartItem{}, "cart_id = ?", cartID).Error
}

// GetRawCartItems returns populated CartItems for a given userID — used by the order service during checkout.
func (r *CartRepository) GetRawCartItems(ctx context.Context, userID uuid.UUID) (uuid.UUID, []CartItem, error) {
	cart, err := r.GetOrCreateCart(ctx, userID)
	if err != nil {
		return uuid.Nil, nil, err
	}

	var items []CartItem
	if err := r.db.WithContext(ctx).
		Preload("ProductVariant").
		Preload("ProductVariant.Product").
		Where("cart_id = ?", cart.ID).
		Find(&items).Error; err != nil {
		return uuid.Nil, nil, err
	}

	return cart.ID, items, nil
}
