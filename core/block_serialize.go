package core

import (
	"bytes"
	"encoding/binary"
)

func (block *Block) Serialize() ([]byte, error) {
	var buf bytes.Buffer

	// Serialize BlockHeader
	err := binary.Write(&buf, binary.LittleEndian, block.Version)
	if err != nil {
		return nil, err
	}
	err = binary.Write(&buf, binary.LittleEndian, block.PrevBlock)
	if err != nil {
		return nil, err
	}
	err = binary.Write(&buf, binary.LittleEndian, block.MerkleRoot)
	if err != nil {
		return nil, err
	}
	err = binary.Write(&buf, binary.LittleEndian, block.Timestamp.Unix())
	if err != nil {
		return nil, err
	}
	err = binary.Write(&buf, binary.LittleEndian, block.Bits)
	if err != nil {
		return nil, err
	}
	err = binary.Write(&buf, binary.LittleEndian, block.Nonce)
	if err != nil {
		return nil, err
	}

	// Serialize TxCount
	err = binary.Write(&buf, binary.LittleEndian, block.TxCount)
	if err != nil {
		return nil, err
	}

	// Serialize Transactions
	for _, tx := range block.Txs {
		txData, err := tx.Serialize()
		if err != nil {
			return nil, err
		}
		_, err = buf.Write(txData)
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}
