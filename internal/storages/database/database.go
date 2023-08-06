package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/BelyaevEI/shortener/internal/logger"
	"github.com/BelyaevEI/shortener/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type database struct {
	DBpath string
	db     *sql.DB
	log    *logger.Logger
}

func New(DBpath string, log *logger.Logger) *database {
	db, err := sql.Open("pgx", DBpath)
	if err != nil {
		log.Log.Error(err)
	}

	_, err = db.Exec("create table IF NOT EXISTS storage_urls(userID bigint NOT NULL, short text not null, long text not null, deleted boolean default false)")
	if err != nil {
		log.Log.Error("Error create tabele", err)
		return nil
	}

	return &database{
		DBpath: DBpath,
		db:     db,
		log:    log,
	}
}

func (d *database) Save(ctx context.Context, url1, url2 string, userID uint32) error {
	_, err := d.db.ExecContext(ctx, "insert into storage_urls(userID, short, long) values ($1, $2, $3)", userID, url1, url2)
	return err
}

func (d *database) GetShortURL(ctx context.Context, inputURL string) (string, error) {

	var (
		foundURL string
		err      error
	)

	row := d.db.QueryRowContext(ctx, "select short from storage_urls where long=$1", inputURL)
	if err = row.Scan(&foundURL); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return "", err
		}
		return "", nil
	}
	return foundURL, nil
}

func (d *database) GetOriginURL(ctx context.Context, inputURL string) (string, bool, error) {

	var (
		foundURL string
		err      error
		deleted  bool
	)

	row := d.db.QueryRowContext(ctx, "select long, deleted from storage_urls where short=$1", inputURL)
	if err = row.Scan(&foundURL, &deleted); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return "", false, err
		}
		return "", false, nil
	}

	if deleted {
		return foundURL, true, nil
	}

	return foundURL, false, nil
}

func (d *database) Ping(ctx context.Context) error {
	if err := d.db.Ping(); err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}

func (d *database) GetUrlsUser(ctx context.Context, userID uint32) ([]models.StorageURL, error) {
	storageURLS := make([]models.StorageURL, 0)

	rows, err := d.db.QueryContext(ctx, "SELECT userID, short, long from storage_urls WHERE userID=$1", userID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var store models.StorageURL
		err = rows.Scan(&store.UserID, &store.ShortURL, &store.OriginalURL)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				return nil, err
			}
			return nil, nil
		}
		storageURLS = append(storageURLS, store)
	}

	// проверяем на ошибки
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return storageURLS, nil

}

func (d *database) UpdateDeletedFlag(ctx context.Context, data []string, userID uint32) error {
	// соберём данные для создания запроса с групповой вставкой
	var values []string
	var args []any

	args = append(args, userID)

	for i, line := range data {
		count := i + 2
		params := fmt.Sprintf("short = $%d", count)
		values = append(values, params)
		args = append(args, line)
	}

	// составляем строку запроса
	query := `UPDATE storage_urls
				set deleted = true
				WHERE userID =$1
				AND (` + strings.Join(values, "OR") + `)`

	_, err := d.db.ExecContext(ctx, query, args)
	return err
}
