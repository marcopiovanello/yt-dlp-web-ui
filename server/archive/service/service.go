package service

import (
	"context"

	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/archive/data"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/archive/domain"
)

type Service struct {
	repository domain.Repository
}

func New(repository domain.Repository) domain.Service {
	return &Service{
		repository: repository,
	}
}

// Archive implements domain.Service.
func (s *Service) Archive(ctx context.Context, entity *domain.ArchiveEntry) error {
	return s.repository.Archive(ctx, &data.ArchiveEntry{
		Id:        entity.Id,
		Title:     entity.Title,
		Path:      entity.Path,
		Thumbnail: entity.Thumbnail,
		Source:    entity.Source,
		Metadata:  entity.Metadata,
		CreatedAt: entity.CreatedAt,
	})
}

// HardDelete implements domain.Service.
func (s *Service) HardDelete(ctx context.Context, id string) (*domain.ArchiveEntry, error) {
	res, err := s.repository.HardDelete(ctx, id)
	if err != nil {
		return nil, err
	}

	return &domain.ArchiveEntry{
		Id:        res.Id,
		Title:     res.Title,
		Path:      res.Path,
		Thumbnail: res.Thumbnail,
		Source:    res.Source,
		Metadata:  res.Metadata,
		CreatedAt: res.CreatedAt,
	}, nil
}

// SoftDelete implements domain.Service.
func (s *Service) SoftDelete(ctx context.Context, id string) (*domain.ArchiveEntry, error) {
	res, err := s.repository.SoftDelete(ctx, id)
	if err != nil {
		return nil, err
	}

	return &domain.ArchiveEntry{
		Id:        res.Id,
		Title:     res.Title,
		Path:      res.Path,
		Thumbnail: res.Thumbnail,
		Source:    res.Source,
		Metadata:  res.Metadata,
		CreatedAt: res.CreatedAt,
	}, nil
}

// List implements domain.Service.
func (s *Service) List(
	ctx context.Context,
	startRowId int,
	limit int,
) (*domain.PaginatedResponse[[]domain.ArchiveEntry], error) {
	res, err := s.repository.List(ctx, startRowId, limit)
	if err != nil {
		return nil, err
	}

	entities := make([]domain.ArchiveEntry, len(*res))

	for i, model := range *res {
		entities[i] = domain.ArchiveEntry{
			Id:        model.Id,
			Title:     model.Title,
			Path:      model.Path,
			Thumbnail: model.Thumbnail,
			Source:    model.Source,
			Metadata:  model.Metadata,
			CreatedAt: model.CreatedAt,
		}
	}

	var (
		first int64
		next  int64
	)

	if len(entities) > 0 {
		first, err = s.repository.GetCursor(ctx, entities[0].Id)
		if err != nil {
			return nil, err
		}

		next, err = s.repository.GetCursor(ctx, entities[len(entities)-1].Id)
		if err != nil {
			return nil, err
		}
	}

	return &domain.PaginatedResponse[[]domain.ArchiveEntry]{
		First: first,
		Next:  next,
		Data:  entities,
	}, nil
}

// GetCursor implements domain.Service.
func (s *Service) GetCursor(ctx context.Context, id string) (int64, error) {
	return s.repository.GetCursor(ctx, id)
}
