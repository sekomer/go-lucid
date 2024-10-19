package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"go-lucid/node"
	"log"
	"time"

	gorpc "github.com/libp2p/go-libp2p-gorpc"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/protocol"
	ping "github.com/libp2p/go-libp2p/p2p/protocol/ping"
	"golang.org/x/exp/rand"
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
	log.Println("dev node starting...")

	var node node.Node
	node.CreateHost(nil, c)
	node.InitPeers()
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

	log.Printf("connection address of this node is: %s/p2p/%s\n", node.Host.Addrs()[0], node.Host.ID())

	// broadcast into topic every 5 seconds
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for range ticker.C {
			msg := struct {
				Message string
			}{
				Message: "Hello, World!",
			}
			data, _ := json.Marshal(msg)
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
			if msg.ReceivedFrom == node.Host.ID() {
				continue
			}

			if err != nil {
				fmt.Println("[ Read Error ]", err)
			}
			fmt.Println("[ Message From ]", msg.ReceivedFrom)
			x := string(msg.Data)

			fmt.Println("[ COW MEAT ]", x)
			fmt.Print('\n')
		}
	}()

	log.Printf("Hello World, hosts ID is %s\n", node.Host.ID())
	log.Printf("connection address of this node is: %s/p2p/%s\n", node.Host.Addrs()[0], node.Host.ID())

	// ping
	pingService := &ping.PingService{Host: node.Host}
	node.Host.SetStreamHandler(ping.ID, pingService.PingHandler)

	// go func() {
	// 	ticker := time.NewTicker(1 * time.Second)

	// 	for range ticker.C {
	// 		peers := node.Host.Network().Peers()
	// 		if len(peers) > 0 {
	// 			peer := peers[rand.Intn(len(peers))]

	// 			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	// 			ch := pingService.Ping(ctx, peer)
	// 			<-ch

	// 			// timeout
	// 			if ctx.Err() == context.DeadlineExceeded {
	// 				log.Printf("ping %s timeout\n", peer)
	// 			} else {
	// 				log.Printf("ping %s success\n", peer)
	// 			}

	// 			cancel()
	// 		}
	// 	}
	// }()

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

			err = rpcClient.Call(peer, "PingService", "Ping", args, &reply)
			if err != nil {
				fmt.Println("rpc call error:", err)
			}
			fmt.Println("rpc call reply:", reply.Data)
		}
	}()

	select {}
}

func Start(c *node.FullNodeConfig) {
	main(c)
}
