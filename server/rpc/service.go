package rpc

import (
	"errors"
	"log/slog"

	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/formats"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal/downloaders"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal/kv"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal/livestream"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/internal/queue"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/playlist"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/sys"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/updater"
)

type Service struct {
	db *kv.Store
	mq *queue.MessageQueue
	lm *livestream.Monitor
}

type Running []internal.ProcessSnapshot
type Pending []string

type NoArgs struct{}

// Exec spawns a Process.
// The result of the execution is the newly spawned process Id.
func (s *Service) Exec(args internal.DownloadRequest, result *string) error {
	d := downloaders.NewGenericDownload(args.URL, args.Params)
	d.SetOutput(internal.DownloadOutput{
		Path:     args.Path,
		Filename: args.Rename,
	})

	s.db.Set(d)
	s.mq.Publish(d)

	*result = d.GetId()
	return nil
}

// Exec spawns a Process.
// The result of the execution is the newly spawned process Id.
func (s *Service) ExecPlaylist(args internal.DownloadRequest, result *string) error {
	err := playlist.PlaylistDetect(args, s.mq, s.db)
	if err != nil {
		return err
	}

	*result = ""
	return nil
}

// TODO: docs
func (s *Service) ExecLivestream(args internal.DownloadRequest, result *string) error {
	s.lm.Add(args.URL)

	*result = args.URL
	return nil
}

// TODO: docs
func (s *Service) ProgressLivestream(args NoArgs, result *livestream.LiveStreamStatus) error {
	*result = s.lm.Status()
	return nil
}

// TODO: docs
func (s *Service) KillLivestream(args string, result *struct{}) error {
	slog.Info("killing livestream", slog.String("url", args))

	err := s.lm.Remove(args)
	if err != nil {
		slog.Error("failed killing livestream", slog.String("url", args), slog.Any("err", err))
		return err
	}

	return nil
}

// TODO: docs
func (s *Service) KillAllLivestream(args NoArgs, result *struct{}) error {
	return s.lm.RemoveAll()
}

// Progess retrieves the Progress of a specific Process given its Id
func (s *Service) Progess(args internal.DownloadRequest, progress *internal.DownloadProgress) error {
	dl, err := s.db.Get(args.Id)
	if err != nil {
		return err
	}

	*progress = dl.Status().Progress
	return nil
}

// Progess retrieves available format for a given resource
func (s *Service) Formats(args internal.DownloadRequest, meta *formats.Metadata) error {
	var err error

	metadata, err := formats.ParseURL(args.URL)
	if err != nil && metadata == nil {
		return err
	}

	if metadata.IsPlaylist() {
		go playlist.PlaylistDetect(args, s.mq, s.db)
	}

	*meta = *metadata
	return nil
}

// Pending retrieves a slice of all Pending/Running processes ids
func (s *Service) Pending(args NoArgs, pending *Pending) error {
	*pending = *s.db.Keys()
	return nil
}

// Running retrieves a slice of all Processes progress
func (s *Service) Running(args NoArgs, running *Running) error {
	*running = *s.db.All()
	return nil
}

// Kill kills a process given its id and remove it from the memoryDB
func (s *Service) Kill(args string, killed *string) error {
	slog.Info("Trying killing process with id", slog.String("id", args))

	download, err := s.db.Get(args)
	if err != nil {
		return err
	}

	if download == nil {
		return errors.New("nil process")
	}

	if err := download.Stop(); err != nil {
		slog.Info("failed killing process", slog.String("id", download.GetId()), slog.Any("err", err))
		return err
	}

	s.db.Delete(download.GetId())
	slog.Info("succesfully killed process", slog.String("id", download.GetId()))

	return nil
}

// KillAll kills all process unconditionally and removes them from
// the memory db
func (s *Service) KillAll(args NoArgs, killed *string) error {
	slog.Info("Killing all spawned processes")

	var (
		keys       = s.db.Keys()
		removeFunc = func(d downloaders.Downloader) error {
			defer s.db.Delete(d.GetId())
			return d.Stop()
		}
	)

	for _, key := range *keys {
		dl, err := s.db.Get(key)
		if err != nil {
			return err
		}

		if dl == nil {
			s.db.Delete(key)
			continue
		}

		if err := removeFunc(dl); err != nil {
			slog.Info(
				"failed killing process",
				slog.String("id", dl.GetId()),
				slog.Any("err", err),
			)
			continue
		}

		slog.Info("succesfully killed process", slog.String("id", dl.GetId()))
	}

	return nil
}

// Remove a process from the db rendering it unusable if active
func (s *Service) Clear(args string, killed *string) error {
	slog.Info("Clearing process with id", slog.String("id", args))
	s.db.Delete(args)
	return nil
}

// Removes completed processes
func (s *Service) ClearCompleted(cleared *string) error {
	var (
		keys       = s.db.Keys()
		removeFunc = func(d downloaders.Downloader) error {
			defer s.db.Delete(d.GetId())

			if !d.IsCompleted() {
				return nil
			}

			return d.Stop()
		}
	)

	for _, key := range *keys {
		proc, err := s.db.Get(key)
		if err != nil {
			return err
		}

		if err := removeFunc(proc); err != nil {
			return err
		}
	}

	return nil
}

// FreeSpace gets the available from package sys util
func (s *Service) FreeSpace(args NoArgs, free *uint64) error {
	freeSpace, err := sys.FreeSpace()
	if err != nil {
		return err
	}

	*free = freeSpace
	return err
}

// Return a flattned tree of the download directory
func (s *Service) DirectoryTree(args NoArgs, tree *[]string) error {
	dfsTree, err := sys.DirectoryTree()

	if err != nil {
		*tree = nil
		return err
	}

	if dfsTree != nil {
		*tree = *dfsTree
	}

	return nil
}

// Updates the yt-dlp binary using its builtin function
func (s *Service) UpdateExecutable(args NoArgs, updated *bool) error {
	slog.Info("Updating yt-dlp executable to the latest release")

	if err := updater.UpdateExecutable(); err != nil {
		slog.Error("Failed updating yt-dlp")
		*updated = false
		return err
	}

	*updated = true
	slog.Info("Succesfully updated yt-dlp")

	return nil
}
