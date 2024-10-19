package rpc

import (
	"log"

	gorpc "github.com/libp2p/go-libp2p-gorpc"
	host "github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/protocol"
)

type RpcApi struct {
	Name          string
	ProtocolId    protocol.ID
	Version       string
	Service       interface{}
	Public        bool
	Authenticated bool
}

type RpcServer struct {
	Endpoint   string
	Api        []RpcApi
	Log        *log.Logger
	AppVersion string
}

func NewRpcServer(endpoint string, log *log.Logger, appVersion string) *RpcServer {
	return &RpcServer{
		Endpoint:   endpoint,
		Api:        []RpcApi{},
		Log:        log,
		AppVersion: appVersion,
	}
}

func (s *RpcServer) AddApi(api RpcApi) {
	s.Api = append(s.Api, api)
}

func StartRpcClient(host *host.Host, protocolId protocol.ID) *gorpc.Client {
	return gorpc.NewClient(*host, protocolId)
}
