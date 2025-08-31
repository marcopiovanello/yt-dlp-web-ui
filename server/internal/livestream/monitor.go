package livestream

import (
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal/kv"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal/queue"
	bolt "go.etcd.io/bbolt"
)

var bucket = []byte("livestreams")

type Monitor struct {
	db      *bolt.DB
	store   *kv.Store              // where the just started livestream will be published
	mq      *queue.MessageQueue    // where the just started livestream will be published
	streams map[string]*LiveStream // keeps track of the livestreams
	done    chan *LiveStream       // to signal individual processes completition
}

func NewMonitor(mq *queue.MessageQueue, store *kv.Store, db *bolt.DB) *Monitor {
	return &Monitor{
		mq:      mq,
		db:      db,
		store:   store,
		streams: make(map[string]*LiveStream),
		done:    make(chan *LiveStream),
	}
}

// Detect each livestream completition, if done detach it from the monitor.
func (m *Monitor) Schedule() {
	for l := range m.done {
		delete(m.streams, l.url)

		m.db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket(bucket)
			return b.Delete([]byte(l.url))
		})
	}
}

func (m *Monitor) Add(url string) {
	ls := New(url, m.done, m.mq, m.store)

	go ls.Start()
	m.streams[url] = ls

	m.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		return b.Put([]byte(url), []byte{})
	})
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
		status[k] = Status{
			Status:   v.status,
			WaitTime: v.waitTime,
			LiveDate: v.liveDate,
		}
	}

	return status
}

// Restore a saved state and resume the monitored livestreams
func (m *Monitor) Restore() error {
	return m.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		return b.ForEach(func(k, v []byte) error {
			m.Add(string(k))
			return nil
		})
	})
}
