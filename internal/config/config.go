package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config is type configuration of service.
type Config struct {
	Telegram TelegramCfg `yaml: telegram`
	Server   ServerCfg   `yaml: server`
}

// TelegramCfg is type telegram configuration.
type TelegramCfg struct {
	BotKey string `yaml: botkey`
}

// ServerCfg is type server configuration.
type ServerCfg struct {
	Enabled       bool   `yaml: enabled`
	Addr          string `yaml: addr`
	Port          string `yaml: port`
	AdminLogin    string `yaml: adminlogin`
	AdminPassword string `yaml: adminpassword`
	JWTSecret     string `yaml: jwtsecret`
}

// New create empty Config.
func New() *Config {
	return &Config{}
}

func NewFromFile(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file %q: %w", path, err)
	}

	cfg := Config{}
	if err := yaml.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("decode data from file %q: %w", path, err)
	}

	return &cfg, nil
}

func (c *Config) WithTelegramBotKey(key string) *Config {
	c.Telegram.BotKey = key

	return c
}
