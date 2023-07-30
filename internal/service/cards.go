package service

import (
	"context"
	"time"

	"github.com/wintermonth2298/xp-loyalty/internal/model"
)

type CardsRepo interface {
	GetStatic(ctx context.Context, ids []string) (model.CardsStatic, error)
	CreateStatic(ctx context.Context, card model.CardStatic) error
	DeleteStatic(ctx context.Context, ids []string) (err error)
	GetStaticByPool(ctx context.Context, pool string) (model.CardsStatic, error)

	GetCardsByOwnerPool(ctx context.Context, ownerUsername, pool string) (model.Cards, error)
	DeleteUsersPendingDailyCards(ctx context.Context, username string) error
	GetCardsByOwner(ctx context.Context, ownerUsername string) (model.Cards, error)
	Get(ctx context.Context, id string) (model.Card, error)
	Create(ctx context.Context, card model.Card) error
	Update(ctx context.Context, card model.Card) error
	ViewCard(ctx context.Context, id string) error
}

type AwardsRepo interface {
	Add(ctx context.Context, award model.UserPrize) error
	GetByUsername(ctx context.Context, username string) ([]model.UserPrize, error)
}

type CardsService struct {
	cardsRepo  CardsRepo
	awardsRepo AwardsRepo
}

func NewCardsStaticService(cardsStaticRepo CardsRepo, awardsRepo AwardsRepo) *CardsService {
	return &CardsService{cardsRepo: cardsStaticRepo, awardsRepo: awardsRepo}
}

func (s *CardsService) ViewCard(ctx context.Context, cardID string) error {
	return s.cardsRepo.ViewCard(ctx, cardID)
}

func (s *CardsService) Create(ctx context.Context, card model.Card) error {
	card.ID = ""
	return s.cardsRepo.Create(ctx, card)
}

func (s *CardsService) GetFormattedCards(ctx context.Context, ownerUsername string) (model.Cards, model.Cards, error) {
	pendingCards := make(model.Cards, 0, 50)
	doneCards := make(model.Cards, 0, 50)

	dailyCards, err := s.cardsRepo.GetCardsByOwnerPool(ctx, ownerUsername, model.PoolDaily)
	if err != nil {
		return model.Cards{}, model.Cards{}, err
	}

	constCards, err := s.cardsRepo.GetCardsByOwnerPool(ctx, ownerUsername, model.PoolConst)
	if err != nil {
		return model.Cards{}, model.Cards{}, err
	}

	for _, c := range dailyCards {
		if c.Done == 0 {
			pendingCards = append(pendingCards, c)
		} else {
			doneCards = append(doneCards, c)
		}
	}

	getPending := func(chain model.Cards) model.Card {
		for i := 0; i < len(chain); i++ {
			if chain[i].Done == 0 {
				return chain[i]
			}
		}

		return chain[len(chain)-1]
	}

	if len(constCards) == 0 {
		return pendingCards, doneCards, nil
	}

	chain := make(model.Cards, 0, 20)
	name := constCards[0].Static.ChainName
	for _, c := range constCards {
		if c.Done != 0 {
			for i := 0; i < c.Done; i++ {
				if c.Static.Type == model.TypeOptions {
					c.OptDoneNum = c.History[i]
				}
				doneCards = append(doneCards, c)
			}
		}

		if c.Static.ChainName == name {
			chain = append(chain, c)
			continue
		}

		pendingCards = append(pendingCards, getPending(chain))

		name = c.Static.ChainName
		chain = chain[:0:20]
		chain = append(chain, c)
	}
	// process last chain
	pendingCards = append(pendingCards, getPending(chain))

	return pendingCards, doneCards, nil
}

