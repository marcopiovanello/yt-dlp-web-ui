package downloaders

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/common"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/config"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal/pipes"
)

type LiveStreamDownloader struct {
	progress internal.DownloadProgress

	proc *os.Process

	logConsumer LogConsumer

	pipes []pipes.Pipe

	// embedded
	DownloaderBase
}

func NewLiveStreamDownloader(url string, pipes []pipes.Pipe) Downloader {
	l := &LiveStreamDownloader{
		logConsumer: NewFFMpegLogConsumer(),
		pipes:       pipes,
	}
	// in base
	l.Id = uuid.NewString()
	l.URL = url
	return l
}

func (l *LiveStreamDownloader) Start() error {
	l.SetPending(true)

	baseParams := []string{
		l.URL,
		"--newline",
		"--no-colors",
		"--no-playlist",
		"--no-exec",
	}

	params := append(baseParams, "-o", "-")

	cmd := exec.Command(config.Instance().Paths.DownloaderPath, params...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// stdout = media stream
	media, err := cmd.StdoutPipe()
	if err != nil {
		slog.Error("failed to get media stdout", slog.Any("err", err))
		panic(err)
	}

	// stderr = log/progress
	stderr, err := cmd.StderrPipe()
	if err != nil {
		slog.Error("failed to get stderr pipe", slog.Any("err", err))
		panic(err)
	}

	if err := cmd.Start(); err != nil {
		slog.Error("failed to start yt-dlp process", slog.Any("err", err))
		panic(err)
	}

	l.proc = cmd.Process

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		l.Complete()
		cancel()
	}()

	// --- costruisci pipeline ---
	reader := io.Reader(media)
	for _, pipe := range l.pipes {
		nr, err := pipe.Connect(reader)
		if err != nil {
			slog.Error("pipe failed", slog.String("pipe", pipe.Name()), slog.Any("err", err))
			return err
		}
		reader = nr
	}

	// --- fallback: se nessun FileWriter, scrivi su file ---
	if !l.hasFileWriter() {
		go func() {
			filepath.Join(
				config.Instance().Paths.DownloadPath,
				fmt.Sprintf("%s (live) %s.mp4", l.Id, time.Now().Format(time.ANSIC)),
			)

			defaultPath := filepath.Join(config.Instance().Paths.DownloadPath)
			f, err := os.Create(defaultPath)
			if err != nil {
				slog.Error("failed to create fallback file", slog.Any("err", err))
				return
			}
			defer f.Close()

			_, err = io.Copy(f, reader)
			if err != nil {
				slog.Error("copy error", slog.Any("err", err))
			}
			slog.Info("download saved", slog.String("path", defaultPath))
		}()
	}

	// --- logs consumer ---
	logs := make(chan []byte)
	go produceLogs(stderr, logs)
	go consumeLogs(ctx, logs, l.logConsumer, l)

	l.progress.Status = internal.StatusLiveStream

	return cmd.Wait()
}

func (l *LiveStreamDownloader) Stop() error {
	defer func() {
		l.progress.Status = internal.StatusCompleted
		l.Complete()
	}()
	// yt-dlp uses multiple child process the parent process
	// has been spawned with setPgid = true. To properly kill
	// all subprocesses a SIGTERM need to be sent to the correct
	// process group
	if l.proc == nil {
		return errors.New("*os.Process not set")
	}

	pgid, err := syscall.Getpgid(l.proc.Pid)
	if err != nil {
		return err
	}
	if err := syscall.Kill(-pgid, syscall.SIGTERM); err != nil {
		return err
	}

	return nil
}

func (l *LiveStreamDownloader) Status() *internal.ProcessSnapshot {
	return &internal.ProcessSnapshot{
		Id:             l.Id,
		Info:           l.Metadata,
		Progress:       l.progress,
		DownloaderName: "livestream",
	}
}

func (l *LiveStreamDownloader) UpdateSavedFilePath(p string) {}

func (l *LiveStreamDownloader) SetOutput(o internal.DownloadOutput)     {}
func (l *LiveStreamDownloader) SetProgress(p internal.DownloadProgress) { l.progress = p }

func (l *LiveStreamDownloader) SetMetadata(fetcher func(url string) (*common.DownloadMetadata, error)) {
	l.FetchMetadata(fetcher)
}

func (l *LiveStreamDownloader) SetPending(p bool) {
	l.Pending = p
}

func (l *LiveStreamDownloader) GetId() string  { return l.Id }
func (l *LiveStreamDownloader) GetUrl() string { return l.URL }

func (l *LiveStreamDownloader) RestoreFromSnapshot(snap *internal.ProcessSnapshot) error {
	if snap == nil {
		return errors.New("cannot restore nil snapshot")
	}

	s := *snap

	l.Id = s.Id
	l.URL = s.Info.URL
	l.Metadata = s.Info
	l.progress = s.Progress

	return nil
}

func (l *LiveStreamDownloader) IsCompleted() bool { return l.Completed }

func (l *LiveStreamDownloader) hasFileWriter() bool {
	return slices.ContainsFunc(l.pipes, func(p pipes.Pipe) bool {
		return p.Name() == "file-writer"
	})
}
