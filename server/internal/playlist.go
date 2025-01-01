package internal

import (
	"encoding/json"
	"errors"
	"log/slog"
	"os/exec"
	"slices"
	"strings"
	"time"

	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/config"
)

type metadata struct {
	Entries       []DownloadInfo `json:"entries"`
	Count         int            `json:"playlist_count"`
	PlaylistTitle string         `json:"title"`
	Type          string         `json:"_type"`
}

func PlaylistDetect(req DownloadRequest, mq *MessageQueue, db *MemoryDB) error {
	var (
		downloader = config.Instance().DownloaderPath
		cmd        = exec.Command(downloader, req.URL, "--flat-playlist", "-J")
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	var m metadata

	if err := cmd.Start(); err != nil {
		return err
	}

	slog.Info("decoding playlist metadata", slog.String("url", req.URL))

	if err := json.NewDecoder(stdout).Decode(&m); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	slog.Info("decoded playlist metadata", slog.String("url", req.URL))

	if m.Type == "" {
		return errors.New("probably not a valid URL")
	}

	if m.Type == "playlist" {
		entries := slices.CompactFunc(slices.Compact(m.Entries), func(a DownloadInfo, b DownloadInfo) bool {
			return a.URL == b.URL
		})

		slog.Info("playlist detected", slog.String("url", req.URL), slog.Int("count", len(entries)))

		for i, meta := range entries {
			// detect playlist title from metadata since each playlist entry will be
			// treated as an individual download
			req.Rename = strings.Replace(
				req.Rename,
				"%(playlist_title)s",
				m.PlaylistTitle,
				1,
			)

			//XXX: it's idiotic but it works: virtually delay the creation time
			meta.CreatedAt = time.Now().Add(time.Millisecond * time.Duration(i*10))

			proc := &Process{
				Url:      meta.URL,
				Progress: DownloadProgress{},
				Output:   DownloadOutput{Filename: req.Rename},
				Info:     meta,
				Params:   req.Params,
			}

			proc.Info.URL = meta.URL

			time.Sleep(time.Millisecond)

			db.Set(proc)
			mq.Publish(proc)
		}
	}

	proc := &Process{
		Url:    req.URL,
		Params: req.Params,
	}

	db.Set(proc)
	mq.Publish(proc)
	slog.Info("sending new process to message queue", slog.String("url", proc.Url))

	return cmd.Wait()
}
