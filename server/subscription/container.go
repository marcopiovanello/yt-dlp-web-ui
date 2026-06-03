package subscription

import (
	"database/sql"

	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/subscription/domain"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/subscription/task"
)

func Container(db *sql.DB, runner task.TaskRunner) domain.RestHandler {
	var (
		r = provideRepository(db)
		s = provideService(r, runner)
		h = provideHandler(s)
	)
	return h
}
