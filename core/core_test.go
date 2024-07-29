package core

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestHelloName(t *testing.T) {
	t.Parallel()

	genScript := []byte("test genesis script")

	tx := RawTransaction{
		Version:   1,
		TxInCount: 1,
		TxIn: []TxIn{
			{
				Coinbase: true,
				PreviousOutput: OutPoint{
					Hash:  Hash{0x01, 0x02, 0x03, 0x04},
					Index: 0,
				},
				Sequence:        0,
				ScriptLength:    uint32(len(genScript)),
				SignatureScript: genScript,
			},
		},
		TxOutCount: 1,
		TxOut: []TxOut{
			{
				Value:        100000000,
				ScriptLength: uint32(len(genScript)),
				PkScript:     genScript,
			},
		},
		LockTime: time.Now().Unix(),
	}

	j, err := json.Marshal(tx)
	if err != nil {
		t.Fatal(err)
	}

	ser, err := tx.Serialize()
	if err != nil {
		t.Fatal(err)
	}

	deserTx, err := Deserialize(ser)
	if err != nil {
		t.Fatal(err)
	}

	if cmp.Equal(tx, *deserTx) == false {
		t.Fatal("Transactions are not equal")
	}

	t.Log("json serialized tx size:", len(j))
	t.Log("custom serialized tx size:", len(ser))
}
