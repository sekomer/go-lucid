package block

import (
	"github.com/libp2p/go-libp2p/core/protocol"
)

// Service constants
const (
	Version     = "1.0.0"
	ServiceName = "BlockService"
	ProtocolID  = protocol.ID("p2p/rpc/block")

	MethodBlock       = "GetBlock"
	MethodBlockHeader = "GetBlockHeader"
)

// PingArgs type
type BlockRpcArgs struct {
	Method string
	Args   []any
}

// PingReply type
type BlockRpcReply struct {
	Success bool
	Error   string
	Result  []byte
}
