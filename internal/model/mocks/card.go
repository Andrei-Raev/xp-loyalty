package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/wintermonth2298/xp-loyalty/internal/model"
)

type CardsRepositoryMock struct {
	mock.Mock
}

func (m *CardsRepositoryMock) ViewCard(ctx context.Context, cardID string) error {
	panic("not implemented") // TODO: Implement
}

func (m *CardsRepositoryMock) GetCardsByOwnerPool(ctx context.Context, ownerUsername, pool string) (model.Cards, error) {
	args := m.Called(ctx, ownerUsername, pool)
	return args.Get(0).(model.Cards), args.Error(1)
}

func (m *CardsRepositoryMock) GetStatic(ctx context.Context, ids []string) (model.CardsStatic, error) {
	panic("not implemented") // TODO: Implement
}

func (m *CardsRepositoryMock) CreateStatic(ctx context.Context, card model.CardStatic) error {
	panic("not implemented") // TODO: Implement
}

func (m *CardsRepositoryMock) DeleteStatic(ctx context.Context, ids []string) (err error) {
	panic("not implemented") // TODO: Implement
}

func (m *CardsRepositoryMock) GetStaticByPool(ctx context.Context, pool string) (model.CardsStatic, error) {
	args := m.Called(ctx, pool)
	return args.Get(0).(model.CardsStatic), args.Error(1)
}

func (m *CardsRepositoryMock) DeleteUsersPendingDailyCards(ctx context.Context, username string) error {
	panic("not implemented") // TODO: Implement
}

func (m *CardsRepositoryMock) GetCardsByOwner(ctx context.Context, ownerUsername string) (model.Cards, error) {
	args := m.Called(ctx, ownerUsername)
	return args.Get(0).(model.Cards), args.Error(1)
}

func (m *CardsRepositoryMock) Get(ctx context.Context, id string) (model.Card, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Card), args.Error(1)
}

func (m *CardsRepositoryMock) Create(ctx context.Context, card model.Card) error {
	args := m.Called(ctx, card)
	return args.Error(0)
}

func (m *CardsRepositoryMock) Update(ctx context.Context, card model.Card) error {
	args := m.Called(ctx, card)
	return args.Error(0)
}
