package core

type Lumen uint64

const OneLumen = Lumen(1e8)

type OutPoint struct {
	Hash  Hash
	Index uint32
}

type TxIn struct {
	Coinbase        bool
	PreviousOutput  OutPoint
	ScriptLength    uint32
	SignatureScript []byte
	Sequence        uint32
}

type TxOut struct {
	Value        Lumen
	ScriptLength uint32
	PkScript     []byte
}

type RawTransaction struct {
	Version    int32
	TxInCount  uint32
	TxIn       []TxIn
	TxOutCount uint32
	TxOut      []TxOut
	LockTime   int64
}

func (tx *RawTransaction) AddTxIn(in TxIn) {
	tx.TxIn = append(tx.TxIn, in)
	tx.TxInCount++
}

func (tx *RawTransaction) AddTxOut(out TxOut) {
	tx.TxOut = append(tx.TxOut, out)
	tx.TxOutCount++
}
