package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Redis    RedisConfig    `yaml:"redis"`
	Telegram TelegramConfig `yaml:"telegram"`
	Server   ServerConfig   `yaml:"server"`

	OpenRouterAPIKey string // из окружения
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Username string // из окружения
	Password string // из окружения
	DB       int    `yaml:"db"`
}

type TelegramConfig struct {
	Token string // из окружения
}

type ServerConfig struct {
	Port string `yaml:"port"`
}

// Load загружает конфиг из YAML и дополняет переменными окружения
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Подгружаем чувствительные данные из ENV
	cfg.Redis.Username = os.Getenv("REDIS_USERNAME")
	cfg.Redis.Password = os.Getenv("REDIS_PASSWORD")
	cfg.Telegram.Token = os.Getenv("TELEGRAM_APITOKEN")
	cfg.OpenRouterAPIKey = os.Getenv("OPENROUTER_API_KEY")

	if cfg.Server.Port == "" {
		cfg.Server.Port = "8080"
	}

	return &cfg, nil
}

