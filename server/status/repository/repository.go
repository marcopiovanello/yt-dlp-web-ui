package repository

import (
	"context"
	"slices"

	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal/kv"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/status/domain"
)

type Repository struct {
	mdb *kv.Store
}

// DownloadSpeed implements domain.Repository.
func (r *Repository) DownloadSpeed(ctx context.Context) int64 {
	processes := r.mdb.All()

	var downloadSpeed float64

	for _, p := range *processes {
		downloadSpeed += p.Progress.Speed
	}

	return int64(downloadSpeed)
}

// Completed implements domain.Repository.
func (r *Repository) Completed(ctx context.Context) int {
	processes := r.mdb.All()

	completed := slices.DeleteFunc(*processes, func(p internal.ProcessSnapshot) bool {
		return p.Progress.Status != internal.StatusCompleted
	})

	return len(completed)
}

// Downloading implements domain.Repository.
func (r *Repository) Downloading(ctx context.Context) int {
	processes := r.mdb.All()

	downloading := slices.DeleteFunc(*processes, func(p internal.ProcessSnapshot) bool {
		return p.Progress.Status != internal.StatusDownloading
	})

	return len(downloading)
}

// Pending implements domain.Repository.
func (r *Repository) Pending(ctx context.Context) int {
	processes := r.mdb.All()

	pending := slices.DeleteFunc(*processes, func(p internal.ProcessSnapshot) bool {
		return p.Progress.Status != internal.StatusPending
	})

	return len(pending)
}

func New(mdb *kv.Store) domain.Repository {
	return &Repository{
		mdb: mdb,
	}
}
