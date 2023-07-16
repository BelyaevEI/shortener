package database

import (
	"database/sql"

	"github.com/BelyaevEI/shortener/internal/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type database struct {
	DBpath string
	db     *sql.DB
}

func New(DBpath string, logger *logger.Logger) *database {
	db, err := sql.Open("pgx", DBpath)
	if err != nil {
		// log.Fatal(err)
		logger.Log.Error(err)
	}

	_, err = db.Exec("create table IF NOT EXISTS storage_urls(short text not null, long text not null)")
	if err != nil {
		// log.Fatal(err)
		logger.Log.Error("Error create tabele", err)
		return nil
	}

	return &database{
		DBpath: DBpath,
		db:     db,
	}
}

func (d *database) Save(url1, url2 string) error {
	_, err := d.db.Exec("insert into storage_urls(short, long) values ($1, $2)", url1, url2)
	if err != nil {
		return err
	}
	return nil
}

func (d *database) Get(inputURL string) string {

	var foundURL string

	row1 := d.db.QueryRow("select long from storage_urls where short=$1", inputURL)
	if err := row1.Scan(&foundURL); err == nil {
		return foundURL
	}
	row2 := d.db.QueryRow("select short from storage_urls where long=$1", inputURL)
	if err := row2.Scan(&foundURL); err == nil {
		return foundURL
	}
	return ""
}

func (d *database) Ping() error {
	if err := d.db.Ping(); err != nil {
		return err
	}
	return nil
}
