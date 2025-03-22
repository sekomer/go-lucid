package block

import (
	"context"
	"encoding/json"
	"go-lucid/core/block"
	"go-lucid/rpc"

	"github.com/libp2p/go-libp2p/core/host"
)

type BlockService struct {
	rpc.BaseService
}

func NewBlockService(host host.Host) *BlockService {
	return &BlockService{
		BaseService: *rpc.NewBaseService(
			host,
			ServiceName,
			ProtocolID,
			Version,
		),
	}
}

// Define the RPC method
func (s *BlockService) GetBlock(ctx context.Context, args *GetBlockRpcArgs, reply *GetBlockRpcReply) error {
	// todo: implement and get real block from db

	replyBlock := block.Block{
		BlockHeader: block.BlockHeader{
			Version:   1,
			Height:    123456789,
			PrevBlock: []byte("prev block"),
			Nonce:     99999999,
		},
	}

	buf, err := json.Marshal(replyBlock)
	if err != nil {
		return err
	}
	reply.Result = buf

	return nil
}
