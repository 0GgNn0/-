package models

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB(dbPath string) error {
	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	DB.Exec("PRAGMA journal_mode=WAL")
	DB.Exec("PRAGMA foreign_keys=ON")

	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS links (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			title      TEXT NOT NULL,
			url        TEXT NOT NULL,
			sort_order REAL NOT NULL DEFAULT 0,
			icon       TEXT DEFAULT '',
			created_at TEXT DEFAULT (datetime('now','localtime'))
		);
		CREATE TABLE IF NOT EXISTS settings (
			key   TEXT PRIMARY KEY,
			value TEXT NOT NULL
		);
	`)
	return err
}
