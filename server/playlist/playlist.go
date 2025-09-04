package playlist

import (
	"encoding/json"
	"errors"
	"log/slog"
	"os/exec"
	"slices"
	"strings"
	"time"

	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/common"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/config"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal/downloaders"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal/kv"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal/queue"
)

func PlaylistDetect(req internal.DownloadRequest, mq *queue.MessageQueue, db *kv.Store) error {
	params := append(req.Params, "--flat-playlist", "-J")
	urlWithParams := append([]string{req.URL}, params...)

	var (
		downloader = config.Instance().Paths.DownloaderPath
		cmd        = exec.Command(downloader, urlWithParams...)
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	var m Metadata

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

	if m.IsPlaylist() {
		entries := slices.CompactFunc(slices.Compact(m.Entries), func(a common.DownloadMetadata, b common.DownloadMetadata) bool {
			return a.URL == b.URL
		})

		entries = slices.DeleteFunc(entries, func(e common.DownloadMetadata) bool {
			return strings.Contains(e.URL, "list=")
		})

		slog.Info("playlist detected", slog.String("url", req.URL), slog.Int("count", len(entries)))

		if err := ApplyModifiers(&entries, req.Params); err != nil {
			return err
		}

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

			downloader := downloaders.NewGenericDownload(meta.URL, req.Params)
			downloader.SetOutput(internal.DownloadOutput{Filename: req.Rename})
			// downloader.SetMetadata(meta)

			db.Set(downloader)
			mq.Publish(downloader)
		}

		return nil
	}

	d := downloaders.NewGenericDownload(req.URL, req.Params)

	db.Set(d)
	mq.Publish(d)
	slog.Info("sending new process to message queue", slog.String("url", d.GetUrl()))

	return cmd.Wait()
}
