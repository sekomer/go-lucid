package node

import (
	"log"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
)

func (n *Node) ResourceDebug() {
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		n.Dht.Host().Network().ResourceManager().ViewSystem(func(rs network.ResourceScope) error {
			log.Println("memory:", rs.Stat().Memory)
			log.Println("cons in:", rs.Stat().NumConnsInbound)
			log.Println("cons out:", rs.Stat().NumConnsOutbound)
			log.Println("strm in:", rs.Stat().NumStreamsInbound)
			log.Println("strm out:", rs.Stat().NumStreamsOutbound)
			log.Println("fds:", rs.Stat().NumFD)
			return nil
		})
	}
}
