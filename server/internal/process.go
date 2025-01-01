package internal

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"regexp"
	"slices"
	"syscall"

	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/archiver"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/config"
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

const (
	StatusPending = iota
	StatusDownloading
	StatusCompleted
	StatusErrored
)

// Process descriptor
type Process struct {
	Id         string
	Url        string
	Livestream bool
	Params     []string
	Info       DownloadInfo
	Progress   DownloadProgress
	Output     DownloadOutput
	proc       *os.Process
}

// Starts spawns/forks a new yt-dlp process and parse its stdout.
// The process is spawned to outputting a custom progress text that
// Resembles a JSON Object in order to Unmarshal it later.
// This approach is anyhow not perfect: quotes are not escaped properly.
// Each process is not identified by its PID but by a UUIDv4
func (p *Process) Start() {
	// escape bash variable escaping and command piping, you'll never know
	// what they might come with...
	p.Params = slices.DeleteFunc(p.Params, func(e string) bool {
		match, _ := regexp.MatchString(`(\$\{)|(\&\&)`, e)
		return match
	})

	p.Params = slices.DeleteFunc(p.Params, func(e string) bool {
		return e == ""
	})

	out := DownloadOutput{
		Path:     config.Instance().DownloadPath,
		Filename: "%(title)s.%(ext)s",
	}

	if p.Output.Path != "" {
		out.Path = p.Output.Path
	}

	if p.Output.Filename != "" {
		out.Filename = p.Output.Filename
	}

	buildFilename(&p.Output)

	templateReplacer := strings.NewReplacer("\n", "", "\t", "", " ", "")

	baseParams := []string{
		strings.Split(p.Url, "?list")[0], //no playlist
		"--newline",
		"--no-colors",
		"--no-playlist",
		"--progress-template",
		templateReplacer.Replace(downloadTemplate),
		"--progress-template",
		templateReplacer.Replace(postprocessTemplate),
	}

	// if user asked to manually override the output path...
	if !(slices.Contains(p.Params, "-P") || slices.Contains(p.Params, "--paths")) {
		p.Params = append(p.Params, "-o")
		p.Params = append(p.Params, fmt.Sprintf("%s/%s", out.Path, out.Filename))
	}

	params := append(baseParams, p.Params...)

	slog.Info("requesting download", slog.String("url", p.Url), slog.Any("params", params))

	cmd := exec.Command(config.Instance().DownloaderPath, params...)
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

	p.proc = cmd.Process

	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		stdout.Close()
		p.Complete()
		cancel()
	}()

	logs := make(chan []byte)
	go produceLogs(stdout, logs)
	go p.consumeLogs(ctx, logs)

	go p.detectYtDlpErrors(stderr)

	cmd.Wait()
}

func produceLogs(r io.Reader, logs chan<- []byte) {
	go func() {
		scanner := bufio.NewScanner(r)

		for scanner.Scan() {
			logs <- scanner.Bytes()
		}
	}()
}

func (p *Process) consumeLogs(ctx context.Context, logs <-chan []byte) {
	for {
		select {
		case <-ctx.Done():
			slog.Info("detaching from yt-dlp stdout",
				slog.String("id", p.getShortId()),
				slog.String("url", p.Url),
			)
			return
		case entry := <-logs:
			p.parseLogEntry(entry)
		}
	}
}

func (p *Process) parseLogEntry(entry []byte) {
	var progress ProgressTemplate
	var postprocess PostprocessTemplate

	if err := json.Unmarshal(entry, &progress); err == nil {
		p.Progress = DownloadProgress{
			Status:     StatusDownloading,
			Percentage: progress.Percentage,
			Speed:      progress.Speed,
			ETA:        progress.Eta,
		}

		slog.Info("progress",
			slog.String("id", p.getShortId()),
			slog.String("url", p.Url),
			slog.String("percentage", progress.Percentage),
		)
	}

	if err := json.Unmarshal(entry, &postprocess); err == nil {
		p.Output.SavedFilePath = postprocess.FilePath

		// slog.Info("postprocess",
		// 	slog.String("id", p.getShortId()),
		// 	slog.String("url", p.Url),
		// 	slog.String("filepath", postprocess.FilePath),
		// )
	}
}

