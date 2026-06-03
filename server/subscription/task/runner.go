package task

import (
	"bytes"
	"context"
	"log/slog"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"

	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/archive"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/config"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/subscription/domain"
	"github.com/robfig/cron/v3"
)

type TaskRunner interface {
	Submit(subcription *domain.Subscription) error
	Spawner(ctx context.Context)
	StopTask(id string) error
	Recoverer()
}

type monitorTask struct {
	Done         chan struct{}
	Schedule     cron.Schedule
	Subscription *domain.Subscription
}

type CronTaskRunner struct {
	mq *internal.MessageQueue
	db *internal.MemoryDB

	tasks  chan monitorTask
	errors chan error

	running map[string]*monitorTask
}

func NewCronTaskRunner(mq *internal.MessageQueue, db *internal.MemoryDB) TaskRunner {
	return &CronTaskRunner{
		mq:      mq,
		db:      db,
		tasks:   make(chan monitorTask),
		errors:  make(chan error),
		running: make(map[string]*monitorTask),
	}
}

var argsSplitterRe = regexp.MustCompile(`(?mi)[^\s"']+|"([^"]*)"|'([^']*)'`)

func (t *CronTaskRunner) Submit(subcription *domain.Subscription) error {
	schedule, err := cron.ParseStandard(subcription.CronExpr)
	if err != nil {
		return err
	}

	job := monitorTask{
		Done:         make(chan struct{}),
		Schedule:     schedule,
		Subscription: subcription,
	}

	t.tasks <- job

	return nil
}

// Handles the entire lifecylce of a monitor job.
func (t *CronTaskRunner) Spawner(ctx context.Context) {
	for req := range t.tasks {
		t.running[req.Subscription.Id] = &req // keep track of the current job

		go func() {
			ctx, cancel := context.WithCancel(ctx) // inject into the job's context a cancellation singal
			fetcherEvents := t.doFetch(ctx, &req)  // retrieve the channel of events of the job

			for {
				select {
				case <-req.Done:
					slog.Info("stopping cron job and removing schedule", slog.String("url", req.Subscription.URL))
					cancel()
					return
				case <-fetcherEvents:
					slog.Info("finished monitoring channel", slog.String("url", req.Subscription.URL))
				}
			}
		}()
	}
}

// Stop a currently scheduled job
func (t *CronTaskRunner) StopTask(id string) error {
	task := t.running[id]
	if task != nil {
		t.running[id].Done <- struct{}{}
		delete(t.running, id)
	}
	return nil
}

// Start a fetcher and notify on a channel when a fetcher has completed
func (t *CronTaskRunner) doFetch(ctx context.Context, req *monitorTask) <-chan struct{} {
	completed := make(chan struct{})

	// generator func
	go func() {
		for {
			sleepFor := t.fetcher(ctx, req)
			completed <- struct{}{}

			time.Sleep(sleepFor)
		}
	}()

	return completed
}

// Perform the retrieval of the latest video of the channel.
// Returns a time.Duration containing the amount of time to the next schedule.
func (t *CronTaskRunner) fetcher(ctx context.Context, req *monitorTask) time.Duration {
	slog.Info("fetching latest video for channel", slog.String("channel", req.Subscription.URL))

	nextSchedule := time.Until(req.Schedule.Next(time.Now()))

	cmd := exec.CommandContext(
		ctx,
		config.Instance().DownloaderPath,
		"-I1",
		"--flat-playlist",
		"--print", "webpage_url",
		req.Subscription.URL,
	)

	stdout, err := cmd.Output()
	if err != nil {
		t.errors <- err
		return time.Duration(0)
	}

	latestVideoURL := string(bytes.Trim(stdout, "\n"))

	// if the download exists there's not point in sending it into the message queue.
	exists, err := archive.DownloadExists(ctx, latestVideoURL)
	if exists && err == nil {
		return nextSchedule
	}

	p := &internal.Process{
		Url: latestVideoURL,
		Params: append(
			argsSplitterRe.FindAllString(req.Subscription.Params, 1),
			[]string{
				"--break-on-existing",
				"--download-archive",
				filepath.Join(config.Instance().Dir(), "archive.txt"),
			}...),
		AutoRemove: true,
	}

	t.db.Set(p)     // give it an id
	t.mq.Publish(p) // send it to the message queue waiting to be processed

	slog.Info(
		"cron task runner next schedule",
		slog.String("url", req.Subscription.URL),
		slog.Any("duration", nextSchedule),
	)

	return nextSchedule
}

func (t *CronTaskRunner) Recoverer() {
	panic("unimplemented")
}
