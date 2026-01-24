package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config is type configuration of service.
type Config struct {
	Telegram TelegramCfg `yaml:"telegram"`
	Server   ServerCfg   `yaml:"server"`
	Database DatabaseCfg `yaml:"database"`
}

// DatabaseCfg is type database configuration.
type DatabaseCfg struct {
	Type string `yaml:"type"` // "memory" or "sqlite"
	Path string `yaml:"path"` // path to sqlite file
}

// TelegramCfg is type telegram configuration.
type TelegramCfg struct {
	BotKey string `yaml:"botkey"`
}

// ServerCfg is type server configuration.
type ServerCfg struct {
	Enabled       bool   `yaml:"enabled"`
	Addr          string `yaml:"addr"`
	Port          string `yaml:"port"`
	AdminLogin    string `yaml:"adminlogin"`
	AdminPassword string `yaml:"adminpassword"`
	JWTSecret     string `yaml:"jwtsecret"`
}

// New create empty Config.
func New() *Config {
	return &Config{}
}

func NewFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file %q: %w", path, err)
	}

	cfg := Config{}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("decode data from file %q: %w", path, err)
	}

	cfg.loadFromEnv()

	return &cfg, nil
}

// LoadEnvFile loads environment variables from a .env file.
// It does not override existing environment variables.
func LoadEnvFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // .env file is optional
		}

		return fmt.Errorf("open env file: %w", err)
	}

	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, found := strings.Cut(line, "=")
		if !found {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		// don't override existing env vars
		if os.Getenv(key) == "" {
			_ = os.Setenv(key, value)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan env file: %w", err)
	}

	return nil
}

func (c *Config) WithTelegramBotKey(key string) *Config {
	c.Telegram.BotKey = key

	return c
}

// loadFromEnv overrides config values with environment variables.
func (c *Config) loadFromEnv() {
	if v := os.Getenv("TELEGRAM_BOT_KEY"); v != "" {
		c.Telegram.BotKey = v
	}

	if v := os.Getenv("ADMIN_LOGIN"); v != "" {
		c.Server.AdminLogin = v
	}

	if v := os.Getenv("ADMIN_PASSWORD"); v != "" {
		c.Server.AdminPassword = v
	}

	if v := os.Getenv("JWT_SECRET"); v != "" {
		c.Server.JWTSecret = v
	}
}
