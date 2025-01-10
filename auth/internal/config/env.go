package config

import (
	"github.com/juxue97/common"
)

type Config struct {
	ApiUrl  string
	Version string
	Addr    string
	Env     string
	DB      pgDBConfig
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
		// mail: mailConfig{
		// 	exp:       time.Hour * 24 * time.Duration(common.GetInt("MAIL_EXPIRATION_DAYS", 3)),
		// 	fromEmail: common.GetString("FROM_EMAIL", ""),
		// 	sendGrid: sendGridConfig{
		// 		apiKey: common.GetString("SENDGRID_API_KEY", ""),
		// 	},
		// 	mailTrap: mailTrapConfig{
		// 		apiKey: common.GetString("MAILTRAP_API_KEY", ""),
		// 	},
		// },
		DB: pgDBConfig{
			Addr:         common.GetString("DB_ADDR", "postgres://juxue:veryStrongPassword@localhost:3000/auth?sslmode=disable"),
			MaxOpenConns: common.GetInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns: common.GetInt("DB_MAX_IDLE_CONNS", 25),
			MaxIdleTime:  common.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}
}
