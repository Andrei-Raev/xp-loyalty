package model

import (
	"time"
)

type (
	Role int
)

const (
	RoleUser Role = iota + 1
	RoleAdmin
	RoleModerator
)

type User struct {
	CredentialsSecure
	Nickname             string      `json:"nickname"`
	AvatarURL            string      `json:"avatar_url"`
	XPoints              int         `json:"XPoints"`
	RegistrationTime     time.Time   `json:"registration_time"`
	LastDailyCardsUpdate time.Time   `json:"last_daily_cards_update"`
	Prizes               []UserPrize `json:"prizes"`
}

type UserPrize struct {
	URL           string `json:"url"`
	OwnerUsername string `json:"owner_username"`
	Available     bool   `json:"available"`
}
