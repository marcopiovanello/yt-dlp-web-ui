package archive

import (
	"database/sql"

	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/archive/domain"
)

func Container(db *sql.DB) (domain.RestHandler, domain.Service) {
	var (
		r = provideRepository(db)
		s = provideService(r)
		h = provideHandler(s)
	)
	return h, s
}
