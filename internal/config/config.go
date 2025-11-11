package config

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	Jwt          Jwt        `split_words:"true" required:"true"`
	HttpServer   Server     `split_words:"true" required:"true"`
	PostgresqlDb PostgreSQL `split_words:"true" required:"true"`
}

type PostgreSQL struct {
	Host     string `split_words:"true" required:"true"`
	Port     string `split_words:"true" required:"true"`
	Name     string `split_words:"true" required:"true"`
	UserName string `split_words:"true" required:"true"`
	Password string `split_words:"true" required:"true"`
	Schema   string `split_words:"true" required:"true"`
}

type Server struct {
	Address    string `split_words:"true" required:"true"`
	CertPath   string `split_words:"true" required:"true"`
	KeyPath    string `split_words:"true" required:"true"`
	CaCertPath string `split_words:"true" required:"true"`
}

type Jwt struct {
	AccessTokenExpiry  int64 `split_words:"true" default:"15"` //min
	RefreshTokenExpiry int64 `split_words:"true" default:"60"`
}

func LoadAppConfig() (*AppConfig, error) {
	if err := godotenv.Load("internal/config/.env"); err != nil {
		log.Printf("No .env file found: %v", err)
	}
	var appCfg AppConfig
	err := envconfig.Process("", &appCfg)
	if err != nil {
		return nil, err
	}
	return &appCfg, nil
}
