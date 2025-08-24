package twitch

import (
	"context"
	"encoding/gob"
	"fmt"
	"iter"
	"log/slog"
	"maps"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/config"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal"
)

type Monitor struct {
	liveChannel           chan *StreamInfo
	monitored             map[string]*Client
	lastState             map[string]bool
	mu                    sync.RWMutex
	authenticationManager *AuthenticationManager
}

func NewMonitor(authenticationManager *AuthenticationManager) *Monitor {
	return &Monitor{
		liveChannel:           make(chan *StreamInfo, 16),
		monitored:             make(map[string]*Client),
		lastState:             make(map[string]bool),
		authenticationManager: authenticationManager,
	}
}

func (m *Monitor) Add(user string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.monitored[user] = NewTwitchClient(m.authenticationManager)
	slog.Info("added user to twitch monitor", slog.String("user", user))
}

func (m *Monitor) Monitor(ctx context.Context, interval time.Duration, handler func(url string) error) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.mu.RLock()
			for user, client := range m.monitored {
				u := user
				c := client

				go func() {
					if err := c.PollStream(u, m.liveChannel); err != nil {
						slog.Error("polling failed", slog.String("user", u), slog.Any("err", err))
					}
				}()
			}
			m.mu.RUnlock()

		case stream := <-m.liveChannel:
			wasLive := m.lastState[stream.UserName]
			if stream.IsLive && !wasLive {
				slog.Info("stream went live", slog.String("user", stream.UserName))
				if err := handler(fmt.Sprintf("https://www.twitch.tv/%s", stream.UserName)); err != nil {
					slog.Error("handler failed", slog.String("user", stream.UserName), slog.Any("err", err))
				}
			}
			m.lastState[stream.UserName] = stream.IsLive

		case <-ctx.Done():
			slog.Info("stopping twitch monitor")
			return
		}
	}
}

func (m *Monitor) GetMonitoredUsers() iter.Seq[string] {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return maps.Keys(m.monitored)
}

func DEFAULT_DOWNLOAD_HANDLER(db *internal.MemoryDB, mq *internal.MessageQueue) func(url string) error {
	return func(url string) error {
		p := &internal.Process{
			Url:        url,
			Livestream: true,
			Params:     []string{"--downloader", "ffmpeg", "--no-part"},
		}
		db.Set(p)
		mq.Publish(p)
		return nil
	}
}

func (m *Monitor) Persist() error {
	filename := filepath.Join(config.Instance().SessionFilePath, "twitch-monitor.dat")

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	users := make([]string, 0, len(m.monitored))

	for user := range m.monitored {
		users = append(users, user)
	}

	return enc.Encode(users)
}

func (m *Monitor) Restore() error {
	filename := filepath.Join(config.Instance().SessionFilePath, "twitch-monitor.dat")

	f, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()

	dec := gob.NewDecoder(f)
	var users []string
	if err := dec.Decode(&users); err != nil {
		return err
	}

	m.monitored = make(map[string]*Client)
	for _, user := range users {
		m.monitored[user] = NewTwitchClient(m.authenticationManager)
	}

	return nil
}
