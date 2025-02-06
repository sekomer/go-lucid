package block

import (
	"go-lucid/core/transaction"

	"gorm.io/gorm"
)

type BlockModel struct {
	gorm.Model
	BlockHeader
	Hash    []byte
	TxCount uint32
	Txs     []transaction.RawTransactionModel `gorm:"foreignKey:BlockID"`
}
