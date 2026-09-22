package models

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitDB(dbPath string) error {
	var err error
	// DSN 级 PRAGMA：对连接池中每条新连接生效（busy_timeout/foreign_keys 是连接级设置）
	dsn := dbPath + "?_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=synchronous(NORMAL)"
	DB, err = sql.Open("sqlite", dsn)
	if err != nil {
		return err
	}

	// SQLite 单写者：限制连接数避免 SQLITE_BUSY；WAL 下读写仍可并行
	DB.SetMaxOpenConns(1)

	DB.Exec("PRAGMA journal_mode=WAL")

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
			content_type TEXT NOT NULL CHECK(content_type IN ('image', 'video', 'audio', 'mixed')),
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
			media_type TEXT NOT NULL CHECK(media_type IN ('image', 'video', 'audio')),
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
		CREATE TABLE IF NOT EXISTS experiences (
			id          INTEGER PRIMARY KEY AUTOINCREMENT,
			period      TEXT NOT NULL,
			title       TEXT NOT NULL,
			description TEXT DEFAULT '',
			sort_order  REAL NOT NULL DEFAULT 0,
			created_at  TEXT DEFAULT (datetime('now','localtime'))
		);
	`)
	if err != nil {
		return err
	}

	// 外键索引（SQLite 不自动为 FK 建索引）
	if _, err = DB.Exec(`
		CREATE INDEX IF NOT EXISTS idx_chat_messages_session ON chat_messages(session_id, id);
		CREATE INDEX IF NOT EXISTS idx_work_media_work ON work_media(work_id, sort_order);
		CREATE INDEX IF NOT EXISTS idx_work_media_file ON work_media(file_id);
		CREATE INDEX IF NOT EXISTS idx_files_category ON files(category_id);
		CREATE INDEX IF NOT EXISTS idx_works_category ON works(category_id);
	`); err != nil {
		return err
	}

	// Migrate: add 'audio' to works.content_type and work_media.media_type CHECK constraints
	migrateCheckConstraints()

	return nil
}

func migrateCheckConstraints() {
	// Check if works table has old CHECK constraint (without 'audio')
	var ck string
	err := DB.QueryRow("SELECT sql FROM sqlite_master WHERE type='table' AND name='works'").Scan(&ck)
	if err == nil && !containsAudio(ck) {
		// Recreate works table with updated CHECK
		DB.Exec("ALTER TABLE works RENAME TO works_old")
		DB.Exec(`CREATE TABLE works (
			id           INTEGER PRIMARY KEY AUTOINCREMENT,
			category_id  INTEGER,
			title        TEXT NOT NULL,
			description  TEXT DEFAULT '',
			cover_url    TEXT DEFAULT '',
			content_type TEXT NOT NULL CHECK(content_type IN ('image', 'video', 'audio', 'mixed')),
			external_url TEXT DEFAULT '',
			sort_order   REAL NOT NULL DEFAULT 0,
			created_at   TEXT DEFAULT (datetime('now','localtime')),
			FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL
		)`)
		DB.Exec("INSERT INTO works SELECT * FROM works_old")
		DB.Exec("DROP TABLE works_old")
	}

	// Check work_media table
	err = DB.QueryRow("SELECT sql FROM sqlite_master WHERE type='table' AND name='work_media'").Scan(&ck)
	if err == nil && !containsAudio(ck) {
		DB.Exec("ALTER TABLE work_media RENAME TO work_media_old")
		DB.Exec(`CREATE TABLE work_media (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			work_id    INTEGER NOT NULL,
			file_id    INTEGER,
			media_url  TEXT NOT NULL,
			media_type TEXT NOT NULL CHECK(media_type IN ('image', 'video', 'audio')),
			sort_order REAL NOT NULL DEFAULT 0,
			FOREIGN KEY (work_id) REFERENCES works(id) ON DELETE CASCADE,
			FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE SET NULL
		)`)
		DB.Exec("INSERT INTO work_media SELECT * FROM work_media_old")
		DB.Exec("DROP TABLE work_media_old")
	}
}

func containsAudio(s string) bool {
	return len(s) > 0 && (contains(s, "'audio'") || contains(s, "\"audio\""))
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
