package database

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type database struct {
	DBpath string
	db     *sql.DB
}

func New(DBpath string) *database {
	db, err := sql.Open("pgx", DBpath)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("create table IF NOT EXISTS storage_urls(short text not null, long text not null)")
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return &database{DBpath: DBpath,
		db: db}
}

func (d *database) Save(url1, url2 string) error {
	panic("No Implement")
}

func (d *database) Get(inputURL string) string {
	panic("No Implement")
}

func (d *database) Ping() error {
	if err := d.db.Ping(); err != nil {
		return err
	}
	return nil
}
