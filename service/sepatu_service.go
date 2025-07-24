package service

import (
	"context"
	"github.com/google/uuid"
	"learnapirest/model"
	"learnapirest/repository"
)

func CreateSepatu(ctx context.Context, sepatus *model.Sepatu) error {
	sepatus.ID = uuid.New()
	return repository.CreateSepatu(ctx, sepatus)
}

func GetSepatu(ctx context.Context) ([]model.Sepatu, error) {
	return repository.GetAllSepatu(ctx)
}

func DeleteSepatu(ctx context.Context, delSepatu *model.DeleteSepatu) error {
	return repository.DeleteSepatu(ctx, delSepatu.ID)
}

func UpdateSepatu(ctx context.Context, sepatu *model.UpdateSepatu, id uuid.UUID) error {
	return repository.UpdateSepatuByID(ctx, sepatu, id)
}
