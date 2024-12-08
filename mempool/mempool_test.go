package mempool_test

import (
	"go-lucid/core"
	"go-lucid/mempool"
	"go-lucid/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMempoolSize(t *testing.T) {
	mempool := mempool.GetMempool()

	for i := 0; i < 100; i++ {
		tx := &core.RawTransaction{}
		tx.Hash = util.GenerateRandomBytes(32)
		mempool.AddTx(tx)
	}

	assert.Equal(t, mempool.Size(), int32(100))

	mempool.Clear()
}

func TestMempoolAddRemove(t *testing.T) {
	mempool := mempool.GetMempool()

	txHash := util.GenerateRandomBytes(32)
	tx := &core.RawTransaction{Hash: txHash}

	mempool.AddTx(tx)
	assert.Equal(t, mempool.Size(), int32(1))

	mempool.RemoveTx(string(txHash))
	assert.Equal(t, mempool.Size(), int32(0))

	mempool.Clear()
}

func TestMempoolGetTx(t *testing.T) {
	mempool := mempool.GetMempool()

	txHash := util.GenerateRandomBytes(32)
	tx := &core.RawTransaction{Hash: txHash}

	mempool.AddTx(tx)

	retrievedTx := mempool.GetTx(string(txHash))
	assert.Equal(t, *tx, *retrievedTx)

	mempool.Clear()
}

func TestMempoolClear(t *testing.T) {
	mempool := mempool.GetMempool()

	for i := 0; i < 100; i++ {
		tx := &core.RawTransaction{}
		tx.Hash = util.GenerateRandomBytes(32)
		mempool.AddTx(tx)
	}

	mempool.Clear()
	assert.Equal(t, mempool.Size(), int32(0))

	mempool.Clear()
}

func TestMempoolGetAllTxs(t *testing.T) {
	mempool := mempool.GetMempool()
	txHashes := make([]string, 0, 100)

	for i := 0; i < 100; i++ {
		txHash := util.GenerateRandomBytes(32)
		txHashes = append(txHashes, string(txHash))
		tx := &core.RawTransaction{Hash: txHash}
		mempool.AddTx(tx)
	}

	txs := mempool.GetTxs()
	assert.Equal(t, len(txs), 100)

	hashMap := make(map[string]bool, 100)
	for _, tx := range txs {
		hashMap[string(tx.Hash)] = true
	}

	for _, txHash := range txHashes {
		found := hashMap[txHash]
		assert.True(t, found)
	}

	mempool.Clear()
}
