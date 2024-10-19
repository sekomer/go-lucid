package cmd

import (
	"context"
	"fmt"
	"go-lucid/node"
	"log"
	"os"
	"time"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	"golang.org/x/exp/rand"

	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/routing"
	ping "github.com/libp2p/go-libp2p/p2p/protocol/ping"
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

	var port int
	if c.Node.Debug.Enabled {
		port = rand.Int()%1000 + 10000
	} else {
		port = c.Node.P2p.ListenPort
	}

	var initialPeers []string
	initialPeers = append(initialPeers, c.Node.Peers...)

	if c.Node.Debug.Enabled {
		initialPeers = make([]string, 1)
		initialPeers[0] = c.Node.Debug.Peer
	}

	log.Println("initial peers:", initialPeers)

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

	// go StartRpcServer(h2)

	ps, err := pubsub.NewGossipSub(context.Background(), h2)
	if err != nil {
		panic(err)
	}

	th, err := ps.Join("test")
	if err != nil {
		panic(err)
	}

	log.Printf("connection address of this node is: %s/p2p/%s\n", h2.Addrs()[0], h2.ID())

	if os.Getenv("RESOURCE_DEBUG") == "true" {
		go func() {
			ticker := time.NewTicker(5 * time.Second)
			for range ticker.C {
				idht.Host().Network().ResourceManager().ViewSystem(func(rs network.ResourceScope) error {
					fmt.Println("memory:", rs.Stat().Memory)
					fmt.Println("cons in:", rs.Stat().NumConnsInbound)
					fmt.Println("cons out:", rs.Stat().NumConnsOutbound)
					fmt.Println("strm in:", rs.Stat().NumStreamsInbound)
					fmt.Println("strm out:", rs.Stat().NumStreamsOutbound)
					fmt.Println("fds:", rs.Stat().NumFD)
					return nil
				})
			}
		}()
	}

	// broadcast into topic every 5 seconds
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for range ticker.C {
			data := []byte(fmt.Sprintf("hello world from %s", h2.ID()))
			err := th.Publish(context.Background(), data)
			if err != nil {
				fmt.Println("[ Broadcast Error ]", err)
			}
		}
	}()

	thSub, err := th.Subscribe()
	if err != nil {
		panic(err)
	}

	// read from topic
	go func() {
		for {
			msg, err := thSub.Next(context.Background())

			// skip self message
			if msg.ReceivedFrom == h2.ID() {
				continue
			}

			if err != nil {
				fmt.Println("[ Read Error ]", err)
			}
			fmt.Println("[ Message From ]", msg.ReceivedFrom)
			fmt.Println("[Cow Meat]", string(msg.Data))
			fmt.Print('\n')
		}
	}()

	// connect to the initial peers
	for _, p := range initialPeers {
		addrInfo, _ := peer.AddrInfoFromString(p)
		err = h2.Connect(context.Background(), *addrInfo)
		if err != nil {
			panic(err)
		}
	}

	log.Printf("Hello World, hosts ID is %s\n", h2.ID())
	log.Printf("connection address of this node is: %s/p2p/%s\n", h2.Addrs()[0], h2.ID())

	go func() {
		ticker := time.NewTicker(5 * time.Second)

		for range ticker.C {
			pxr := h2.Peerstore().PeersWithAddrs()
			fmt.Println("peerstore len:", pxr.Len())
			peers := h2.Network().Peers()
			fmt.Println("active peers", len(peers))

			fmt.Print("\n\n")

			// idht.RefreshRoutingTable()
		}
	}()

	// ping
	pingService := &ping.PingService{Host: h2}
	h2.SetStreamHandler(ping.ID, pingService.PingHandler)

	go func() {
		ticker := time.NewTicker(1 * time.Second)

		for range ticker.C {
			peers := h2.Network().Peers()
			if len(peers) > 0 {
				peer := peers[rand.Intn(len(peers))]

				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

				ch := pingService.Ping(ctx, peer)
				<-ch

				// timeout
				if ctx.Err() == context.DeadlineExceeded {
					log.Printf("ping %s timeout\n", peer)
				} else {
					// log.Printf("ping %s success\n", peer)
				}

				cancel()
			}
		}
	}()

	h2.SetStreamHandler(ping.ID, func(s network.Stream) {
		defer s.Close()
		// log.Println("!!!!!!!!!!!!!ping handler")
	})

	// rpcClient := StartRpcClient(h2)

	// go func() {
	// 	ticker := time.NewTicker(5 * time.Second)
	// 	for range ticker.C {
	// 		peers := h2.Network().Peers()
	// 		if len(peers) < 1 {
	// 			continue
	// 		}

	// 		peer := peers[rand.Intn(len(peers))]

	// 		var reply PingReply
	// 		var args PingArgs

	// 		b := make([]byte, 32)
	// 		_, err := rand.Read(b)
	// 		if err != nil {
	// 			panic(err)
	// 		}

	// 		args.Data = b

	// 		rpcClient.Call(peer, "PingService", "Ping", args, &reply)
	// 		fmt.Println("rpc call reply:", reply.Data)
	// 	}
	// }()

	select {}
}

func Start(c *node.FullNodeConfig) {
	main(c)
}
