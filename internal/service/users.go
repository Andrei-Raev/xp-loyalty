package service

import (
	"context"
	"time"

	"github.com/Andrei-Raev/xp-loyalty/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user model.User) error
	GetByUsername(ctx context.Context, username string) (model.User, error)
	GetAll(ctx context.Context) ([]model.User, error)
	Update(ctx context.Context, user model.User) error
}

type UserService struct {
	userRepo   UserRepository
	awardsRepo AwardsRepo
	imagesRepo ImageRepository
}

func NewUserService(userRepo UserRepository, awardsRepo AwardsRepo, imagesRepo ImageRepository) *UserService {
	return &UserService{
		userRepo:   userRepo,
		awardsRepo: awardsRepo,
		imagesRepo: imagesRepo,
	}
}

func (s UserService) Create(ctx context.Context, id, username, avatarURL, nickname string, role model.Role) error {
	u := model.User{
		CredentialsSecure: model.CredentialsSecure{
			ID:       id,
			Username: username,
			Role:     role,
		},
		Nickname:             nickname,
		AvatarURL:            avatarURL,
		XPoints:              0,
		RegistrationTime:     time.Now(),
		LastDailyCardsUpdate: time.Now().Add(-24 * time.Hour),
	}

	return s.userRepo.Create(ctx, u)
}

func (s UserService) Prizes(ctx context.Context, username string) ([]model.UserPrize, error) {
	allPrizes, err := s.imagesRepo.GetPrizes(ctx)
	if err != nil {
		return nil, err
	}

	userPrizes, err := s.awardsRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	for _, p1 := range allPrizes {
		isInUserPrizes := false
		for _, p2 := range userPrizes {
			if p1.URL == p2.URL {
				isInUserPrizes = true
				break
			}
		}
		if !isInUserPrizes {
			convertedPrize := model.UserPrize{
				URL:           p1.URL,
				Available:     false,
				OwnerUsername: username,
			}
			userPrizes = append(userPrizes, convertedPrize)
		}
	}

	return userPrizes, nil
}

func (s UserService) GetByUsername(ctx context.Context, username string) (model.User, error) {
	return s.userRepo.GetByUsername(ctx, username)
}

func (s UserService) GetAll(ctx context.Context) ([]model.User, error) {
	return s.userRepo.GetAll(ctx)
}

func (s UserService) Update(ctx context.Context, user model.User) error {
	return s.userRepo.Update(ctx, user)
}
