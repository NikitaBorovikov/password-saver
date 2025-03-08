package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		Postgres    Postgres
		EncryptKeys EncryptKeys
		Http        Http `yml:"http"`
	}

	Postgres struct {
		URL string `env:"PG_URL"`
	}

	Http struct {
		Port       string `yml:"port"`
		SessionKey string `env:"SESSION_KEY"`
	}

	EncryptKeys struct {
		EncPasswordKey string `env:"PASSWORD_ENC_KEY"`
		EncServiceKey  string `env:"SERVICE_ENC_KEY"`
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
	if err := godotenv.Load(); err != nil {
		return err
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
