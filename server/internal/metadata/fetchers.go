package metadata

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"os/exec"
	"syscall"
	"time"

	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/common"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/config"
)

func DefaultFetcher(url string) (*common.DownloadMetadata, error) {
	cmd := exec.Command(config.Instance().DownloaderPath, url, "-J")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	meta := common.DownloadMetadata{
		URL:       url,
		CreatedAt: time.Now(),
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	var bufferedStderr bytes.Buffer

	go func() {
		io.Copy(&bufferedStderr, stderr)
	}()

	slog.Info("retrieving metadata", slog.String("url", url))

	if err := json.NewDecoder(stdout).Decode(&meta); err != nil {
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		return nil, errors.New(bufferedStderr.String())
	}

	return &meta, nil
}
