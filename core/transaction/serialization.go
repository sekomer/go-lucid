package transaction

import (
	"bytes"
	"encoding/gob"
)

// * RawTransaction

func (tx *RawTransaction) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(tx)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (tx *RawTransaction) Deserialize(data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(tx)
}

func (tx *RawTransaction) SerializeWithoutHash() ([]byte, error) {
	tempTx := *tx
	tempTx.Hash = nil
	return tempTx.Serialize()
}

// * RawTransactionModel
