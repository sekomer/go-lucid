package core

import "gorm.io/gorm"

type Lumen uint64

const OneLumen = Lumen(1e8)

type OutPoint struct {
	Hash  []byte `gorm:"index"` // Use string for compatibility with database fields
	Index uint32
}

type TxIn struct {
	gorm.Model
	TransactionID   uint32 // Foreign key to associate with RawTransaction
	Coinbase        bool
	PreviousOutput  OutPoint `gorm:"embedded;embeddedPrefix:prev_"`
	ScriptLength    uint32
	SignatureScript []byte
	Sequence        uint32
}

type TxOut struct {
	gorm.Model
	TransactionID uint32 // Foreign key to associate with RawTransaction
	Value         Lumen
	ScriptLength  uint32
	PkScript      []byte
}

type RawTransaction struct {
	gorm.Model
	Hash       []byte `gorm:"index"`
	Version    int32
	TxInCount  uint32
	TxOutCount uint32
	TxIns      []TxIn  `gorm:"foreignKey:TransactionID"` // One-to-Many relationship
	TxOuts     []TxOut `gorm:"foreignKey:TransactionID"` // One-to-Many relationship
	LockTime   int64
}

func (tx *RawTransaction) AddTxIn(in TxIn) {
	tx.TxIns = append(tx.TxIns, in)
}

func (tx *RawTransaction) AddTxOut(out TxOut) {
	tx.TxOuts = append(tx.TxOuts, out)
}
