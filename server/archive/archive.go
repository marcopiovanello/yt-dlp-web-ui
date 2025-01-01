package archive

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/archive/domain"
)

// alias type
// TODO: remove after refactoring
type Service = domain.Service
type Entity = domain.ArchiveEntry

func ApplyRouter(db *sql.DB) func(chi.Router) {
	handler, _ := Container(db)
	return handler.ApplyRouter()
}
