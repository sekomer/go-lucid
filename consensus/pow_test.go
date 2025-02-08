package consensus_test

import (
	"bytes"
	"context"
	"encoding/hex"
	"go-lucid/consensus"
	"go-lucid/core/block"
	"go-lucid/core/transaction"
	"go-lucid/util"
	"testing"
	"time"
)

func TestMineBlock(t *testing.T) {
	t.Parallel()

	previousBlock := &block.Block{
		BlockHeader: block.BlockHeader{
			Version:    0,
			Height:     0,
			PrevBlock:  []byte("prev"),
			MerkleRoot: []byte("merkle"),
			Timestamp:  time.Unix(0, 0),
			Bits:       0x1e0ffff0,
			Nonce:      22,
		},
	}
	previousBlock.Hash, _ = previousBlock.GetHash()
	transactions := []transaction.RawTransaction{
		{
			Hash:       util.HashToBytes("a105afd81d4eca54972c6e0db6720c9b6cb1894f8d4fd9494825bb07ab1a4590"),
			BlockID:    42,
			Version:    0,
			TxInCount:  1,
			TxOutCount: 1,
			TxIns: []transaction.TxIn{
				{
					TransactionID: 42,
					Coinbase:      true,
					PreviousOutput: transaction.OutPoint{
						Hash:  []byte("prev"),
						Index: 0,
					},
					SignatureScript: []byte("sig"),
					Sequence:        0,
				},
			},
			TxOuts: []transaction.TxOut{
				{
					TransactionID: 42,
					Value:         100,
					PkScript:      []byte("pk"),
				},
			},
			LockTime: 0,
		},
	}

	transaction := transactions[0]
	txHash, _ := transaction.GetHash()
	if !bytes.Equal(txHash, transaction.Hash) {
		t.Fatal("wrong tx hash!")
	}

	minedBlock, err := consensus.MineBlock(transactions, previousBlock, context.Background())
	if err != nil {
		t.Fatal(err)
	}

	t.Log("minedBlock:", hex.EncodeToString(minedBlock.Hash))
}
