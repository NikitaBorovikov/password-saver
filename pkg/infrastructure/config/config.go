package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	configYAMLPath    = "config/config.yml"
	configProdENVPath = ".prod.env"
	configDevENVPath  = ".dev.env"
)

type (
	Config struct {
		Postgres    Postgres
		EncryptKeys EncryptKeys
		Http        Http       `yaml:"http"`
		RateLimits  RateLimits `yaml:"rate_limits"`
		Telegram    Telegram
	}

	Telegram struct {
		Token string `env:"TG_TOKEN"`
	}

	Postgres struct {
		Host     string `env:"PG_HOST"`
		Port     int    `env:"PG_PORT"`
		User     string `env:"PG_USER"`
		Password string `env:"PG_PASSWORD"`
		Name     string `env:"PG_NAME"`
	}

	Http struct {
		Port              string        `yaml:"port"`
		MiddlewareTimeout time.Duration `yaml:"middleware_timeout"`
		ReadTimeout       time.Duration `yaml:"read_timeout"`
		IdleTimeout       time.Duration `yaml:"idle_timeout"`
		SessionKey        string        `env:"SESSION_KEY"`
		SessionName       string        `env:"SESSION_NAME"`
	}

	RateLimits struct {
		Auth        int `yaml:"auth"`
		CloseRoutes int `yaml:"close_routes"`
		OpenRoutes  int `yaml:"open_routes"`
	}

	EncryptKeys struct {
		EncPasswordKey string `env:"PASSWORD_ENC_KEY"`
		EncServiceKey  string `env:"SERVICE_ENC_KEY"`
		EncLoginKey    string `env:"LOGIN_ENC_KEY"`
	}
)

func Init() (*Config, error) {
	var cfg Config

	if err := cleanenv.ReadConfig(configYAMLPath, &cfg); err != nil {
		return nil, err
	}

	if err := readFromEnv(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func readFromEnv(cfg *Config) error {
	env := os.Getenv("APP_ENV")

	switch env {
	case "prod":
		if err := cleanenv.ReadConfig(configProdENVPath, cfg); err != nil {
			return err
		}
	default:
		if err := cleanenv.ReadConfig(configDevENVPath, cfg); err != nil {
			return err
		}
	}
	return nil
}
