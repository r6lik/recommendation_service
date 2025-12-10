package postgresql

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/r6lik/recommendation_service/internal/adapters/config"
)

type Database struct {
	*sqlx.DB
}

func NewDatabase(cfg *config.Config) (*Database, error) {
	db, err := sqlx.Connect("postgres", cfg.Database.DSN())
	if err != nil {
		return nil, errors.New("failed to connect to database")
	}
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetime) * time.Minute)

	return &Database{db}, nil
}
