package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/google/uuid"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/config"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal/downloaders"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal/kv"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal/livestream"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal/queue"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/playlist"

	bolt "go.etcd.io/bbolt"
)

type Service struct {
	mdb *kv.Store
	db  *bolt.DB
	mq  *queue.MessageQueue
	lm  *livestream.Monitor
}

func NewService(
	mdb *kv.Store,
	db *bolt.DB,
	mq *queue.MessageQueue,
	lm *livestream.Monitor,
) *Service {
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("templates"))
		return err
	})
	return &Service{
		mdb: mdb,
		db:  db,
		mq:  mq,
		lm:  lm,
	}
}

func (s *Service) Exec(req internal.DownloadRequest) (string, error) {
	d := downloaders.NewGenericDownload(req.URL, req.Params)
	d.SetOutput(internal.DownloadOutput{
		Path:     req.Path,
		Filename: req.Rename,
	})

	id := s.mdb.Set(d)
	s.mq.Publish(d)

	return id, nil
}

func (s *Service) ExecPlaylist(req internal.DownloadRequest) error {
	return playlist.PlaylistDetect(req, s.mq, s.mdb)
}

func (s *Service) ExecLivestream(req internal.DownloadRequest) {
	s.lm.Add(req.URL)
}

func (s *Service) Running(ctx context.Context) (*[]internal.ProcessSnapshot, error) {
	select {
	case <-ctx.Done():
		return nil, context.Canceled
	default:
		return s.mdb.All(), nil
	}
}

func (s *Service) GetCookies(ctx context.Context) ([]byte, error) {
	fd, err := os.Open("cookies.txt")
	if err != nil {
		return nil, err
	}

	defer fd.Close()

	cookies, err := io.ReadAll(fd)
	if err != nil {
		return nil, err
	}

	return cookies, nil
}

func (s *Service) SetCookies(ctx context.Context, cookies string) error {
	fd, err := os.Create("cookies.txt")
	if err != nil {
		return err
	}

	defer fd.Close()
	fd.WriteString(cookies)

	return nil
}

func (s *Service) SaveTemplate(ctx context.Context, template *internal.CustomTemplate) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("templates"))
		v, err := json.Marshal(template)
		if err != nil {
			return err
		}
		return b.Put([]byte(uuid.NewString()), v)
	})
}

func (s *Service) GetTemplates(ctx context.Context) (*[]internal.CustomTemplate, error) {
	templates := make([]internal.CustomTemplate, 0)

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("templates"))
		if b == nil {
			return nil // bucket vuoto, restituisco lista vuota
		}

		return b.ForEach(func(k, v []byte) error {
			var t internal.CustomTemplate
			if err := json.Unmarshal(v, &t); err != nil {
				return err
			}
			templates = append(templates, t)
			return nil
		})
	})

	if err != nil {
		return nil, err
	}

	return &templates, nil
}

func (s *Service) UpdateTemplate(ctx context.Context, t *internal.CustomTemplate) (*internal.CustomTemplate, error) {
	data, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	err = s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("templates"))
		if b == nil {
			return fmt.Errorf("bucket templates not found")
		}
		return b.Put([]byte(t.Id), data)
	})

	if err != nil {
		return nil, err
	}

	return t, nil
}

func (s *Service) DeleteTemplate(ctx context.Context, id string) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("templates"))
		return b.Delete([]byte(id))
	})
}

func (s *Service) GetVersion(ctx context.Context) (string, string, error) {
	//TODO: load from realease properties file, or anything else outside code
	const CURRENT_RPC_VERSION = "3.2.6"

	result := make(chan string, 1)

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	cmd := exec.CommandContext(ctx, config.Instance().Paths.DownloaderPath, "--version")
	go func() {
		stdout, _ := cmd.Output()
		result <- string(stdout)
	}()

	select {
	case <-ctx.Done():
		return CURRENT_RPC_VERSION, "", errors.New("requesting yt-dlp version took too long")
	case res := <-result:
		return CURRENT_RPC_VERSION, res, nil
	}
}
