package model

import (
	"encoding/json"
	"math/rand"
	"time"
)

const (
	TypeOrdinary string = "ordinary"
	TypeProgress string = "progress"
	TypeOptions  string = "options"

	GoalBuyFood        string = "food"
	GoalBuyDrink       string = "drink"
	GoalPlayMore       string = "play"
	GoalSocialActivity string = "social"

	PoolDaily string = "daily"
	PoolConst string = "const"
)

type CardStatic struct {
	ID               string       `json:"id"`
	Title            string       `json:"title"`
	ShortDescription string       `json:"short_description"`
	LongDescription  string       `json:"long_description"`
	Goal             string       `json:"goal"`
	CreatedAt        time.Time    `json:"created_at"`
	Type             string       `json:"type"`
	Pool             string       `json:"pool"`
	BackgroundURL    string       `json:"background_url"`
	ChainName        string       `json:"chain_name"`
	ChainOrder       int          `json:"chain_order"`
	OrdSettings      *OrdSettings `json:"ordinary_settings,omitempty"`
	PrgSettings      *PrgSettings `json:"progress_settings,omitempty"`
	OptSettings      *OptSettings `json:"options_settings,omitempty"`
}

type OrdSettings struct {
	Award Award `json:"award"`
}

type PrgSettings struct {
	Award       Award `json:"award"`
	MaxProgress int   `json:"max_progress"`
}

type OptSettings struct {
	Awards  []Award   `json:"awards"`
	Options []float32 `json:"options"`
}

type Award struct {
	XPoints       int    `json:"XPoints"`
	Prize         string `json:"prize"`
	PrizeImageURL string `json:"prize_image_url"`
}

func (card CardStatic) Card(ownerUsername string) Card {
	c := Card{Static: card}
	c.OwnerUsername = ownerUsername
	c.Done = 0
	switch card.Type {
	case TypeProgress:
		c.Progress = 0
	case TypeOptions:
		c.History = []int{}
	}

	return c
}

type CardsStatic []CardStatic

func (cards CardsStatic) Unique() []CardStatic {
	exist := make(map[string]bool)
	unique := make([]CardStatic, 0, len(cards))

	for _, c := range cards {
		if _, v := exist[c.ID]; !v {
			exist[c.ID] = true
			unique = append(unique, c)
		}
	}
	return unique
}

func (cards CardsStatic) Random(num int, uniqueGoal bool) []CardStatic {
	if len(cards) < num {
		return []CardStatic{}
	}

	uniqueRandomCards := make([]CardStatic, 0, num)

	indexes := make(map[int]struct{})
	goals := make(map[string]struct{})
	for len(uniqueRandomCards) < num {
		i := rand.Intn(len(cards))

		if _, ok := indexes[i]; ok {
			continue
		}

		if uniqueGoal {
			if _, ok := goals[cards[i].Goal]; ok {
				continue
			}
			goals[cards[i].Goal] = struct{}{}
		}

		indexes[i] = struct{}{}
		uniqueRandomCards = append(uniqueRandomCards, cards[i])
	}
	return uniqueRandomCards
}

func (c *CardStatic) MarshalJSON() ([]byte, error) {
	type Alias CardStatic

	var chainName *string
	var chainOrder *int
	if c.Pool == PoolConst {
		chainName = &c.ChainName
		chainOrder = &c.ChainOrder
	}

	return json.Marshal(&struct {
		ChainName  *string `json:"chain_name,omitempty"`
		ChainOrder *int    `json:"chain_order,omitempty"`
		*Alias
	}{
		ChainName:  chainName,
		ChainOrder: chainOrder,
		Alias:      (*Alias)(c),
	})
}
