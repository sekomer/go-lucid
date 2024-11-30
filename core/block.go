package core

import "time"

type Hash [32]byte

type BlockHeader struct {
	Version    int32
	PrevBlock  Hash
	MerkleRoot Hash
	Timestamp  time.Time
	Bits       uint32
	Nonce      uint32
}

type Block struct {
	BlockHeader
	TxCount uint32
	Txs     []RawTransaction
}
