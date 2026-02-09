package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Dsn      string `env:"MONGO_DB_DSN" env-required:"true"`
	Dbname   string `env:"DB_NAME" env-required:"true"`
	DbUser   string `env:"DB_USER" env-required:"true"`
	DbPass   string `env:"DB_PASS" env-required:"true"`
	HostPort string `env:"HOST_PORT" env-required:"true"`
	Kafka
	Telegram
}

type Kafka struct {
	Brokers []string `env:"KAFKA_BROKERS" env-required:"true"`
	Topic   string   `env:"KAFKA_TOPIC" env-required:"true"`
}

type Telegram struct {
	Token  string `env:"TELEGRAM_BOT_TOKEN"`
	ChatId string `env:"TELEGRAM_CHAT_ID"`
}

func MustLoad() *Config {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	return &cfg
}
