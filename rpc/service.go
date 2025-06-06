package rpc

import (
	gorpc "github.com/libp2p/go-libp2p-gorpc"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/protocol"
)

type RpcService any

type BaseService struct {
	Name       string
	ProtocolID protocol.ID
	Version    string
	Client     *gorpc.Client
	Service    RpcService
}

func NewBaseClient(host host.Host, protocolID protocol.ID, opts ...gorpc.ClientOption) *gorpc.Client {
	return gorpc.NewClient(host, protocolID, opts...)
}

func NewBaseService(host host.Host, name string, protocolID protocol.ID, version string) *BaseService {
	return &BaseService{
		Name:       name,
		ProtocolID: protocolID,
		Version:    version,
		Client:     NewBaseClient(host, protocolID),
	}
}
