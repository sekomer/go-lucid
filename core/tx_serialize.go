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
	for _, in := range tx.TxIn {
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
	for _, out := range tx.TxOut {
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
	tx.TxIn = make([]TxIn, tx.TxInCount)
	for i := 0; i < int(tx.TxInCount); i++ {
		err = binary.Read(buf, binary.LittleEndian, &tx.TxIn[i].Coinbase)
		if err != nil {
			return err
		}
		err = binary.Read(buf, binary.LittleEndian, &tx.TxIn[i].PreviousOutput.Hash)
		if err != nil {
			return err
		}
		err = binary.Read(buf, binary.LittleEndian, &tx.TxIn[i].PreviousOutput.Index)
		if err != nil {
			return err
		}
		err = binary.Read(buf, binary.LittleEndian, &tx.TxIn[i].ScriptLength)
		if err != nil {
			return err
		}
		tx.TxIn[i].SignatureScript = make([]byte, tx.TxIn[i].ScriptLength)
		err = binary.Read(buf, binary.LittleEndian, &tx.TxIn[i].SignatureScript)
		if err != nil {
			return err
		}
		err = binary.Read(buf, binary.LittleEndian, &tx.TxIn[i].Sequence)
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
	tx.TxOut = make([]TxOut, tx.TxOutCount)
	for i := 0; i < int(tx.TxOutCount); i++ {
		err = binary.Read(buf, binary.LittleEndian, &tx.TxOut[i].Value)
		if err != nil {
			return err
		}
		err = binary.Read(buf, binary.LittleEndian, &tx.TxOut[i].ScriptLength)
		if err != nil {
			return err
		}
		tx.TxOut[i].PkScript = make([]byte, tx.TxOut[i].ScriptLength)
		err = binary.Read(buf, binary.LittleEndian, &tx.TxOut[i].PkScript)
		if err != nil {
			return err
		}
	}

	return nil
}

func Deserialize(data []byte) (*RawTransaction, error) {
	buf := bytes.NewBuffer(data)
	tx := &RawTransaction{}

	err := binary.Read(buf, binary.LittleEndian, &tx.Version)
	if err != nil {
		return nil, err
	}
	err = deserializeTxIn(buf, tx)
	if err != nil {
		return nil, err
	}
	err = deserializeTxOut(buf, tx)
	if err != nil {
		return nil, err
	}
	err = binary.Read(buf, binary.LittleEndian, &tx.LockTime)
	if err != nil {
		return nil, err
	}

	return tx, nil
}
