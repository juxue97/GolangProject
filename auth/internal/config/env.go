package config

import (
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
}

type pgDBConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

const version = "0.0.1" // set into .env

var Configs Config

func init() {
	Configs = Config{
		Version: version,
		ApiUrl:  common.GetString("API_URL", "localhost:8000"),
		Addr:    common.GetString("API_ADDR", ":8000"),
		Env:     common.GetString("ENV", "development"),
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
	}
}
