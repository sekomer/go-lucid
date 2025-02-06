package p2p

import (
	"context"
	"go-lucid/core/block"
	"go-lucid/p2p"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
)

type BlockService struct {
	*p2p.BaseService
}

var _ p2p.P2PService = (*BlockService)(nil)

func NewBlockService(h host.Host, ps *pubsub.PubSub) (*BlockService, error) {
	base, err := p2p.NewBaseService(h, ps, BlockServiceName)
	if err != nil {
		return nil, err
	}
	return &BlockService{
		BaseService: base,
	}, nil
}

func (s *BlockService) Broadcast(ctx context.Context, block block.Block) error {
	data, err := block.Serialize()
	if err != nil {
		return err
	}
	return s.BaseService.Publish(ctx, data)
}
