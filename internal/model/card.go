package model

import "encoding/json"

const (
	StatusPending string = "pending"
	StatusDone    string = "done"
)

type Card struct {
	ID            string     `json:"id"`
	OwnerUsername string     `json:"owner_username"`
	Static        CardStatic `json:"static"`
	Done          int        `json:"done"`
	Progress      int        `json:"progress"`
	IsViewed      bool       `json:"is_viewed"`
	History       []int      `json:"history"`
	OptDoneNum    int        `json:"opt_done_num"`
}

type Cards []Card

func (cards Cards) ContainsStatic(card CardStatic) bool {
	for _, c := range cards {
		if c.Static.ID == card.ID {
			return true
		}
	}
	return false
}

func (c *Card) MarshalJSON() ([]byte, error) {
	type Alias Card

	var progress *int
	if c.Static.Type == TypeProgress {
		progress = &c.Progress
	}

	var history []int
	if c.Static.Type == TypeOptions {
		history = c.History
		if len(history) == 0 {
			history = []int{}
		}
	}

	return json.Marshal(&struct {
		Progress *int  `json:"progress,omitempty"`
		History  []int `json:"history"`
		*Alias
	}{
		Progress: progress,
		History:  history,
		Alias:    (*Alias)(c),
	})
}
