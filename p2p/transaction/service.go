package p2p

import (
	"context"
	"go-lucid/core"
	"go-lucid/p2p"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
)

type TransactionService struct {
	*p2p.BaseService
}

var _ p2p.P2PService = (*TransactionService)(nil)

func NewTransactionService(h host.Host, ps *pubsub.PubSub) (*TransactionService, error) {
	base, err := p2p.NewBaseService(h, ps, TransactionServiceName)
	if err != nil {
		return nil, err
	}
	return &TransactionService{
		BaseService: base,
	}, nil
}

func (s *TransactionService) Broadcast(ctx context.Context, tx core.RawTransaction) error {
	data, err := tx.Serialize()
	if err != nil {
		return err
	}
	return s.BaseService.Publish(ctx, data)
}
