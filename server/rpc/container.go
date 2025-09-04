package rpc

import (
	"github.com/go-chi/chi/v5"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/config"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal/kv"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal/livestream"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal/queue"
	middlewares "github.com/marcopiovanello/yt-dlp-web-ui/v3/server/middleware"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/openid"
)

// Dependency injection container.
func Container(db *kv.Store, mq *queue.MessageQueue, lm *livestream.Monitor) *Service {
	return &Service{
		db: db,
		mq: mq,
		lm: lm,
	}
}

// RPC service must be registered before applying this router!
func ApplyRouter() func(chi.Router) {
	return func(r chi.Router) {
		if config.Instance().Authentication.RequireAuth {
			r.Use(middlewares.Authenticated)
		}
		if config.Instance().OpenId.UseOpenId {
			r.Use(openid.Middleware)
		}
		r.Get("/ws", WebSocket)
		r.Post("/http", Post)
	}
}
