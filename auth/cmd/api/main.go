package main

import (
	"github.com/juxue97/auth/internal/db"
	"github.com/juxue97/auth/internal/repository"
	"github.com/juxue97/common"
)

const version = "0.0.1" // set into .env

//	@title			Auth API
//	@description	This is a authentication backend server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

// @license.name				Apache 2.0
// @license.url				http://www.apache.org/licenses/LICENSE-2.0.html
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	// Consume .env here
	cfg := config{
		version: version,
		apiUrl:  common.GetString("API_URL", "localhost:8000"),
		addr:    common.GetString("API_ADDR", ":8000"),
		env:     common.GetString("ENV", "development"),
		db: db.PgDBConfig{
			Addr:         common.GetString("DB_ADDR", "postgres://cibai:sohai@localhost:3000/social?sslmode=disable"),
			MaxOpenConns: common.GetInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns: common.GetInt("DB_MAX_IDLE_CONNS", 25),
			MaxIdleTime:  common.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	logger := common.Logger

	db, err := db.New(cfg.db.Addr, cfg.db.MaxOpenConns, cfg.db.MaxIdleConns, cfg.db.MaxIdleTime)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("Pg database connection established")

	store := repository.NewRepository(db)

	app := &application{
		config: cfg,
		store:  store,
		logger: logger,
	}
	mux := app.mount()

	logger.Fatal(app.run(mux))
}
