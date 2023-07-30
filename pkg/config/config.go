package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	ServerPort        string `json:"server_port"`
	ModeratorUsername string `json:"moderator_username"`
	ModeratorPassword string `json:"moderator_password"`
	Mongo             Mongo  `json:"mongo"`
	SQL               SQL    `json:"sql"`
	User              User   `json:"user"`
}

type Mongo struct {
	URI  string `json:"uri"`
	Name string `json:"name"`
}

type SQL struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
}

type User struct {
	DailyCardsNum int  `json:"daily_cards_num"`
	UniqueGoals   bool `json:"unique_goals"`
}

func New(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %w", err)
	}

	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("error decoding json: %w", err)
	}

	if err := file.Close(); err != nil {
		return nil, err
	}

	return &cfg, nil
}
