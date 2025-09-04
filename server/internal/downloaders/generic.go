package downloaders

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"slices"
	"strings"
	"syscall"

	"github.com/google/uuid"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/common"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/config"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal"
)

const downloadTemplate = `download:
{
	"eta":%(progress.eta)s,
	"percentage":"%(progress._percent_str)s",
	"speed":%(progress.speed)s
}`

// filename not returning the correct extension after postprocess
const postprocessTemplate = `postprocess:
{
	"filepath":"%(info.filepath)s"
}
`

type GenericDownloader struct {
	Params []string

	AutoRemove bool

	progress internal.DownloadProgress
	output   internal.DownloadOutput

	proc *os.Process

	logConsumer LogConsumer

	// embedded
	DownloaderBase
}

func NewGenericDownload(url string, params []string) Downloader {
	g := &GenericDownloader{
		logConsumer: NewJSONLogConsumer(),
	}
	// in base
	g.Id = uuid.NewString()
	g.URL = url
	return g
}

func (g *GenericDownloader) Start() error {
	g.SetPending(true)

	g.Params = argsSanitizer(g.Params)

	out := internal.DownloadOutput{
		Path:     config.Instance().Paths.DownloadPath,
		Filename: "%(title)s.%(ext)s",
	}

	if g.output.Path != "" {
		out.Path = g.output.Path
	}

	if g.output.Filename != "" {
		out.Filename = g.output.Filename
	}

	buildFilename(&g.output)

	templateReplacer := strings.NewReplacer("\n", "", "\t", "", " ", "")

	baseParams := []string{
		strings.Split(g.URL, "?list")[0], //no playlist
		"--newline",
		"--no-colors",
		"--no-playlist",
		"--progress-template",
		templateReplacer.Replace(downloadTemplate),
		"--progress-template",
		templateReplacer.Replace(postprocessTemplate),
		"--no-exec",
	}

	// if user asked to manually override the output path...
	if !(slices.Contains(g.Params, "-P") || slices.Contains(g.Params, "--paths")) {
		g.Params = append(g.Params, "-o")
		g.Params = append(g.Params, fmt.Sprintf("%s/%s", out.Path, out.Filename))
	}

	params := append(baseParams, g.Params...)

	slog.Info("requesting download", slog.String("url", g.URL), slog.Any("params", params))

	cmd := exec.Command(config.Instance().Paths.DownloaderPath, params...)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		slog.Error("failed to get a stdout pipe", slog.Any("err", err))
		panic(err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		slog.Error("failed to get a stderr pipe", slog.Any("err", err))
		panic(err)
	}

	if err := cmd.Start(); err != nil {
		slog.Error("failed to start yt-dlp process", slog.Any("err", err))
		panic(err)
	}

	g.proc = cmd.Process

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		stdout.Close()
		g.Complete()
		cancel()
	}()

	logs := make(chan []byte)
	go produceLogs(stdout, logs)
	go consumeLogs(ctx, logs, g.logConsumer, g)

	go printYtDlpErrors(stderr, g.Id, g.URL)

	g.SetPending(false)
	return cmd.Wait()
}

func (g *GenericDownloader) Stop() error {
	defer func() {
		g.progress.Status = internal.StatusCompleted
		g.Complete()
	}()
	// yt-dlp uses multiple child process the parent process
	// has been spawned with setPgid = true. To properly kill
	// all subprocesses a SIGTERM need to be sent to the correct
	// process group
	if g.proc == nil {
		return errors.New("*os.Process not set")
	}

	pgid, err := syscall.Getpgid(g.proc.Pid)
	if err != nil {
		return err
	}
	if err := syscall.Kill(-pgid, syscall.SIGTERM); err != nil {
		return err
	}

	return nil
}

func (g *GenericDownloader) Status() *internal.ProcessSnapshot {
	return &internal.ProcessSnapshot{
		Id:             g.Id,
		Info:           g.Metadata,
		Progress:       g.progress,
		Output:         g.output,
		Params:         g.Params,
		DownloaderName: "generic",
	}
}

func (g *GenericDownloader) UpdateSavedFilePath(p string) { g.output.SavedFilePath = p }

func (g *GenericDownloader) SetOutput(o internal.DownloadOutput)     { g.output = o }
func (g *GenericDownloader) SetProgress(p internal.DownloadProgress) { g.progress = p }

func (g *GenericDownloader) SetMetadata(fetcher func(url string) (*common.DownloadMetadata, error)) {
	g.FetchMetadata(fetcher)
}

func (g *GenericDownloader) SetPending(p bool) {
	g.Pending = p
}

func (g *GenericDownloader) GetId() string  { return g.Id }
func (g *GenericDownloader) GetUrl() string { return g.URL }

func (g *GenericDownloader) RestoreFromSnapshot(snap *internal.ProcessSnapshot) error {
	if snap == nil {
		return errors.New("cannot restore nil snapshot")
	}

	s := *snap

	g.Id = s.Id
	g.URL = s.Info.URL
	g.Metadata = s.Info
	g.progress = s.Progress
	g.output = s.Output
	g.Params = s.Params

	return nil
}

func (g *GenericDownloader) IsCompleted() bool { return g.Completed }
