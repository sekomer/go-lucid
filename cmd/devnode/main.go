package cmd

import (
	"context"
	"go-lucid/core/transaction"
	"go-lucid/node"
	tx_p2p "go-lucid/p2p/transaction"
	"go-lucid/rpc/ping"
	"go-lucid/service/health"
	"log"
	"time"

	gorpc "github.com/libp2p/go-libp2p-gorpc"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/protocol"
)

var rpcProtocolID = protocol.ID("/p2p/rpc/ping")

func StartRpcClient(client host.Host) *gorpc.Client {
	return gorpc.NewClient(client, rpcProtocolID)
}

func main(c *node.FullNodeConfig) {
	log.Println("dev node starting...")

	n := node.CreateHost(nil, c)
	n.InitPeers()
	defer n.Close()

	log.Printf("Hello World, hosts ID is %s\n", n.Host.ID())
	log.Printf("connection address of this node is: %s/p2p/%s\n", n.Host.Addrs()[0], n.Host.ID())

	pingService := ping.NewPingService(n.Host)
	err := n.Rpc.RegisterService(pingService, ping.ProtocolID)
	if err != nil {
		panic(err)
	}

	healthService := health.NewHealthService(n.Host)
	go func() {
		err := healthService.Start(context.Background())
		if err != nil {
			log.Println("error starting health service:", err)
		}
	}()

	ps, err := pubsub.NewGossipSub(context.Background(), n.Host)
	if err != nil {
		panic(err)
	}
	transactionService, err := tx_p2p.NewTransactionService(n.Host, ps)
	if err != nil {
		panic(err)
	}
	ch, err := transactionService.Subscribe(context.Background())
	if err != nil {
		panic(err)
	}
	go func() {
		for range time.Tick(1 * time.Second) {
			log.Println("broadcasting block... devnode")
			err := transactionService.Broadcast(context.Background(), transaction.RawTransactionModel{Version: 32})
			if err != nil {
				log.Println("error broadcasting block:", err)
			}
		}
	}()
	go func() {
		for msg := range ch {
			tx := transaction.RawTransactionModel{}
			err := tx.Deserialize(msg.Payload)
			if err != nil {
				log.Println("error deserializing block:", err)
				continue
			}
			log.Printf("from: %s, block: %+v\n", msg.From, tx)
		}
	}()

	select {}
}

func Start(c *node.FullNodeConfig) {
	main(c)
}
