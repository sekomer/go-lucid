package transaction

import (
	"go-lucid/core"

	"gorm.io/gorm"
)

type TxInModel struct {
	gorm.Model
	TransactionID   uint32 // Foreign key to associate with RawTransaction
	Coinbase        bool
	PreviousOutput  OutPoint `gorm:"embedded;embeddedPrefix:prev_"`
	SignatureScript []byte
	Sequence        uint32
}

type TxOutModel struct {
	gorm.Model
	TransactionID uint32 // Foreign key to associate with RawTransaction
	Value         core.Lumen
	PkScript      []byte
}

type RawTransactionModel struct {
	gorm.Model
	Hash       []byte `gorm:"index"`
	Version    int32
	TxInCount  uint32
	TxIns      []TxInModel `gorm:"foreignKey:TransactionID"` // One-to-Many relationship
	TxOutCount uint32
	TxOuts     []TxOutModel `gorm:"foreignKey:TransactionID"` // One-to-Many relationship
	BlockID    uint32       // Foreign key to associate with Block
	LockTime   int64
}
