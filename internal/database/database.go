package database

import (
	"database/sql"

	"github.com/BelyaevEI/shortener/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect(cfg config.Parameters) (*sql.DB, error) {

	db, err := sql.Open("pgx", cfg.DBStoragePath)
	if err != nil {
		return nil, err
	}
	return db, nil

}
