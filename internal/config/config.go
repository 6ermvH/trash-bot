package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Telegram TelegramCfg `yaml: telegram`
	Server   ServerCfg   `yaml: server`
}

type TelegramCfg struct {
	BotKey string `yaml: botkey`
}

type ServerCfg struct {
	Enabled       bool   `yaml: enabled`
	Addr          string `yaml: addr`
	Port          string `yaml: port`
	AdminLogin    string `yaml: adminlogin`
	AdminPassword string `yaml: adminpassword`
}

func New() *Config {
	return &Config{}
}

func NewFromFile(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open file %q: %v", path, err)
	}

	cfg := Config{}
	if err := yaml.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("decode data from file %q: %v", path, err)
	}

	return &cfg, nil
}
