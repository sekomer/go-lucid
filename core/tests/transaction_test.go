package core_test

import (
	"encoding/json"
	"go-lucid/config"
	"go-lucid/core"
	"go-lucid/core/transaction"
	"go-lucid/database"
	"go-lucid/util"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestSerializeDeserialize(t *testing.T) {
	t.Parallel()

	config.MustReadConfig("../../config/fullnode.yaml")

	tx := transaction.RawTransaction{
		Hash:      []byte(util.GenerateRandomBytes(core.HASH_LEN)),
		Version:   int32(util.RandRange(1, 100)),
		TxInCount: 1,
		TxIns: []transaction.TxIn{
			{
				Coinbase: bool(util.RandRange(0, 1) == 0),
				PreviousOutput: transaction.OutPoint{
					Hash:  []byte(util.GenerateRandomBytes(core.HASH_LEN)),
					Index: uint32(util.RandRange(0, 100)),
				},
				ScriptSig:     util.GenerateRandomBytes(core.HASH_LEN),
				Sequence:      uint32(util.RandRange(0, 100)),
				TransactionID: 0,
			},
		},
		TxOutCount: 1,
		TxOuts: []transaction.TxOut{
			{
				Value:         core.Lumen(util.RandRange(0, 10000000000)),
				PkScript:      util.GenerateRandomBytes(core.HASH_LEN),
				TransactionID: 0,
			},
		},
		LockTime: time.Now().Unix(),
	}

	_, err := json.Marshal(tx)
	if err != nil {
		t.Log("marshal error:", err)
		t.Fatal(err)
	}

	ser, err := tx.Serialize()
	if err != nil {
		t.Log("serialize error:", err)
		t.Fatal(err)
	}

	deserTx := &transaction.RawTransaction{}
	err = deserTx.Deserialize(ser)

	if err != nil {
		t.Log("deserialize error:", err)
		t.Fatal(err)
	}

	if cmp.Equal(tx, *deserTx) == false {
		t.Log(cmp.Diff(tx, *deserTx))
		t.Fatal("Transactions are not equal")
	}

	t.Log("Test passed")
}

func TestTransaction(t *testing.T) {
	t.Parallel()

	config.MustReadConfig("../../config/fullnode.yaml")

	db := database.GetTestDB()
	if db == nil {
		t.Fatal("db is nil")
	}

	rawTx := &transaction.RawTransactionModel{}
	rawTx.Hash = []byte(util.GenerateRandomBytes(core.HASH_LEN))
	rawTx.Version = int32(util.RandRange(1, 100))
	rawTx.TxInCount = 1
	rawTx.TxOutCount = 1
	rawTx.TxIns = []transaction.TxInModel{
		{
			Coinbase: bool(util.RandRange(0, 1) == 0),
			PreviousOutput: transaction.OutPoint{
				Index: uint32(util.RandRange(1, 100)),
				Hash:  []byte(util.GenerateRandomBytes(core.HASH_LEN)),
			},
			Sequence:      uint32(util.RandRange(1, 100)),
			TransactionID: 0,
			ScriptSig:     util.GenerateRandomBytes(core.HASH_LEN),
		},
	}
	rawTx.TxOuts = []transaction.TxOutModel{
		{
			Value:         core.Lumen(util.RandRange(1, 1000000)),
			TransactionID: 0,
			PkScript:      util.GenerateRandomBytes(core.HASH_LEN),
		},
	}

	db.Create(rawTx)

	dbTx := &transaction.RawTransactionModel{}
	db.Preload("TxIns").Preload("TxOuts").First(dbTx)

	if cmp.Equal(rawTx, dbTx) == false {
		t.Log("diff:", cmp.Diff(rawTx, dbTx))
		t.Fatal("Transactions are not equal")
	}

	t.Log("db test finished")
}
