package service

import (
	"context"
	"log"

	"github.com/libp2p/go-libp2p/core/host"
)

// Service defines the interface for a generic service in the system.
type Service interface {
	// Start initiates the service with the given context.
	// It returns an error if the service fails to start.
	Start(context.Context) error

	// Stop terminates the service with the given context.
	// It returns an error if the service fails to stop gracefully.
	Stop(context.Context) error

	// Name returns the identifier of the service.
	Name() string
}

type BaseService struct {
	name string
	host host.Host
	log  *log.Logger
}

func NewBaseService(h host.Host, name string) *BaseService {
	return &BaseService{
		host: h,
		name: name,
	}
}

func (bs *BaseService) Start(ctx context.Context) error {
	return nil
}

func (bs *BaseService) Stop(ctx context.Context) error {
	return nil
}

func (bs *BaseService) Name() string {
	return bs.name
}

func (bs *BaseService) GetHost() host.Host {
	return bs.host
}

func (bs *BaseService) GetLogger() *log.Logger {
	if bs.log == nil {
		bs.log = log.Default()
	}
	return bs.log
}
