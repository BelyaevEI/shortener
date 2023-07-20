package database

import (
	"database/sql"
	"errors"

	"github.com/BelyaevEI/shortener/internal/logger"
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

	_, err = db.Exec("create table IF NOT EXISTS storage_urls(short text not null, long text not null)")
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

func (d *database) Save(url1, url2 string) error {
	_, err := d.db.Exec("insert into storage_urls(short, long) values ($1, $2)", url1, url2)
	return err
}

func (d *database) GetShortURL(inputURL string) (s string, e error) {

	var (
		foundURL string
		err      error
	)

	row := d.db.QueryRow("select short from storage_urls where long=$1", inputURL)
	if err = row.Scan(&foundURL); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return "", err
		}
		return "", nil
	}
	return foundURL, nil
}

func (d *database) GetOriginURL(inputURL string) (s string, e error) {

	var (
		foundURL string
		err      error
	)

	row := d.db.QueryRow("select long from storage_urls where short=$1", inputURL)
	if err = row.Scan(&foundURL); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return "", err
		}
		return "", nil
	}
	return foundURL, nil
}

func (d *database) Ping() error {
	if err := d.db.Ping(); err != nil {
		return err
	}
	return nil
}
