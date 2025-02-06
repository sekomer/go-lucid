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

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/crypto"
	"golang.org/x/exp/rand"
)

func main(c *node.FullNodeConfig) {
	log.Println("bootnode starting...")

	priv, _, err := crypto.GenerateEd25519Key(rand.New(rand.NewSource(c.Node.Debug.Seed)))
	if err != nil {
		panic(err)
	}

	n := node.CreateHost(priv, c)
	defer n.Close()

	n.InitPeers()

	log.Printf("Hello World, hosts ID is %s\n", n.Host.ID())
	log.Printf("connection address of this node is: %s/p2p/%s\n", n.Host.Addrs()[0], n.Host.ID())

	pingService := ping.NewPingService(n.Host)
	err = n.Rpc.RegisterService(pingService, ping.ProtocolID)
	if err != nil {
		panic(err)
	}

	healthService := health.NewHealthService(n.Host)
	go func() {
		err := healthService.Start(context.Background())
		if err != nil {
			log.Printf("Error starting health service: %v", err)
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
			log.Println("broadcasting block... bootnode")
			err := transactionService.Broadcast(context.Background(), transaction.RawTransactionModel{
				Version: 32,
			})
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
