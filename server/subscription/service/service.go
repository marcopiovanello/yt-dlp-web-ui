package service

import (
	"context"
	"errors"
	"math"

	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/subscription/data"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/subscription/domain"
	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/subscription/task"
	"github.com/robfig/cron/v3"
)

type Service struct {
	r      domain.Repository
	runner task.TaskRunner
}

func New(r domain.Repository, runner task.TaskRunner) domain.Service {
	s := &Service{
		r:      r,
		runner: runner,
	}

	// very crude recoverer
	initial, _ := s.List(context.Background(), 0, math.MaxInt)
	if initial != nil {
		for _, v := range initial.Data {
			s.runner.Submit(&v)
		}
	}

	return s
}

func fromDB(model *data.Subscription) domain.Subscription {
	return domain.Subscription{
		Id:       model.Id,
		URL:      model.URL,
		Params:   model.Params,
		CronExpr: model.CronExpr,
	}
}

func toDB(dto *domain.Subscription) data.Subscription {
	return data.Subscription{
		Id:       dto.Id,
		URL:      dto.URL,
		Params:   dto.Params,
		CronExpr: dto.CronExpr,
	}
}

// Delete implements domain.Service.
func (s *Service) Delete(ctx context.Context, id string) error {
	s.runner.StopTask(id)
	return s.r.Delete(ctx, id)
}

// GetCursor implements domain.Service.
func (s *Service) GetCursor(ctx context.Context, id string) (int64, error) {
	return s.r.GetCursor(ctx, id)
}

// List implements domain.Service.
func (s *Service) List(ctx context.Context, start int64, limit int) (
	*domain.PaginatedResponse[[]domain.Subscription],
	error,
) {
	dbSubs, err := s.r.List(ctx, start, limit)
	if err != nil {
		return nil, err
	}

	subs := make([]domain.Subscription, len(*dbSubs))

	for i, v := range *dbSubs {
		subs[i] = fromDB(&v)
	}

	var (
		first int64
		next  int64
	)

	if len(subs) > 0 {
		first, err = s.r.GetCursor(ctx, subs[0].Id)
		if err != nil {
			return nil, err
		}

		next, err = s.r.GetCursor(ctx, subs[len(subs)-1].Id)
		if err != nil {
			return nil, err
		}
	}

	return &domain.PaginatedResponse[[]domain.Subscription]{
		First: first,
		Next:  next,
		Data:  subs,
	}, nil
}

// Submit implements domain.Service.
func (s *Service) Submit(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error) {
	if sub.CronExpr == "" {
		sub.CronExpr = "*/5 * * * *"
	}

	_, err := cron.ParseStandard(sub.CronExpr)
	if err != nil {
		return nil, errors.Join(errors.New("failed parsing cron expression"), err)
	}

	subDB, err := s.r.Submit(ctx, &data.Subscription{
		URL:      sub.URL,
		Params:   sub.Params,
		CronExpr: sub.CronExpr,
	})

	retval := fromDB(subDB)

	if err := s.runner.Submit(sub); err != nil {
		return nil, err
	}

	return &retval, err
}

// UpdateByExample implements domain.Service.
func (s *Service) UpdateByExample(ctx context.Context, example *domain.Subscription) error {
	_, err := cron.ParseStandard(example.CronExpr)
	if err != nil {
		return errors.Join(errors.New("failed parsing cron expression"), err)
	}

	e := toDB(example)

	return s.r.UpdateByExample(ctx, &e)
}
