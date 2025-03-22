package core_test

import (
	"go-lucid/core"
	"go-lucid/core/block"
	"go-lucid/core/transaction"
	"go-lucid/database"
	"go-lucid/util"
	"log"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestBlockChain(t *testing.T) {
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

	for i := range 100 {
		rawBlock := &block.BlockModel{}
		rawBlock.Hash = []byte(util.GenerateRandomBytes(core.HASH_LEN))
		rawBlock.Height = uint32(i)
		rawBlock.PrevBlock = []byte(util.GenerateRandomBytes(core.HASH_LEN))
		rawBlock.MerkleRoot = []byte(util.GenerateRandomBytes(core.HASH_LEN))
		rawBlock.Timestamp = time.Now()
		rawBlock.Bits = uint32(i)
		rawBlock.Nonce = uint32(i)
		rawBlock.Version = int32(i)
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
						Coinbase: true,
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

		if i > 0 {
			prevBlock := &block.BlockModel{}
			db.Preload("Txs").Preload("Txs.TxIns").Preload("Txs.TxOuts").Where("height=?", i-1).First(prevBlock)
			rawBlock.PrevBlock = prevBlock.Hash
		}

		db.Create(rawBlock)
	}

	dbBlocks := []*block.BlockModel{}
	db.Preload("Txs").Preload("Txs.TxIns").Preload("Txs.TxOuts").Find(&dbBlocks)

	for i := range len(dbBlocks) {
		if i == 0 {
			continue
		}

		log.Printf("dbBlocks[%d].PrevBlock: %x\n", i, dbBlocks[i].PrevBlock)
		log.Printf("dbBlocks[%d].Hash: %x\n", i, dbBlocks[i].Hash)

		if cmp.Equal(dbBlocks[i].PrevBlock, dbBlocks[i-1].Hash) == false {
			t.Fatal("PrevBlock is not equal")
		}
	}
}
