package simple

import (
	"context"
	"time"

	"go-lucid/core/block"
	"go-lucid/core/transaction"
	"go-lucid/miner/core"
	"go-lucid/state"
	"go-lucid/util/hexutil"

	"github.com/holiman/uint256"
)

func (m *SimpleMiner) Mine(transactions []transaction.RawTransaction, previousBlock *block.Block, state *state.State, ctx context.Context) (*block.Block, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		newBlock := &block.Block{
			BlockHeader: block.BlockHeader{
				PrevBlock:  previousBlock.Hash,
				Timestamp:  time.Unix(state.GetCurrentTimestamp(), 0),
				Version:    state.GetVersion(),
				Height:     state.GetCurrentBlockHeight() + 1,
				MerkleRoot: []byte{},
				Nonce:      0,
				Bits:       0,
			},
			TxCount: uint32(len(transactions)),
			Txs:     transactions,
		}

		target := state.GetTarget()
		for {
			newBlock.Nonce++
			blockHash, err := newBlock.GetHash()
			if err != nil {
				return nil, err
			}
			encoded := hexutil.EncodeToPrefixedHex(blockHash)
			blockHashUint256, err := uint256.FromHex(encoded)
			if err != nil {
				return nil, err
			}

			if target.Gt(blockHashUint256) {
				newBlock.Hash = blockHash
				return newBlock, nil
			}
		}
	}
}

func (m *SimpleMiner) GetResources() *core.Resources {
	return m.Resources
}

func (m *SimpleMiner) SetResources(resources *core.Resources) {
	m.Resources = resources
}

func NewSimpleMiner() core.MinerInterface {
	return &SimpleMiner{
		Resources: &core.Resources{},
	}
}
