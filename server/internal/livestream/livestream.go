package livestream

import (
	"bufio"
	"errors"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/config"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal/downloaders"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal/kv"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal/pipes"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal/queue"
)

const (
	waiting = iota
	inProgress
	completed
	errored
)

// Defines a generic livestream.
// A livestream is identified by its url.
type LiveStream struct {
	url          string
	proc         *os.Process        // used to manually kill the yt-dlp process
	status       int                // whether is monitoring or completed
	done         chan *LiveStream   // where to signal the completition
	waitTimeChan chan time.Duration // time to livestream start
	waitTime     time.Duration
	liveDate     time.Time

	mq    *queue.MessageQueue
	store *kv.Store
}

func New(url string, done chan *LiveStream, mq *queue.MessageQueue, store *kv.Store) *LiveStream {
	return &LiveStream{
		url:          url,
		done:         done,
		status:       waiting,
		waitTime:     time.Second * 0,
		waitTimeChan: make(chan time.Duration),
		mq:           mq,
		store:        store,
	}
}

// Start the livestream monitoring process, once completion signals on the done channel
func (l *LiveStream) Start() error {
	cmd := exec.Command(
		config.Instance().DownloaderPath,
		l.url,
		"--wait-for-video", "30", // wait for the stream to be live and recheck every 10 secs
		"--no-colors", // no ansi color fuzz
		"--simulate",
		"--newline",
		"--paths", config.Instance().DownloadPath,
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		l.status = errored
		return err
	}
	defer stdout.Close()

	if err := cmd.Start(); err != nil {
		l.status = errored
		return err
	}

	l.proc = cmd.Process
	l.status = waiting

	// Start monitoring when the livestream is goin to be live.
	// If already live do nothing.
	go l.monitorStartTime(stdout)

	// Wait to the simulated download process to finish.
	cmd.Wait()

	// Set the job as completed and notify the parent the completion.
	l.status = completed
	l.done <- l

	// Send the started livestream to the message queue! :D

	//TODO: add pipes
	d := downloaders.NewLiveStreamDownloader(l.url, []pipes.Pipe{})

	l.store.Set(d)
	l.mq.Publish(d)

	return nil
}

func (l *LiveStream) monitorStartTime(r io.Reader) {
	// yt-dlp shows the time in the stdout
	scanner := bufio.NewScanner(r)

	defer func() {
		l.status = inProgress
		close(l.waitTimeChan)
	}()

	// however the time to live is not shown in a new line (and atm there's nothing to do about)
	// use a custom split funciton to set the line separator to \r instead of \r\n or \n
	scanner.Split(stdoutSplitFunc)

	waitTimeScanner := func() {
		for scanner.Scan() {
			// l.log <- scanner.Bytes()

			// if this substring is in the current line the download is starting,
			// no need to monitor the time to live.
			//TODO: silly
			if !strings.Contains(scanner.Text(), "Remaining time until next attempt") {
				return
			}

			parts := strings.Split(scanner.Text(), ": ")
			if len(parts) < 2 {
				continue
			}

			startsIn := parts[1]
			parsed, err := parseTimeSpan(startsIn)
			if err != nil {
				continue
			}

			l.liveDate = parsed

			//TODO: check if using channels is stupid or not
			// l.waitTimeChan <- time.Until(start)
			l.waitTime = time.Until(parsed)
		}
	}

	scanner.Scan()

	for !strings.Contains(scanner.Text(), "Waiting for") {
		scanner.Scan()
	}

	waitTimeScanner()
}

func (l *LiveStream) WaitTime() <-chan time.Duration {
	return l.waitTimeChan
}

// Kills a livestream process and signal its completition
func (l *LiveStream) Kill() error {
	l.done <- l

	if l.proc != nil {
		return l.proc.Kill()
	}

	return errors.New("nil yt-dlp process")
}

// Parse the timespan returned from yt-dlp (time to live)
//
//	parsed := parseTimeSpan("76:12:15")
//	fmt.Println(parsed) // 2024-07-21 13:59:59.634781 +0200 CEST
func parseTimeSpan(timeStr string) (time.Time, error) {
	parts := strings.Split(timeStr, ":")

	hh, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Time{}, err
	}
	mm, err := strconv.Atoi(parts[1])
	if err != nil {
		return time.Time{}, err
	}
	ss, err := strconv.Atoi(parts[2])
	if err != nil {
		return time.Time{}, err
	}

	dd := 0

	if hh > 24 {
		dd = hh / 24
		hh = hh % 24
	}

	start := time.Now()
	start = start.AddDate(0, 0, dd)
	start = start.Add(time.Duration(hh) * time.Hour)
	start = start.Add(time.Duration(mm) * time.Minute)
	start = start.Add(time.Duration(ss) * time.Second)

	return start, nil
}
