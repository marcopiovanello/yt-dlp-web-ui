package kv

import (
	"encoding/gob"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/config"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal/downloaders"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal/queue"
)

var memDbEvents = make(chan downloaders.Downloader, runtime.NumCPU())

// In-Memory Thread-Safe Key-Value Storage with optional persistence
type Store struct {
	table map[string]downloaders.Downloader
	mu    sync.RWMutex
}

func NewStore() *Store {
	return &Store{
		table: make(map[string]downloaders.Downloader),
	}
}

// Get a process pointer given its id
func (m *Store) Get(id string) (downloaders.Downloader, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	entry, ok := m.table[id]
	if !ok {
		return nil, errors.New("no process found for the given key")
	}

	return entry, nil
}

// Store a pointer of a process and return its id
func (m *Store) Set(d downloaders.Downloader) string {
	m.mu.Lock()
	m.table[d.GetId()] = d
	m.mu.Unlock()

	return d.GetId()
}

// Removes a process progress, given the process id
func (m *Store) Delete(id string) {
	m.mu.Lock()
	delete(m.table, id)
	m.mu.Unlock()
}

func (m *Store) Keys() *[]string {
	var running []string

	m.mu.RLock()
	defer m.mu.RUnlock()

	for id := range m.table {
		running = append(running, id)
	}

	return &running
}

// Returns a slice of all currently stored processes progess
func (m *Store) All() *[]internal.ProcessSnapshot {
	running := []internal.ProcessSnapshot{}

	m.mu.RLock()
	for _, v := range m.table {
		running = append(running, *(v.Status()))
	}
	m.mu.RUnlock()

	return &running
}

// Persist the database in a single file named "session.dat"
func (m *Store) Persist() error {
	running := m.All()

	sf := filepath.Join(config.Instance().SessionFilePath, "session.dat")

	fd, err := os.Create(sf)
	if err != nil {
		return errors.Join(errors.New("failed to persist session"), err)
	}

	m.mu.RLock()
	defer m.mu.RUnlock()
	session := Session{Processes: *running}

	if err := gob.NewEncoder(fd).Encode(session); err != nil {
		return errors.Join(errors.New("failed to persist session"), err)
	}

	return nil
}

// Restore a persisted state
func (m *Store) Restore(mq *queue.MessageQueue) {
	sf := filepath.Join(config.Instance().SessionFilePath, "session.dat")

	fd, err := os.Open(sf)
	if err != nil {
		return
	}

	var session Session

	if err := gob.NewDecoder(fd).Decode(&session); err != nil {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	for _, snap := range session.Processes {
		var restored downloaders.Downloader

		if snap.DownloaderName == "generic" {
			d := downloaders.NewGenericDownload("", []string{})
			err := d.RestoreFromSnapshot(&snap)
			if err != nil {
				continue
			}
			restored = d

			m.table[snap.Id] = restored

			if !restored.(*downloaders.GenericDownloader).DownloaderBase.Completed {
				mq.Publish(restored)
			}
		}
	}
}

func (m *Store) EventListener() {
	for p := range memDbEvents {
		if p.Status().DownloaderName == "livestream" {
			slog.Info("compacting Store", slog.String("id", p.GetId()))
			m.Delete(p.GetId())
		}
	}
}
