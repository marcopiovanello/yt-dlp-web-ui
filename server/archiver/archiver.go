package archiver

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/archive"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/config"
)

var (
	ch             = make(chan *Message, 1)
	archiveService archive.Service
)

type Message = archive.Entity

func Register(db *sql.DB) {
	_, s := archive.Container(db)
	archiveService = s
}

func init() {
	go func() {
		for m := range ch {
			slog.Info(
				"archiving completed download",
				slog.String("title", m.Title),
				slog.String("source", m.Source),
			)
			archiveService.Archive(context.Background(), m)
		}
	}()
}

func Publish(m *Message) {
	if config.Instance().AutoArchive {
		ch <- m
	}
}
