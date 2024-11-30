package ping

import (
	"context"
	"go-lucid/rpc"

	"github.com/libp2p/go-libp2p/core/host"
)

type PingService struct {
	rpc.BaseService
}

func NewPingService(host host.Host) *PingService {
	return &PingService{
		BaseService: *rpc.NewBaseService(
			host,
			ServiceName,
			ProtocolID,
			Version,
		),
	}
}

// Define the RPC method
func (s *PingService) Ping(ctx context.Context, args *PingArgs, reply *PingReply) error {
	reply.Data = args.Data
	return nil
}
