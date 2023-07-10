package database

import (
	"database/sql"

	"github.com/BelyaevEI/shortener/internal/config"
)

func Connect(cfg config.Parameters) (*sql.DB, error) {

	db, err := sql.Open("postgres", cfg.DBStoragePath)
	if err != nil {
		return nil, err
	}
	return db, nil
}
