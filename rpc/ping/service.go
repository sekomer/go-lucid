package ping

import (
	"context"
	"go-lucid/rpc"

	"github.com/libp2p/go-libp2p/core/protocol"
)

type RpcArgs struct {
	Data []byte
}

type RpcReply struct {
	Data []byte
}

type PingService struct{}

func (t *PingService) Ping(ctx context.Context, argType RpcArgs, replyType *RpcReply) error {
	replyType.Data = argType.Data
	return nil
}

func CreatePingService() rpc.RpcApi {
	svc := PingService{}

	pingApi := rpc.RpcApi{
		ProtocolId:    protocol.ID("/p2p/rpc/ping"),
		Version:       "1.0",
		Service:       &svc,
		Public:        true,
		Authenticated: false,
	}

	return pingApi
}
