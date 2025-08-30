package downloaders

import (
	"encoding/json"
	"log/slog"
	"strings"

	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal"
)

type LogConsumer interface {
	GetName() string
	ParseLogEntry(entry []byte, downloader Downloader)
}

type JSONLogConsumer struct{}

func NewJSONLogConsumer() LogConsumer {
	return &JSONLogConsumer{}
}

func (j *JSONLogConsumer) GetName() string { return "json-log-consumer" }

func (j *JSONLogConsumer) ParseLogEntry(entry []byte, d Downloader) {
	var progress internal.ProgressTemplate
	var postprocess internal.PostprocessTemplate

	if err := json.Unmarshal(entry, &progress); err == nil {
		d.SetProgress(internal.DownloadProgress{
			Status:     internal.StatusDownloading,
			Percentage: progress.Percentage,
			Speed:      progress.Speed,
			ETA:        progress.Eta,
		})

		slog.Info("progress",
			slog.String("id", j.GetShortId(d.GetId())),
			slog.String("url", d.GetUrl()),
			slog.String("percentage", progress.Percentage),
		)
	}

	if err := json.Unmarshal(entry, &postprocess); err == nil {
		d.UpdateSavedFilePath(postprocess.FilePath)
	}
}

func (j *JSONLogConsumer) GetShortId(id string) string {
	return strings.Split(id, "-")[0]
}

//TODO: split in different files

type FFMpegLogConsumer struct{}

func NewFFMpegLogConsumer() LogConsumer {
	return &JSONLogConsumer{}
}

func (f *FFMpegLogConsumer) GetName() string { return "ffmpeg-log-consumer" }

func (f *FFMpegLogConsumer) ParseLogEntry(entry []byte, d Downloader) {
	slog.Info("ffmpeg output",
		slog.String("id", d.GetId()),
		slog.String("url", d.GetUrl()),
		slog.String("output", string(entry)),
	)
}
