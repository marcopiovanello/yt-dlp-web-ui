package status

import (
	"github.com/go-chi/chi/v5"
	"github.com/marcopiovanello/yt-dlp-web-ui/v4/server/internal/kv"
	"github.com/marcopiovanello/yt-dlp-web-ui/v4/server/status/repository"
	"github.com/marcopiovanello/yt-dlp-web-ui/v4/server/status/rest"
	"github.com/marcopiovanello/yt-dlp-web-ui/v4/server/status/service"
)

func ApplyRouter(mdb *kv.Store) func(chi.Router) {
	var (
		r = repository.New(mdb)
		s = service.New(r, nil)
		h = rest.New(s)
	)

	return func(r chi.Router) {
		r.Get("/", h.Status())
	}
}
