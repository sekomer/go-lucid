package block

import (
	"go-lucid/core/transaction"
	"time"
)

type BlockHeader struct {
	Version    int32
	Height     uint32 `gorm:"unique;not null;index"`
	PrevBlock  []byte `gorm:"not null"`
	MerkleRoot []byte
	Timestamp  time.Time
	Bits       uint32
	Nonce      uint32
}

type Block struct {
	BlockHeader
	Hash    []byte `gorm:"unique;not null;index"`
	TxCount uint32
	Txs     []transaction.RawTransaction
}
