package rest

import (
	"sync"
)

var (
	service *Service
	handler *Handler

	serviceOnce sync.Once
	handlerOnce sync.Once
)

func ProvideService(args *ContainerArgs) *Service {
	serviceOnce.Do(func() {
		service = NewService(args.MDB, args.DB, args.MQ, args.LM)
	})
	return service
}

func ProvideHandler(svc *Service) *Handler {
	handlerOnce.Do(func() {
		handler = &Handler{
			service: svc,
		}
	})
	return handler
}
