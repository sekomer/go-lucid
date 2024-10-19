package node

import (
	"context"
	"fmt"
	"go-lucid/rpc"
	"go-lucid/rpc/ping"
	"log"
	"time"

	"github.com/libp2p/go-libp2p"
	gorpc "github.com/libp2p/go-libp2p-gorpc"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/libp2p/go-libp2p/core/routing"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	libp2ptls "github.com/libp2p/go-libp2p/p2p/security/tls"
	"golang.org/x/exp/rand"
)

type Node struct {
	Host   host.Host
	Dht    *dht.IpfsDHT
	Rpc    rpc.RpcServer
	config *FullNodeConfig
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
	n.config = c

	if priv == nil {
		var err error
		priv, _, err = crypto.GenerateEd25519Key(rand.New(rand.NewSource(rand.Uint64())))
		if err != nil {
			panic(err)
		}
	}

	var idht *dht.IpfsDHT

	var port int
	switch c.Node.Type {
	default:
		panic("unknown node type")
	case DevNode:
		port = rand.Int()%1000 + 10000
	case BootNode:
		port = c.Node.Debug.Port
	case FullNode:
		port = c.Node.P2p.ListenPort
	}

	connmgr, err := connmgr.NewConnManager(
		c.Node.P2p.MinPeers,
		c.Node.P2p.MaxPeers,
		connmgr.WithGracePeriod(time.Duration(c.Node.P2p.GracePeriod)*time.Second),
	)
	if err != nil {
		panic(err)
	}

	h2, err := libp2p.New(
		libp2p.Identity(priv),
		libp2p.ListenAddrStrings(
			fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port),
			fmt.Sprintf("/ip4/0.0.0.0/udp/%d/quic", port),
		),
		libp2p.Security(libp2ptls.ID, libp2ptls.New),
		libp2p.Security(noise.ID, noise.New),
		libp2p.DefaultTransports,
		libp2p.ConnectionManager(connmgr),
		libp2p.NATPortMap(),
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			idht, err = dht.New(context.Background(), h, dht.Mode(dht.ModeAutoServer))
			return idht, err
		}),
		libp2p.EnableNATService(),
	)
	if err != nil {
		panic(err)
	}

	if n.config.Node.Debug.Enabled {
		// go peerstoreDebug(h2)
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

func peerstoreDebug(host host.Host) {
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		pxr := host.Peerstore().PeersWithAddrs()
		log.Println("peerstore len:", pxr.Len())
		peers := host.Network().Peers()
		log.Println("active peers:", len(peers))
		fmt.Print("\n\n")
	}
}

func (n *Node) InitPeers() {
	if n.config.Node.Type == "boot" {
		return
	}

	initialPeers := n.config.Node.Peers
	if n.config.Node.Debug.Enabled {
		initialPeers = []string{n.config.Node.Debug.Peer}
	}

	for _, p := range initialPeers {
		addrInfo, _ := peer.AddrInfoFromString(p)
		err := n.Host.Connect(context.Background(), *addrInfo)
		if err != nil {
			// TODO: handle gracefully
			panic(err)
		}
	}
}