func (s *CardsService) Update(ctx context.Context, id string, progress int, doneOption float32) (string, int, error) {
	card, err := s.cardsRepo.Get(ctx, id)
	if err != nil {
		return "", 0, err
	}

	var gotAward model.Award

	var XPoints int
	switch card.Static.Type {
	case model.TypeOrdinary:
		XPoints = card.Static.OrdSettings.Award.XPoints
		card.Done += 1
		gotAward = card.Static.OrdSettings.Award
	case model.TypeProgress:
		card.Progress += progress
		if card.Progress >= card.Static.PrgSettings.MaxProgress {
			XPoints = card.Static.PrgSettings.Award.XPoints
			card.Progress = 0
			card.Done += 1
		}
		gotAward = card.Static.PrgSettings.Award
	case model.TypeOptions:
		if doneOption < card.Static.OptSettings.Options[0] {
			return "", 0, nil
		}

		card.Done += 1

		awards := card.Static.OptSettings.Awards
		options := card.Static.OptSettings.Options
		for i := 0; i < len(options)-1; i++ {
			if doneOption >= options[i] && doneOption < options[i+1] {
				XPoints = awards[i].XPoints
				card.History = append(card.History, i)
				gotAward = awards[i]
				break
			}
		}
		if doneOption >= options[len(options)-1] {
			XPoints = awards[len(awards)-1].XPoints
			gotAward = awards[len(awards)-1]
			card.History = append(card.History, len(options)-1)
		}
	}

	a := model.UserPrize{
		URL:           gotAward.PrizeImageURL,
		OwnerUsername: card.OwnerUsername,
		Available:     true,
	}

	if err := s.awardsRepo.Add(ctx, a); err != nil {
		return "", 0, err
	}

	err = s.cardsRepo.Update(ctx, card)
	return card.OwnerUsername, XPoints, err
}

func (s *CardsService) UpdateDailyCards(ctx context.Context, users []model.User, dcardsnum int, uniqueGoal bool) (time.Time, []int, error) {
	updatedUsers := make([]int, 0, 5)
	now := time.Now()
	nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// fetch available daily cards
	dailyStaticCards, err := s.cardsRepo.GetStaticByPool(ctx, model.PoolDaily)
	if err != nil {
		return time.Time{}, nil, err
	}

	for i, u := range users {
		t := u.LastDailyCardsUpdate.Local()
		tDate := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

		if !nowDate.After(tDate) {
			continue
		}

		newCards := dailyStaticCards.Random(dcardsnum, uniqueGoal)
		if len(newCards) == 0 {
			return time.Time{}, nil, model.ErrNoRandomCards
		}

		err := s.cardsRepo.DeleteUsersPendingDailyCards(ctx, u.Username)
		if err != nil {
			return time.Time{}, nil, err
		}

		for _, c := range newCards {
			if err := s.cardsRepo.Create(ctx, c.Card(u.Username)); err != nil {
				return time.Time{}, nil, err
			}
		}
		updatedUsers = append(updatedUsers, i)
	}

	// return nowDate to update LastDailyCardsUpdate for users
	return nowDate, updatedUsers, nil
}

func (s *CardsService) UpdateConstCards(ctx context.Context, users []model.User) error {
	constStaticCards, err := s.cardsRepo.GetStaticByPool(ctx, model.PoolConst)
	if err != nil {
		return err
	}

	for _, u := range users {
		userCards, err := s.cardsRepo.GetCardsByOwner(ctx, u.Username)
		if err != nil {
			return err
		}

		for i, c := range constStaticCards {
			if userCards.ContainsStatic(c) {
				continue
			}
			if err := s.cardsRepo.Create(ctx, constStaticCards[i].Card(u.Username)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *CardsService) CreateStatic(ctx context.Context, card model.CardStatic) error {
	card.ID = ""
	card.CreatedAt = time.Now()
	return s.cardsRepo.CreateStatic(ctx, card)
}

func (s *CardsService) GetStatic(ctx context.Context, ids []string) (model.CardsStatic, error) {
	return s.cardsRepo.GetStatic(ctx, ids)
}

func (s *CardsService) DeleteStatic(ctx context.Context, ids []string) (err error) {
	return s.cardsRepo.DeleteStatic(ctx, ids)
}
