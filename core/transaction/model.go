package transaction

import (
	"go-lucid/core"

	"gorm.io/gorm"
)

type TxInModel struct {
	gorm.Model
	TransactionID  uint32 // Foreign key to associate with RawTransaction
	Coinbase       bool
	PreviousOutput OutPoint `gorm:"embedded;embeddedPrefix:prev_"`
	ScriptSig      []byte
	Sequence       uint32
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

func (t TxInModel) ToDomain() TxIn {
	return TxIn{
		TransactionID:  t.TransactionID,
		Coinbase:       t.Coinbase,
		PreviousOutput: t.PreviousOutput,
		ScriptSig:      t.ScriptSig,
		Sequence:       t.Sequence,
	}
}

func (t TxOutModel) ToDomain() TxOut {
	return TxOut{
		TransactionID: t.TransactionID,
		Value:         t.Value,
		PkScript:      t.PkScript,
	}
}

func (t RawTransactionModel) ToDomain() RawTransaction {
	txIns := make([]TxIn, len(t.TxIns))
	for i, txIn := range t.TxIns {
		txIns[i] = txIn.ToDomain()
	}
	txOuts := make([]TxOut, len(t.TxOuts))
	for i, txOut := range t.TxOuts {
		txOuts[i] = txOut.ToDomain()
	}
	return RawTransaction{
		Hash:      t.Hash,
		Version:   t.Version,
		TxInCount: t.TxInCount,
		TxIns:     txIns,
		TxOuts:    txOuts,
		BlockID:   t.BlockID,
		LockTime:  t.LockTime,
	}
}
