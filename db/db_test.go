package db

import (
	"go-lucid/core"
	"go-lucid/util"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDB(t *testing.T) {
	t.Parallel()

	t.Log("db test starting...")

	db := GetDB(":memory:")
	if db == nil {
		t.Fatal("db is nil")
	}
	db.AutoMigrate(&core.RawTransaction{}, &core.TxIn{}, &core.TxOut{})

	rawTx := &core.RawTransaction{}
	rawTx.Hash = util.GenerateRandomBytes(util.RandRange(32, 64))
	rawTx.Version = 1
	rawTx.TxInCount = 1
	rawTx.TxOutCount = 1
	rawTx.TxIns = []core.TxIn{
		{
			Coinbase: true,
			PreviousOutput: core.OutPoint{
				Index: uint32(util.RandRange(1, 100)),
				Hash:  util.GenerateRandomBytes(util.RandRange(1, 32)),
			},
			SignatureScript: util.GenerateRandomBytes(util.RandRange(32, 64)),
		},
	}
	rawTx.TxOuts = []core.TxOut{
		{
			Value:    core.Lumen(util.RandRange(1, 1000000)),
			PkScript: util.GenerateRandomBytes(util.RandRange(32, 64)),
		},
	}

	db.Create(rawTx)

	dbTx := &core.RawTransaction{}
	db.Preload("TxIns").Preload("TxOuts").First(dbTx)

	if cmp.Equal(rawTx, dbTx) == false {
		t.Fatal("Transactions are not equal")
	}

	t.Log("db test finished")
}
