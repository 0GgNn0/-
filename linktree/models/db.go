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
		CREATE TABLE IF NOT EXISTS categories (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			name       TEXT NOT NULL,
			type       TEXT NOT NULL CHECK(type IN ('file', 'work')),
			sort_order REAL NOT NULL DEFAULT 0,
			created_at TEXT DEFAULT (datetime('now','localtime'))
		);
		CREATE TABLE IF NOT EXISTS files (
			id             INTEGER PRIMARY KEY AUTOINCREMENT,
			category_id    INTEGER,
			original_name  TEXT NOT NULL,
			stored_name    TEXT NOT NULL,
			mime_type      TEXT NOT NULL,
			size           INTEGER NOT NULL,
			is_public      INTEGER NOT NULL DEFAULT 1,
			download_count INTEGER NOT NULL DEFAULT 0,
			sort_order     REAL NOT NULL DEFAULT 0,
			created_at     TEXT DEFAULT (datetime('now','localtime')),
			FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
		);
		CREATE TABLE IF NOT EXISTS works (
			id           INTEGER PRIMARY KEY AUTOINCREMENT,
			category_id  INTEGER,
			title        TEXT NOT NULL,
			description  TEXT DEFAULT '',
			cover_url    TEXT DEFAULT '',
			content_type TEXT NOT NULL CHECK(content_type IN ('image', 'video', 'mixed')),
			external_url TEXT DEFAULT '',
			sort_order   REAL NOT NULL DEFAULT 0,
			created_at   TEXT DEFAULT (datetime('now','localtime')),
			FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
		);
		CREATE TABLE IF NOT EXISTS work_media (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			work_id    INTEGER NOT NULL,
			file_id    INTEGER,
			media_url  TEXT NOT NULL,
			media_type TEXT NOT NULL CHECK(media_type IN ('image', 'video')),
			sort_order REAL NOT NULL DEFAULT 0,
			FOREIGN KEY (work_id) REFERENCES works(id) ON DELETE CASCADE,
			FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE SET NULL
		);
		CREATE TABLE IF NOT EXISTS chat_sessions (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			title      TEXT NOT NULL DEFAULT '新对话',
			model      TEXT DEFAULT '',
			created_at TEXT DEFAULT (datetime('now','localtime')),
			updated_at TEXT DEFAULT (datetime('now','localtime'))
		);
		CREATE TABLE IF NOT EXISTS chat_messages (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			session_id INTEGER NOT NULL,
			role       TEXT NOT NULL CHECK(role IN ('user', 'assistant', 'system')),
			content    TEXT NOT NULL,
			created_at TEXT DEFAULT (datetime('now','localtime')),
			FOREIGN KEY (session_id) REFERENCES chat_sessions(id) ON DELETE CASCADE
		);
	`)
	return err
}
