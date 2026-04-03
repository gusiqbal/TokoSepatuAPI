package repository

import (
	"context"
	"learnapirest/model"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ISepatuRepository interface {
	CreateSepatu(ctx context.Context, sepatu *model.Sepatu) error
	GetAllSepatu(ctx context.Context) ([]model.Sepatu, error)
	DeleteSepatu(ctx context.Context, id uuid.UUID) error
	UpdateSepatuByID(ctx context.Context, sepatuUpdate *model.UpdateSepatu, id uuid.UUID) error
}

type SepatuRepoSitory struct {
	db *gorm.DB
}

func NewSepatuRepo(db *gorm.DB) *SepatuRepoSitory {
	return &SepatuRepoSitory{
		db: db,
	}
}

func (r *SepatuRepoSitory) CreateSepatu(ctx context.Context, sepatu *model.Sepatu) error {
	sepatu.ID = uuid.New()
	sepatu.CreatedAt = time.Now().Unix()
	sepatu.LastUpdatedAt = time.Now().Unix()
	return r.db.WithContext(ctx).Create(&sepatu).Error
}

func (r *SepatuRepoSitory) GetAllSepatu(ctx context.Context) ([]model.Sepatu, error) {
	var sepatus []model.Sepatu
	err := r.db.WithContext(ctx).Find(&sepatus).Error
	return sepatus, err
}

func (r *SepatuRepoSitory) DeleteSepatu(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.Sepatu{}, "id = ?", id).Error
}

func (r *SepatuRepoSitory) UpdateSepatuByID(ctx context.Context, sepatuUpdate *model.UpdateSepatu, id uuid.UUID) error {
	var sepatus model.Sepatu
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
