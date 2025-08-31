package subscription

import (
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/subscription/domain"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/subscription/task"

	bolt "go.etcd.io/bbolt"
)

func Container(db *bolt.DB, runner task.TaskRunner) domain.RestHandler {
	var (
		r = provideRepository(db)
		s = provideService(r, runner)
		h = provideHandler(s)
	)
	return h
}
