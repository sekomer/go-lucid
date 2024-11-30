package core

import (
	"bytes"
	"encoding/binary"
)

func serializeTxIn(buf *bytes.Buffer, tx *RawTransaction) error {
	err := binary.Write(buf, binary.LittleEndian, tx.TxInCount)
	if err != nil {
		return err
	}
	for _, in := range tx.TxIns {
		err = binary.Write(buf, binary.LittleEndian, in.Coinbase)
		if err != nil {
			return err
		}
		err = binary.Write(buf, binary.LittleEndian, in.PreviousOutput.Hash)
		if err != nil {
			return err
		}
		err = binary.Write(buf, binary.LittleEndian, in.PreviousOutput.Index)
		if err != nil {
			return err
		}
		err = binary.Write(buf, binary.LittleEndian, in.ScriptLength)
		if err != nil {
			return err
		}
		err = binary.Write(buf, binary.LittleEndian, in.SignatureScript)
		if err != nil {
			return err
		}
		err = binary.Write(buf, binary.LittleEndian, in.Sequence)
		if err != nil {
			return err
		}
	}

	return nil
}

func serializeTxOut(buf *bytes.Buffer, tx *RawTransaction) error {
	err := binary.Write(buf, binary.LittleEndian, tx.TxOutCount)
	if err != nil {
		return err
	}
	for _, out := range tx.TxOuts {
		err = binary.Write(buf, binary.LittleEndian, out.Value)
		if err != nil {
			return err
		}
		err = binary.Write(buf, binary.LittleEndian, out.ScriptLength)
		if err != nil {
			return err
		}
		err = binary.Write(buf, binary.LittleEndian, out.PkScript)
		if err != nil {
			return err
		}
	}

	return nil
}

func (tx *RawTransaction) Serialize() ([]byte, error) {
	var buf bytes.Buffer

	err := binary.Write(&buf, binary.LittleEndian, tx.Version)
	if err != nil {
		return nil, err
	}
	err = serializeTxIn(&buf, tx)
	if err != nil {
		return nil, err
	}
	err = serializeTxOut(&buf, tx)
	if err != nil {
		return nil, err
	}
	err = binary.Write(&buf, binary.LittleEndian, tx.LockTime)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func deserializeTxIn(buf *bytes.Buffer, tx *RawTransaction) error {
	err := binary.Read(buf, binary.LittleEndian, &tx.TxInCount)
	if err != nil {
		return err
	}
	tx.TxIns = make([]TxIn, tx.TxInCount)
	for i := 0; i < int(tx.TxInCount); i++ {
		err = binary.Read(buf, binary.LittleEndian, &tx.TxIns[i].Coinbase)
		if err != nil {
			return err
		}
		err = binary.Read(buf, binary.LittleEndian, &tx.TxIns[i].PreviousOutput.Hash)
		if err != nil {
			return err
		}
		err = binary.Read(buf, binary.LittleEndian, &tx.TxIns[i].PreviousOutput.Index)
		if err != nil {
			return err
		}
		err = binary.Read(buf, binary.LittleEndian, &tx.TxIns[i].ScriptLength)
		if err != nil {
			return err
		}
		tx.TxIns[i].SignatureScript = make([]byte, tx.TxIns[i].ScriptLength)
		err = binary.Read(buf, binary.LittleEndian, &tx.TxIns[i].SignatureScript)
		if err != nil {
			return err
		}
		err = binary.Read(buf, binary.LittleEndian, &tx.TxIns[i].Sequence)
		if err != nil {
			return err
		}
	}

	return nil
}

func deserializeTxOut(buf *bytes.Buffer, tx *RawTransaction) error {
	err := binary.Read(buf, binary.LittleEndian, &tx.TxOutCount)
	if err != nil {
		return err
	}
	tx.TxOuts = make([]TxOut, tx.TxOutCount)
	for i := 0; i < int(tx.TxOutCount); i++ {
		err = binary.Read(buf, binary.LittleEndian, &tx.TxOuts[i].Value)
		if err != nil {
			return err
		}
		err = binary.Read(buf, binary.LittleEndian, &tx.TxOuts[i].ScriptLength)
		if err != nil {
			return err
		}
		tx.TxOuts[i].PkScript = make([]byte, tx.TxOuts[i].ScriptLength)
		err = binary.Read(buf, binary.LittleEndian, &tx.TxOuts[i].PkScript)
		if err != nil {
			return err
		}
	}

	return nil
}

func (tx *RawTransaction) Deserialize(data []byte) error {
	buf := bytes.NewBuffer(data)

	err := binary.Read(buf, binary.LittleEndian, &tx.Version)
	if err != nil {
		return err
	}
	err = deserializeTxIn(buf, tx)
	if err != nil {
		return err
	}
	err = deserializeTxOut(buf, tx)
	if err != nil {
		return err
	}
	err = binary.Read(buf, binary.LittleEndian, &tx.LockTime)
	if err != nil {
		return err
	}

	return nil
}
