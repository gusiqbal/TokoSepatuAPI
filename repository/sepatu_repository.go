package repository

import (
	"context"
	"github.com/google/uuid"
	"learnapirest/config"
	"learnapirest/model"
	"time"
)

func CreateSepatu(ctx context.Context, sepatu *model.Sepatu) error {
	sepatu.ID = uuid.New()
	sepatu.CreatedAt = time.Now().Unix()
	sepatu.LastUpdatedAt = time.Now().Unix()
	return config.DB.WithContext(ctx).Create(&sepatu).Error
}

func GetAllSepatu(ctx context.Context) ([]model.Sepatu, error) {
	var sepatus []model.Sepatu
	err := config.DB.WithContext(ctx).Find(&sepatus).Error
	return sepatus, err
}

func DeleteSepatu(ctx context.Context, id uuid.UUID) error {
	return config.DB.WithContext(ctx).Delete(&model.Sepatu{}, "id = ?", id).Error
}

func UpdateSepatuByID(ctx context.Context, sepatuUpdate *model.UpdateSepatu, id uuid.UUID) error {
	var sepatus model.Sepatu
	if err := config.DB.WithContext(ctx).First(&sepatus, "id = ?", id).Error; err != nil {
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

	return config.DB.Save(&sepatus).Error
}
