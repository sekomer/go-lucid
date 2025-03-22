package block

import (
	"go-lucid/core/transaction"
	"time"

	"gorm.io/gorm"
)

type BlockHeaderModel struct {
	Version    int32
	Height     uint32 `gorm:"unique;not null;index"`
	PrevBlock  []byte `gorm:"not null"`
	MerkleRoot []byte
	Timestamp  time.Time
	Bits       uint32
	Nonce      uint32
}

type BlockModel struct {
	gorm.Model
	BlockHeaderModel
	Hash    []byte `gorm:"index"`
	TxCount uint32
	Txs     []transaction.RawTransactionModel `gorm:"foreignKey:BlockID"`
}

func (b BlockHeaderModel) ToDomain() BlockHeader {
	return BlockHeader{
		Version:    b.Version,
		Height:     b.Height,
		PrevBlock:  b.PrevBlock,
		MerkleRoot: b.MerkleRoot,
		Timestamp:  b.Timestamp,
		Bits:       b.Bits,
		Nonce:      b.Nonce,
	}
}

func (b BlockModel) ToDomain() Block {
	txs := make([]transaction.RawTransaction, len(b.Txs))
	for i, tx := range b.Txs {
		txs[i] = tx.ToDomain()
	}

	return Block{
		BlockHeader: b.BlockHeaderModel.ToDomain(),
		Hash:        b.Hash,
		TxCount:     b.TxCount,
		Txs:         txs,
	}
}
