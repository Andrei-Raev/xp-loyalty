package service

import (
	"context"

	"github.com/wintermonth2298/xp-loyalty/internal/model"
)

type AdminsRepository interface {
	Create(ctx context.Context, admin model.Admin) error
	GetByUsername(ctx context.Context, username string) (model.Admin, error)
}

type AdminsService struct {
	adminsRepo AdminsRepository
}

func NewAdminsService(adminRepo AdminsRepository) *AdminsService {
	return &AdminsService{adminsRepo: adminRepo}
}

func (s *AdminsService) Create(ctx context.Context, admin model.Admin) error {
	return s.adminsRepo.Create(ctx, admin)
}
