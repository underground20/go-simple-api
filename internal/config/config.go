package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Dsn    string `env:"MONGO_DB_DSN" env-required:"true"`
	Dbname string `env:"DB_NAME" env-required:"true"`
	DbUser string `env:"DB_USER" env-required:"true"`
	DbPass string `env:"DB_PASS" env-required:"true"`
}

func MustLoad() *Config {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	return &cfg
}
