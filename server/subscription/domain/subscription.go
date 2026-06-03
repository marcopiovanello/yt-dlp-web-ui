package domain

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/subscription/data"
)

type Subscription struct {
	Id       string `json:"id"`
	URL      string `json:"url"`
	Params   string `json:"params"`
	CronExpr string `json:"cron_expression"`
}

type PaginatedResponse[T any] struct {
	First int64 `json:"first"`
	Next  int64 `json:"next"`
	Data  T     `json:"data"`
}

type Repository interface {
	Submit(ctx context.Context, sub *data.Subscription) (*data.Subscription, error)
	List(ctx context.Context, start int64, limit int) (*[]data.Subscription, error)
	UpdateByExample(ctx context.Context, example *data.Subscription) error
	Delete(ctx context.Context, id string) error
	GetCursor(ctx context.Context, id string) (int64, error)
}

type Service interface {
	Submit(ctx context.Context, sub *Subscription) (*Subscription, error)
	List(ctx context.Context, start int64, limit int) (*PaginatedResponse[[]Subscription], error)
	UpdateByExample(ctx context.Context, example *Subscription) error
	Delete(ctx context.Context, id string) error
	GetCursor(ctx context.Context, id string) (int64, error)
}

type RestHandler interface {
	Submit() http.HandlerFunc
	List() http.HandlerFunc
	UpdateByExample() http.HandlerFunc
	Delete() http.HandlerFunc
	GetCursor() http.HandlerFunc
	ApplyRouter() func(chi.Router)
}
