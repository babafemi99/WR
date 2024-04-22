package config

import (
	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
	"log"
)

const (
	AppSrvName = "core-api"
)

type Config struct {
	ServiceName string
	DataBaseUrl string `env:"DSN"`
	RedisURL    string `env:"REDIS_URL"`
	Port        int16  `env:"PORT" envDefault:"7019"`

	EmailPassword string `env:"SMTP_PASSWORD"`
	EmailPort     string `env:"SMTP_PORT"`
	EmailHost     string `env:"SMTP_HOST"`
	EmailUsername string `env:"SMTP_USERNAME"`
	EmailFromAddr string `env:"EMAIL_FROM_ADDR" envDefault:"info@weddingregistry.com"`
}

// New loads the environment variable, parses them to the Config struct
// and returns an instance of Config
func New() *Config {

	if loadErr := godotenv.Load(".env"); loadErr != nil {
		log.Printf("[Env]: unable to load .env file %v", loadErr)
	}

	var cfg Config
	if parseErr := env.Parse(&cfg); parseErr != nil {
		log.Fatalf("[Env]: failed to parse environment variables: %v", parseErr)
	}
	cfg.ServiceName = AppSrvName

	return &cfg
}
