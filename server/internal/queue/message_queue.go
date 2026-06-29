package queue

import (
	"context"
	"errors"
	"log/slog"

	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/config"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal/downloaders"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal/metadata"
)

type MessageQueue struct {
	concurrency   int
	downloadQueue chan downloaders.Downloader
	metadataQueue chan downloaders.Downloader
	ctx           context.Context
	cancel        context.CancelFunc
}

func NewMessageQueue() (*MessageQueue, error) {
	qs := config.Instance().Server.QueueSize
	if qs <= 0 {
		return nil, errors.New("invalid queue size")
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &MessageQueue{
		concurrency:   qs,
		downloadQueue: make(chan downloaders.Downloader, qs*2),
		metadataQueue: make(chan downloaders.Downloader, qs*4),
		ctx:           ctx,
		cancel:        cancel,
	}, nil
}

// Publish download job
func (m *MessageQueue) Publish(d downloaders.Downloader) {
	d.SetPending(true)

	select {
	case m.downloadQueue <- d:
		slog.Info("published download", slog.String("id", d.GetId()))
	case <-m.ctx.Done():
		slog.Warn("queue stopped, dropping download", slog.String("id", d.GetId()))
	}
}

// Workers: download + metadata
func (m *MessageQueue) SetupConsumers() {
	// N parallel workers for downloadQueue
	for i := 1; i < m.concurrency+1; i++ {
		go m.downloadWorker(i)
	}

	// 1 serial worker for metadata
	go m.metadataWorker()
}

// Worker dei download
func (m *MessageQueue) downloadWorker(workerId int) {
	slog.Info("download worker spawned", slog.Int("worker", workerId))

	for {
		select {
		case <-m.ctx.Done():
			return
		case p := <-m.downloadQueue:
			if p == nil {
				continue
			}
			if p.IsCompleted() {
				continue
			}

			slog.Info("download worker starting download",
				slog.Int("worker", workerId),
				slog.String("id", p.GetId()),
			)

			m.metadataQueue <- p
			slog.Info("queued for metadata", slog.String("id", p.GetId()))

			p.Start()
		}
	}
}

func (m *MessageQueue) metadataWorker() {
	slog.Info("metadata worker spawned", slog.Int("worker", 1))

	for {
		select {
		case <-m.ctx.Done():
			return
		case p := <-m.metadataQueue:
			if p == nil {
				continue
			}

			slog.Info("metadata worker started",
				slog.String("id", p.GetId()),
			)

			if p.IsCompleted() {
				slog.Warn("metadata skipped, illegal state",
					slog.String("id", p.GetId()),
				)
				continue
			}

			slog.Info("metadata worker started", slog.String("id", p.GetId()))
			p.SetMetadata(metadata.DefaultFetcher)
		}
	}
}

func (m *MessageQueue) Stop() {
	m.cancel()
	close(m.downloadQueue)
	close(m.metadataQueue)
}
