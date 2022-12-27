package sqlite

import (
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type DB struct {
	DSN            string
	MigrationsPath string
	db             *sql.DB
}

func NewDB(dsn string, migrationsPath string) *DB {
	return &DB{
		DSN:            dsn,
		MigrationsPath: migrationsPath,
	}
}

func (db *DB) Open() error {
	if db.DSN == "" {
		return errors.New("database dsn required.")
	}

	if db.DSN == ":memory:" {
		return errors.New("db should be persistent")
	}

	dbSQL, err := sql.Open("sqlite3", db.DSN)
	if err != nil {
		return err
	}
	db.db = dbSQL

	if _, err := db.db.Exec(`PRAGMA journal_mode = wal;`); err != nil {
		return err
	}

	driver, err := sqlite3.WithInstance(dbSQL, &sqlite3.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		db.MigrationsPath,
		"students",
		driver,
	)
	if err != nil {
		return err
	}

	return m.Up()
}
