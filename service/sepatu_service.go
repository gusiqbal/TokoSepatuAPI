package service

import (
	"context"
	"learnapirest/model"
	"learnapirest/repository"

	"github.com/google/uuid"
)

type ISepatuService interface {
	CreateSepatu(ctx context.Context, sepatus *model.Sepatu) error
	GetSepatu(ctx context.Context) ([]model.Sepatu, error)
	DeleteSepatu(ctx context.Context, delSepatu *model.DeleteSepatu) error
	UpdateSepatu(ctx context.Context, sepatu *model.UpdateSepatu, id uuid.UUID) error
}

type SepatuService struct {
	repo repository.ISepatuRepository
}

func NewSepatuService(repo repository.ISepatuRepository) *SepatuService {
	return &SepatuService{
		repo: repo,
	}
}

func (s *SepatuService) CreateSepatu(ctx context.Context, sepatus *model.Sepatu) error {
	sepatus.ID = uuid.New()
	return s.repo.CreateSepatu(ctx, sepatus)
}

func (s *SepatuService) GetSepatu(ctx context.Context) ([]model.Sepatu, error) {
	return s.repo.GetAllSepatu(ctx)
}

func (s *SepatuService) DeleteSepatu(ctx context.Context, delSepatu *model.DeleteSepatu) error {
	return s.repo.DeleteSepatu(ctx, delSepatu.ID)
}

func (s *SepatuService) UpdateSepatu(ctx context.Context, sepatu *model.UpdateSepatu, id uuid.UUID) error {
	return s.repo.UpdateSepatuByID(ctx, sepatu, id)
}
