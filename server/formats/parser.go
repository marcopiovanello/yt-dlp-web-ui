package formats

import (
	"encoding/json"
	"log/slog"
	"os/exec"
	"sync"

	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/config"
)

func ParseURL(url string) (*Metadata, error) {
	cmd := exec.Command(config.Instance().Paths.DownloaderPath, url, "-J")

	stdout, err := cmd.Output()
	if err != nil {
		slog.Error("failed to retrieve metadata", slog.String("err", err.Error()))
		return nil, err
	}

	slog.Info(
		"retrieving metadata",
		slog.String("caller", "getFormats"),
		slog.String("url", url),
	)

	info := &Metadata{URL: url}
	best := &Format{}

	var (
		wg            sync.WaitGroup
		decodingError error
	)

	wg.Add(2)

	go func() {
		decodingError = json.Unmarshal(stdout, &info)
		wg.Done()
	}()

	go func() {
		decodingError = json.Unmarshal(stdout, &best)
		wg.Done()
	}()

	wg.Wait()

	if decodingError != nil {
		return nil, err
	}

	info.Best = *best

	return info, nil
}
