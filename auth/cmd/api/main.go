package main

import (
	"fmt"

	"github.com/juxue97/auth/internal/db"
	"github.com/juxue97/auth/internal/repository"
	"github.com/juxue97/common"
)

const version = "0.0.1" // set into .env

func main() {
	// Consume .env here
	cfg := config{
		version: version,
		url:     common.GetString("URL", "http://localhost"),
		addr:    common.GetString("ADDR", "8000"),
		env:     common.GetString("ENV", "development"),
		db: db.PgDBConfig{
			Addr:         common.GetString("DB_ADDR", "postgres://cibai:sohai@localhost:3000/social?sslmode=disable"),
			MaxOpenConns: common.GetInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns: common.GetInt("DB_MAX_IDLE_CONNS", 25),
			MaxIdleTime:  common.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	logger := common.NewLogger()
	fmt.Println(cfg.db.Addr)
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
