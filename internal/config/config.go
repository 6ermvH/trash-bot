package config

import (
	"os"
)

type Config struct {
	Redis            RedisConfig
	Telegram         TelegramConfig
	OpenRouterAPIKey string
	Server           ServerConfig
}

type RedisConfig struct {
	Addr     string
	Username string
	Password string
	DB       int
}

type TelegramConfig struct {
	Token string
}

type ServerConfig struct {
	Port string
}

// Load reads environment variables into Config
func Load() (*Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return &Config{
		Redis: RedisConfig{
			Addr:     os.Getenv("REDIS_ADDR"),
			Username: os.Getenv("REDIS_USERNAME"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       1,
		},
		Telegram: TelegramConfig{
			Token: os.Getenv("TELEGRAM_APITOKEN"),
		},
		OpenRouterAPIKey: os.Getenv("OPENROUTER_API_KEY"),
		Server:           ServerConfig{Port: port},
	}, nil
}
