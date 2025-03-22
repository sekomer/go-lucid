package core_test

import (
	"go-lucid/core"
	"go-lucid/core/block"
	"go-lucid/core/transaction"
	"go-lucid/database"
	"go-lucid/util"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestBlock(t *testing.T) {
	t.Parallel()

	t.Log("db test starting...")

	db := database.GetTestDB()
	if db == nil {
		t.Fatal("db is nil")
	}
	db.AutoMigrate(
		&transaction.RawTransactionModel{},
		&transaction.TxInModel{},
		&transaction.TxOutModel{},
		&block.BlockModel{},
	)

	rawBlock := &block.BlockModel{}
	rawBlock.ID = uint(util.RandRange(1, 100))
	rawBlock.Hash = []byte(util.GenerateRandomBytes(core.HASH_LEN))
	rawBlock.Height = uint32(util.RandRange(1, 100))
	rawBlock.PrevBlock = []byte(util.GenerateRandomBytes(core.HASH_LEN))
	rawBlock.MerkleRoot = []byte(util.GenerateRandomBytes(core.HASH_LEN))
	rawBlock.Timestamp = time.Now()
	rawBlock.Bits = uint32(util.RandRange(1, 100))
	rawBlock.Nonce = uint32(util.RandRange(1, 100))
	rawBlock.Version = int32(util.RandRange(1, 100))
	rawBlock.TxCount = uint32(1)
	rawBlock.Txs = []transaction.RawTransactionModel{
		{
			Hash:       []byte(util.GenerateRandomBytes(core.HASH_LEN)),
			Version:    int32(util.RandRange(1, 100)),
			TxInCount:  uint32(1),
			TxOutCount: uint32(1),
			BlockID:    uint32(rawBlock.ID),
			LockTime:   int64(util.RandRange(1, 100)),
			TxIns: []transaction.TxInModel{
				{
					Coinbase: bool(util.RandRange(0, 1) == 0),
					PreviousOutput: transaction.OutPoint{
						Index: uint32(util.RandRange(1, 100)),
						Hash:  []byte(util.GenerateRandomBytes(core.HASH_LEN)),
					},
					Sequence:      uint32(util.RandRange(1, 100)),
					TransactionID: 0,
					ScriptSig:     []byte(util.GenerateRandomBytes(core.PKS_LEN)),
				},
			},
			TxOuts: []transaction.TxOutModel{
				{
					Value:         core.Lumen(util.RandRange(1, 1000000)),
					TransactionID: 0,
					PkScript:      []byte(util.GenerateRandomBytes(core.PKS_LEN)),
				},
			},
		},
	}
	rawBlock.CreatedAt = time.Now()
	rawBlock.UpdatedAt = time.Now()

	db.Create(rawBlock)

	dbBlock := &block.BlockModel{}
	db.Preload("Txs").Preload("Txs.TxIns").Preload("Txs.TxOuts").First(dbBlock)

	if cmp.Equal(rawBlock, dbBlock) == false {
		t.Log("diff:", cmp.Diff(rawBlock, dbBlock))
		t.Fatal("Blocks are not equal")
	}
}
