package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	TelegramToken  string `json:"telegram_token"`
	TelegramChatID string `json:"telegram_chat_id"`
}

var cachedConfig *Config

func LoadConfig() Config {
	if cachedConfig != nil {
		return *cachedConfig
	}

	data, err := os.ReadFile("config/nodeos_config.json")
	if err != nil {
		return Config{}
	}

	cfg := Config{}
	json.Unmarshal(data, &cfg)
	cachedConfig = &cfg
	return cfg
}

func InvalidateCache() {
	cachedConfig = nil
}
