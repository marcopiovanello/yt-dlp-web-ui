package service

import (
	"context"
	"sync"

	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/rest"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/status/domain"
)

type Service struct {
	repository     domain.Repository
	utilityService *rest.Service
}

// Version implements domain.Status.
func (s *Service) Status(ctx context.Context) (*domain.Status, error) {
	// rpcVersion, downloaderVersion, err := s.utilityService.GetVersion(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	var (
		wg          sync.WaitGroup
		pending     int
		downloading int
		completed   int
		speed       int64
		// version     = fmt.Sprintf("RPC: %s yt-dlp: %s", rpcVersion, downloaderVersion)
	)

	wg.Add(4)

	go func() {
		pending = s.repository.Pending(ctx)
		wg.Done()
	}()

	go func() {
		downloading = s.repository.Downloading(ctx)
		wg.Done()
	}()

	go func() {
		completed = s.repository.Completed(ctx)
		wg.Done()
	}()

	go func() {
		speed = s.repository.DownloadSpeed(ctx)
		wg.Done()
	}()

	wg.Wait()

	return &domain.Status{
		Downloading:   downloading,
		Pending:       pending,
		Completed:     completed,
		DownloadSpeed: int(speed),
	}, nil
}

func New(repository domain.Repository, utilityService *rest.Service) domain.Service {
	return &Service{
		repository:     repository,
		utilityService: utilityService,
	}
}
