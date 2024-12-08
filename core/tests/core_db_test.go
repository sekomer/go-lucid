package core_test

import (
	"go-lucid/core"
	"go-lucid/database"
	"go-lucid/util"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDB(t *testing.T) {
	t.Parallel()
	t.Log("db test starting...")

	db := database.GetDB(":memory:")
	if db == nil {
		t.Fatal("db is nil")
	}
	db.AutoMigrate(
		&core.RawTransaction{},
		&core.TxIn{},
		&core.TxOut{},
	)

	rawTx := &core.RawTransaction{}
	rawTx.Hash = []byte(util.GenerateRandomBytes(core.HASH_LEN))
	rawTx.Version = int32(util.RandRange(1, 100))
	rawTx.TxInCount = 1
	rawTx.TxIns = []core.TxIn{
		{
			Coinbase: bool(util.RandRange(0, 1) == 0),
			PreviousOutput: core.OutPoint{
				Index: uint32(util.RandRange(1, 100)),
				Hash:  []byte(util.GenerateRandomBytes(core.HASH_LEN)),
			},
			Sequence:        uint32(util.RandRange(1, 100)),
			TransactionID:   0,
			ScriptLength:    32,
			SignatureScript: util.GenerateRandomBytes(32),
		},
	}
	rawTx.TxOutCount = 1
	rawTx.TxOuts = []core.TxOut{
		{
			Value:         core.Lumen(util.RandRange(1, 1000000)),
			TransactionID: 0,
			ScriptLength:  32,
			PkScript:      util.GenerateRandomBytes(32),
		},
	}

	db.Create(rawTx)

	xx := &core.RawTransaction{}
	db.First(xx)
	t.Log("rawtransactionsfromdatabase:", xx)

	dbTx := &core.RawTransaction{}
	db.Preload("TxIns").Preload("TxOuts").First(dbTx)

	if cmp.Equal(rawTx, dbTx) == false {
		t.Log("diff:", cmp.Diff(rawTx, dbTx))
		t.Fatal("Transactions are not equal")
	}

	t.Log("db test finished")
}
