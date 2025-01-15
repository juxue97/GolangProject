package config

import (
	"log"
	"time"

	"github.com/juxue97/auth/internal/types"
	"github.com/juxue97/common"
)

type Config struct {
	ApiUrl      string
	Version     string
	Addr        string
	Env         string
	DB          pgDBConfig
	Mail        types.MailConfig
	FrontendURL string
	Auth        types.AuthConfig
	RedisCfg    types.RedisConfig
	RateLimit   types.RateLimitConfig
}

type pgDBConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

const version = "0.0.1" // set into .env

func GetConfig(path string) *Config {
	err := common.EnvInit(path)
	if err != nil {
		log.Fatal(err)
	}

	configs := &Config{
		Version: version,
		ApiUrl:  common.GetString("API_URL", "localhost:8000"),
		Addr:    common.GetString("API_ADDR", ":8000"),
		Env:     common.GetString("API_ENV", "development"),
		Mail: types.MailConfig{
			Exp:       time.Hour * 24 * time.Duration(common.GetInt("MAIL_EXPIRATION_DAYS", 3)),
			FromEmail: common.GetString("FROM_EMAIL", ""),
			SendGrid: types.SendGridConfig{
				ApiKey: common.GetString("SENDGRID_API_KEY", ""),
			},
			MailTrap: types.MailTrapConfig{
				ApiKey: common.GetString("MAILTRAP_API_KEY", ""),
			},
		},
		DB: pgDBConfig{
			Addr:         common.GetString("DB_ADDR", "postgres://juxue:veryStrongPassword@localhost:3000/auth?sslmode=disable"),
			MaxOpenConns: common.GetInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns: common.GetInt("DB_MAX_IDLE_CONNS", 25),
			MaxIdleTime:  common.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		FrontendURL: common.GetString("FRONTEND_URL", "http://localhost:4000"),
		Auth: types.AuthConfig{
			Token: types.TokenConfig{
				Secret: common.GetString("AUTH_TOKEN_SECRET", "secret"),
				Exp:    time.Hour * 24 * time.Duration(common.GetInt("AUTH_TOKEN_EXPIRATION_DAYS", 3)),
				Iss:    common.GetString("AUTH_TOKEN_ISSUER", "MehNohNahSuperAuth"),
				Aud:    common.GetString("AUTH_TOKEN_AUDIENCE", "MehNohNahSuperAuth"),
			},
		},
		RedisCfg: types.RedisConfig{
			Addr:     common.GetString("REDIS_ADDR", "localhost:6379"),
			Password: common.GetString("REDIS_PASSWORD", ""),
			DB:       common.GetInt("REDIS_DB", 0),
			TTL:      time.Minute * time.Duration(common.GetInt("REDIS_TTL", 1)),
			Enabled:  common.GetBool("REDIS_ENABLED", false),
		},
		RateLimit: types.RateLimitConfig{
			DB:      common.GetInt("RATE_LIMIT_DB", 1),
			Limit:   common.GetInt("RATE_LIMIT", 20),
			Window:  time.Minute * time.Duration(common.GetInt("RATE_LIMIT_WINDOW", 1)),
			Enabled: common.GetBool("RATE_LIMIT_ENABLED", false),
		},
	}
	return configs
}