func (p *Process) detectYtDlpErrors(r io.Reader) {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		slog.Error("yt-dlp process error",
			slog.String("id", p.getShortId()),
			slog.String("url", p.Url),
			slog.String("err", scanner.Text()),
		)
	}
}

// Keep process in the memoryDB but marks it as complete
// Convention: All completed processes has progress -1
// and speed 0 bps.
func (p *Process) Complete() {
	// auto archive
	// TODO: it's not that deterministic :/
	if p.Progress.Percentage == "" && p.Progress.Speed == 0 {
		var serializedMetadata bytes.Buffer

		json.NewEncoder(&serializedMetadata).Encode(p.Info)

		archiver.Publish(&archiver.Message{
			Id:        p.Id,
			Path:      p.Output.SavedFilePath,
			Title:     p.Info.Title,
			Thumbnail: p.Info.Thumbnail,
			Source:    p.Url,
			Metadata:  serializedMetadata.String(),
			CreatedAt: p.Info.CreatedAt,
		})
	}

	p.Progress = DownloadProgress{
		Status:     StatusCompleted,
		Percentage: "-1",
		Speed:      0,
		ETA:        0,
	}

	// for safety, if the filename is not set, set it with original function
	if p.Output.SavedFilePath == "" {
		p.GetFileName(&p.Output)
	}

	slog.Info("finished",
		slog.String("id", p.getShortId()),
		slog.String("url", p.Url),
	)
}

// Kill a process and remove it from the memory
func (p *Process) Kill() error {
	defer func() {
		p.Progress.Status = StatusCompleted
	}()
	// yt-dlp uses multiple child process the parent process
	// has been spawned with setPgid = true. To properly kill
	// all subprocesses a SIGTERM need to be sent to the correct
	// process group
	if p.proc == nil {
		return errors.New("*os.Process not set")
	}

	pgid, err := syscall.Getpgid(p.proc.Pid)
	if err != nil {
		return err
	}
	if err := syscall.Kill(-pgid, syscall.SIGTERM); err != nil {
		return err
	}

	return nil
}

func (p *Process) GetFileName(o *DownloadOutput) error {
	cmd := exec.Command(
		config.Instance().DownloaderPath,
		"--print", "filename",
		"-o", fmt.Sprintf("%s/%s", o.Path, o.Filename),
		p.Url,
	)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	out, err := cmd.Output()
	if err != nil {
		return err
	}

	p.Output.SavedFilePath = strings.Trim(string(out), "\n")
	return nil
}

func (p *Process) SetPending() {
	// Since video's title isn't available yet, fill in with the URL.
	p.Info = DownloadInfo{
		URL:       p.Url,
		Title:     p.Url,
		CreatedAt: time.Now(),
	}
	p.Progress.Status = StatusPending
}

func (p *Process) SetMetadata() error {
	cmd := exec.Command(config.Instance().DownloaderPath, p.Url, "-J")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		slog.Error("failed to connect to stdout",
			slog.String("id", p.getShortId()),
			slog.String("url", p.Url),
			slog.String("err", err.Error()),
		)
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		slog.Error("failed to connect to stderr",
			slog.String("id", p.getShortId()),
			slog.String("url", p.Url),
			slog.String("err", err.Error()),
		)
		return err
	}

	info := DownloadInfo{
		URL:       p.Url,
		CreatedAt: time.Now(),
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	var bufferedStderr bytes.Buffer

	go func() {
		io.Copy(&bufferedStderr, stderr)
	}()

	slog.Info("retrieving metadata",
		slog.String("id", p.getShortId()),
		slog.String("url", p.Url),
	)

	if err := json.NewDecoder(stdout).Decode(&info); err != nil {
		return err
	}

	p.Info = info
	p.Progress.Status = StatusPending

	if err := cmd.Wait(); err != nil {
		return errors.New(bufferedStderr.String())
	}

	return nil
}

func (p *Process) getShortId() string { return strings.Split(p.Id, "-")[0] }

func buildFilename(o *DownloadOutput) {
	if o.Filename != "" && strings.Contains(o.Filename, ".%(ext)s") {
		o.Filename += ".%(ext)s"
	}

	o.Filename = strings.Replace(
		o.Filename,
		".%(ext)s.%(ext)s",
		".%(ext)s",
		1,
	)
}
