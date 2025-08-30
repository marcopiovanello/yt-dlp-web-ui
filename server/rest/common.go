package rest

import (
	"database/sql"

	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal/kv"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal/queue"
)

type ContainerArgs struct {
	DB  *sql.DB
	MDB *kv.Store
	MQ  *queue.MessageQueue
}
