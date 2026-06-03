package rest

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/google/uuid"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/config"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal/livestream"
)

type Service struct {
	mdb *internal.MemoryDB
	db  *sql.DB
	mq  *internal.MessageQueue
	lm  *livestream.Monitor
}

func (s *Service) Exec(req internal.DownloadRequest) (string, error) {
	p := &internal.Process{
		Url:    req.URL,
		Params: req.Params,
		Output: internal.DownloadOutput{
			Path:     req.Path,
			Filename: req.Rename,
		},
	}

	id := s.mdb.Set(p)
	s.mq.Publish(p)

	return id, nil
}

func (s *Service) ExecPlaylist(req internal.DownloadRequest) error {
	return internal.PlaylistDetect(req, s.mq, s.mdb)
}

func (s *Service) ExecLivestream(req internal.DownloadRequest) {
	s.lm.Add(req.URL)
}

func (s *Service) Running(ctx context.Context) (*[]internal.ProcessResponse, error) {
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
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return err
	}

	defer conn.Close()

	_, err = conn.ExecContext(
		ctx,
		"INSERT INTO templates (id, name, content) VALUES (?, ?, ?)",
		uuid.NewString(),
		template.Name,
		template.Content,
	)

	return err
}

func (s *Service) GetTemplates(ctx context.Context) (*[]internal.CustomTemplate, error) {
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	rows, err := conn.QueryContext(ctx, "SELECT * FROM templates")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	templates := make([]internal.CustomTemplate, 0)

	for rows.Next() {
		t := internal.CustomTemplate{}

		err := rows.Scan(&t.Id, &t.Name, &t.Content)
		if err != nil {
			return nil, err
		}

		templates = append(templates, t)
	}

	return &templates, nil
}

func (s *Service) UpdateTemplate(ctx context.Context, t *internal.CustomTemplate) (*internal.CustomTemplate, error) {
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	_, err = conn.ExecContext(ctx, "UPDATE templates SET name = ?, content = ? WHERE id = ?", t.Name, t.Content, t.Id)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (s *Service) DeleteTemplate(ctx context.Context, id string) error {
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return err
	}

	defer conn.Close()

	_, err = conn.ExecContext(ctx, "DELETE FROM templates WHERE id = ?", id)

	return err
}

func (s *Service) GetVersion(ctx context.Context) (string, string, error) {
	//TODO: load from realease properties file, or anything else outside code
	const CURRENT_RPC_VERSION = "3.2.5"

	result := make(chan string, 1)

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	cmd := exec.CommandContext(ctx, config.Instance().DownloaderPath, "--version")
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
