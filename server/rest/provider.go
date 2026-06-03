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
		service = &Service{
			mdb: args.MDB,
			db:  args.DB,
			mq:  args.MQ,
		}
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
