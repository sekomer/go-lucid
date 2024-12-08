package core_test

import (
	"encoding/json"
	"go-lucid/core"
	"go-lucid/util"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestSerializeDeserialize(t *testing.T) {
	t.Parallel()

	tx := core.RawTransaction{
		Hash:      []byte(util.GenerateRandomBytes(core.HASH_LEN)),
		Version:   int32(util.RandRange(1, 100)),
		TxInCount: 1,
		TxIns: []core.TxIn{
			{
				Coinbase: bool(util.RandRange(0, 1) == 0),
				PreviousOutput: core.OutPoint{
					Hash:  []byte(util.GenerateRandomBytes(core.HASH_LEN)),
					Index: uint32(util.RandRange(0, 100)),
				},
				ScriptLength:    64,
				SignatureScript: util.GenerateRandomBytes(64),
				Sequence:        uint32(util.RandRange(0, 100)),
				TransactionID:   0,
			},
		},
		TxOutCount: 1,
		TxOuts: []core.TxOut{
			{
				Value:         core.Lumen(util.RandRange(0, 10000000000)),
				ScriptLength:  32,
				PkScript:      util.GenerateRandomBytes(32),
				TransactionID: 0,
			},
		},
		LockTime: time.Now().Unix(),
	}

	j, err := json.Marshal(tx)
	if err != nil {
		t.Log("marshal error:", err)
		t.Fatal(err)
	}
	t.Log("json serialized tx struct:", string(j))

	ser, err := tx.Serialize()
	if err != nil {
		t.Log("serialize error:", err)
		t.Fatal(err)
	}
	t.Log("custom serialized tx struct:", ser)

	deserTx := &core.RawTransaction{}
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
