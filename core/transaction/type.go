package transaction

type Lumen uint64

const OneLumen = Lumen(1e8)

type TxIn struct {
	TransactionID   uint32
	Coinbase        bool
	PreviousOutput  OutPoint
	SignatureScript []byte
	Sequence        uint32
}

type TxOut struct {
	TransactionID uint32
	Value         Lumen
	PkScript      []byte
}

type RawTransaction struct {
	Hash       []byte
	Version    int32
	TxInCount  uint32
	TxIns      []TxIn
	TxOutCount uint32
	TxOuts     []TxOut
	BlockID    uint32
	LockTime   int64
}

type OutPoint struct {
	Hash  []byte
	Index uint32
}
