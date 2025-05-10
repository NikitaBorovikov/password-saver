package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		Postgres    Postgres
		EncryptKeys EncryptKeys
		Http        Http       `yaml:"http"`
		RateLimits  RateLimits `yaml:"rate_limits"`
	}

	Postgres struct {
		Host     string `env:"PG_HOST"`
		Port     int64  `env:"PG_PORT"`
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

	if err := cleanenv.ReadConfig("config/config.yml", &cfg); err != nil {
		return nil, err
	}

	if err := ReadFromEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func ReadFromEnv(cfg *Config) error {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	//if not prod, then we take the data from the .env.dev file.
	//if prod, then from the environment variables
	if env == "dev" {
		if err := godotenv.Load(".env.dev"); err != nil {
			return err
		}
	}

	if err := cleanenv.ReadEnv(&cfg.Postgres); err != nil {
		return err
	}

	if err := cleanenv.ReadEnv(&cfg.Http); err != nil {
		return err
	}

	err := cleanenv.ReadEnv(&cfg.EncryptKeys)
	return err
}
