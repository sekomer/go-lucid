package block

import (
	"context"
	"encoding/json"
	"errors"
	"go-lucid/core/block"
	"go-lucid/database"
	"go-lucid/rpc"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"gorm.io/gorm"
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

// GetBlockHeader returns the block header for the given block height
func (s *BlockService) GetBlockHeader(ctx context.Context, args *BlockRpcArgs, reply *BlockRpcReply) error {
	// todo: implement and get real block header from db

	h := block.BlockHeader{
		Version:    1,
		Height:     123456789,
		PrevBlock:  []byte("prev block"),
		Nonce:      99999999,
		MerkleRoot: []byte("merkle root"),
		Timestamp:  time.Now(),
		Bits:       123456789,
	}
	buf, err := json.Marshal(h)
	if err != nil {
		return err
	}
	reply.Result = buf

	return nil
}

// GetBlock returns the block for the given block height
func (s *BlockService) GetBlock(ctx context.Context, args *BlockRpcArgs, reply *BlockRpcReply) error {
	db := database.GetDB()

	blockNumber := args.Args[0].(int64)

	block := block.BlockModel{}
	err := db.Where("height = ?", blockNumber).First(&block).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			reply.Success = false
			reply.Error = "block not found"
			reply.Result = nil
			return nil
		} else {
			panic("TODO unexpected error: " + err.Error())
		}
	}

	buf, err := json.Marshal(block)
	if err != nil {
		return err
	}

	reply.Success = true
	reply.Error = ""
	reply.Result = buf

	return nil
}
