package rpc

import (
	"log"

	gorpc "github.com/libp2p/go-libp2p-gorpc"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/protocol"
)

type RpcServer struct {
	host   host.Host
	logger *log.Logger
}

// todo: maybe move this under host
func NewRpcServer(host host.Host, logger *log.Logger) *RpcServer {
	return &RpcServer{
		host:   host,
		logger: logger,
	}
}

func (s *RpcServer) RegisterService(service interface{}, protocolID protocol.ID) error {
	rpcHost := gorpc.NewServer(s.host, protocolID)
	err := rpcHost.Register(service)
	if err != nil {
		return err
	}
	s.logger.Printf("Registered RPC service: %s\n", protocolID)
	return nil
}
