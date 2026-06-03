package domain

import (
	"context"
	"net/http"
)

type Status struct {
	Downloading   int `json:"downloading"`
	Pending       int `json:"pending"`
	Completed     int `json:"completed"`
	DownloadSpeed int `json:"download_speed"`
}

type Repository interface {
	Pending(ctx context.Context) int
	Completed(ctx context.Context) int
	Downloading(ctx context.Context) int
	DownloadSpeed(ctx context.Context) int64
}

type Service interface {
	Status(ctx context.Context) (*Status, error)
}

type RestHandler interface {
	Status() http.HandlerFunc
}
