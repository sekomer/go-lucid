package block

import (
	"go-lucid/core/transaction"
	"time"
)

type BlockHeader struct {
	Version    int32
	Height     uint32
	PrevBlock  []byte
	MerkleRoot []byte
	Timestamp  time.Time
	Bits       uint32
	Nonce      uint32
}

type Block struct {
	BlockHeader
	Hash    []byte
	TxCount uint32
	Txs     []transaction.RawTransaction
}
