package archiver

import (
	"context"
	"database/sql"
	"log/slog"

	evbus "github.com/asaskevich/EventBus"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/archive"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/config"
)

const QueueName = "process:archive"

var (
	eventBus       = evbus.New()
	archiveService archive.Service
)

type Message = archive.Entity

func Register(db *sql.DB) {
	_, s := archive.Container(db)
	archiveService = s
}

func init() {
	eventBus.Subscribe(QueueName, func(m *Message) {
		slog.Info(
			"archiving completed download",
			slog.String("title", m.Title),
			slog.String("source", m.Source),
		)
		archiveService.Archive(context.Background(), m)
	})
}

func Publish(m *Message) {
	if config.Instance().AutoArchive {
		eventBus.Publish(QueueName, m)
	}
}
