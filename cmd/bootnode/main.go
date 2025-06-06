package cmd

import (
	"context"
	"encoding/json"
	block_core "go-lucid/core/block"
	"go-lucid/core/transaction"
	"go-lucid/database"
	"go-lucid/node"
	tx_p2p "go-lucid/p2p/transaction"
	"go-lucid/rpc/block"
	"go-lucid/rpc/ping"
	"go-lucid/service/health"
	"log"
	"time"

	"github.com/libp2p/go-libp2p/core/crypto"
	"golang.org/x/exp/rand"
)

func main(c *node.FullNodeConfig) {
	log.Println("bootnode starting...")

	database.InitDB("/tmp/bootnode.db")

	priv, _, err := crypto.GenerateEd25519Key(rand.New(rand.NewSource(c.Node.Debug.Seed)))
	if err != nil {
		panic(err)
	}

	n := node.CreateHost(priv, c)
	defer n.Close()

	n.InitPeers()

	log.Printf("[lucid-go], hosts ID is %s\n", n.Host.ID())
	log.Printf("connection address of this node is: %s/p2p/%s\n", n.Host.Addrs()[0], n.Host.ID())

	pingService := ping.NewPingService(n.Host)
	err = n.Rpc.RegisterService(pingService, ping.ProtocolID)
	if err != nil {
		panic(err)
	}

	blockService := block.NewBlockService(n.Host)
	err = n.Rpc.RegisterService(blockService, block.ProtocolID)
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

	transactionService, err := tx_p2p.NewTransactionService(n.Host, n.PubSub)
	if err != nil {
		panic(err)
	}
	ch, err := transactionService.Subscribe(context.Background())
	if err != nil {
		panic(err)
	}
	go func() {
		block_number := uint32(0)
		for range time.Tick(3 * time.Second) {
			err := transactionService.Broadcast(context.Background(), transaction.RawTransaction{
				Version:   31,
				BlockID:   block_number,
				Hash:      []byte("bootnode hash"),
				TxInCount: 9999999,
			})
			if err != nil {
				log.Println("error broadcasting block:", err)
			}
			block_number++
		}
	}()
	go func() {
		for msg := range ch {
			tx := transaction.RawTransaction{}
			err := tx.Deserialize(msg.Payload)
			if err != nil {
				log.Println("error deserializing block:", err)
				continue
			}
		}
	}()

	go func() {
		for range time.Tick(10 * time.Second) {
			peers := n.Host.Network().Peers()
			log.Println("bootnode pubsub peers:", peers)
		}
	}()

	go func() {
		for {
			time.Sleep(3 * time.Second)
			blockClient := block.NewBlockClient(n.Host)

			blockRpcArgs := block.BlockRpcArgs{
				Method: "GetBlock",
				Args:   []any{1},
			}
			reply := block.BlockRpcReply{}

			if len(n.Host.Network().Peers()) == 0 {
				log.Println("no peers found")
				continue
			}

			err := blockClient.Call(context.Background(), n.Host.Network().Peers()[0], "GetBlock", &blockRpcArgs, &reply)
			if err != nil {
				log.Println("error calling get block rpc:", err)
				continue
			}

			if !reply.Success {
				log.Println("xo boot success:", reply.Success)
				log.Println("xo boot error:", reply.Error)
				log.Println("xo boot result:", reply.Result)
			}

			block := block_core.Block{}
			json.Unmarshal(reply.Result, &block)
			log.Println("init test: bootnode block:", block)
		}
	}()

	select {}
}

func Start(c *node.FullNodeConfig) {
	main(c)
}
