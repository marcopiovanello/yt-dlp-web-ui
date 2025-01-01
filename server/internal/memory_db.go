package internal

import (
	"encoding/gob"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/config"
)

// In-Memory Thread-Safe Key-Value Storage with optional persistence
type MemoryDB struct {
	table map[string]*Process
	mu    sync.RWMutex
}

func NewMemoryDB() *MemoryDB {
	return &MemoryDB{
		table: make(map[string]*Process),
	}
}

// Get a process pointer given its id
func (m *MemoryDB) Get(id string) (*Process, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	entry, ok := m.table[id]
	if !ok {
		return nil, errors.New("no process found for the given key")
	}

	return entry, nil
}

// Store a pointer of a process and return its id
func (m *MemoryDB) Set(process *Process) string {
	id := uuid.NewString()

	m.mu.Lock()
	process.Id = id
	m.table[id] = process
	m.mu.Unlock()

	return id
}

// Removes a process progress, given the process id
func (m *MemoryDB) Delete(id string) {
	m.mu.Lock()
	delete(m.table, id)
	m.mu.Unlock()
}

func (m *MemoryDB) Keys() *[]string {
	var running []string

	m.mu.RLock()
	defer m.mu.RUnlock()

	for id := range m.table {
		running = append(running, id)
	}

	return &running
}

// Returns a slice of all currently stored processes progess
func (m *MemoryDB) All() *[]ProcessResponse {
	running := []ProcessResponse{}

	m.mu.RLock()
	for k, v := range m.table {
		running = append(running, ProcessResponse{
			Id:       k,
			Info:     v.Info,
			Progress: v.Progress,
			Output:   v.Output,
			Params:   v.Params,
		})
	}
	m.mu.RUnlock()

	return &running
}

// Persist the database in a single file named "session.dat"
func (m *MemoryDB) Persist() error {
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
func (m *MemoryDB) Restore(mq *MessageQueue) {
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

	for _, proc := range session.Processes {
		restored := &Process{
			Id:       proc.Id,
			Url:      proc.Info.URL,
			Info:     proc.Info,
			Progress: proc.Progress,
			Output:   proc.Output,
			Params:   proc.Params,
		}

		m.table[proc.Id] = restored

		if restored.Progress.Status != StatusCompleted {
			mq.Publish(restored)
		}
	}
}
