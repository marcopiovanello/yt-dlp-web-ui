package status

import (
	"github.com/go-chi/chi/v5"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/status/repository"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/status/rest"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/status/service"
)

func ApplyRouter(mdb *internal.MemoryDB) func(chi.Router) {
	var (
		r = repository.New(mdb)
		s = service.New(r, nil) //TODO: nil, wtf?
		h = rest.New(s)
	)

	return func(r chi.Router) {
		r.Get("/", h.Status())
	}
}
