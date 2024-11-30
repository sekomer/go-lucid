package ping

import "github.com/libp2p/go-libp2p/core/protocol"

// Service constants
const (
	Version     = "1.0.0"
	ServiceName = "PingService"
	ProtocolID  = protocol.ID("p2p/rpc/ping")
)

// PingArgs type
type PingArgs struct {
	Data []byte
}

// PingReply type
type PingReply struct {
	Data []byte
}
