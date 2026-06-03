package repository

import (
	"context"
	"database/sql"
	"os"

	"github.com/google/uuid"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/archive/data"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/archive/domain"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) domain.Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Archive(ctx context.Context, entry *data.ArchiveEntry) error {
	conn, err := r.db.Conn(ctx)
	if err != nil {
		return err
	}

	defer conn.Close()

	_, err = conn.ExecContext(
		ctx,
		"INSERT INTO archive (id, title, path, thumbnail, source, metadata, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		uuid.NewString(),
		entry.Title,
		entry.Path,
		entry.Thumbnail,
		entry.Source,
		entry.Metadata,
		entry.CreatedAt,
	)

	return err
}

func (r *Repository) SoftDelete(ctx context.Context, id string) (*data.ArchiveEntry, error) {
	conn, err := r.db.Conn(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	tx, err := conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var model data.ArchiveEntry

	row := tx.QueryRowContext(ctx, "SELECT * FROM archive WHERE id = ?", id)

	if err := row.Scan(
		&model.Id,
		&model.Title,
		&model.Path,
		&model.Thumbnail,
		&model.Source,
		&model.Metadata,
		&model.CreatedAt,
	); err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM archive WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &model, nil
}

func (r *Repository) HardDelete(ctx context.Context, id string) (*data.ArchiveEntry, error) {
	entry, err := r.SoftDelete(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := os.Remove(entry.Path); err != nil {
		return nil, err
	}

	return entry, nil
}

func (r *Repository) List(ctx context.Context, startRowId int, limit int) (*[]data.ArchiveEntry, error) {
	conn, err := r.db.Conn(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	var entries []data.ArchiveEntry

	// cursor based pagination
	rows, err := conn.QueryContext(ctx, "SELECT rowid, * FROM archive WHERE rowid > ? LIMIT ?", startRowId, limit)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var rowId int64
		var entry data.ArchiveEntry

		if err := rows.Scan(
			&rowId,
			&entry.Id,
			&entry.Title,
			&entry.Path,
			&entry.Thumbnail,
			&entry.Source,
			&entry.Metadata,
			&entry.CreatedAt,
		); err != nil {
			return &entries, err
		}

		entries = append(entries, entry)
	}

	return &entries, err
}

func (r *Repository) GetCursor(ctx context.Context, id string) (int64, error) {
	conn, err := r.db.Conn(ctx)
	if err != nil {
		return -1, err
	}
	defer conn.Close()

	row := conn.QueryRowContext(ctx, "SELECT rowid FROM archive WHERE id = ?", id)

	var rowId int64

	if err := row.Scan(&rowId); err != nil {
		return -1, err
	}

	return rowId, nil
}
