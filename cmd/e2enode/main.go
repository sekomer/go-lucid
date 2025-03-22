package cmd

import (
	"context"
	"fmt"
	"go-lucid/node"
	"log"
	"time"

	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	"golang.org/x/exp/rand"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/routing"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	libp2ptls "github.com/libp2p/go-libp2p/p2p/security/tls"

	gorpc "github.com/libp2p/go-libp2p-gorpc"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/protocol"
)

var rpcProtocolID = protocol.ID("/p2p/rpc/ping")

type PingArgs struct {
	Data []byte
}
type PingReply struct {
	Data []byte
}
type PingService struct{}

func (t *PingService) Ping(ctx context.Context, argType PingArgs, replyType *PingReply) error {
	log.Println("Received a Ping call")
	replyType.Data = argType.Data
	return nil
}

func StartRpcServer(host host.Host) {
	rpcHost := gorpc.NewServer(host, rpcProtocolID)

	svc := PingService{}
	err := rpcHost.Register(&svc)
	if err != nil {
		panic(err)
	}
}

func StartRpcClient(client host.Host) *gorpc.Client {
	return gorpc.NewClient(client, rpcProtocolID)
}

func main(c *node.FullNodeConfig) {
	log.Println("normal node starting...")

	port := c.Node.P2p.ListenPort

	priv, _, err := crypto.GenerateEd25519Key(rand.New(rand.NewSource(rand.Uint64())))
	if err != nil {
		panic(err)
	}

	var idht *dht.IpfsDHT

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
		libp2p.Security(noise.ID, noise.New),
		libp2p.DefaultTransports,
		libp2p.ConnectionManager(connmgr),
		libp2p.NATPortMap(),
		libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
			idht, err = dht.New(context.Background(), h, dht.Mode(dht.ModeAutoServer))
			return idht, err
		}),
		libp2p.EnableNATService(),
		libp2p.EnableAutoNATv2(),
	)

	if err != nil {
		panic(err)
	}
	defer h2.Close()

	log.Printf("Hello World, hosts ID is %s\n", h2.ID())
	log.Printf("connection address of this node is: %s/p2p/%s\n", h2.Addrs()[0], h2.ID())

	select {}
}

func Start(c *node.FullNodeConfig) {
	main(c)
}
