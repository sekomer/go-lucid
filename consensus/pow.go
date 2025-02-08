package consensus

import (
	"context"
	"encoding/hex"
	"fmt"
	"go-lucid/core/block"
	"go-lucid/core/transaction"
	"go-lucid/mempool"
	"go-lucid/miner"
	"go-lucid/state"
)

func SelectTransactions(mempool *mempool.Mempool, maxBlockSize int) (selected []*transaction.RawTransaction) {
	currentSize := 0

	for _, tx := range mempool.GetTxs() {
		ser, _ := tx.Serialize() // txs are already verified at this point
		txSize := len(ser)
		if currentSize+txSize > maxBlockSize {
			break
		}
		selected = append(selected, tx)
		currentSize += txSize
	}

	return selected
}

func MineBlock(
	transactions []transaction.RawTransaction,
	previousBlock *block.Block,
	ctx context.Context,
) (*block.Block, error) {
	state := state.GetState()
	miner := miner.GetGlobalMiner()
	newBlock, err := miner.Mine(transactions, previousBlock, state, ctx)
	if err != nil {
		return nil, err
	}

	nonce := newBlock.BlockHeader.Nonce
	blockHash, err := newBlock.GetHash()
	if err != nil {
		return nil, err
	}
	fmt.Println("nonce", nonce)
	hexhash := hex.EncodeToString(blockHash)
	fmt.Println("blockHash", hexhash)

	return newBlock, nil
}
