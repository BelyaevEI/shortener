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
		log.Log.Error("Error create table", err)
		return nil
	}

	return &database{
		DBpath: DBpath,
		db:     db,
		log:    log,
	}
}

func (d *database) Save(ctx context.Context, url1, url2 string, userID uint32) error {
	_, err := d.db.ExecContext(ctx, "INSERT INTO storage_urls(userID, short, long) values ($1, $2, $3)", userID, url1, url2)
	return err
}

func (d *database) GetShortURL(ctx context.Context, inputURL string, log *logger.Logger) (string, error) {

	var (
		foundURL string
		err      error
	)

	log.Log.Info(inputURL)

	row := d.db.QueryRowContext(ctx, "SELECT short FROM storage_urls where long=$1", inputURL)
	if err = row.Scan(&foundURL); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return "", err
		}
		return "", nil
	}

	log.Log.Info(foundURL)

	return foundURL, nil
}

func (d *database) GetOriginURL(ctx context.Context, inputURL string, log *logger.Logger) (string, bool, error) {

	var (
		foundURL string
		err      error
		deleted  bool
	)

	log.Log.Info(inputURL)

	row := d.db.QueryRowContext(ctx, "SELECT long, deleted FROM storage_urls where short=$1", inputURL)
	if err = row.Scan(&foundURL, &deleted); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return "", false, err
		}
		return "", false, nil
	}

	log.Log.Info(foundURL)

	return foundURL, deleted, nil
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

func (d *database) UpdateDeletedFlag(ctx context.Context, data []string, userID uint32, log *logger.Logger) error {
	// соберём данные для создания запроса с групповой вставкой
	var (
		values []string
		args   []any
		query  string
	)

	args = append(args, userID)

	for i, line := range data {
		count := i + 2
		params := fmt.Sprintf("short = $%d", count)
		values = append(values, params)
		args = append(args, line)
	}

	// составляем строку запроса
	if len(data) < 2 {
		query = "UPDATE storage_urls SET deleted = true WHERE userID = $1 AND short = $2"
	} else {
		query = "UPDATE storage_urls SET deleted = true WHERE userID = $1 AND (" + strings.Join(values, " OR ") + ")"
	}
	fmt.Println(query, args)
	sql, err := d.db.ExecContext(ctx, query, args)
	log.Log.Infoln(sql.RowsAffected())
	return err

}
