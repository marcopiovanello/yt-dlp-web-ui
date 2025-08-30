package livestream

import (
	"encoding/gob"
	"log/slog"
	"maps"
	"os"
	"path/filepath"

	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/config"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal/kv"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal/queue"
)

type Monitor struct {
	db      *kv.Store              // where the just started livestream will be published
	mq      *queue.MessageQueue    // where the just started livestream will be published
	streams map[string]*LiveStream // keeps track of the livestreams
	done    chan *LiveStream       // to signal individual processes completition
}

func NewMonitor(mq *queue.MessageQueue, db *kv.Store) *Monitor {
	return &Monitor{
		mq:      mq,
		db:      db,
		streams: make(map[string]*LiveStream),
		done:    make(chan *LiveStream),
	}
}

// Detect each livestream completition, if done detach it from the monitor.
func (m *Monitor) Schedule() {
	for l := range m.done {
		delete(m.streams, l.url)
	}
}

func (m *Monitor) Add(url string) {
	ls := New(url, m.done, m.mq, m.db)

	go ls.Start()
	m.streams[url] = ls
}

func (m *Monitor) Remove(url string) error {
	return m.streams[url].Kill()
}

func (m *Monitor) RemoveAll() error {
	for _, v := range m.streams {
		if err := v.Kill(); err != nil {
			return err
		}
	}
	return nil
}

func (m *Monitor) Status() LiveStreamStatus {
	status := make(LiveStreamStatus)

	for k, v := range m.streams {
		// wt, ok := <-v.WaitTime()
		// if !ok {
		// 	continue
		// }

		status[k] = Status{
			Status:   v.status,
			WaitTime: v.waitTime,
			LiveDate: v.liveDate,
		}
	}

	return status
}

// Persist the monitor current state to a file.
// The file is located in the configured config directory
func (m *Monitor) Persist() error {
	fd, err := os.Create(filepath.Join(config.Instance().SessionFilePath, "livestreams.dat"))
	if err != nil {
		return err
	}

	defer fd.Close()

	slog.Debug("persisting livestream monitor state")

	var toPersist []string
	for url := range maps.Keys(m.streams) {
		toPersist = append(toPersist, url)
	}

	return gob.NewEncoder(fd).Encode(toPersist)
}

// Restore a saved state and resume the monitored livestreams
func (m *Monitor) Restore() error {
	fd, err := os.Open(filepath.Join(config.Instance().SessionFilePath, "livestreams.dat"))
	if err != nil {
		return err
	}

	defer fd.Close()

	var toRestore []string

	if err := gob.NewDecoder(fd).Decode(&toRestore); err != nil {
		return err
	}

	for _, url := range toRestore {
		m.Add(url)
	}

	slog.Debug("restored livestream monitor state")

	return nil
}
