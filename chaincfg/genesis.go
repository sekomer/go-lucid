package chaincfg

import (
	"go-lucid/core"
	"time"
)

// GenesisBlock defines the genesis block of the block chain which serves as the
// public transaction ledger for the lumen network.
var GenesisBlock = core.Block{
	BlockHeader: core.BlockHeader{
		Version:   1,
		PrevBlock: core.Hash{
			// TODO: put magic here
		},
		MerkleRoot: core.Hash{
			// TODO: put magic here
		},
		Timestamp: time.Unix(0, 0),
		Bits:      0x00FFFFFF,
		Nonce:     0,
	},
	TxCount: 1,
	Txs:     []core.RawTransaction{
		// TODO: put magic here
	},
}
