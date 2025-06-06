package block

import (
	"context"
	"go-lucid/rpc"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

type BlockClient struct {
	client *rpc.BaseService
}

func NewBlockClient(host host.Host) *BlockClient {
	return &BlockClient{
		client: rpc.NewBaseService(
			host,
			ServiceName,
			ProtocolID,
			Version,
		),
	}
}

// Block rpc call to peer
func (c *BlockClient) Call(ctx context.Context, peer peer.ID, method string, args *BlockRpcArgs, reply *BlockRpcReply) error {
	return c.client.Client.Call(
		peer,
		c.client.Name,
		method,
		args,
		reply,
	)
}
