package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/BelyaevEI/shortener/internal/logger"
	"github.com/BelyaevEI/shortener/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type database struct {
	DBpath string
	db     *sql.DB
	log    *logger.Logger
}

// Create a new storage using database
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

// Save short and origin urls to database
func (d *database) Save(ctx context.Context, url1, url2 string, userID uint32) error {
	_, err := d.db.ExecContext(ctx, "INSERT INTO storage_urls(userID, short, long) values ($1, $2, $3)", userID, url1, url2)
	return err
}

// Find short url
func (d *database) GetShortURL(ctx context.Context, inputURL string) (string, error) {

	var (
		foundURL string
		err      error
	)

	d.log.Log.Info(inputURL)

	row := d.db.QueryRowContext(ctx, "SELECT short FROM storage_urls where long=$1", inputURL)
	if err = row.Scan(&foundURL); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return "", err
		}
		return "", nil
	}

	d.log.Log.Info(foundURL)

	return foundURL, nil
}

// Find origin url
func (d *database) GetOriginURL(ctx context.Context, inputURL string) (string, bool, error) {

	var (
		foundURL string
		err      error
		deleted  bool
	)

	d.log.Log.Info(inputURL)

	row := d.db.QueryRowContext(ctx, "SELECT long, deleted FROM storage_urls where short=$1", inputURL)
	if err = row.Scan(&foundURL, &deleted); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return "", false, err
		}
		return "", false, nil
	}

	d.log.Log.Info(foundURL)

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

// Find all user url
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

// Update deleted flag in datatabse
func (d *database) UpdateDeletedFlag(data models.DeleteURL) {

	sql, err := d.db.Exec("UPDATE storage_urls SET deleted = true WHERE userID = $1 AND short = $2", data.UserID, data.ShortURL)
	if err != nil {
		d.log.Log.Error(err)
		return
	}

	d.log.Log.Infoln(sql.RowsAffected())
}

func (d *database) GetStatistic() models.Statistic {
	var stat models.Statistic

	res := d.db.QueryRow("SELECT count(distinct userID), count(*) FROM storage_urls")
	err := res.Scan(&stat.Users, &stat.Urls)
	if err != nil {
		return stat
	}
	return stat
}
