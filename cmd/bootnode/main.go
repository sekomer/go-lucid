package cmd

import (
	"context"
	"fmt"
	"go-lucid/node"
	"log"
	"time"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"golang.org/x/exp/rand"

	"github.com/libp2p/go-libp2p/core/crypto"

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
	log.Println("bootnode starting...")

	var initialPeers []string = make([]string, 0)

	var err error
	var priv crypto.PrivKey

	priv, _, err = crypto.GenerateEd25519Key(rand.New(rand.NewSource(c.Node.Debug.Seed)))
	if err != nil {
		panic(err)
	}

	var node node.Node
	node.CreateHost(priv, c)
	defer node.Close()

	go StartRpcServer(node.Host)

	ps, err := pubsub.NewGossipSub(context.Background(), node.Host)
	if err != nil {
		panic(err)
	}

	th, err := ps.Join("test")
	if err != nil {
		panic(err)
	}

	// broadcast into topic every 5 seconds
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for range ticker.C {
			err := th.Publish(context.Background(), []byte("Hello World"))
			if err != nil {
				log.Println("Broadcast Error:", err)
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
			if msg.ReceivedFrom == node.Host.ID() {
				continue
			}
			if err != nil {
				log.Println("Read Error:", err)
			}

			log.Println("Message From", msg.ReceivedFrom)
			log.Println("Read Topic", string(msg.Data))
			fmt.Print('\n')
		}
	}()

	// connect to the initial peers
	for _, p := range initialPeers {
		addrInfo, _ := peer.AddrInfoFromString(p)
		err = node.Host.Connect(context.Background(), *addrInfo)
		if err != nil {
			panic(err)
		}
	}

	log.Printf("Hello World, hosts ID is %s\n", node.Host.ID())
	log.Printf("connection address of this node is: %s/p2p/%s\n", node.Host.Addrs()[0], node.Host.ID())

	// ping
	pingService := &ping.PingService{Host: node.Host}
	node.Host.SetStreamHandler(ping.ID, pingService.PingHandler)

	go func() {
		ticker := time.NewTicker(1 * time.Second)

		for range ticker.C {
			peers := node.Host.Network().Peers()
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

	node.Host.SetStreamHandler(ping.ID, func(s network.Stream) {
		defer s.Close()
		// log.Println("!!!!!!!!!!!!!ping handler")
	})

	rpcClient := StartRpcClient(node.Host)

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for range ticker.C {
			peers := node.Host.Network().Peers()
			if len(peers) < 1 {
				continue
			}

			peer := peers[rand.Intn(len(peers))]

			var reply PingReply
			var args PingArgs

			b := make([]byte, 32)
			_, err := rand.Read(b)
			if err != nil {
				panic(err)
			}

			args.Data = b

			rpcClient.Call(peer, "PingService", "Ping", args, &reply)
			log.Println("rpc call reply:", reply.Data)
		}
	}()

	select {}
}

func Start(c *node.FullNodeConfig) {
	main(c)
}
