package dbutil

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"

	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/config"
)

var lockFilePath = filepath.Join(config.Instance().Dir(), ".db.lock")

// Run the table migration
func Migrate(ctx context.Context, db *sql.DB) error {
	conn, err := db.Conn(ctx)
	if err != nil {
		return err
	}

	defer func() {
		conn.Close()
		createLockFile()
	}()

	if _, err := db.ExecContext(
		ctx,
		`CREATE TABLE IF NOT EXISTS templates (
			id CHAR(36) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			content TEXT NOT NULL
		)`,
	); err != nil {
		return err
	}

	if _, err := db.ExecContext(
		ctx,
		`CREATE TABLE IF NOT EXISTS archive (
			id CHAR(36) PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			path VARCHAR(255) NOT NULL,
			thumbnail TEXT,
			source VARCHAR(255),
			metadata TEXT,
			created_at DATETIME
		)`,
	); err != nil {
		return err
	}

	if _, err := db.ExecContext(
		ctx,
		`CREATE TABLE IF NOT EXISTS subscriptions (
			id CHAR(36) PRIMARY KEY,
			url VARCHAR(2048) UNIQUE NOT NULL,
			params TEXT NOT NULL,
			cron TEXT
		)`,
	); err != nil {
		return err
	}

	if lockFileExists() {
		return nil
	}

	db.ExecContext(
		ctx,
		`INSERT INTO templates (id, name, content) VALUES
			($1, $2, $3),
			($4, $5, $6);`,
		"0", "default", "--no-mtime",
		"1", "audio only", "-x",
	)

	return nil
}

func createLockFile() { os.Create(lockFilePath) }

func lockFileExists() bool {
	_, err := os.Stat(lockFilePath)
	return os.IsExist(err)
}
