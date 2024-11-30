package ping

import (
	"context"
	"go-lucid/rpc"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

type PingClient struct {
	client *rpc.BaseService
}

func NewPingClient(host host.Host) *PingClient {
	return &PingClient{
		client: rpc.NewBaseService(
			host,
			ServiceName,
			ProtocolID,
			Version,
		),
	}
}

// Ping to peer
func (c *PingClient) Call(ctx context.Context, peer peer.ID, method string, args *PingArgs, reply *PingReply) error {
	return c.client.Client.Call(
		peer,
		c.client.Name,
		method,
		args,
		reply,
	)
}
