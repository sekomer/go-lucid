package chaincfg

import (
	"bytes"
	"go-lucid/core"
	"go-lucid/core/block"
	"go-lucid/core/transaction"
	"go-lucid/util"
	"time"
)

// This is the genesis block of the Lucid protocol.
var GenesisBlock = block.Block{
	BlockHeader: block.BlockHeader{
		Height:     0,
		Version:    1,
		PrevBlock:  util.PadOrTrimTo32Bytes([]byte{0x42}),
		MerkleRoot: util.PadOrTrimTo32Bytes([]byte{0x42}),
		Timestamp:  time.Unix(42, 0),
		Bits:       0x42,
		Nonce:      0x42,
	},
	Hash:    util.HashToBytes("eb4aff2f304b7d9fad63d9bcf95f7067c32b2206d1278c92d7a385f2a74b56f3"),
	TxCount: 1,
	Txs: []transaction.RawTransaction{
		{
			Version:    1,
			TxInCount:  1,
			TxOutCount: 1,
			TxIns: []transaction.TxIn{
				{
					TransactionID: 0x42,
					Coinbase:      true,
					PreviousOutput: transaction.OutPoint{
						Hash:  util.PadOrTrimTo32Bytes(bytes.Repeat([]byte{0x42}, core.HASH_LEN)),
						Index: 0x42,
					},
					SignatureScript: util.PadOrTrimTo32Bytes(bytes.Repeat([]byte{0x42}, core.HASH_LEN)),
					Sequence:        0x42,
				},
			},
			TxOuts: []transaction.TxOut{
				{
					TransactionID: 0x42,
					Value:         0x42,
					PkScript:      util.PadOrTrimTo32Bytes(bytes.Repeat([]byte{0x42}, core.HASH_LEN)),
				},
			},
			LockTime: 0x42,
		},
	},
}
