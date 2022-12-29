package sqlite

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/mattn/go-sqlite3"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

//go:embed subjects_data.sql
var data string

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

	if _, err := db.db.Exec(`PRAGMA foreign_keys = on;`); err != nil {
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

	if err := m.Up(); err != nil {
		return err
	}

	return db.populateSubjects()
}

func (db *DB) populateSubjects() error {
	conn, err := db.db.Conn(context.Background())
	if err != nil {
		return err
	}
	defer conn.Close()

	resp, err := conn.QueryContext(context.Background(), `SELECT COUNT(*) FROM subjects`)
	if err != nil {
		return err
	}
	defer resp.Close()

	var n int
	if err := resp.Scan(&n); err != nil {
		return err
	}

	if n != 0 {
		return nil
	}

	_, err = conn.ExecContext(context.Background(), data)
	return err
}
