package node

import (
	"context"
	"fmt"
	"go-lucid/rpc"
	"go-lucid/rpc/ping"
	"log"
	"os"
	"time"

	"github.com/libp2p/go-libp2p"
	gorpc "github.com/libp2p/go-libp2p-gorpc"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/core/routing"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	libp2ptls "github.com/libp2p/go-libp2p/p2p/security/tls"
)

type Node struct {
	Host host.Host
	Rpc  rpc.RpcServer
}

func (n *Node) RegisterRpcServices() {
	for _, api := range n.Rpc.Api {
		rpcHost := gorpc.NewServer(n.Host, api.ProtocolId)
		err := rpcHost.Register(api.Service)
		if err != nil {
			panic(err)
		}

		n.Rpc.Log.Printf("Registered service %s", api.ProtocolId)
	}
}

func (n *Node) CreateHost(priv crypto.PrivKey, c *FullNodeConfig) {
	var port int
	var idht *dht.IpfsDHT

	if c.Node.Debug.Enabled {
		port = c.Node.Debug.Port
	} else {
		port = c.Node.P2p.ListenPort
	}

	connmgr, err := connmgr.NewConnManager(
		c.Node.P2p.MinPeers, // Lowwater
		c.Node.P2p.MaxPeers, // HighWater,
		connmgr.WithGracePeriod(time.Duration(c.Node.P2p.GracePeriod)*time.Second),
	)
	if err != nil {
		panic(err)
	}

	h2, err := libp2p.New(
		// Use the keypair we generated
		libp2p.Identity(priv),
		// Multiple listen addresses
		libp2p.ListenAddrStrings(
			fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port),      // regular tcp connections
			fmt.Sprintf("/ip4/0.0.0.0/udp/%d/quic", port), // a UDP endpoint for the QUIC transport
		),
		// support TLS connections
		libp2p.Security(libp2ptls.ID, libp2ptls.New),
		// support noise connection
		libp2p.Security(noise.ID, noise.New),
		// support any other default transports (TCP)
		libp2p.DefaultTransports,
		// Let's prevent our peer from having too many
		// connections by attaching a connection manager.
		libp2p.ConnectionManager(connmgr),
		// Attempt to open ports using uPNP for NATed hosts.
		libp2p.NATPortMap(),
		// Let this host use the DHT to find other hosts
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			idht, err = dht.New(context.Background(), h, dht.Mode(dht.ModeServer))
			return idht, err
		}),
		libp2p.EnableNATService(),
	)
	if err != nil {
		panic(err)
	}
	if true || os.Getenv("RESOURCE_DEBUG") == "true" {
		go resourceDebug(idht)
		go peerstoreDebug(h2, idht)
	}

	n.Host = h2
}

func (n *Node) Close() {
	n.Host.Close()
}

func (n *Node) StartPingRpc(protocolId protocol.ID, server *ping.PingService) {
	rpcHost := gorpc.NewServer(n.Host, protocolId)
	err := rpcHost.Register(server)
	if err != nil {
		panic(err)
	}
}

func resourceDebug(idht *dht.IpfsDHT) {
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		idht.Host().Network().ResourceManager().ViewSystem(func(rs network.ResourceScope) error {
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

func peerstoreDebug(host host.Host, idht *dht.IpfsDHT) {
	ticker := time.NewTicker(5 * time.Second)

	for range ticker.C {
		pxr := host.Peerstore().PeersWithAddrs()
		log.Println("peerstore len:", pxr.Len())
		peers := host.Network().Peers()
		log.Println("active peers:", len(peers))
		fmt.Print("\n\n")

		idht.RefreshRoutingTable()
	}
}
