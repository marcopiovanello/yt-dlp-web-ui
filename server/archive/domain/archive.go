package domain

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/archive/data"
)

type ArchiveEntry struct {
	Id        string    `json:"id"`
	Title     string    `json:"title"`
	Path      string    `json:"path"`
	Thumbnail string    `json:"thumbnail"`
	Source    string    `json:"source"`
	Metadata  string    `json:"metadata"`
	CreatedAt time.Time `json:"created_at"`
}

type PaginatedResponse[T any] struct {
	First int64 `json:"first"`
	Next  int64 `json:"next"`
	Data  T     `json:"data"`
}

type Repository interface {
	Archive(ctx context.Context, model *data.ArchiveEntry) error
	SoftDelete(ctx context.Context, id string) (*data.ArchiveEntry, error)
	HardDelete(ctx context.Context, id string) (*data.ArchiveEntry, error)
	List(ctx context.Context, startRowId int, limit int) (*[]data.ArchiveEntry, error)
	GetCursor(ctx context.Context, id string) (int64, error)
}

type Service interface {
	Archive(ctx context.Context, entity *ArchiveEntry) error
	SoftDelete(ctx context.Context, id string) (*ArchiveEntry, error)
	HardDelete(ctx context.Context, id string) (*ArchiveEntry, error)
	List(ctx context.Context, startRowId int, limit int) (*PaginatedResponse[[]ArchiveEntry], error)
	GetCursor(ctx context.Context, id string) (int64, error)
}

type RestHandler interface {
	List() http.HandlerFunc
	Archive() http.HandlerFunc
	SoftDelete() http.HandlerFunc
	HardDelete() http.HandlerFunc
	GetCursor() http.HandlerFunc
	ApplyRouter() func(chi.Router)
}
