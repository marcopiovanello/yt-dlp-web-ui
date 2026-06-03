package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/subscription/data"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/subscription/domain"
)

type Repository struct {
	db *sql.DB
}

// Delete implements domain.Repository.
func (r *Repository) Delete(ctx context.Context, id string) error {
	conn, err := r.db.Conn(ctx)
	if err != nil {
		return err
	}

	defer conn.Close()

	_, err = conn.ExecContext(ctx, "DELETE FROM subscriptions WHERE id = ?", id)

	return err
}

// GetCursor implements domain.Repository.
func (r *Repository) GetCursor(ctx context.Context, id string) (int64, error) {
	conn, err := r.db.Conn(ctx)
	if err != nil {
		return -1, err
	}

	defer conn.Close()

	row := conn.QueryRowContext(ctx, "SELECT rowid FROM subscriptions WHERE id = ?", id)

	var rowId int64

	if err := row.Scan(&rowId); err != nil {
		return -1, err
	}

	return rowId, nil
}

// List implements domain.Repository.
func (r *Repository) List(ctx context.Context, start int64, limit int) (*[]data.Subscription, error) {
	conn, err := r.db.Conn(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	var elements []data.Subscription

	rows, err := conn.QueryContext(ctx, "SELECT rowid, * FROM subscriptions WHERE rowid > ? LIMIT ?", start, limit)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var rowId int64
		var element data.Subscription

		if err := rows.Scan(
			&rowId,
			&element.Id,
			&element.URL,
			&element.Params,
			&element.CronExpr,
		); err != nil {
			return &elements, err
		}

		elements = append(elements, element)
	}

	return &elements, nil
}

// Submit implements domain.Repository.
func (r *Repository) Submit(ctx context.Context, sub *data.Subscription) (*data.Subscription, error) {
	conn, err := r.db.Conn(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	_, err = conn.ExecContext(
		ctx,
		"INSERT INTO subscriptions (id, url, params, cron) VALUES (?, ?, ?, ?)",
		uuid.NewString(),
		sub.URL,
		sub.Params,
		sub.CronExpr,
	)

	return sub, err
}

// UpdateByExample implements domain.Repository.
func (r *Repository) UpdateByExample(ctx context.Context, example *data.Subscription) error {
	conn, err := r.db.Conn(ctx)
	if err != nil {
		return err
	}

	defer conn.Close()

	_, err = conn.ExecContext(
		ctx,
		"UPDATE subscriptions SET url = ?, params = ?, cron = ? WHERE id = ? OR url = ?",
		example.URL,
		example.Params,
		example.CronExpr,
		example.Id,
		example.URL,
	)

	return err
}

func New(db *sql.DB) domain.Repository {
	return &Repository{
		db: db,
	}
}
