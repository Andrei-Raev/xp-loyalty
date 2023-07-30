package model

const (
	TypeImageAvatar    string = "avatar"
	TypeImagePrize     string = "prize"
	TypeCardBackground string = "card"
)

type Image struct {
	ID   string `json:"id"`
	URL  string `json:"url"`
	Type string `json:"type"`
}
