package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/juxue97/auth/internal/config"
	"github.com/juxue97/common"
	_ "github.com/lib/pq"
)

var (
	PgDB *sql.DB
	err  error
)

func init() {
	PgDB, err = NewPgClient(config.Configs.DB.Addr, config.Configs.DB.MaxOpenConns, config.Configs.DB.MaxIdleConns, config.Configs.DB.MaxIdleTime)
	if err != nil {
		common.Logger.Fatal(err)
	}

	// defer PgDB.Close()
	common.Logger.Info("Pg database connection established")
}

func NewPgClient(addr string, maxOpenConns int, maxIdleConns int, maxIdleTime string) (*sql.DB, error) {
	db, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, err
	}
	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}
