package rest

import (
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal/kv"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal/livestream"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal/queue"

	bolt "go.etcd.io/bbolt"
)

type ContainerArgs struct {
	DB  *bolt.DB
	MDB *kv.Store
	MQ  *queue.MessageQueue
	LM  *livestream.Monitor
}
